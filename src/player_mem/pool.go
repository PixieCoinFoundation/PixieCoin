package player_mem

import (
	"appcfg"
	. "github.com/jolestar/go-commons-pool"
	. "logger"
)

var (
	playerCommonPool *ObjectPool
)

func init() {
	if appcfg.GetServerType() != "" {
		return
	}

	config := &ObjectPoolConfig{
		Lifo:                           false,
		MaxTotal:                       -1,
		MaxIdle:                        100,
		MinIdle:                        DefaultMinIdle,
		MaxWaitMillis:                  DefaultMaxWaitMillis,
		MinEvictableIdleTimeMillis:     DefaultMinEvictableIdleTimeMillis,
		SoftMinEvictableIdleTimeMillis: DefaultSoftMinEvictableIdleTimeMillis,
		NumTestsPerEvictionRun:         0,
		EvictionPolicyName:             DefaultEvictionPolicyName,
		TestOnCreate:                   false,
		TestOnBorrow:                   false,
		TestOnReturn:                   false,
		TestWhileIdle:                  false,
		TimeBetweenEvictionRunsMillis:  int64(-1),
		BlockWhenExhausted:             false,
	}
	playerCommonPool = NewObjectPool(NewPooledObjectFactorySimple(
		func() (interface{}, error) {
			return &GFPlayer{}, nil
		}), config)
}

func GetPoolInfo() (int, int, int) {
	return playerCommonPool.GetNumActive(), playerCommonPool.GetNumIdle(), playerCommonPool.GetDestroyedCount()
}

func GetGFPlayer() *GFPlayer {
	if p, err := playerCommonPool.BorrowObject(); err != nil {
		Err(err)
		return &GFPlayer{}
	} else {
		return p.(*GFPlayer)
	}
}

func PutGFPlayer(item *GFPlayer) {
	playerCommonPool.ReturnObject(item)
}
