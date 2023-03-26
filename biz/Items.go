package biz

import (
	"encoding/json"
	"fmt"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
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
	}

	return item, nil
}

func getItemInfoKey(itemId int32) string {
	return fmt.Sprintf("iteminfo-pack-%d", itemId)
}
