package tags

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/khicago/got/util/procast"

	"github.com/adjust/redismq"
	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"
)

// TagUpdateManager manages the state and timing of tag updates.
type TagUpdateManager[EntityType any] struct {
	mu      sync.Mutex
	running bool
	status  string

	producer Producer
	consumer Consumer[*redismq.Package]

	runningCtx context.Context
	cancel     context.CancelFunc

	fnHandleTagUpdate func(ctx context.Context, message TagUpdateMessage[EntityType]) error
}

// NewTagUpdateManager creates a new TagUpdateManager.
func NewTagUpdateManager[EntityType any](
	producer Producer, consumer Consumer[*redismq.Package],
	fnHandleTagUpdate func(ctx context.Context, message TagUpdateMessage[EntityType]) error,
) *TagUpdateManager[EntityType] {
	if producer == nil {
		panic("producer cannot be nil")
	}
	if consumer == nil {
		panic("consumer cannot be nil")
	}
	if fnHandleTagUpdate == nil {
		panic("fnHandleTagUpdate cannot be nil")
	}

	return &TagUpdateManager[EntityType]{
		producer:          producer,
		consumer:          consumer,
		fnHandleTagUpdate: fnHandleTagUpdate,
	}
}

func (m *TagUpdateManager[EntityType]) Start(ctx context.Context) *TagUpdateManager[EntityType] {
	m.mu.Lock()
	defer m.mu.Unlock()
	log, ctx := wlog.ByCtxAndCache(ctx, "TagUpdateManager")
	if m.running {
		log.Warnf("Manager is already running, do nothing")
		return m
	}
	m.running = true
	m.runningCtx, m.cancel = context.WithCancel(ctx)
	m.status = "running"

	go m.run(m.runningCtx)
	log.Infof("Manager started")
	return m
}

func (m *TagUpdateManager[EntityType]) Stop(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	log := wlog.ByCtx(ctx, "TagUpdateManager")
	if !m.running {
		log.Warnf("Manager is not running, do nothing")
		return
	}

	if m.runningCtx != nil {
		m.cancel()
		m.running = false
		m.status = "stopped"

		log.Infof("Manager stopped")
		m.runningCtx = nil
	}
}

func (m *TagUpdateManager[EntityType]) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func (m *TagUpdateManager[EntityType]) GetStatus() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

func (m *TagUpdateManager[EntityType]) Put(ctx context.Context, message TagUpdateMessage[EntityType]) error {
	log := wlog.ByCtx(ctx, "TagUpdateManager.Put")
	data, err := json.Marshal(message)
	if err != nil {
		return irr.Wrap(err, "failed to marshal tag update message").LogError(log)
	}

	if err = m.producer.Put(ctx, string(data)); err != nil {
		return irr.Wrap(err, "failed to enqueue tag update message").LogError(log)
	}

	log.Infof("Successfully enqueued tag update message")
	return nil
}

// run starts the tag update worker to process messages from the queue.
func (m *TagUpdateManager[EntityType]) run(ctx context.Context) {
	log, ctx := wlog.ByCtxAndCache(ctx, "TagUpdateWorker", uuid.NewString())
	defer procast.Recover(func(err error) {
		log.Errorf("recovered from panic: %v", err)
	})

	var (
		packages []*redismq.Package
		err      error
		retries  int

		unackedRetryInterval = time.Second
	)

	for {
		select {
		case <-ctx.Done():
			log.Infof("Shutting down tag update worker")
			m.Stop(ctx)
			return
		default:
			// 处理所有未确认的包
			hasUnacked := true
			for hasUnacked {
				var err error
				hasUnacked, err = m.handleUnackedPackages(ctx)
				if err != nil {
					log.Errorf("处理未确认的包时出错: %v", err)
					time.Sleep(unackedRetryInterval)
				}
			}

			// Fetch a batch of messages from the queue
			packages, err = m.consumer.MGet(ctx, 10) // 获取最多10个消息
			if err != nil {
				retries++
				if retries > MaxRetryAttempts {
					log.Errorf("Failed to get messages from queue，已达到最大重试次数: %v", err)
					time.Sleep(time.Second * 5)
					retries = 0
				} else {
					log.Warnf("Failed to get messages from queue，正在重试 (%d/%d): %v", retries, MaxRetryAttempts, err)
					time.Sleep(time.Second)
				}
				continue
			}
			retries = 0

			// 如果没有获取到消息，等待一段时间后继续
			if len(packages) == 0 {
				time.Sleep(time.Second)
				continue
			}

			// 处理获取到的所有消息
			for _, pkg := range packages {
				if err = m.HandlePackage(ctx, m.consumer, pkg); err != nil {
					log.Errorf("Handle message failed: %v", err)
				}
			}
		}
	}
}

