package config

import (
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config group
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "설정 관리 명령어 그룹",
}

func init() {

}
