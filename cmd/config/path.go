package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

// configRootCmd sets or shows the root folder
var configRootCmd = &cobra.Command{
	Use:   "path [Root PATH]",
	Short: "데이터 분석의 기준이 되는 루트 폴더를 설정하거나 조회합니다",
	Long: `인자를 주면 해당 경로를 루트 폴더로 저장(config.yaml).
인자 없이 실행하면 현재 설정된 루트 폴더를 출력합니다.
--clear 플래그로 설정을 기본값(현재 디렉터리)으로 초기화합니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clear, _ := cmd.Flags().GetBool("clear")
		if clear {
			// TODO: clear logic
			fmt.Println("루트 폴더 설정이 초기화되었습니다. (기본값: 현재 디렉터리)")
		} else if len(args) > 0 {
			// TODO: set logic
			fmt.Printf("⚙️ 루트 폴더를 '%s'로 설정했습니다.\n", args[0])
			fmt.Println("다음 단계: tori bundle")
		} else {
			// TODO: get logic
			fmt.Println("현재 루트 폴더: /path/to/root")
		}
		return nil
	},
}

func init() {
	configRootCmd.Flags().Bool("clear", false, "설정을 초기화합니다")
	ConfigCmd.AddCommand(configRootCmd)
}
