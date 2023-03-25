package logic

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/url"
	"sync"
)

var (
	mu  sync.Mutex
	cli *redis.Redis
)

const (
	delayQueue = "sts:dq:timeOutObjectUrl"
	usedUrl    = "sts:store:usedUrl"
)

// checkSingletonRedis singleton redis pattern
func checkSingletonRedis(redisConf *redis.RedisConf) {
	if cli == nil {
		mu.Lock()
		if cli == nil {
			cli = redis.MustNewRedis(*redisConf)
		}
		mu.Unlock()
	}
}

// addToUsedUrl 对于单个url使用
func addToUsedUrl(redisConf *redis.RedisConf, nUrl string) {
	checkSingletonRedis(redisConf)
	pUrl, _ := url.Parse(nUrl)
	_, _ = cli.Sadd(usedUrl, pUrl.Path)
}

// addUrlsToUsedUrl 对于string数组使用
func addUrlsToUsedUrl(redisConf *redis.RedisConf, urls []string) {
	checkSingletonRedis(redisConf)
	for _, val := range urls {
		nUrl, _ := url.Parse(val)
		_, _ = cli.Sadd(usedUrl, nUrl.Path)
	}
}

// addUnExistUrlToUsedUrl 当更新某项url才使用，较为耗费性能
func addUnExistUrlToUsedUrl(redisConf *redis.RedisConf, urls []string) {
	checkSingletonRedis(redisConf)
	for _, val := range urls {
		pUrl, _ := url.Parse(val)
		_, err := cli.Zrank(delayQueue, pUrl.Path)
		// err != nil 时，代表该元素不存在，所以不加入我们的等待集合中
		if err == nil {
			_, _ = cli.Sadd(usedUrl, pUrl.Path)
		}
	}
}
