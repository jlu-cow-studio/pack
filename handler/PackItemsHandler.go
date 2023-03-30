package handler

import (
	"context"
	"log"

	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	"github.com/jlu-cow-studio/common/dal/rpc/pack"
	"github.com/jlu-cow-studio/pack/biz"
	"github.com/sanity-io/litter"
)

func (h *Handler) PackItems(ctx context.Context, req *pack.PackItemsForFeedReq) (res *pack.PackItemsForFeedRes, erro error) {

	res = &pack.PackItemsForFeedRes{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
		ItemList: []*feed_service.ItemForFeed{},
	}

	log.Printf("[Pack] pack items req : %v", litter.Sdump(req))

	ids := req.GetItemIdList()

	for _, id := range ids {
		if itemRedis, err := biz.GetItemInfoForFeed(id); err != nil {
			res.Base.Code = "400"
			res.Base.Message = err.Error()
			return
		} else {
			itemInfo := &feed_service.ItemForFeed{
				Id:                 int32(itemRedis.ID),
				Name:               itemRedis.Name,
				Description:        itemRedis.Description,
				Category:           itemRedis.Category,
				Price:              itemRedis.Price,
				Stock:              int32(itemRedis.Stock),
				Province:           itemRedis.Province,
				City:               itemRedis.City,
				District:           itemRedis.District,
				ImageUrl:           itemRedis.ImageURL,
				UserId:             int32(itemRedis.UserID),
				UserType:           itemRedis.UserType,
				SpecificAttributes: itemRedis.SpecificAttributes,
				Uid:                int32(itemRedis.UID),
				Username:           itemRedis.Username,
				Uprovince:          itemRedis.UProvince,
				Ucity:              itemRedis.UCity,
				Udistrict:          itemRedis.UDistrict,
				Urole:              itemRedis.URole,
			}
			res.ItemList = append(res.ItemList, itemInfo)
		}

	}

	res.Base.Code = "200"

	log.Printf("[Pack] pack items res : %v", litter.Sdump(res))
	return
}
