package tags_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/adjust/redismq"
	"github.com/alicebob/miniredis/v2"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/internal/utils/cache"
	"github.com/bagaking/memorianexus/pkg/tags"
	"github.com/khicago/got/util/typer"
	"github.com/stretchr/testify/assert"
)

// 定义 EntityType 枚举
type EntityType int

const (
	EntityTypeBlock EntityType = iota
	EntityTypePost
	EntityTypeComment
)

// MockMQ 是 Producer 和 Consumer 接口的模拟实现
type MockMQ struct {
	lst []string
}

func (m *MockMQ) Put(ctx context.Context, payload string) error {
	m.lst = append(m.lst, payload)
	return nil
}

func (m *MockMQ) Get(ctx context.Context) (*redismq.Package, error) {
	if len(m.lst) == 0 {
		return nil, nil
	}
	msg := m.lst[0]
	m.lst = m.lst[1:]
	return &redismq.Package{Payload: msg}, nil
}

func (m *MockMQ) MGet(ctx context.Context, count int) ([]*redismq.Package, error) {
	if len(m.lst) == 0 {
		return []*redismq.Package{}, nil
	}
	var ret []string
	if count >= len(m.lst) {
		ret, m.lst = m.lst, []string{}
	} else {
		ret, m.lst = m.lst[0:count], m.lst[count:]
	}
	return typer.SliceMap(ret, func(s string) *redismq.Package {
		return &redismq.Package{Payload: s}
	}), nil
}

func (m *MockMQ) Ack(ctx context.Context, pkg *redismq.Package) error {
	return nil
}

func (m *MockMQ) Fail(ctx context.Context, pkg *redismq.Package) error {
	m.lst = append(m.lst, pkg.Payload)
	return nil
}

func (m *MockMQ) Requeue(ctx context.Context, pkg *redismq.Package) error {
	m.lst = append(m.lst, pkg.Payload)
	return nil
}

// MockTagRepository 是 TagRepository 接口的模拟实现
type MockTagRepository struct {
	tagsByUser    map[utils.UInt64][]string
	usersByTag    map[string][]utils.UInt64
	tagsByEntity  map[utils.UInt64][]string
	entitiesByTag map[string]map[EntityType][]utils.UInt64
}

func NewMockTagRepository() *MockTagRepository {
	repo := &MockTagRepository{
		tagsByUser:    make(map[utils.UInt64][]string),
		usersByTag:    make(map[string][]utils.UInt64),
		tagsByEntity:  make(map[utils.UInt64][]string),
		entitiesByTag: make(map[string]map[EntityType][]utils.UInt64),
	}
	repo.tagsByUser[12345] = []string{"tag1", "tag2"}
	repo.usersByTag["tag1"] = []utils.UInt64{12345}
	repo.usersByTag["tag2"] = []utils.UInt64{12345}
	repo.entitiesByTag["tag1"] = map[EntityType][]utils.UInt64{EntityTypeBlock: {100001, 100002}}
	repo.entitiesByTag["tag2"] = map[EntityType][]utils.UInt64{EntityTypeBlock: {100001}}
	repo.tagsByEntity[100001] = []string{"tag1", "tag2"}
	repo.tagsByEntity[100002] = []string{"tag1"}
	return repo
}

func (m *MockTagRepository) GetTag(ctx context.Context, tag string) (*tags.Tag[EntityType], error) {
	return nil, nil
}

func (m *MockTagRepository) GetUsersByTag(ctx context.Context, tag string) ([]utils.UInt64, error) {
	return m.usersByTag[tag], nil
}

func (m *MockTagRepository) GetTagsByUser(ctx context.Context, userID utils.UInt64) ([]string, error) {
	return m.tagsByUser[userID], nil
}

func (m *MockTagRepository) GetTagsByEntity(ctx context.Context, entityID utils.UInt64) ([]string, error) {
	return m.tagsByEntity[entityID], nil
}

