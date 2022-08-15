package common

import (
	"github.com/wooyang2018/corechain/example/pb"
	"github.com/wooyang2018/corechain/example/utils"
)

// MakeTransactionID 计算交易ID，包括签名。
func MakeTransactionID(tx *pb.Transaction) ([]byte, error) {
	return utils.MakeTxId(tx)
}

// MakeTxDigestHash 计算交易哈希，不包括签名。
func MakeTxDigestHash(tx *pb.Transaction) ([]byte, error) {
	return utils.MakeTxDigestHash(tx)
}
