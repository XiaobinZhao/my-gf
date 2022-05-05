package token

import (
	"context"
	"myapp/internal/errorCode"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/container/gmap"

	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	// instances is the instances map for management
	// for multiple Token instance by name.
	instances = gmap.NewStrAnyMap(true)
)

// Instance returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Instance(name ...string) *MyToken {
	key := "default"
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		timeout := g.Cfg().MustGet(context.Background(), "token.timeout", CacheTimeout).Int()
		cacheMode := g.Cfg().MustGet(context.Background(), "token.cacheMode", CacheModeCache).Int()
		token := &MyToken{
			Timeout:   timeout,
			CacheMode: cacheMode,
		}
		return token
	}).(*MyToken)
}

type MyToken struct {
	Timeout   int // 超时时间 （毫秒）
	CacheMode int // 缓存类型：内存:1，redis:2
}

type MyRequestToken struct {
	UserKey string
	Uuid    string
	Token   string
}

type MyCacheToken struct {
	Token         string
	UserKey       string
	Uuid          string
	UserData      interface{}
	CreatedAt     int64 // Token 生成的时间
	NextFreshTime int64 // 下次token刷新时间, =0，一次性token
}

// 获取请求header中的token
// 只支持http header中，Authorization: Bearer 类型
func GetRequestToken(r *ghttp.Request) (token string, err error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			return "", errorCode.NewMyErr(r.Context(), errorCode.AuthHeaderInvalidError, authHeader)
		}
		if parts[0] != "Bearer" {
			return "", errorCode.NewMyErr(r.Context(), errorCode.AuthHeaderInvalidError, authHeader)
		}
		return parts[1], nil

	}
	return "", nil
}

// 验证Token
func (m *MyToken) ValidToken(ctx context.Context, token string) (myCacheToken *MyCacheToken, err error) {
	if token == "" {
		return nil, errorCode.NewMyErr(ctx, errorCode.Unauthorized)
	}

	decryptToken, err := m.decryptToken(ctx, token)
	if err != nil {
		return nil, errorCode.MyWrapCode(ctx, errorCode.AuthorizedFailed, err)
	}

	userCache, err := m.getAndFreshCacheToken(ctx, decryptToken.UserKey)
	if err != nil {
		return nil, errorCode.MyWrapCode(ctx, errorCode.AuthorizedFailed, err)
	}

	if decryptToken.Uuid != userCache.Uuid {
		return nil, errorCode.NewMyErr(ctx, errorCode.AuthorizedFailed)
	}

	return userCache, nil
}

/**
* @description 加密生成token.
               token的生成规则是base64(gaes.Encrypt(base64(userKey)+TokenDelimiter+uuid)): 其中TokenDelimiter默认为_;
               为什么要base64(userKey)，因为可能userKey包含_; 标准base64是使用 `数字`+`大小写字母`+`/`+`+`以及`=`组成
               解释：为什么还要对token进行加解密？答：加密因为token携带了userKey信息，且便于过滤掉不合法token；
* @param userKey 用户的标识，一般使用用户名称或者用户的uuid
* @param uuid 可以使用外部提供的uuid，如果为空，会重新生成
**/
func (m *MyToken) EncryptToken(ctx context.Context, userKey string, uuid string) (*MyRequestToken, error) {
	if userKey == "" {
		return nil, errorCode.NewMyErr(ctx, errorCode.TokenKeyEmpty)
	}

	if uuid == "" {
		// 重新生成uuid, 使用UUID
		uuid = guid.S()
	}

	tokenStr := gbase64.EncodeToString([]byte(userKey)) + TokenDelimiter + uuid

	token, err := gaes.Encrypt([]byte(tokenStr), []byte(EncryptKey))
	if err != nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
	}
	return &MyRequestToken{
		UserKey: userKey,
		Uuid:    uuid,
		Token:   gbase64.EncodeToString(token),
	}, nil
}

/**
* @description 解密token。token的生成规则是base64(gaes.Encrypt(base64(userKey)+TokenDelimiter+uuid))

* @param userKey 用户的标识，一般使用用户名称或者用户的uuid
* @param uuid 可以使用外部提供的uuid，如果为空，会重新生成
**/
func (m *MyToken) decryptToken(ctx context.Context, token string) (tokenDecrypted *MyRequestToken, err error) {
	if token == "" {
		return nil, errorCode.NewMyErr(ctx, errorCode.TokenEmpty)
	}
	token64, err := gbase64.Decode([]byte(token))
	if err != nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.TokenInvalidError, token, err)
	}
	decryptTokenStr, err := gaes.Decrypt(token64, []byte(EncryptKey))
	if err != nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.TokenInvalidError, token, err)
	}
	tokenArray := gstr.Split(string(decryptTokenStr), TokenDelimiter)
	if len(tokenArray) < 2 {
		return nil, errorCode.NewMyErr(ctx, errorCode.TokenInvalidError, token, err)
	}
	userKey, err := gbase64.Decode([]byte(tokenArray[0]))
	if err != nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.TokenInvalidError, token, err)
	}

	return &MyRequestToken{string(userKey), tokenArray[1], token}, nil
}

func (m *MyToken) getAndFreshCacheToken(ctx context.Context, userKey string) (*MyCacheToken, error) {
	cacheKey := CacheKeyPrefix + userKey

	cacheToken, err := m.getCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	nowTime := gtime.Now().TimestampMilli()

	// 需要进行缓存超时时间刷新
	// cacheToken.NextFreshTime == 0, 表明是一个一次性的token
	if gconv.Int64(cacheToken.NextFreshTime) == 0 || nowTime > gconv.Int64(cacheToken.NextFreshTime) {
		cacheToken.CreatedAt = gtime.Now().TimestampMilli()
		cacheToken.NextFreshTime = gtime.Now().TimestampMilli() + gconv.Int64(CacheMaxRefresh)
		m.setCache(ctx, cacheKey, cacheToken)
	}
	return cacheToken, nil
}

func (m *MyToken) GenerateToken(ctx context.Context, userKey string, data interface{}) (*MyCacheToken, error) {
	myRequestToke, err := m.EncryptToken(ctx, userKey, "")
	if err != nil {
		return nil, err
	}

	cacheKey := CacheKeyPrefix + userKey
	nowTime := gtime.Now().TimestampMilli()
	myCacheToken := &MyCacheToken{
		Token:         myRequestToke.Token,
		Uuid:          myRequestToke.Uuid,
		UserKey:       userKey,
		UserData:      data,
		CreatedAt:     nowTime,
		NextFreshTime: nowTime + gconv.Int64(CacheMaxRefresh),
	}

	err = m.setCache(ctx, cacheKey, myCacheToken)
	if err != nil {
		return nil, errorCode.NewMyErr(ctx, errorCode.MyInternalError, err)
	}

	return myCacheToken, nil
}