func (m *MockTagRepository) GetEntitiesByTag(ctx context.Context, userID utils.UInt64, tag string, entityType EntityType) ([]utils.UInt64, error) {
	if entities, ok := m.entitiesByTag[tag]; ok {
		return entities[entityType], nil
	}
	return nil, nil
}

func TestMain(m *testing.M) {
	redisServer, err := miniredis.Run()
	if err != nil {
		panic(fmt.Sprintf("Failed to start miniredis: %v", err))
	}

	defer redisServer.Close()
	fmt.Println("Redis server started at", redisServer.Addr())
	// 初始化缓存
	cache.Init(redisServer.Addr())

	// 运行测试
	code := m.Run()

	// 退出测试
	os.Exit(code)
}

// 测试 TagUpdateManager 的启动和停止
func TestTagUpdateManager_StartStop(t *testing.T) {
	mq := new(MockMQ)
	manager := tags.NewTagUpdateManager(mq, mq, func(ctx context.Context, message tags.TagUpdateMessage[EntityType]) error {
		return nil
	})

	ctx := context.TODO()

	// 测试启动
	manager.Start(ctx)
	assert.True(t, manager.IsRunning())
	assert.Equal(t, "running", manager.GetStatus())

	time.Sleep(time.Second * 2)

	// 测试停止
	manager.Stop(ctx)
	assert.False(t, manager.IsRunning())
	assert.Equal(t, "stopped", manager.GetStatus())
}

// 测试 TagUpdateManager 的 Put 方法
func TestTagUpdateManager_Put(t *testing.T) {
	mq := new(MockMQ)
	manager := tags.NewTagUpdateManager(mq, mq, func(ctx context.Context, message tags.TagUpdateMessage[EntityType]) error {
		return nil
	})

	ctx := context.TODO()
	message := tags.TagUpdateMessage[EntityType]{Action: tags.EventInvalidUser, UserID: 12345}
	err := manager.Put(ctx, message)
	assert.NoError(t, err)
	assert.NotEmpty(t, mq.lst)
}

// 测试 TagUpdateManager 的 HandlePackage 方法
func TestTagUpdateManager_HandlePackage(t *testing.T) {
	mq := new(MockMQ)
	manager := tags.NewTagUpdateManager(mq, mq, func(ctx context.Context, message tags.TagUpdateMessage[EntityType]) error {
		return nil
	})

	ctx := context.TODO()
	message := tags.TagUpdateMessage[EntityType]{Action: tags.EventInvalidUser, UserID: 12345}
	payload, _ := json.Marshal(message)
	pkg := &redismq.Package{Payload: string(payload)}

	err := manager.HandlePackage(ctx, mq, pkg)
	assert.NoError(t, err)
	assert.Empty(t, mq.lst)
}

// 测试 TagService 的 GetTagsByUser 方法
func TestTagService_GetTagsByUser(t *testing.T) {
	repo := NewMockTagRepository()
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.User2Tags.MustBuild(12345)
	tags, err := cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, tags)

	// 设置缓存
	tags, err = service.GetTagsByUser(ctx, 12345)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)

	time.Sleep(time.Second)

	// 验证缓存是否存在
	tags, err = cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)
}

// 测试 TagService 的 InvalidateUserCache 方法
func TestTagService_InvalidateUserCache(t *testing.T) {
	repo := NewMockTagRepository()
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.User2Tags.MustBuild(12345)
	tags, err := cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, tags)

	// 设置缓存
	tags, err = service.GetTagsByUser(ctx, 12345)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)

	// 验证缓存是否存在
	tags, err = cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)

	// 调用 InvalidateUserCache 方法
	err = service.InvalidateUserCache(ctx, 12345, true)
	assert.NoError(t, err)

	// 验证缓存是否被清除
	tags, err = cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, tags)
}

