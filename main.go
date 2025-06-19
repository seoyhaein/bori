package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"github.com/seoyhaein/bori/cmd"
	"github.com/seoyhaein/bori/grpcv1"
	globallog "github.com/seoyhaein/bori/log"
	"github.com/seoyhaein/go-grpc-kit/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var logger = globallog.Log

func main() {
	// 인증서가 없으면 생성
	if !fileExists("ca.crt") || !fileExists("ca.key") {
		if err := GenerateTLSCerts("my-service", "my-client", 365*24*time.Hour); err != nil {
			log.Fatalf("TLS cert 생성 실패: %v", err)
		}
		log.Println("✅ TLS 인증서 생성 완료")
	}

	// 1) gRPC 서버를 비동기 시작
	srv, err := grpcv1.StartPeerNodeAsync(
		"bori",
		":50051",     // 자신의 주소
		[]string{},   // client 주소들
		"server.crt", // 서버 인증서
		"server.key", // 서버 개인키
		"ca.crt",     // CA 인증서
	)
	if err != nil {
		log.Fatalf("StartPeerNodeAsync 에러: %v", err)
	}
	logger.Info("▶ PeerNode gRPC 서버가 :50051에서 시작되었습니다.")

	// 2) REPL 을 별도 고루틴에서 실행
	ctx := context.Background()
	shell := ishell.New()
	shell.SetPrompt("tori> ")
	shell.Println("TORI REPL (type help for commands)")
	registerCommandsRecursively(ctx, shell, cmd.RootCmd, nil)
	shell.AddCmd(&ishell.Cmd{
		Name: "exit", Help: "종료",
		Func: func(c *ishell.Context) {
			c.Println("bye!")
			c.Stop()
		},
	})
	go shell.Run()

	// 3) 시그널 대기 및 graceful shutdown
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-sigCtx.Done()
	logger.Info("◀ 시그널 수신. 서버 종료 중...")

	// 4) gRPC 서버 종료
	srv.GracefulStop()
	logger.Info("✅ PeerNode 서버 정상 종료")
}

// path는 부모 명령들("bundle","list" 등)을 순서대로 담은 슬라이스입니다.
func registerCommandsRecursively(ctx context.Context, shell *ishell.Shell, cCmd *cobra.Command, path []string) {
	for _, sub := range cCmd.Commands() {
		newPath := append(path, sub.Name())
		name := strings.Join(newPath, " ")

		isCmd := &ishell.Cmd{
			Name:      name,
			Help:      sub.Short,
			LongHelp:  sub.Long,
			Completer: makeCompleter(sub),
			Func: func(c *ishell.Context) {
				// "bundle list --verbose" 같은 전체 커맨드를 다시 Fields 로 쪼갭니다
				base := strings.Fields(c.Cmd.Name)
				args := append(base, c.Args...)
				if err := cmd.Execute(ctx, args); err != nil {
					c.Printf("error: %v\n", err)
				}
			},
		}

		shell.AddCmd(isCmd)
		registerCommandsRecursively(ctx, shell, sub, newPath)
	}
}

// Cobra 커맨드의 서브커맨드 이름과 플래그 이름을 제안하는 Completer // TODO 확인해야함.
func makeCompleter(cCmd *cobra.Command) func([]string) []string {
	return func(args []string) []string {
		suggestions := []string{}

		// 1) 서브커맨드 이름
		for _, sc := range cCmd.Commands() {
			if sc.IsAvailableCommand() {
				suggestions = append(suggestions, sc.Name())
			}
		}
		// 2) 플래그(단축 + 전체)
		cCmd.Flags().VisitAll(func(f *pflag.Flag) {
			suggestions = append(suggestions, "--"+f.Name)
			if f.Shorthand != "" {
				suggestions = append(suggestions, "-"+f.Shorthand)
			}
		})
		return suggestions
	}
}

// GenerateTLSCerts generates a self-signed CA, plus server and client certs,
// then writes them out as PEM files.
// - serverCN: 서버 인증서의 CommonName (예: "server.example.com")
// - clientCN: 클라이언트 인증서의 CommonName (예: "client-node")
// - validFor: 인증서 유효기간 (예: 365*24*time.Hour)
func GenerateTLSCerts(serverCN, clientCN string, validFor time.Duration) error {
	// 1) Self-signed CA
	caCert, caKey, err := utils.GenerateSelfSignedCA(validFor)
	if err != nil {
		return fmt.Errorf("generate CA failed: %w", err)
	}

	// 2) Server cert
	serverTLS, err := utils.GenerateServerCert(caCert, caKey, serverCN, validFor)
	if err != nil {
		return fmt.Errorf("generate server cert failed: %w", err)
	}

	// 3) Client cert
	clientTLS, err := utils.GenerateClientCert(caCert, caKey, clientCN, validFor)
	if err != nil {
		return fmt.Errorf("generate client cert failed: %w", err)
	}

	// 한 번에 저장할 파일 리스트 구성
	certs := []struct {
		DER      []byte
		Key      interface{} // *rsa.PrivateKey
		CertPath string
		KeyPath  string
	}{
		{DER: caCert.Raw, Key: caKey, CertPath: "ca.crt", KeyPath: "ca.key"},
		{DER: serverTLS.Certificate[0], Key: serverTLS.PrivateKey, CertPath: "server.crt", KeyPath: "server.key"},
		{DER: clientTLS.Certificate[0], Key: clientTLS.PrivateKey, CertPath: "client.crt", KeyPath: "client.key"},
	}

	// 4) 파일로 저장
	for _, c := range certs {
		if err := utils.SavePEM(c.CertPath, c.KeyPath, c.DER, c.Key.(*rsa.PrivateKey)); err != nil {
			return fmt.Errorf("save %s/%s failed: %w", c.CertPath, c.KeyPath, err)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
