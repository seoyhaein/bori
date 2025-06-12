package cmd

import (
	"context"
	cfg "github.com/seoyhaein/cliTester/cmd/config"
	"github.com/seoyhaein/cliTester/cmd/query"
	"github.com/seoyhaein/cliTester/state"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tori",
	Short: "TORI CLI 도구",
	Long: `Tidy Omics & Raw-data Integrator (TORI)는 대용량 오믹스 데이터의
스캔·그룹핑·번들링·전송·조회 기능을 제공합니다.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(ctx context.Context) error {
	RootCmd.SetContext(ctx)
	if err := RootCmd.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	// 글로벌(공통) 설정 플래그를 여기에 추가할 수 있습니다.
	// e.g. RootCmd.PersistentFlags().BoolP("debug", "d", false, "디버그 모드 활성화")
	// 명령 정렬 비활성화
	cobra.EnableCommandSorting = false
	// 비활성화: Cobra 기본 completion 명령 제거
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	// 커스텀 completion 명령 등록
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(bundleCmd)
	RootCmd.AddCommand(query.QueryCmd)
	RootCmd.AddCommand(cfg.ConfigCmd)
	RootCmd.AddCommand(completionCmd)
	// 기존 help 서브커맨드 제거
	RootCmd.SetHelpCommand(helpCmd)
	RootCmd.AddCommand(adminCmd)

	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if state.Ready() {
			cmd.Root().Use = "tori(server ready)"
		} else {
			cmd.Root().Use = "tori(server not ready)"
		}
	}
}
