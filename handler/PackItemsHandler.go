package handler

import (
	"context"

	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/pack"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/pack/biz"
)

func PackItemsHandler(ctx context.Context, req *pack.PackItemsReq) (res *pack.PackItemsRes, erro error) {

	res = &pack.PackItemsRes{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
		ItemList: []*product_core.ItemInfo{},
	}

	ids := req.GetItemIdList()

	for _, id := range ids {
		if itemRedis, err := biz.GetItemInfo(id); err != nil {
			res.Base.Code = "400"
			res.Base.Message = err.Error()
			return
		} else {
			itemInfo := &product_core.ItemInfo{
				ItemId:             itemRedis.ID,
				Name:               itemRedis.Name,
				Description:        itemRedis.Description,
				Category:           itemRedis.Category,
				Price:              itemRedis.Price,
				Stock:              itemRedis.Stock,
				Province:           itemRedis.Province,
				City:               itemRedis.City,
				District:           itemRedis.District,
				ImageUrl:           itemRedis.ImageURL,
				UserId:             itemRedis.UserID,
				UserType:           itemRedis.UserType,
				SpecificAttributes: itemRedis.SpecificAttr,
			}
			res.ItemList = append(res.ItemList, itemInfo)
		}

	}

	res.Base.Code = "200"
	return
}
