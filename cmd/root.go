package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "socketgen",
	Short: "SocketGen is a CLI tool for generating WebSocket packet dispatchers",
	Long: `SocketGen automates the creation of message routing (Dispatcher) and handler interfaces 
based on Protobuf definitions for Go, TypeScript, Python, and C#.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
