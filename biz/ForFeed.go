package biz

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/sanity-io/litter"
)

const (
	ItemCacheTTL      = time.Hour // 缓存时间1小时
	TableName_ForFeed = "item_user_for_feed"
	KeyPrefix_ForFeed = "iteminfo-pack"
)

func GetItemInfoForFeed(itemId int32) (*redis_model.ItemForFeed, error) {
	var item *redis_model.ItemForFeed

	cacheKey := getItemForFeedKey(itemId)

	if cmd := redis.DB.Exists(cacheKey); cmd.Err() != nil || cmd.Val() == 0 {
		//redis中不存在，从mysql中获取
		itemMysql := &mysql_model.ItemForFeed{}
		if tx := mysql.GetDBConn().Table(TableName_ForFeed).Where("id = ?", itemId).First(itemMysql); tx.Error != nil {
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
		item = new(redis_model.ItemForFeed)
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
	cacheKey := getItemForFeedKey(itemMysql.ID)
	if tx := mysql.GetDBConn().Table(TableName_ForFeed).Where("id = ?", itemMysql.ID).First(itemMysql); tx.Error != nil {
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

func DeleteItemForFeed(itemKey int32) error {
	log.Println("delete item cache, ", itemKey)
	return redis.DB.Del(getItemForFeedKey(itemKey)).Err()
}

func getItemForFeedKey(itemId int32) string {
	return fmt.Sprintf("%v-%d", KeyPrefix_ForFeed, itemId)
}

func ClearCacheForFeed() error {
	// 获取所有与通配符匹配的键
	keys, err := redis.DB.Keys(KeyPrefix_ForFeed + "*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		// 如果没有找到匹配的键，则直接返回
		return nil
	}
	// 删除所有匹配的键
	_, err = redis.DB.Del(keys...).Result()
	if err != nil {
		return err
	}
	return nil
}
