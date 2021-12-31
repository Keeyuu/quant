package request

const (
	RECOMMENDED = "recommend"
	DEFAULT     = "default"
	DISTANCE    = "distance"
	SCORE       = "score"
	CASHBACK    = "cashback"
)

type SearchStoreReq struct {
	Size       int64      `json:"size"`
	From       int64      `json:"from"`
	Where      Where      `json:"where"`
	Coordinate Coordinate `json:"coordinate"`
	Sort       Sort       `json:"sort"`
}

type Where struct {
	Keyword   string   `json:"keyword"`
	Zone      string   `json:"zone"`
	Distance  int64    `json:"distance"`
	CashBack  bool     `json:"cash_back"`
	Score     float64  `json:"score"`
	Category  string   `json:"category"`
	Country   string   `json:"country"`
	StoreType string   `json:"store_type"`
	Status    []string `json:"status"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Sort struct {
	Keyword  string `json:"keyword"`
	Sequence string `json:"sequence"`
}
