package resource

import (
	"context"
	"fmt"
	"net"
	"rulex/rulexrpc"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
	"google.golang.org/grpc"
)

const (
	DEFAULT_TRANSPORT = "tcp"
)

//
type grpcConfig struct {
	Port       uint16             `json:"port" validate:"required"`
	DataModels []typex.XDataModel `json:"dataModels"`
}

type RulexRpcServer struct {
	grpcInEndResource *GrpcInEndResource
	rulexrpc.UnimplementedRulexRpcServer
}

//
// Resource interface
//
type GrpcInEndResource struct {
	typex.XStatus
	rulexServer *RulexRpcServer
	rpcServer   *grpc.Server
}

//
func NewGrpcInEndResource(inEndId string, e typex.RuleX) typex.XResource {
	g := GrpcInEndResource{}
	g.PointId = inEndId
	g.RuleEngine = e
	return &g
}

//
func (g *GrpcInEndResource) Start() error {
	inEnd := g.RuleEngine.GetInEnd(g.PointId)
	config := inEnd.Config
	var mainConfig grpcConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}

	listener, err := net.Listen(DEFAULT_TRANSPORT, fmt.Sprintf(":%d", mainConfig.Port))
	if err != nil {
		return err
	}
	// Important !!!
	g.rpcServer = grpc.NewServer()
	g.rulexServer = new(RulexRpcServer)
	g.rulexServer.grpcInEndResource = g
	//
	rulexrpc.RegisterRulexRpcServer(g.rpcServer, g.rulexServer)
	go func(c context.Context) {
		log.Info("RulexRpc resource started on", listener.Addr())
		g.rpcServer.Serve(listener)
	}(context.Background())

	return nil
}

//
func (g *GrpcInEndResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

//
func (g *GrpcInEndResource) Stop() {
	if g.rpcServer != nil {
		g.rpcServer.Stop()
	}

}
func (g *GrpcInEndResource) Reload() {

}
func (g *GrpcInEndResource) Pause() {

}
func (g *GrpcInEndResource) Status() typex.ResourceState {
	return typex.UP
}

func (g *GrpcInEndResource) Register(inEndId string) error {
	g.PointId = inEndId
	return nil
}

func (g *GrpcInEndResource) Test(inEndId string) bool {
	return true
}

func (g *GrpcInEndResource) Enabled() bool {
	return true
}

func (g *GrpcInEndResource) Details() *typex.InEnd {
	return g.RuleEngine.GetInEnd(g.PointId)
}
func (m *GrpcInEndResource) OnStreamApproached(data string) error {
	return nil
}
func (*GrpcInEndResource) Driver() typex.XExternalDriver {
	return nil
}

//
func (r *RulexRpcServer) Work(ctx context.Context, in *rulexrpc.Data) (*rulexrpc.Response, error) {
	r.grpcInEndResource.RuleEngine.Work(
		r.grpcInEndResource.RuleEngine.GetInEnd(r.grpcInEndResource.PointId),
		in.Value,
	)
	return &rulexrpc.Response{
		Code:    0,
		Message: "OK",
	}, nil
}
func (*RulexRpcServer) Configs() []typex.XConfig {
	return []typex.XConfig{}
}
