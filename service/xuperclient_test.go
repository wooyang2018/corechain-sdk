package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/wooyang2018/corechain-sdk/account"
	"github.com/wooyang2018/corechain-sdk/common/config"
	"github.com/wooyang2018/corechain/example/pb"
)

func TestNewXClient(t *testing.T) {
	x, err := New("127.0.0.1:9999")
	if err != nil {
		t.Error("New xuperclient asser failed:", err)
	}
	x.Close()

	_, err = New("127.0.0.1:9999", WithConfigFile("./conf/sdk.yaml"))
	if err == nil {
		t.Error("New xuperclient asser failed:", err)
	}

	x, err = New("127.0.0.1:9999", WithGrpcGZIP())
	if err != nil {
		t.Error("New xuperclient asser failed:", err)
	}
	x.Close()
	_, err = New("127.0.0.1:9999", WithGrpcGZIP(), WithGrpcTLS("aaa", "aaa", "aaa", "aaa"))
	if err == nil {
		t.Error("New xuperclient asser failed:", err)
	}
}

func String(t *pb.Transaction) string {
	b, err := json.Marshal(*t)
	if err != nil {
		return fmt.Sprintf("%+v", *t)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *t)
	}
	return out.String()
}

func newClient() *XClient {
	if testNode == "" {
		return &XClient{
			xc:  &MockXClient{},
			ec:  &MockEClient{},
			esc: &MockESClient{},
		}
	}
	return xclient
}

func TestXClient_WatchBlockEvent(t *testing.T) {
	client := newClient()
	watcher, err := client.WatchBlockEvent()
	if err != nil {
		t.Error(err)
	}

	go func() {
		select {
		case _, ok := <-watcher.FilteredBlockChan:
			if !ok {
				t.Error("unexpected closed channel")
			}
		case <-time.After(5 * time.Second):
			t.Error("timed out waiting for block event")
		}
	}()

	time.Sleep(time.Second * 5)
}

func TestXClient_QueryTxByID(t *testing.T) {
	txID := "b1ae1868c4e46651657b5aa9be20ab284d36161c9cc311787a9e81e391dc2bed"
	client := newClient()

	tx, err := client.QueryTxByID(txID)
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("tx:%v\n", tx)
}

func TestXClient_QueryBlockByHeight(t *testing.T) {
	client := newClient()
	height := int64(213822)

	block, err := client.QueryBlockByHeight(height)
	if err != nil {
		t.Fatalf("err:%v\n", err.Error())
	}

	fmt.Printf("block:%v\n", block)
}

func TestXClient_QueryBlockByID(t *testing.T) {
	client := newClient()
	blockID := "b4213754f1f645a9bbcf5cfe13de65d019dd843b76929782d78045fff3cceace"

	block, err := client.QueryBlockByID(blockID)
	if err != nil {
		t.Fatalf("err:%v\n", err.Error())
	}

	fmt.Printf("block:%+v\n", block)
}

func TestXClient_QueryAccountAcl(t *testing.T) { //todo
	client := newClient()
	account := "XC1111111111111111@xuper"

	acl, err := client.QueryAccountACL(account)
	if err != nil {
		t.Fatalf("err:%v\n", err.Error())
	}

	fmt.Printf("acl:%v\n", acl)
}

func TestXClient_QueryMethodAcl(t *testing.T) { // todo
	name := "test0611"
	method := "testUint256Event"
	client := newClient()

	acl, err := client.QueryMethodACL(name, method)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", acl)
}

func TestXClient_QueryAccountContracts(t *testing.T) {
	account := "XC1111111111111111@xuper"
	client := newClient()

	contracts, err := client.QueryAccountContracts(account)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	for _, contract := range contracts {
		fmt.Printf("%v\n", contract)
	}
}

func TestXClient_QueryAddressContracts(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	contracts, err := client.QueryAddressContracts(address)
	if err != nil {
		t.Fatalf("err:%s\n", err)
	}

	for _, contract := range contracts {
		fmt.Printf("%v\n", contract)
	}
}

