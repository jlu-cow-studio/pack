package biz

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/sanity-io/litter"
)

const (
	ItemCacheTTL = 3600 // 缓存时间1小时
)

func GetItemInfo(itemId int32) (*redis_model.Item, error) {
	var item *redis_model.Item

	cacheKey := getItemInfoKey(itemId)

	if cmd := redis.DB.Exists(cacheKey); cmd.Err() != nil || cmd.Val() == 0 {
		//redis中不存在，从mysql中获取
		itemMysql := &mysql_model.Item{}
		if tx := mysql.GetDBConn().Table("items").Where("id = ?", itemId).First(itemMysql); tx.Error != nil {
			return nil, tx.Error
		}
		item = itemMysql.ToRedis()
		log.Printf("get item info from mysql\n")
		litter.Dump(item)
		if itembyte, err := json.Marshal(item); err != nil {
			return nil, err
		} else if setcmd := redis.DB.Set(cacheKey, string(itembyte), ItemCacheTTL); setcmd.Err() != nil {
			return nil, err
		}
	} else {
		strcmd := redis.DB.Get(cacheKey)
		if strcmd.Err() != nil {
			return nil, strcmd.Err()
		}
		if err := json.Unmarshal([]byte(strcmd.Val()), item); err != nil {
			return nil, err
		}
		log.Printf("get item info from redis\n")
		litter.Dump(item)
	}
	log.Printf("item info \n")
	litter.Dump(item)

	return item, nil
}

func UpdateItemInfo(itemMysql *mysql_model.Item) error {
	log.Println("update item info cache: ")
	litter.Dump(itemMysql)
	cacheKey := getItemInfoKey(itemMysql.ID)
	if tx := mysql.GetDBConn().Table("items").Where("id = ?", itemMysql.ID).First(itemMysql); tx.Error != nil {
		return tx.Error
	}
	item := itemMysql.ToRedis()
	if itembyte, err := json.Marshal(item); err != nil {
		return err
	} else if setcmd := redis.DB.Set(cacheKey, string(itembyte), ItemCacheTTL); setcmd.Err() != nil {
		return err
	}
	return nil
}

func DeleteItemInfo(itemKey int32) error {
	log.Println("delete item cache, ", itemKey)
	return redis.DB.Del(getItemInfoKey(itemKey)).Err()
}

func getItemInfoKey(itemId int32) string {
	return fmt.Sprintf("iteminfo-pack-%d", itemId)
}
