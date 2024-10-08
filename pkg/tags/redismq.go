package tags

import (
	"context"

	"github.com/adjust/redismq"
	"github.com/khicago/irr"
)

// RedisMQProducer is a RedisMQ implementation of the Producer interface.
type RedisMQProducer struct {
	queue *redismq.Queue
}

var (
	_ Producer                   = (*RedisMQProducer)(nil)
	_ Consumer[*redismq.Package] = (*RedisMQConsumer)(nil)
)

// NewRedisMQProducer creates a new RedisMQProducer.
func NewRedisMQProducer(queue *redismq.Queue) *RedisMQProducer {
	return &RedisMQProducer{queue: queue}
}

// Put puts a message into the queue.
func (p *RedisMQProducer) Put(ctx context.Context, payload string) error {
	return p.queue.Put(payload)
}

// RedisMQConsumer is a RedisMQ implementation of the Consumer interface.
type RedisMQConsumer struct {
	consumer *redismq.Consumer
}

// NewRedisMQConsumer creates a new RedisMQConsumer.
func NewRedisMQConsumer(queue *redismq.Queue, name string) (*RedisMQConsumer, error) {
	consumer, err := queue.AddConsumer(name)
	if err != nil {
		return nil, irr.Wrap(err, "failed to add consumer")
	}
	return &RedisMQConsumer{consumer: consumer}, nil
}

// Get gets a single message from the queue.
func (c *RedisMQConsumer) Get(ctx context.Context) (*redismq.Package, error) {
	pkg, err := c.consumer.Get()
	if err != nil {
		return nil, irr.Wrap(err, "failed to get message")
	}
	return pkg, nil
}

// MGet gets multiple pkgs from the queue.
func (c *RedisMQConsumer) MGet(ctx context.Context, count int) ([]*redismq.Package, error) {
	packages, err := c.consumer.MultiGet(count)
	if err != nil {
		return nil, irr.Wrap(err, "failed to get pkgs")
	}
	return packages, nil
}

// Ack acknowledges a pkg.
func (c *RedisMQConsumer) Ack(ctx context.Context, pkg *redismq.Package) error {
	return pkg.Ack()
}

// Fail marks a pkg as failed.
func (c *RedisMQConsumer) Fail(ctx context.Context, pkg *redismq.Package) error {
	return pkg.Fail()
}

// Requeue requeues a pkg.
func (c *RedisMQConsumer) Requeue(ctx context.Context, pkg *redismq.Package) error {
	return pkg.Requeue()
}

// GetUnacked gets an unacked package from the queue.
func (c *RedisMQConsumer) GetUnacked(ctx context.Context) (*redismq.Package, error) {
	unacked, err := c.consumer.GetUnacked()
	if err != nil {
		return nil, irr.Wrap(err, "获取未确认的包失败")
	}
	return unacked, nil
}
