package main

import (
	"fmt"
	"time"

	"github.com/wooyang2018/corechain-sdk/service"
)

func main() {
	testEvent()
}

func testEvent() error {
	// 创建节点客户端。
	client, err := service.New("127.0.0.1:37101")
	if err != nil {
		return err
	}

	// 监听时间，返回 Watcher，通过 Watche 中的 channel 获取block。
	watcher, err := client.WatchBlockEvent(service.WithSkipEmplyTx())
	if err != nil {
		return err
	}

	go func() {
		for {
			b, ok := <-watcher.FilteredBlockChan
			if !ok {
				fmt.Println("watch block event channel closed.")
				return
			}
			fmt.Printf("%+v\n", b)
		}
	}()

	time.Sleep(time.Second * 3)

	// 关闭监听。
	watcher.Close()
	client.Close()
	return nil
}
