package cmd

import (
	"github.com/spf13/cobra"
)

var peerCmd = &cobra.Command{
	Use:   "peer",
	Short: "P2P PeerNode 관리",
}

var peerConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "모든 피어에 연결",
	RunE: func(_ *cobra.Command, _ []string) error {
		// 이 코드 블록은 REPL 외에도 일반 CLI(비-REPL)에서
		// `tori peer connect` 로 호출됐을 때 실행됩니다.
		// TODO 일단 주석 처리함.
		// conns := globalPeerNode.ConnectAll(context.Background())
		// fmt.Printf("Connected to %d peers\n", len(conns))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(peerCmd)
	peerCmd.AddCommand(peerConnectCmd)
}
