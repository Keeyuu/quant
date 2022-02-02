package main

import (
	"app/httphandler"
	"app/infrastructure/config"
	"app/infrastructure/debug"
	"app/infrastructure/net"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var apis = []net.Api{
	{http.MethodPost, "/debug", httphandler.HandleDebug, false},
}

func main() {

	debug.PprofRouter()

	// start gin
	net.NewHttp(config.Get().Http).Run(apis)
}
