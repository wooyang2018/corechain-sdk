package native

import (
	"context"
	"net/url"
	"time"

	"github.com/wooyang2018/corechain-sdk/code"
	"github.com/wooyang2018/corechain-sdk/exec"
	"github.com/wooyang2018/corechain/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	_ protos.NativeCodeServer = (*nativeCodeService)(nil)
)

type nativeCodeService struct {
	contract  code.Contract
	rpcClient *grpc.ClientConn
	lastping  time.Time
}

func newNativeCodeService(chainAddr string, contract code.Contract) *nativeCodeService {
	uri, err := url.Parse(chainAddr)
	if err != nil {
		panic(err)
	}
	switch uri.Scheme {
	case "tcp":
	default:
		panic("unsupported protocol " + uri.Scheme)
	}
	conn, err := grpc.Dial(uri.Host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &nativeCodeService{
		contract:  contract,
		rpcClient: conn,
		lastping:  time.Now(),
	}
}

func (s *nativeCodeService) bridgeCall(method string, request proto.Message, response proto.Message) error {
	// NOTE sync with contract.proto's package name
	fullmethod := "/protos.Syscall/" + method
	return s.rpcClient.Invoke(context.Background(), fullmethod, request, response)
}

func (s *nativeCodeService) Call(ctx context.Context, request *protos.NativeCallRequest) (*protos.NativeCallResponse, error) {
	exec.RunContract(request.GetCtxid(), s.contract, s.bridgeCall)
	return new(protos.NativeCallResponse), nil
}

func (s *nativeCodeService) Ping(ctx context.Context, request *protos.PingRequest) (*protos.PingResponse, error) {
	s.lastping = time.Now()
	return &protos.PingResponse{}, nil
}

func (s *nativeCodeService) LastpingTime() time.Time {
	return s.lastping
}

func (s *nativeCodeService) Close() error {
	return nil
}
