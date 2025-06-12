### dependency
- go get github.com/abiosoft/ishell/v2@latest
- go get github.com/AlecAivazis/survey/v2
- go get github.com/spf13/cobra@latest

### 시작하기
- cobra 익히기.
- github.com/olekukonko/tablewriter 익히기
- github.com/AlecAivazis/survey/v2 익히기
- briandowns/spinner 익히기
- fatih/color 익히기

### cobra

- 설치하기

```aiignore
# Go 1.17 이상 필요
go install github.com/spf13/cobra-cli@latest
```

- main.go
```aiignore

package main

import (
	"context"
	"github.com/abiosoft/ishell/v2"
	"github.com/seoyhaein/cliTester/cmd"
	"github.com/seoyhaein/cliTester/state"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 컨텍스트 설정 & 시그널 핸들
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// REPL 실행
	if len(os.Args) == 1 {
		runREPL(ctx)
	} else {
		if err := cmd.Execute(ctx); err != nil {
			log.Fatalf("Command failed: %v", err)
		}
	}
}

func runREPL(ctx context.Context) {
	shell := ishell.New()

	// 프롬프트 갱신 함수
	updatePrompt := func() {
		if state.Ready() {
			shell.SetPrompt("tori(server ready)> ")
		} else {
			shell.SetPrompt("tori(server not ready)> ")
		}
	}

	// 최초 프롬프트 설정
	updatePrompt()

	// Cobra 커맨드들 자동 등록 + 래퍼로 프롬프트 갱신
	for _, cobraCmd := range cmd.RootCmd.Commands() {
		ccmd := cobraCmd // 클로저 캡처 주의
		shell.AddCmd(&ishell.Cmd{
			Name:    ccmd.Name(),
			Aliases: ccmd.Aliases,
			Help:    ccmd.Short,
			Func: func(c *ishell.Context) {
				// 사용자가 친 전체 라인 → cobra에 세팅
				//args := strings.Fields(c.Line)
				//cmd.RootCmd.SetArgs(args)
				cmd.RootCmd.SetArgs(c.RawArgs)
				if _, err := cmd.RootCmd.ExecuteC(); err != nil {
					c.Err(err)
				}
				// 명령이 끝난 뒤 프롬프트 재설정
				updatePrompt()
			},
		})
	}

	// exit/quit 커맨드
	shell.AddCmd(&ishell.Cmd{
		Name: "exit", Help: "REPL 종료", Func: func(c *ishell.Context) {
			c.Println("Bye!")
			shell.Stop()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "quit", Help: "REPL 종료", Func: func(c *ishell.Context) {
			c.Println("Bye!")
			shell.Stop()
		},
	})

	// 컨텍스트 취소 시에도 REPL 종료
	go func() {
		<-ctx.Done()
		shell.Stop()
	}()

	shell.Run()
}


```
### 생각하기
- https://www.notion.so/context-20d86e6dd1d1807eaceffdd2acbd9ff4
- https://www.notion.so/20d86e6dd1d180b39cc8cd0522230090

- completion 관련
- https://chatgpt.com/share/68453f6a-c5cc-800d-9420-fcf0d616fb8f
