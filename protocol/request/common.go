package request

const (
	Signature = "signature"
	Timestamp = "timestamp"
	Secret    = "secret"
)

type ClientCommonParam struct {
	Di        string `form:"di" binding:"required"`
	Gi        string `form:"gi"`
	Pn        string `form:"pn" binding:"required"`
	Vn        string `form:"vn" binding:"required"`
	Sv        string `form:"sv" binding:"required"`
	Pf        string `form:"pf" binding:"required"`
	Ui        string `form:"ui"`
	Ts        string `form:"ts" binding:"required"`
	Tz        string `form:"tz" binding:"required"`
	Cc        string `form:"cc" binding:"required"`
	Ul        string `form:"ul" binding:"required"`
	Rs        string `form:"rs" binding:"required"`
	Signature string `form:"signature" binding:"required"`
}

func (p ClientCommonParam) Convert2Map() map[string]interface{} {
	pMap := make(map[string]interface{}, 11)
	pMap["di"] = p.Di
	if p.Gi != "" {
		pMap["gi"] = p.Gi
	}
	pMap["pn"] = p.Pn
	pMap["vn"] = p.Vn
	pMap["sv"] = p.Sv
	pMap["pf"] = p.Pf
	if p.Ui != "" {
		pMap["ui"] = p.Ui
	}
	pMap["ts"] = p.Ts
	pMap["tz"] = p.Tz
	pMap["cc"] = p.Cc
	pMap["ul"] = p.Ul
	pMap["rs"] = p.Rs
	pMap["signature"] = p.Signature
	return pMap
}
