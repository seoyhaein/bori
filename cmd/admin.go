package cmd

import (
	"github.com/spf13/cobra"
)

// adminCmd provides hidden admin menu
var adminCmd = &cobra.Command{
	Use:    "admin",
	Short:  "관리자 전용 인터랙티브 메뉴를 제공합니다",
	Long:   "",
	Hidden: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// TODO: ID/PW 인증 로직
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 대화형 메뉴 구현
	},
}

func init() {

}
