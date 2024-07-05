package cache

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// CacheKey 用于构建缓存键。
//
// Usage:
//
//	type A struct {
//	    XXX string
//	    YYY int
//	}
//
// var ckb = MustNewCacheKey[A]("key:{xxx}:{yyy}", 10*time.Minute)
//
// ... ckb.build()
type CacheKey[ParamsTable any] struct {
	schema   string
	expireAt time.Duration
	extra    any
}

// MustNewCacheKey 创建一个新的 CacheKey 实例。
func MustNewCacheKey[ParamsTable any](schema string, expireAt time.Duration) *CacheKey[ParamsTable] {
	// 创建一个 ParamsTable 类型的零值实例
	var zeroValue ParamsTable

	ck := &CacheKey[ParamsTable]{
		schema:   schema,
		expireAt: expireAt,
	}
	_, err := ck.Build(zeroValue)
	if err != nil {
		panic(err)
	}
	return ck
}

// Build 方法中应用驼峰转蛇形
func (ckb *CacheKey[ParamsTable]) Build(params ParamsTable) (string, error) {
	var paramsMap map[string]any
	var err error

	// 检查 ParamsTable 是否为结构体或指向结构体的指针
	paramsType := reflect.TypeOf(params)
	if paramsType.Kind() == reflect.Struct || (paramsType.Kind() == reflect.Ptr && paramsType.Elem().Kind() == reflect.Struct) {
		// 如果是结构体或指向结构体的指针，使用现有的 structToMap 方法
		paramsMap, err = structToMap(params)
		if err != nil {
			return "", err
		}
	} else {
		// 如果不是结构体，直接使用一个占位符替换
		// 假设 params 是一个可以直接转换为字符串的类型
		paramsMap = map[string]any{"placeholder": fmt.Sprintf("%v", params)}
	}

	return replacePlaceholders(ckb.schema, paramsMap)
}

func (ckb *CacheKey[ParamsTable]) MustBuild(params ParamsTable) string {
	paramsMap, err := ckb.Build(params)
	if err != nil {
		panic(err)
	}
	return paramsMap
}

// SetExpireAt 设置缓存键的过期时间。
func (ckb *CacheKey[ParamsTable]) SetExpireAt(expireAt time.Duration) {
	ckb.expireAt = expireAt
}

// GetExpireAt 获取缓存键的过期时间。
func (ckb *CacheKey[ParamsTable]) GetExpireAt() time.Duration {
	return ckb.expireAt
}

// WithExtra 设置上下文，返回一个新的 CacheKey[ParamsTable] 实例。
func (ckb *CacheKey[ParamsTable]) WithExtra(extra any) *CacheKey[ParamsTable] {
	return &CacheKey[ParamsTable]{
		schema:   ckb.schema,
		expireAt: ckb.expireAt,
		extra:    extra,
	}
}

func replacePlaceholders(schema string, paramsMap map[string]any) (string, error) {
	var result strings.Builder
	var placeholderName strings.Builder
	inPlaceholder := false

	for _, char := range schema {
		if char == '{' {
			inPlaceholder = true
			placeholderName.Reset()
			continue
		} else if char == '}' {
			inPlaceholder = false
			fieldName := placeholderName.String()
			fieldValue, exists := paramsMap[fieldName]
			if !exists {
				return "", fmt.Errorf("field %s does not exist", fieldName)
			}
			result.WriteString(fmt.Sprintf("%v", fieldValue))
			continue
		}

		if inPlaceholder {
			placeholderName.WriteRune(char)
		} else {
			result.WriteRune(char)
		}
	}

	if inPlaceholder {
		return "", errors.New("malformed schema: unclosed placeholder")
	}

	return result.String(), nil
}

// camelToSnake 将驼峰命名转换为蛇形命名
func camelToSnake(name string) string {
	if name == "" {
		return ""
	}
	var result strings.Builder
	result.Grow(len(name) + 5) // 预留一些空间以减少扩容操作

	for i, char := range name {
		if unicode.IsUpper(char) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// structToMap 使用反射将结构体转换为映射
func structToMap(item any) (map[string]any, error) {
	itemVal := reflect.ValueOf(item)
	if itemVal.Kind() == reflect.Ptr {
		itemVal = itemVal.Elem()
	}

	if itemVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", itemVal.Kind())
	}

	itemType := itemVal.Type()
	result := make(map[string]any)

	for i := 0; i < itemVal.NumField(); i++ {
		field := itemType.Field(i)
		fieldVal := itemVal.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}
		fieldName := camelToSnake(field.Name)
		result[fieldName] = fieldVal.Interface()
	}

	return result, nil
}
