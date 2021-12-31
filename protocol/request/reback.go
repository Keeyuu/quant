package request

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenData struct {
	UserId string `json:"user_id" binding:"required"`
	Token  string `json:"token"`
}

type QrcodeDecodeData struct {
	UserId string `json:"user_id,omitempty"`
	Data   string `json:"data" binding:"required"`
	Token  string `json:"token"`
}

type DecodeData struct {
	Id        string `json:"id"`
	Timestamp string `json:"ts"`
}

type OfflinePopupData struct {
	OrderId []string `json:"order_id"`
}
