package token

const (
	CacheModeCache = 1
	CacheModeRedis = 2

	TokenDelimiter  = "_"
	EncryptKey      = "12345678912345678912345678912345"
	CacheKeyPrefix  = "Token:"
	CacheTimeout    = 10 * 60 * 1000 // 10分钟
	CacheMaxRefresh = CacheTimeout / 2
)
