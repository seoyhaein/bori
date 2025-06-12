package cmd

import (
	"github.com/spf13/cobra"
)

// helpCmd replaces the default help command with a localized Korean version.
var helpCmd = &cobra.Command{
	Use:   "help",           // 한국어로 번역된 명령어 이름
	Short: "명령어 사용법을 보여줍니다", // 한국어로 번역된 설명
	Long: `TORI CLI 도구의 전체 사용법을 출력합니다.
예) tori bundle --verbose`, // 필요에 따라 더 긴 설명 추가
	Run: func(cmd *cobra.Command, args []string) {
		// 기본 help 로직 재활용
		_ = cmd.Root().Help()
	},
}

func init() {
}
