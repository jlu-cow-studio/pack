package biz

import (
	"testing"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/sanity-io/litter"
)

func TestClearCacheForFeed(t *testing.T) {

	discovery.Init()
	redis.Init()
	mysql.Init()

	err := ClearCacheForFeed()
	litter.Dump(err)
}
