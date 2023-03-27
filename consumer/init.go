package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jlu-cow-studio/common/dal/mq"
	"github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	"github.com/jlu-cow-studio/common/model/mq_struct"
	"github.com/jlu-cow-studio/pack/biz"
	"github.com/sanity-io/litter"
	"github.com/segmentio/kafka-go"
)

var ErrChan chan error

func Init() {
	go func() {
		log.Println("listening topic : ", mq.Topic_ItemChange)
		err := mq.Listen(context.Background(), mq.Topic_ItemChange, consumer)
		log.Println("consumer error : ", err.Error())
		close(ErrChan)
		panic(err)
	}()
}

func consumer(msg kafka.Message) error {

	log.Println("get a message, ")
	// litter.Dump(msg.)
	litter.Dump(msg.Headers)
	litter.Dump(msg.Offset)
	litter.Dump(msg.Partition)
	litter.Dump(msg.Topic)
	litter.Dump(msg.Time)
	log.Println("value : ", string(msg.Value))

	itemChange := new(mq_struct.ItemChangeMsg)

	if err := json.Unmarshal(msg.Value, itemChange); err != nil {
		return err
	}

	switch itemChange.Op {
	case mq_struct.ItemOp_Create:
	case mq_struct.ItemOp_Update:
		return biz.UpdateItemInfo(mysql.ItemFromRedis(itemChange.Info))
	case mq_struct.ItemOp_Delete:
		return biz.DeleteItemInfo(itemChange.Info.ID)
	default:
	}

	return nil
}
