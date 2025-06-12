package query

import (
	"fmt"
	"github.com/spf13/cobra"
)

// infoCmd shows metadata of a .pb file
var queryInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Protobuf 파일의 메타데이터를 조회합니다",
	Long: `--file 플래그로 지정한 .pb 파일의 버전, 생성 시각,
포함된 샘플 수 등을 출력합니다. 플래그 없으면 기본 경로를 사용합니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		// TODO: display metadata
		fmt.Printf("파일: %s\n버전: v1.0.0\n생성 시각: ...\n샘플 수: ...\n", file)
		return nil
	},
}

func init() {
	queryInfoCmd.Flags().StringP("file", "f", "blocks.pb", "조회할 .pb 파일 경로")
	QueryCmd.AddCommand(queryInfoCmd)
}
