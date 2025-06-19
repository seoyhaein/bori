package main

import (
	"context"
	"github.com/abiosoft/ishell/v2"
	"github.com/seoyhaein/bori/cmd"
	"github.com/seoyhaein/go-grpc-kit/peernode"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"strings"
)

func main() {
	ctx := context.Background()

	// grpc 연결.
	pn := peernode.NewPeerNode(
		"tori",            // 내 이름
		"localhost:50051", // 내 gRPC 서버 주소
		[]string{"peerA:50051", "peerB:50051"},
		nil, // 서버 옵션
		nil, // 클라이언트 옵션
		/* registerServices... */
	)

	// 1-2) gRPC 서버 시작
	if err := pn.ServerStart(); err != nil {
		log.Fatalf("PeerNode start failed: %v", err)
	}

	shell := ishell.New()
	shell.SetPrompt("tori> ")
	shell.Println("TORI REPL (type help for commands)")

	// Cobra 명령 트리를 읽어와서 ishell에 등록
	registerCommandsRecursively(ctx, shell, cmd.RootCmd, nil)

	// REPL 종료
	shell.AddCmd(&ishell.Cmd{
		Name: "exit", Help: "종료",
		Func: func(c *ishell.Context) {
			c.Println("bye!")
			c.Stop()
		},
	})

	shell.Run()
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