func (m *TagUpdateManager[EntityType]) handleUnackedPackages(ctx context.Context) (bool, error) {
	log := wlog.ByCtx(ctx, "handleUnackedPackages")

	// 获取未确认的包
	unackedPackage, err := m.consumer.GetUnacked(ctx)
	if err != nil {
		// 如果是因为没有未确认的包而返回错误，我们不认为这是一个真正的错误
		if err.Error() == "no unacked Packages found" {
			return false, nil
		}
		return false, irr.Wrap(err, "获取未确认的包失败").LogError(log)
	}

	// 如果没有未确认的包，直接返回
	if unackedPackage == nil {
		return false, nil
	}

	// 处理未确认的包
	if err := m.HandlePackage(ctx, m.consumer, unackedPackage); err != nil {
		log.Errorf("处理未确认的包失败: %v", err)
		// 考虑是否需要重新入队或标记为失败
		if m.shouldRetry(err) {
			if err := m.consumer.Requeue(ctx, unackedPackage); err != nil {
				log.Errorf("重新入队未确认的包失败: %v", err)
			}
		} else {
			if err := m.consumer.Fail(ctx, unackedPackage); err != nil {
				log.Errorf("标记未确认的包为失败状态时出错: %v", err)
			}
		}
		return true, err
	}

	return true, nil
}

func (m *TagUpdateManager[EntityType]) HandlePackage(ctx context.Context, consumer Consumer[*redismq.Package], pkg *redismq.Package) error {
	log := wlog.ByCtx(ctx, "HandlePackage")

	// 添加空值检查
	if consumer == nil {
		return irr.Error("consumer is nil").LogError(log)
	}
	if pkg == nil {
		return irr.Error("package is nil").LogError(log)
	}

	var message TagUpdateMessage[EntityType]
	if err := json.Unmarshal([]byte(pkg.Payload), &message); err != nil {
		if err = consumer.Fail(ctx, pkg); err != nil {
			return irr.Wrap(err, "failed to acknowledge message failed").LogError(log)
		}
		return irr.Wrap(err, "failed to unmarshal message")
	}

	// Process the message
	if err := m.fnHandleTagUpdate(ctx, message); err != nil {
		log.Errorf("Failed to handle message: %v", err)
		// 根据错误类型决定是否重试或标记为失败
		if m.shouldRetry(err) {
			log.Warnf("Requeuing message due to error: %v", err)
			if err = consumer.Requeue(ctx, pkg); err != nil {
				return irr.Wrap(err, "failed to requeue message").LogError(log)
			}
		} else {
			if err = consumer.Fail(ctx, pkg); err != nil {
				return irr.Wrap(err, "failed to acknowledge message failed").LogError(log)
			}
			return irr.Wrap(err, "failed to handle message")
		}
	}

	// Acknowledge the message
	if err := consumer.Ack(ctx, pkg); err != nil {
		return irr.Wrap(err, "failed to acknowledge message").LogError(log)
	}
	return nil
}

// shouldRetry 判断是否应该重试处理消息
func (m *TagUpdateManager[EntityType]) shouldRetry(err error) bool {
	// 根据错误类型或内容决定是���重试
	// 这里可以根据具体需求进行实现
	// 例如：如果错误是临时性的，可以返回 true；否则返回 false
	return true // 示例：默认所有错误都重试
}
