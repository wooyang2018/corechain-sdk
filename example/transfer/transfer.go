package main

import (
	"fmt"

	"github.com/wooyang2018/corechain-sdk/account"
	"github.com/wooyang2018/corechain-sdk/service"
)

func main() {
	akTransfer()
	contractAccountTransfer()
}

// akTransfer 普通账户转账（Ak）示例。
func akTransfer() {
	// 创建或者使用已有账户，此处为新创建一个账户。
	from, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}

	to, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(to.Address)
	fmt.Println(to.Mnemonic)

	// 节点地址。
	node := "127.0.0.1:37101"

	// 创建节点客户端。
	xclient, _ := service.New(node)

	// 转账前查看两个地址余额。
	fmt.Println(xclient.QueryBalance(from.Address))
	fmt.Println(xclient.QueryBalance(to.Address))

	tx, err := xclient.Transfer(from, to.Address, "10")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", tx.Tx.Txid)

	// 转账后查看两个地址余额。
	fmt.Println(xclient.QueryBalance(from.Address))
	fmt.Println(xclient.QueryBalance(to.Address))
}

// contractAccountTransfer 合约账户转账示例。
func contractAccountTransfer() {
	// 创建或者使用已有账户，此处为新创建一个账户。
	me, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	// XC1234567812345678@xuper 为合约账户，如果没有合约账户需要先创建合约账户。
	me.SetContractAccount("XC1234567812345678@xuper")
	fmt.Println(me.Address)
	fmt.Println(me.Mnemonic)
	fmt.Println(me.GetContractAccount())
	fmt.Println(me.GetAuthRequire())

	to, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(to.Address)
	fmt.Println(to.Mnemonic)

	// 节点地址。
	node := "127.0.0.1:37101"
	xclient, _ := service.New(node)

	// 转账前查看两个地址余额。
	fmt.Println(xclient.QueryBalance(me.Address))
	fmt.Println(xclient.QueryBalance(to.Address))

	tx, err := xclient.Transfer(me, "a", "10")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", tx.Tx.Txid)

	// 转账后查看两个地址余额。
	fmt.Println(xclient.QueryBalance(me.GetContractAccount())) // 转账时使用的是合约账户，因此查询余额时也是合约账户。
	fmt.Println(xclient.QueryBalance(to.Address))
}
