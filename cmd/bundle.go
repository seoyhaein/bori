package cmd

import (
	"github.com/spf13/cobra"
)

// bundleCmd scans folders and creates .pb files
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "지정 폴더를 스캔하여 .pb 파일을 생성합니다",
	Long: `설정된 루트 폴더 하위의 모든 샘플 폴더를 검사하고,
rule.json 유효성 검사를 통과한 뒤 하나의 Protobuf 파일(.pb)로 묶습니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 이렇게 꺼내 올 수 있음.
		//ctx := cmd.Context()

		return nil
	},
}

func init() {

}