// 测试 TagService 的 InvalidateTagCache 方法
func TestTagService_InvalidateTagCache(t *testing.T) {
	repo := NewMockTagRepository()
	repo.usersByTag["tag1"] = []utils.UInt64{12345, 67890}
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.Tag2Users.MustBuild("tag1")
	users, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, users)

	// 设置缓存
	users, err = service.GetUsersByTag(ctx, "tag1")
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{12345, 67890}, users)

	// 验证缓存是否存在
	users, err = cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{12345, 67890}, users)

	// 调用 InvalidateTagCache 方法
	err = service.InvalidateTagCache(ctx, "tag1", true)
	assert.NoError(t, err)

	// 验证缓存是否被清除
	users, err = cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, users)
}

// 测试 TagService 的 InvalidateUserTagCache 方法
func TestTagService_InvalidateUserTagCache(t *testing.T) {
	repo := NewMockTagRepository()
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.Entities.MustBuild(tags.ParamUserTagType[EntityType]{UserID: 12345, Tag: "tag1"})
	entities, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, entities)

	// 设置缓存
	entities, err = service.GetEntities(ctx, 12345, "tag1", nil)
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{100001, 100002}, entities)

	// 验证缓存是否存在
	entities, err = cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{100001, 100002}, entities)

	// 调用 InvalidateUserTagCache 方法
	err = service.InvalidateUserTagCache(ctx, 12345, "tag1", true)
	assert.NoError(t, err)

	// 验证缓存是否被清除
	entities, err = cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, entities)
}

// 测试 TagService 的 InvalidateEntityCache 方法
func TestTagService_InvalidateEntityCache(t *testing.T) {
	repo := NewMockTagRepository()
	repo.tagsByEntity[12345] = []string{"tag1", "tag2"}
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.Entity2Tags.MustBuild(12345)
	tags, err := cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, tags)

	// 设置缓存
	tags, err = service.GetTagsOfEntity(ctx, 12345)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)

	// 验证缓存是否存在
	tags, err = cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)

	// 调用 InvalidateEntityCache 方法
	err = service.InvalidateEntityCache(ctx, 12345, true)
	assert.NoError(t, err)

	// 验证缓存是否被清除
	tags, err = cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, tags)
}

// 测试 TagService 的 GetUsersByTag 方法
func TestTagService_GetUsersByTag(t *testing.T) {
	repo := NewMockTagRepository()
	repo.usersByTag["tag1"] = []utils.UInt64{12345, 67890}
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.Tag2Users.MustBuild("tag1")
	users, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, users)

	// 设置缓存
	users, err = service.GetUsersByTag(ctx, "tag1")
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{12345, 67890}, users)

	// 验证缓存是否存在
	users, err = cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{12345, 67890}, users)
}

// 测试 TagService 的 GetEntities 方法
func TestTagService_GetEntities(t *testing.T) {
	repo := NewMockTagRepository()
	repo.entitiesByTag["tag1"] = map[EntityType][]utils.UInt64{EntityTypeBlock: {100001, 100002}, EntityTypePost: {200001}}
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.Entities.MustBuild(tags.ParamUserTagType[EntityType]{UserID: 12345, Tag: "tag1"})
	entities, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, entities)

	// 设置缓存
	entities, err = service.GetEntities(ctx, 12345, "tag1", nil)
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{100001, 100002, 200001}, entities)

	// 验证缓存是否存在
	entities, err = cache.SET().GetAllUInt64s(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []utils.UInt64{100001, 100002, 200001}, entities)
}

// 测试 TagService 的 GetTagsOfEntity 方法
func TestTagService_GetTagsOfEntity(t *testing.T) {
	repo := NewMockTagRepository()
	repo.tagsByEntity[12345] = []string{"tag1", "tag2"}
	mq := new(MockMQ)
	service := tags.NewTagService(repo, []EntityType{EntityTypeBlock, EntityTypePost}, mq, mq)

	ctx := context.TODO()

	// 验证缓存不存在
	cacheKey := service.Schemas.Entity2Tags.MustBuild(12345)
	tags, err := cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Empty(t, tags)

	// 设置缓存
	tags, err = service.GetTagsOfEntity(ctx, 12345)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)

	// 验证缓存是否存在
	tags, err = cache.SET().GetAll(ctx, cacheKey)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)
}
