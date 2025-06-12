package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/abiosoft/ishell/v2"
	"os"
)

func main() {
	shell := ishell.New()
	shell.SetPrompt("tori> ")

	shell.AddCmd(&ishell.Cmd{
		Name: "ask",
		Help: "간단한 입력 프롬프트 예제",
		Func: func(c *ishell.Context) {
			// AskOne에 ishell의 입출력(streams)을 전달
			var name string
			prompt := &survey.Input{
				Message: "이름을 입력하세요:",
			}
			survey.AskOne(prompt, &name,
				survey.WithStdio(os.Stdin, os.Stdout, os.Stderr))

			c.Printf("반갑습니다, %s님!\n", name)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "select",
		Help: "선택 메뉴 예제",
		Func: func(c *ishell.Context) {
			var choice string
			prompt := &survey.Select{
				Message: "원하는 기능을 고르세요:",
				Options: []string{"파이프라인 실행", "상태 확인", "종료"},
			}
			survey.AskOne(prompt, &choice,
				survey.WithStdio(os.Stdin, os.Stdout, os.Stderr))

			c.Println("선택된 항목:", choice)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "종료",
		Func: func(c *ishell.Context) {
			c.Println("bye!")
			c.Stop()
		},
	})

	shell.Run()
}
