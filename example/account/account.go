package main

import (
	"fmt"

	"github.com/wooyang2018/corechain-sdk/account"
	"github.com/wooyang2018/corechain-sdk/service"
)

func testAccount() {
	// 测试创建账户
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		fmt.Printf("create account error: %v\n", err)
	}
	fmt.Println(acc)
	fmt.Println(acc.Mnemonic)

	// 测试恢复账户
	acc, err = account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("RetrieveAccount: to %v\n", acc)

	// 创建账户并存储到文件中
	acc, err = account.CreateAndSaveAccountToFile("./keys", "123", 1, 1)
	if err != nil {
		fmt.Printf("createAndSaveAccountToFile err: %v\n", err)
	}
	fmt.Printf("CreateAndSaveAccountToFile: %v\n", acc)

	// 从文件中恢复账户
	acc, err = account.GetAccountFromFile("keys/", "123")
	if err != nil {
		fmt.Printf("getAccountFromFile err: %v\n", err)
	}
	fmt.Printf("getAccountFromFile: %v\n", acc)
	return
}

//测试合约账户
func testContractAccount() {
	// 通过助记词恢复账户
	account, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("retrieveAccount address: %v\n", account.Address)

	// 创建一个合约账户
	// 合约账户是由 (XC + 16个数字 + @xuper) 组成, 比如 "XC1234567890123456@xuper"
	contractAccount := "XC1234567890123456@xuper"

	xchainClient, err := service.New("127.0.0.1:37101")
	tx, err := xchainClient.CreateContractAccount(account, contractAccount)
	if err != nil {
		fmt.Printf("createContractAccount err:%s\n", err.Error())
	}
	fmt.Println(tx.Tx.Txid)
	return
}

func main() {
	//testAccount()
	testContractAccount()
}
