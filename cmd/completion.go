package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	install       bool
	completionDir string
)

var completionCmd = &cobra.Command{
	Use:    "completion [bash|zsh|fish|powershell]",
	Short:  "Generate or install the autocompletion script for the specified shell",
	Hidden: true,
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),  // 정확히 1개의 인자를 요구
		cobra.OnlyValidArgs, // ValidArgs 에 정의된 값만 허용
	),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]
		binName := filepath.Base(os.Args[0])
		// 터미널로 연결
		var out io.Writer = os.Stdout

		if install {
			// 1) 저장 디렉터리 결정
			dir := completionDir
			if dir == "" {
				dir = defaultCompletionDir(shell)
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}

			// 2) 스크립트 파일로 생성
			filePath := filepath.Join(dir, fmt.Sprintf("%s.%s", binName, shell))
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
			if err != nil {
				return err
			}
			defer f.Close()
			// 파일로 연결
			out = f

			// 생성
			if err := generateCompletion(shell, out); err != nil {
				return err
			}
			fmt.Printf("✅ %s completion script written to %s\n", shell, filePath)

			// 3) 셸 설정 파일에 source 라인 추가
			rcFile := shellRcFile(shell)
			line := fmt.Sprintf("\n# %s autocomplete\ntest -f %s && source %s\n", binName, filePath, filePath)
			if err := appendUniqueLine(rcFile, line); err != nil {
				return err
			}
			fmt.Printf("✅ Added source line to %s\n", rcFile)

			fmt.Println("로그아웃/로그인 또는 'exec $SHELL' 후에 자동 완성이 활성화됩니다.")
			return nil
		}

		// 단순 생성만
		return generateCompletion(shell, os.Stdout)
	},
}

func init() {
	// 플래그 정의
	// --install 을 주면 install 변수에 true, --install 안주면 false 가 저장됨.
	completionCmd.Flags().BoolVarP(&install, "install", "i", false,
		"자동 완성 스크립트를 생성하고 쉘 설정 파일에 설치합니다")
	completionCmd.Flags().StringVar(&completionDir, "dir", "",
		"자동 완성 스크립트를 저장할 디렉터리 (기본: ~/.<shell>_completion.d)")

}

func generateCompletion(shell string, out io.Writer) error {
	switch shell {
	case "bash":
		return RootCmd.GenBashCompletion(out)
	case "zsh":
		return RootCmd.GenZshCompletion(out)
	case "fish":
		return RootCmd.GenFishCompletion(out, true)
	case "powershell":
		return RootCmd.GenPowerShellCompletionWithDesc(out)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}

func appendUniqueLine(filePath, line string) error {
	// 파일이 없으면 생성
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, []byte(line), 0o644); err != nil {
			return err
		}
		return nil
	}
	// 중복 검사
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(line) {
			f.Close()
			return nil // 이미 존재
		}
	}
	f.Close()
	// 추가
	f2, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f2.Close()
	if _, err := f2.WriteString(line); err != nil {
		return err
	}
	return nil
}

func shellRcFile(shell string) string {
	home := homeDir()
	switch shell {
	case "bash":
		return filepath.Join(home, ".bashrc")
	case "zsh":
		return filepath.Join(home, ".zshrc")
	default:
		// fish, powershell 등은 별도 조정 필요
		return filepath.Join(home, ".bashrc")
	}
}

func homeDir() string {
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	return os.Getenv("HOME")
}

func defaultCompletionDir(shell string) string {
	// ~/.<shell>_completion.d 예: ~/.bash_completion.d
	return filepath.Join(homeDir(), fmt.Sprintf(".%s_completion.d", shell))
}
