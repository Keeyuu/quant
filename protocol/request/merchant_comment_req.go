package request

type MerchantCommentUpsertReq struct {
	CommentId string         `json:"comment_id"`
	Score     *MerchantScore `json:"score" binding:"required"`
	Content   string         `json:"content"`
	StoreId   string         `json:"store_id" binding:"required"`
	OrderId   string         `json:"order_id"`
	Media     []string       `json:"media"`
	Anonymous bool           `json:"anonymous"`
}

type MerchantScore struct {
	Overall float64 `json:"overall"`
	Service float64 `json:"service"`
	Env     float64 `json:"env"`
	Price   float64 `json:"price"`
}

type MerchantAdditionalCommentUpsertReq struct {
	CommentId string   `json:"comment_id" binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Media     []string `json:"media"`
}

type MerchantCommentDeleteReq struct {
	CommentId string `json:"comment_id" binding:"required"`
}
