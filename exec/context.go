package exec

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/wooyang2018/corechain-sdk/code"
	"github.com/wooyang2018/corechain/protos"
)

const (
	methodPut                = "PutObject"
	methodGet                = "GetObject"
	methodDelete             = "DeleteObject"
	methodOutput             = "SetOutput"
	methodGetCallArgs        = "GetCallArgs"
	methodTransfer           = "Transfer"
	methodContractCall       = "ContractCall"
	methodCrossContractQuery = "CrossContractQuery"
	methodQueryTx            = "QueryTx"
	methodQueryBlock         = "QueryBlock"
	methodNewIterator        = "NewIterator"
)

var (
	_ code.Context = (*contractContext)(nil)
)

type contractContext struct {
	callArgs       protos.CallArgs
	contractArgs   map[string][]byte
	bridgeCallFunc BridgeCallFunc
	header         protos.SyscallHeader
}

func newContractContext(ctxid int64, bridgeCallFunc BridgeCallFunc) *contractContext {
	return &contractContext{
		contractArgs:   make(map[string][]byte),
		bridgeCallFunc: bridgeCallFunc,
		header: protos.SyscallHeader{
			Ctxid: ctxid,
		},
	}
}

func (c *contractContext) Init() error {
	var request protos.GetCallArgsRequest
	request.Header = &c.header
	err := c.bridgeCallFunc(methodGetCallArgs, &request, &c.callArgs)
	if err != nil {
		return err
	}
	for _, pair := range c.callArgs.GetArgs() {
		c.contractArgs[pair.GetKey()] = pair.GetValue()
	}
	return nil
}

func (c *contractContext) Method() string {
	return c.callArgs.GetMethod()
}

func (c *contractContext) Args() map[string][]byte {
	return c.contractArgs
}

func (c *contractContext) Caller() string {
	caller := c.callArgs.GetCaller()
	//  fall back toInitiator if caller is not set
	if caller == "" {
		caller = c.callArgs.GetInitiator()
	}
	return caller
}

func (c *contractContext) Initiator() string {
	return c.callArgs.Initiator
}

func (c *contractContext) AuthRequire() []string {
	return c.callArgs.AuthRequire
}

func (c *contractContext) PutObject(key, value []byte) error {
	req := &protos.PutRequest{
		Header: &c.header,
		Key:    key,
		Value:  value,
	}
	rep := new(protos.PutResponse)
	return c.bridgeCallFunc(methodPut, req, rep)
}

func (c *contractContext) GetObject(key []byte) ([]byte, error) {
	req := &protos.GetRequest{
		Header: &c.header,
		Key:    key,
	}
	rep := new(protos.GetResponse)
	err := c.bridgeCallFunc(methodGet, req, rep)
	if err != nil {
		return nil, err
	}
	return rep.Value, nil
}

func (c *contractContext) DeleteObject(key []byte) error {
	req := &protos.DeleteRequest{
		Header: &c.header,
		Key:    key,
	}
	rep := new(protos.DeleteResponse)
	return c.bridgeCallFunc(methodDelete, req, rep)
}

func (c *contractContext) NewIterator(start, limit []byte) code.Iterator {
	return newKvIterator(c, start, limit)
}

func (c *contractContext) QueryTx(txid string) (*protos.Transaction, error) {
	req := &protos.QueryTxRequest{
		Header: &c.header,
		Txid:   string(txid),
	}
	resp := new(protos.QueryTxResponse)
	if err := c.bridgeCallFunc(methodQueryTx, req, resp); err != nil {
		return nil, err
	}
	return resp.Tx, nil
}

func (c *contractContext) QueryBlock(blockid string) (*protos.Block, error) {
	req := &protos.QueryBlockRequest{
		Header:  &c.header,
		Blockid: string(blockid),
	}
	resp := new(protos.QueryBlockResponse)
	if err := c.bridgeCallFunc(methodQueryBlock, req, resp); err != nil {
		return nil, err
	}
	return resp.Block, nil
}

func (c *contractContext) Transfer(to string, amount *big.Int) error {
	req := &protos.TransferRequest{
		Header: &c.header,
		To:     to,
		Amount: amount.Text(10),
	}
	rep := new(protos.TransferResponse)
	return c.bridgeCallFunc(methodTransfer, req, rep)
}

func (c *contractContext) TransferAmount() (*big.Int, error) {
	amount, ok := new(big.Int).SetString(c.callArgs.GetTransferAmount(), 10)
	if !ok {
		return nil, errors.New("bad amount:" + c.callArgs.GetTransferAmount())
	}
	return amount, nil
}

func (c *contractContext) Call(module, contract, method string, args map[string][]byte) (*code.Response, error) {
	var argPairs []*protos.ArgPair
	// 在合约里面单次合约调用的map迭代随机因子是确定的，因此这里不需要排序
	for key, value := range args {
		argPairs = append(argPairs, &protos.ArgPair{
			Key:   key,
			Value: value,
		})
	}
	req := &protos.ContractCallRequest{
		Header:   &c.header,
		Module:   module,
		Contract: contract,
		Method:   method,
		Args:     argPairs,
	}
	rep := new(protos.ContractCallResponse)
	err := c.bridgeCallFunc(methodContractCall, req, rep)
	if err != nil {
		return nil, err
	}
	return &code.Response{
		Status:  int(rep.Response.Status),
		Message: rep.Response.Message,
		Body:    rep.Response.Body,
	}, nil
}

func (c *contractContext) CrossQuery(uri string, args map[string][]byte) (*code.Response, error) {
	var argPairs []*protos.ArgPair
	for key, value := range args {
		argPairs = append(argPairs, &protos.ArgPair{
			Key:   key,
			Value: value,
		})
	}
	req := &protos.CrossContractQueryRequest{
		Header: &c.header,
		Uri:    uri,
		Args:   argPairs,
	}
	rep := new(protos.CrossContractQueryResponse)
	err := c.bridgeCallFunc(methodCrossContractQuery, req, rep)
	if err != nil {
		return nil, err
	}
	return &code.Response{
		Status:  int(rep.Response.Status),
		Message: rep.Response.Message,
		Body:    rep.Response.Body,
	}, nil
}

func (c *contractContext) SetOutput(response *code.Response) error {
	req := &protos.SetOutputRequest{
		Header: &c.header,
		Response: &protos.Response{
			Status:  int32(response.Status),
			Message: response.Message,
			Body:    response.Body,
		},
	}
	rep := new(protos.SetOutputResponse)
	err := c.bridgeCallFunc(methodOutput, req, rep)
	if err != nil {
		log.Printf("Setoutput error:%s", err)
	}
	return err
}

func (c *contractContext) Logf(fmtstr string, args ...interface{}) {
	entry := fmt.Sprintf(fmtstr, args...)
	request := &protos.PostLogRequest{
		Header: &c.header,
		Entry:  entry,
	}
	c.bridgeCallFunc("PostLog", request, new(protos.PostLogResponse))
}

func (c *contractContext) EmitEvent(name string, body []byte) error {
	request := &protos.EmitEventRequest{
		Header: &c.header,
		Name:   name,
		Body:   body,
	}
	return c.bridgeCallFunc("EmitEvent", request, new(protos.EmitEventResponse))
}

func (c *contractContext) EmitJSONEvent(name string, body interface{}) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return c.EmitEvent(name, buf)
}
