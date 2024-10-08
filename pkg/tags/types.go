package tags

import (
	"context"

	"github.com/bagaking/memorianexus/internal/utils"
)

type (
	// TagRepository defines the interface for tag repository.
	TagRepository[EntityType any] interface {
		GetTag(ctx context.Context, tag string) (*Tag[EntityType], error)

		GetUsersByTag(ctx context.Context, tag string) ([]utils.UInt64, error)
		GetTagsByUser(ctx context.Context, userID utils.UInt64) ([]string, error)

		GetTagsByEntity(ctx context.Context, entityID utils.UInt64) ([]string, error)
		GetEntitiesByTag(ctx context.Context, userID utils.UInt64, tag string, entityType EntityType) ([]utils.UInt64, error)
	}

	// Producer defines the interface for a message producer.
	Producer interface {
		Put(ctx context.Context, payload string) error
	}

	// Consumer defines the interface for a message consumer.
	Consumer[MessageType any] interface {
		Get(ctx context.Context) (MessageType, error)
		MGet(ctx context.Context, count int) ([]MessageType, error)
		GetUnacked(ctx context.Context) (MessageType, error)

		Ack(ctx context.Context, pkg MessageType) error
		Fail(ctx context.Context, pkg MessageType) error
		Requeue(ctx context.Context, pkg MessageType) error
	}
)
