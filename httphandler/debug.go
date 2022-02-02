package httphandler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type S struct {
	A []float32 `json:"a"`
	B []float32 `json:"b"`
	C []float32 `json:"c"`
}

func HandleDebug(ctx *gin.Context) {
	go func() {
		s := new(S)
		ctx.ShouldBindJSON(s)
		fmt.Println("get a request", s)
	}()
}
