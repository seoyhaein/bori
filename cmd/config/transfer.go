package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

// configTransferCmd sets or shows the transfer mode
var configTransferCmd = &cobra.Command{
	Use:   "transfer [auto|manual]",
	Short: "파일 전송 방식을 설정하거나 조회합니다",
	Long:  "'auto' 또는 'manual'로 서버 전송 모드를 설정합니다. 인자 없이 실행하면 현재 모드를 출력합니다.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			mode := args[0]
			// TODO: save mode
			fmt.Printf("⚙️ 전송 방식을 '%s'로 설정했습니다.\n", mode)
		} else {
			// TODO: show current mode
			fmt.Println("현재 전송 방식: auto")
		}
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(configTransferCmd)
}
