package query

import (
	"github.com/spf13/cobra"
)

// QueryCmd represents the query group
var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "조회 명령어 그룹",
}

func init() {
	// cmd.RootCmd.AddCommand(queryCmd)
}
