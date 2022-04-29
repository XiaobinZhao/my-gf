package token

import (
	"context"
	"myapp/internal/errorCode"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
)

// setCache 设置缓存
func (m *MyToken) setCache(ctx context.Context, cacheKey string, userCache *MyCacheToken) error {
	switch m.CacheMode {
	case CacheModeCache:
		err := gcache.Set(ctx, cacheKey, userCache, gconv.Duration(m.Timeout)*time.Millisecond)
		if err != nil {
			return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
	case CacheModeRedis:
		cacheValueJson, err := gjson.Encode(userCache)
		if err != nil {
			return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
		_, err = g.Redis().Do(ctx, "SETEX", cacheKey, m.Timeout/1000, cacheValueJson) // SETEX 单位为秒
		if err != nil {
			return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
	default:
		return errorCode.NewMyErr(ctx, errorCode.NotSupportedCacheModeError, m.CacheMode)
	}

	return nil
}

// getCache 获取缓存
func (m *MyToken) getCache(ctx context.Context, cacheKey string) (myCacheToken *MyCacheToken, err error) {
	switch m.CacheMode {
	case CacheModeCache:
		cacheValue, err := gcache.Get(ctx, cacheKey)
		if err != nil {
			return nil, errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
		if cacheValue.IsNil() {
			return nil, errorCode.NewMyErr(ctx, errorCode.Unauthorized)
		}
		myCacheToken = &MyCacheToken{}
		err = gconv.Struct(cacheValue, myCacheToken)
		if err != nil {
			return nil, errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
	case CacheModeRedis:
		userCacheJson, err := g.Redis().Do(ctx, "GET", cacheKey)
		if err != nil {
			if err != nil {
				return nil, errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
			}
		}
		if userCacheJson.IsNil() {
			return nil, errorCode.NewMyErr(ctx, errorCode.Unauthorized)
		}

		err = gconv.Struct(userCacheJson, myCacheToken)
		if err != nil {
			return nil, errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
	default:
		return nil, errorCode.NewMyErr(ctx, errorCode.NotSupportedCacheModeError, m.CacheMode)
	}

	return myCacheToken, nil
}

// removeCache 删除缓存
func (m *MyToken) RemoveCache(ctx context.Context, cacheKey string) error {
	switch m.CacheMode {
	case CacheModeCache:
		_, err := gcache.Remove(ctx, cacheKey)
		if err != nil {
			return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
	case CacheModeRedis:
		_, err := g.Redis().Do(ctx, "DEL", cacheKey)
		if err != nil {
			return errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
		}
	default:
		return errorCode.NewMyErr(ctx, errorCode.NotSupportedCacheModeError, m.CacheMode)
	}

	return nil
}