func TestXClient_QueryBalance(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	bal, err := client.QueryBalance(address)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", bal)
}

func TestXClient_QueryBalanceDetail(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	addrBalanceStatus, err := client.QueryBalanceDetail(address)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", addrBalanceStatus)
}

func TestXClient_QuerySystemStatus(t *testing.T) {
	client := newClient()

	reply, err := client.QuerySystemStatus()
	if err != nil {
		t.Fatalf("err:%s\n", err)
	}

	fmt.Printf("reply:%v\n", reply)
}

func TestXClient_QueryBlockChains(t *testing.T) {
	client := newClient()

	bcStatus, err := client.QueryBlockChains()
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	fmt.Printf("bcStatus:%v\n", bcStatus)
}

func TestXClient_QueryBlockChainStatus(t *testing.T) {
	client := newClient()

	bcStatus, err := client.QueryBlockChainStatus()
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("bcStatus:%v\n", bcStatus)
}

func TestXClient_QueryNetURL(t *testing.T) {
	client := newClient()

	rawURL, err := client.QueryNetURL()
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("rawURL:%v\n", rawURL)
}

func TestXClient_QueryAccountByAK(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	resp, err := client.QueryAccountByAK(address)
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("resp:%v\n", resp)
}

var (
	testNode           string = "" // ??????????????? xchain ?????????URL ????????????????????????????????????????????????//10.12.199.82:8701
	wasmCppCounterFile        = "" // ???????????????counter wasm???????????????????????????
)

var (
	// ?????????QZcRXTPwtYo2UtZXNRppEC3Mto2AvmgKZ???????????????????????????????????????????????????
	mnemonic    = "??? ??? ??? ??? ??? ??? ??? ??? ??? ??? ??? ???"
	aks         = []*account.Account{}
	contractAcc = []string{
		"XC0000000000000000@xuper",
		"XC1111111111111111@xuper",
		"XC2222222222222222@xuper",
		"XC3333333333333333@xuper",
		"XC4444444444444444@xuper",
		"XC5555555555555555@xuper",
		"XC6666666666666666@xuper",
		"XC7777777777777777@xuper",
		"XC8888888888888888@xuper",
		"XC9999999999999999@xuper",
	}

	xclient *XClient = nil
)

// ??????10???AK?????????AK???????????????????????????0-9.
// ?????????????????????????????????????????? QZcRXTPwtYo2UtZXNRppEC3Mto2AvmgKZ ????????????????????????100000000000???
// ??????????????? ??? ??? ??? ??? ??? ??? ??? ??? ??? ??? ???
func init() {
	for i := 0; i < len(contractAcc); i++ {
		acc, _ := account.CreateAccount(1, 1)
		aks = append(aks, acc)
		if testNode != "" {
			var err error
			richAcc, _ := account.RetrieveAccount(mnemonic, 1)
			if xclient == nil {
				xclient, err = New(testNode)
				if err != nil {
					panic(err)
				}
			}

			// ??????
			_, err = xclient.Transfer(richAcc, acc.Address, "1000000000")
			if err != nil {
				panic(err)
			}
			_, err = xclient.CreateContractAccount(acc, contractAcc[i])
			if err != nil {
				panic(err)
			}
			// ?????????????????????
			_, err = xclient.Transfer(richAcc, contractAcc[i], "1000000000")
			if err != nil {
				panic(err)
			}
		}
	}
	if testNode != "" {
		fmt.Println("Test real node client start...")
	} else {
		fmt.Println("Test mock client start...")
	}
}

