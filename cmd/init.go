package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// initCmd initializes TORI's state
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "TORI 기록을 초기화하고 처음 상태로 되돌립니다",
	Long: `모든 스캔 및 그룹핑 기록을 삭제하고 TORI를
처음 설치된 상태로 되돌립니다. 실행 전 사용자에게 y/N 확인 프롬프트를 띄웁니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: y/N 확인 및 초기화 로직 구현
		fmt.Println("⚠️ 이 작업을 실행하면 기존 기록이 삭제됩니다. 계속하시겠습니까? (y/N)")
		// ...
		fmt.Println("✅ 초기화가 완료되었습니다.")
		return nil
	},
}

func init() {

}
