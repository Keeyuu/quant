package response

import "app/model"

type MerchantCommentUpsertResp struct {
	Review *MerchantCommentReview `json:"review"`
	Bonus  *MerchantCommentBonus  `json:"bonus"`
}

type MerchantCommentReview struct {
	Id                string                     `json:"id"`
	UserName          string                     `json:"user_name"`
	UserAvatar        string                     `json:"user_avatar"`
	Score             float64                    `json:"score"`
	SubScore          map[string]float64         `json:"sub_score"`
	Content           string                     `json:"content"`
	StoreId           string                     `json:"store_id"`
	OrderId           string                     `json:"order_id"`
	Media             []string                   `json:"media"`
	Recommend         bool                       `json:"recommend"`
	Tags              []string                   `json:"tags"`
	AdditionalComment []*model.AdditionalComment `json:"additional_comment"`
	CreateAt          int64                      `json:"create_at"`
	Status            string                     `json:"status"`
	StoreName         string                     `json:"store_name"`
	StoreLogo         string                     `json:"store_logo"`
	UserId            string                     `json:"user_id"`
	Amount            int64                      `json:"amount"`
	Currency          string                     `json:"currency"`
	Anonymous         bool                       `json:"anonymous"`
}

type MerchantCommentBonus struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
	Min      int64  `json:"min"` // 积分区间最小值
	Max      int64  `json:"max"` // 积分区间最大值
}
