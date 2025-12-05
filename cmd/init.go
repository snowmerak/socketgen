package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project with a basic packet.proto",
	Long:  `Creates a 'packet.proto' file with the standard structure required by SocketGen.`,
	Run: func(cmd *cobra.Command, args []string) {
		content := `syntax = "proto3";
package packet;

option go_package = "./;packet";

// [헤더]: 모든 패킷에 포함될 메타데이터
message Header {
  int64 timestamp = 1;
  string request_id = 2;
}

// [페이로드]: 실제 전송할 데이터들 (사용자가 추가하는 부분)
message LoginReq { string id = 1; string pw = 2; }
message LoginRes { bool success = 1; }
message ChatMsg  { string text = 1; }

// [패킷 래퍼]: 네트워크 전송 단위
message GamePacket {
  Header header = 1;

  // 도구는 이 'oneof'를 파싱하여 분기문을 작성합니다.
  oneof payload {
    LoginReq login_req = 10;
    LoginRes login_res = 11;
    ChatMsg chat_msg = 12;
  }
}
`
		filename := "packet.proto"
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("Error: '%s' already exists.\n", filename)
			return
		}

		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}

		fmt.Printf("Created '%s' with basic structure.\n", filename)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