func TestTransfers(t *testing.T) {
	type Case struct {
		from                    *account.Account
		to                      string
		amount                  string
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:   aks[0],
			to:     aks[1].Address,
			amount: "100",
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			desc: "????????????100????????? fee???",
		},
		{
			from:   aks[0],
			to:     aks[1].Address,
			amount: "100",
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			fee:  big.NewInt(100),
			desc: "???????????? fee 100",
			opts: []RequestOption{WithFee("100")},
		},
		{
			from:           aks[0],
			to:             aks[1].Address,
			amount:         "100",
			hasContractAcc: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			fee:  big.NewInt(100),
			desc: "???????????????????????? fee 100",
			opts: []RequestOption{WithFee("100")},
		},
		{
			from:           aks[0],
			to:             aks[1].Address,
			amount:         "100",
			hasContractAcc: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			fee:                     big.NewInt(100),
			desc:                    "???????????????????????? fee 100",
			opts:                    []RequestOption{WithFee("100"), WithFeeFromAccount()},
			onlyFeeFromContractAccc: true,
		},
		{
			from:           aks[0],
			to:             aks[1].Address,
			amount:         "100",
			hasContractAcc: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			fee:  big.NewInt(100),
			desc: "???????????? fee 100",
			opts: []RequestOption{WithFee("100")},
		},
	}

	for _, c := range cases {
		if testNode != "" {
			// ????????????????????????????????????
			bal, err := xclient.QueryBalance(c.from.Address)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(bal)

			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[0])
			}
			xclient.Transfer(c.from, c.to, c.amount, c.opts...)
			c.from.RemoveContractAccount()
			time.Sleep(time.Millisecond * 500)
			// ????????????????????????????????????
		} else {
			// mock ??????????????????
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[0])
			}
			_, err := xc.Transfer(c.from, c.to, c.amount, c.opts...)
			if err != nil {
				t.Error(err)
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func string2Bigint(v string) *big.Int {
	b, ok := big.NewInt(0).SetString(v, 10)
	if !ok {
		panic("string 2 bigint failed:" + v)
	}
	return b
}

func TestDeployContract(t *testing.T) {

	type Case struct {
		from                    *account.Account
		name                    string
		runtime                 string
		abi                     []byte
		code                    []byte
		args                    map[string]string
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: GoRuntime,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "?????? go native ?????????",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: JavaRuntime,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc:          true,
			onlyFeeFromContractAccc: true,
			opts:                    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			desc:                    "?????? java native ?????????account ?????? fee???fee=10???",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: EvmContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: true,
			desc:           "?????? evm ???????????????????????????????????????????????????",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: CRuntime,
			opts:    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ????????????????????????100",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: CRuntime,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ????????????????????????100???account ??????????????????fee=10???",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock ??????????????????
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			switch c.runtime {
			case JavaRuntime:
				_, err := xc.DeployNativeJavaContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

				_, err = xc.UpgradeNativeContract(c.from, c.name, c.code, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case GoRuntime:
				_, err := xc.DeployNativeGoContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

			case CRuntime:
				_, err := xc.DeployWasmContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

				_, err = xc.UpgradeWasmContract(c.from, c.name, c.code, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case EvmContractModule:
				_, err := xc.DeployEVMContract(c.from, c.name, c.abi, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			default:
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestInvokeContract(t *testing.T) {
	type Case struct {
		from                    *account.Account
		name                    string
		module                  string
		method                  string
		args                    map[string]string
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: NativeContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "?????? native ?????????",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: WasmContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc:          true,
			onlyFeeFromContractAccc: true,
			opts:                    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			desc:                    "?????? wasm ?????????account ?????? fee???fee=10???",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: WasmContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ???????????????????????????????????????????????????",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: WasmContractModule,
			opts: []RequestOption{WithContractInvokeAmount("10"), WithBcname("xuper"),
				WithDesc("haha"), WithOtherAuthRequires([]string{"a", "b"})},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ????????????????????????100",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: EvmContractModule,
			opts:   []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ????????????????????????100???account ??????????????????fee=10???",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock ??????????????????
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			switch c.module {
			case NativeContractModule:
				_, err := xc.InvokeNativeContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

				_, err = xc.QueryNativeContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case WasmContractModule:
				_, err := xc.InvokeWasmContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
				_, err = xc.QueryWasmContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case EvmContractModule:
				_, err := xc.InvokeEVMContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
				_, err = xc.QueryEVMContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			default:
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestACLSet(t *testing.T) {
	type Case struct {
		from                    *account.Account
		name                    string
		method                  string
		acl                     *ACL
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "?????????????????? ACL???",
		},
		{
			from: aks[1],
			acl:  getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "?????????????????? ACL???",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc:          true,
			onlyFeeFromContractAccc: true,
			opts: []RequestOption{WithFeeFromAccount(), WithFee("10"), WithBcname("xuper"),
				WithDesc("haha"), WithOtherAuthRequires([]string{"a", "b"})},
			desc: "???????????? ACL???account ?????? fee???fee=10???",
		},
		{
			from: aks[1],
			acl:  getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: true,
			desc:           "?????????????????? ACL?????????????????????????????????????????????",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			opts:   []RequestOption{},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ????????????????????????100",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			opts:   []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "?????? wasm ????????????????????????100???account ??????????????????fee=10???",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock ??????????????????
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			if c.name != "" {
				_, err := xc.SetMethodACL(c.from, c.name, c.method, c.acl, c.opts...)
				if err != nil {
					t.Error(err)
				}
			} else {
				_, err := xc.SetAccountACL(c.from, c.acl, c.opts...)
				if err != nil {
					t.Error(err)
				}
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestCreateAccount(t *testing.T) {
	type Case struct {
		from                    *account.Account
		account                 string
		acl                     *ACL
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
		hasError                bool
	}

	cases := []Case{
		{
			from:    aks[1],
			account: "",
			acl:     getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasError: true,
			desc:     "??????????????????????????????????????????",
		},
		{
			from:     aks[1],
			account:  "XC1234567812345678@xuper",
			acl:      getDefaultACL("a"),
			hasError: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "?????????????????????from ????????????????????????",
		},
		{
			from:     aks[1],
			account:  "hello",
			acl:      getDefaultACL("a"),
			hasError: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			onlyFeeFromContractAccc: true,
			opts: []RequestOption{WithFeeFromAccount(), WithFee("10"), WithBcname("xuper"),
				WithDesc("haha"), WithOtherAuthRequires([]string{"a", "b"})},
			desc: "?????????????????????????????????????????????",
		},
		{
			from:    aks[1],
			acl:     getDefaultACL("a"),
			account: "XC1234567812345678@xuper",
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: false,
			desc:           "?????????????????? ok???",
		},
		{
			from:    aks[1],
			account: "XC1234567812345678@xuper",
			acl:     getDefaultACL("a"),
			opts:    []RequestOption{},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			desc: "????????????????????????????????????100",
		},
		{
			from:    aks[1],
			account: "XC1234567812345678@xuper",
			acl:     getDefaultACL("a"),
			opts:    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasError: true,
			desc:     "????????????????????????????????????100???account ??????????????????fee=10???",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock ??????????????????
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			tx, err := xc.CreateContractAccount(c.from, c.account, c.opts...)
			if c.hasError {
				if err == nil {
					t.Error("Create contract assert err filed", c.desc)
				}
			} else {
				if err != nil {
					t.Error("Create contract assert err filed", err, c.desc)
				}

				if len(tx.Tx.GetTxid()) == 0 {
					t.Error("Create contract assert tx filed", c.desc)
				}
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestRequest(t *testing.T) {
	r := new(Request)
	r.SetArgs(map[string][]byte{"a": []byte("a")})
	r.SetContractName("counter")
	acc, _ := account.CreateAccount(1, 1)
	r.SetInitiatorAccount(acc)
	r.SetModule("xx")
	r.SetTransferAmount("10")
	r.SetTransferTo("bob")

	if r.contractName != "counter" {
		t.Error("Request set assert failed")
	}

	if r.module != "xx" {
		t.Error("Request set assert failed")
	}

	if r.transferTo != "bob" {
		t.Error("Request set assert failed")
	}

	if r.transferAmount != "10" {
		t.Error("Request set assert failed")
	}

	if r.initiatorAccount.Address != acc.Address {
		t.Error("Request set assert failed")
	}

}
