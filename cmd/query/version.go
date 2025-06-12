package query

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Version is injected at build time
var Version = "v1.0.0"

// versionCmd prints CLI version
var queryVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "CLI 도구의 버전을 출력합니다",
	Long:  "빌드 시점에 설정된 Version 변수를 읽어와 현재 설치된 TORI CLI 버전을 출력합니다.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tori version %s\n", Version)
	},
}

func init() {
	QueryCmd.AddCommand(queryVersionCmd)
}
