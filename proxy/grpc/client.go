package grpc

import (
	"github.com/integration-system/isp-lib/backend"
	"github.com/integration-system/isp-lib/structure"
	log "github.com/integration-system/isp-log"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"isp-gate-service/conf"
	"isp-gate-service/log_code"
	"isp-gate-service/service"
	"net/http"
	"time"
)

type grpcProxy struct {
	client *backend.RxGrpcClient
}

func NewProxy() *grpcProxy {
	return &grpcProxy{client: backend.NewRxGrpcClient(
		backend.WithDialOptions(
			grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(conf.DefaultMaxResponseBodySize))),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(conf.DefaultMaxResponseBodySize))),
		),
		backend.WithDialingErrorHandler(func(err error) {
			log.Errorf(log_code.ErrorClientGrpc, "dialing err: %v", err)
		}),
	)}
}

func (p *grpcProxy) ProxyRequest(ctx *fasthttp.RequestCtx) {
	currentTime := time.Now()

	uri := string(ctx.RequestURI())
	SetHandler(ctx).Complete(ctx, uri, p.client)

	executionTime := time.Since(currentTime) / 1e6

	service.Metrics.UpdateStatusCounter(ctx.Response.StatusCode())
	if ctx.Response.StatusCode() == http.StatusOK {
		service.Metrics.UpdateResponseTime(executionTime)
		service.Metrics.UpdateMethodResponseTime(uri, executionTime)
	}
}

func (p *grpcProxy) Consumer(addr []structure.AddressConfiguration) bool {
	return p.client.ReceiveAddressList(addr)
}

func (p *grpcProxy) Close() {
	p.client.Close()
}
