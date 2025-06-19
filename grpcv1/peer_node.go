package grpcv1

import (
	"context"
	"fmt"
	globallog "github.com/seoyhaein/bori/log"
	"github.com/seoyhaein/go-grpc-kit/client"
	"github.com/seoyhaein/go-grpc-kit/peernode"
	"github.com/seoyhaein/go-grpc-kit/server"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
)

// TODO 파일이름 수정될 수 있음.

var logger = globallog.Log

func StartPeerNode(name, address string, peers []string, certFile, keyFile, caFile string) error {
	// 1) mTLS + 기본 서버 옵션
	grpcOpts := server.DefaultServerOptions(
		server.WithMTLS(certFile, keyFile, caFile),
	)

	// 2) PeerNode 생성
	pn := peernode.NewPeerNode(
		name,
		address,
		peers,
		grpcOpts,
		[]client.Option{},
		server.WithHealthCheck(),
		server.WithReflection(),
	)

	// 3) 서버 시작
	if err := pn.ServerStart(); err != nil {
		return fmt.Errorf("failed to start PeerNode %s: %w", name, err)
	}
	logger.Infof("▶ PeerNode[%s] gRPC 서버 시작됨: %s", name, address)

	// 4) signal.NotifyContext 로 블로킹 + cancel 처리
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // 프로그램 종료 시 clean-up

	<-ctx.Done() // 여기서 Ctrl+C/SIGTERM 대기
	logger.Infof("◀ 시그널 수신. 서버 종료 중...")

	return nil
}

func StartPeerNodeAsync(name, address string, peers []string, certFile, keyFile, caFile string) (*grpc.Server, error) {
	// 1) mTLS + 기본 서버 옵션
	grpcOpts := server.DefaultServerOptions(
		server.WithMTLS(certFile, keyFile, caFile),
	)

	// 2) PeerNode 생성
	pn := peernode.NewPeerNode(
		name,
		address,
		peers,
		grpcOpts,
		[]client.Option{},
		server.WithHealthCheck(),
		server.WithReflection(),
	)

	// 3) 비동기 서버 시작, *grpc.Server 반환
	srv, err := pn.ServerStartAsync()
	if err != nil {
		return nil, fmt.Errorf("failed to start PeerNode %s async: %w", name, err)
	}
	logger.Infof("▶ PeerNode[%s] async gRPC server started at %s", name, address)
	return srv, nil
}
