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

type OrderClient struct {
	OrderService orderRpc.OrderServiceClient
	connection   *grpc.ClientConn
	locker       sync.Mutex
}

var orderClient OrderClient

func GetOrderClient() (*OrderClient, error) {
	if orderClient.connection == nil {
		orderClient.locker.Lock()
		defer orderClient.locker.Unlock()
		if orderClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.Order)
			if err != nil {
				return &orderClient, err
			}
			// 初始化connection
			orderClient.connection = conn
			orderClient.OrderService = orderRpc.NewOrderServiceClient(conn)
		}
	}
	return &orderClient, nil
}

// rebate order服务
var rebateOrderClient RebateOrderClient

type RebateOrderClient struct {
	RebateOrder rebateOrderRpc.RebateOrderServiceClient

	connection *grpc.ClientConn
	locker     sync.Mutex
}

func GetRebateOrderClient() (*RebateOrderClient, error) {
	if rebateOrderClient.connection == nil {
		rebateOrderClient.locker.Lock()
		defer rebateOrderClient.locker.Unlock()
		if rebateOrderClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.RebateOrder)
			if err != nil {
				return &rebateOrderClient, err
			}
			// 初始化connection
			rebateOrderClient.connection = conn
			rebateOrderClient.RebateOrder = rebateOrderRpc.NewRebateOrderServiceClient(conn)
		}
	}
	return &rebateOrderClient, nil
}

// new ums-user service
var userClient UserClient

type UserClient struct {
	UserService userRpc.UserServiceClient

	connection *grpc.ClientConn
	locker     sync.Mutex
}

func GetUserClient() (*UserClient, error) {
	if userClient.connection == nil {
		userClient.locker.Lock()
		defer userClient.locker.Unlock()
		if userClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.User)
			if err != nil {
				return &userClient, err
			}
			// 初始化connection
			userClient.connection = conn
			userClient.UserService = userRpc.NewUserServiceClient(conn)
		}
	}
	return &userClient, nil
}

var balanceClient BalanceClient

type BalanceClient struct {
	Balance    balanceRpc.BalanceServiceClient
	connection *grpc.ClientConn
	locker     sync.Mutex
}

func GetBalanceClient() (*BalanceClient, error) {
	if balanceClient.connection == nil {
		balanceClient.locker.Lock()
		defer balanceClient.locker.Unlock()
		if balanceClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.Balance)
			if err != nil {
				return &balanceClient, err
			}
			// 初始化connection
			balanceClient.connection = conn
			balanceClient.Balance = balanceRpc.NewBalanceServiceClient(conn)
		}
	}
	return &balanceClient, nil
}

var activityClient ActivityServiceClient

type ActivityServiceClient struct {
	Cli activityRpc.ActivityServiceClient

	connection *grpc.ClientConn
	locker     sync.Mutex
}

func GetActivityServiceClient() (*ActivityServiceClient, error) {
	if activityClient.connection == nil {
		activityClient.locker.Lock()
		defer activityClient.locker.Unlock()
		if activityClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.Activity)
			if err != nil {
				return &activityClient, err
			}
			// 初始化connection
			activityClient.connection = conn
			activityClient.Cli = activityRpc.NewActivityServiceClient(conn)

		}
	}
	return &activityClient, nil
}

var dishesTableClient DishesTableServiceClient

type DishesTableServiceClient struct {
	DishesTableServiceClient dishes_table.DishesTableServiceClient
	connection               *grpc.ClientConn
	locker                   sync.Mutex
}

func GetDishesTableServiceClient() (*DishesTableServiceClient, error) {
	if dishesTableClient.connection == nil {
		dishesTableClient.locker.Lock()
		defer dishesTableClient.locker.Unlock()
		if dishesTableClient.connection == nil {
			conn, err := initConnection(config.Get().GRpc.Clients.DishesTable)
			if err != nil {
				return &dishesTableClient, err
			}
			// 初始化connection
			dishesTableClient.connection = conn
			dishesTableClient.DishesTableServiceClient = dishes_table.NewDishesTableServiceClient(conn)
		}
	}
	return &dishesTableClient, nil
}

func initConnection(url string) (connection *grpc.ClientConn, err error) {
	connection, err = grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Error("did not connect", log.String("err", err.Error()))
	}
	return
}
