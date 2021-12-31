package response

import "app/model"

type UserBalanceResp struct {
	Code int32           `json:"code"`
	Body UserBalanceBody `json:"body"`
}

type UserBalanceBody struct {
	Result bool `json:"result"`
}

type SearchStoreRsb struct {
	Data  []model.EsShop `json:"shop_list"`
	Total int            `json:"total"`
}
