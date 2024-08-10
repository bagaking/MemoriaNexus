package tags

import (
	"context"
	"encoding/json"
	"sync"
	"time"

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
	data, err := json.Marshal(message)
	if err != nil {
		return irr.Wrap(err, "failed to marshal tag update message")
	}

	if err = m.producer.Put(ctx, string(data)); err != nil {
		return irr.Wrap(err, "failed to enqueue tag update message")
	}

	return nil
}

// run starts the tag update worker to process messages from the queue.
func (m *TagUpdateManager[EntityType]) run(ctx context.Context) {
	log := wlog.ByCtx(ctx, "TagUpdateWorker")

	var (
		packages []*redismq.Package
		err      error
	)

	for {
		select {
		case <-ctx.Done():
			log.Infof("Shutting down tag update worker")
			m.Stop(ctx)
			return
		default:
			// Fetch a batch of messages from the queue
			packages, err = m.consumer.MGet(ctx, 10)
			if err != nil {
				log.Errorf("Failed to get messages from queue: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, pkg := range packages {
				if err = m.HandlePackage(ctx, m.consumer, pkg); err != nil {
					log.Errorf("Handle message failed: %v", err)
				}
			}
		}
	}
}

func (m *TagUpdateManager[EntityType]) HandlePackage(ctx context.Context, consumer Consumer[*redismq.Package], pkg *redismq.Package) error {
	log := wlog.ByCtx(ctx, "HandlePackage")
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
	// 根据错误类型或内容决定是否重试
	// 这里可以根据具体需求进行实现
	// 例如：如果错误是临时性的，可以返回 true；否则返回 false
	return true // 示例：默认所有错误都重试
}
