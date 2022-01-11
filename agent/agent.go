package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"log"
	"math/rand"
	"net"
	pb "supwisdom.com/diancan/weigh_dc/proto"
	"time"
)

var (
	uniqueIDDict = "0123456789abcdefghijklmnopqrstuvwxyz"
)

//go:generate protoc --go_out=../proto --go_opt=paths=source_relative --go-grpc_out=../proto --go-grpc_opt=paths=source_relative agent.proto

//Agent
type Agent struct {
	agentID    string
	config     *Config
	grpcServer *grpc.Server
	status     string
	startTime  time.Time
	cancel     context.CancelFunc
}

func NewAgent(agentID string, opts ...Option) *Agent {
	a := &Agent{
		agentID: agentID,
		config: &Config{
			heartBeatInterval:        60,
			transExpiredMinutes:      10,
			historyRecordExpiredDays: 30,
			debug:                    false,
		},
		status: pb.AgentStatus_CLOSED.String(),
	}
	for _, opt := range opts {
		opt.apply(a.config)
	}
	return a
}

// NewAgent2
func NewAgent2(agentID string) *Agent {
	return NewAgent(agentID)
}

// agent set option
func (a *Agent) setOption(opt Option) {
	opt.apply(a.config)
}

func (a *Agent) uniqueID(length int) string {
	n := len(uniqueIDDict)
	uniqueIDs := make([]byte, n)
	for i := 0; i < n; i++ {
		pos := rand.Intn(n)
		uniqueIDs[i] = uniqueIDDict[pos]
	}
	return string(uniqueIDs)
}

func (a *Agent) startGRPCServer() error {
	addr := fmt.Sprintf(":%v", a.config.gRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("agentgo grpcserver listen recv remote request err: %v", err)
		return err
	}
	defer listener.Close()
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(a.config.gRPCCACert)) {
		log.Printf("agentgo AppendCertsFromPEM err %v", "can not load ca certificate")
		return errors.New("can't load ca Certificate")
	}
	certificate, err := tls.X509KeyPair([]byte(a.config.gRPCTLSCert), []byte(a.config.gRPCTLSKey))
	if err != nil {
		log.Printf("agentgo X509KeyPair can't load server certificate: %v", err)
		return errors.New("can't load server certificate")
	}
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    certPool,
		Certificates: []tls.Certificate{certificate},
	})
	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute, Timeout: 5 * time.Minute}),
		grpc.ConnectionTimeout(time.Duration(5) * time.Minute)}
	a.grpcServer = grpc.NewServer(opts...)
	pb.RegisterWeighScaleDeviceServer(a.grpcServer, newDeviceServer(a))
	pb.RegisterBindDeviceServer(a.grpcServer, newBindServer(a))
	a.status = pb.AgentStatus_OPEN.String()
	a.startTime = time.Now()
	log.Printf("agentgo Start grpc server addr: %v ...", addr)
	return a.grpcServer.Serve(listener)
}

func (a *Agent) Option(opt Option) {
	opt.apply(a.config)
}

func (a *Agent) uploadTransTask(ctx context.Context) error {
	return nil
}
