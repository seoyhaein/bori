package query

import (
	"fmt"
	"github.com/spf13/cobra"
)

// errorsCmd shows the last bundle errors
var queryErrorsCmd = &cobra.Command{
	Use:   "errors",
	Short: "가장 최근에 발생한 오류의 내역을 조회합니다",
	Long: `마지막 bundle 실행 시 발생한 오류를 테이블 형식으로 출력합니다.
--last 플래그는 기본 동작과 동일합니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: load and display errors
		fmt.Println("최근 발생한 오류 (2건):")
		fmt.Println("| 경로 | 오류 내용 |")
		return nil
	},
}

func init() {
	queryErrorsCmd.Flags().Bool("last", false, "마지막 오류만 조회")
	QueryCmd.AddCommand(queryErrorsCmd)
}
