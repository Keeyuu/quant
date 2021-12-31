package commonutil

import (
	"app/infrastructure/util/log"
	"fmt"
	"runtime/debug"
	"time"
)

/**
 * @Description: TODO
 *
 *
 * @Author: acezhuyf
 * @Since: 2021-05-25
 */

func Retry(retryTimes int, f func() error) error {
	//返回最终结果
	var err error
	for i := 1; i <= retryTimes; i++ {
		err = f()
		if err != nil {
			log.Error(fmt.Sprintf("func [%d] time processing err: %s, retry in %ds", i, err.Error(), 1<<i))
			time.Sleep(time.Second * time.Duration(1<<i))
			continue
		}
		log.Info("func processing successfully")
		break
	}
	return err
}

func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("goroutine panic recover err: %s\n,stack:%s", err, string(debug.Stack())))
		}
	}()
	fn()
}
