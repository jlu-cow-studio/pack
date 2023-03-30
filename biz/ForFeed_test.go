package biz

import (
	"testing"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/discovery"
)

func TestClearCacheForFeed(t *testing.T) {

	discovery.Init()
	redis.Init()
	mysql.Init()
	ClearCacheForFeed()
}
