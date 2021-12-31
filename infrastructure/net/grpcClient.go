package net

import (
	"app/infrastructure/config"
	"app/infrastructure/util/log"
	activityRpc "app/protocol/grpc/activity"
	balanceRpc "app/protocol/grpc/balance"
	dishes_table "app/protocol/grpc/dishes-table"
	orderRpc "app/protocol/grpc/order"
	pmsGwRpc "app/protocol/grpc/pms-gateway"
	rebateOrderRpc "app/protocol/grpc/rebate-order"
	userRpc "app/protocol/grpc/ums-user"
	"sync"

	"google.golang.org/grpc"
)

var pmsClient PMSServiceClient

type PMSServiceClient struct {
	Cli pmsGwRpc.PMSServiceClient

	connection *grpc.ClientConn
	locker     sync.Mutex
}

func GetPmsClient() (*PMSServiceClient, error) {
	if pmsClient.connection == nil {
		pmsClient.locker.Lock()
		defer pmsClient.locker.Unlock()
		if pmsClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.PMS)
			if err != nil {
				return &pmsClient, err
			}
			// 初始化connection
			pmsClient.connection = conn
			pmsClient.Cli = pmsGwRpc.NewPMSServiceClient(conn)

		}
	}
	return &pmsClient, nil
}

func initConnection(url string) (connection *grpc.ClientConn, err error) {
	connection, err = grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Error("did not connect", log.String("err", err.Error()))
	}
	return
}
