package main

import (
	"github.com/jlu-cow-studio/common/dal/mq"
	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/jlu-cow-studio/pack/consumer"
	"github.com/jlu-cow-studio/pack/rpc"
)

func main() {
	discovery.Init()
	redis.Init()
	mysql.Init()
	mq.Init()
	consumer.Init()
	rpc.Init()

	select {
	case err, ok := <-consumer.ErrChan:
		if !ok {
			panic("consumer down! " + err.Error())
		}
	case err, ok := <-rpc.ErrChan:
		if !ok {
			panic("rpc down !" + err.Error())
		}
	}
}
