package cmd

import (
	"fmt"

	"github.com/snowmerak/socketgen/parser"
	"github.com/spf13/cobra"
)

var (
	languages []string
	outDir    string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code for selected languages",
	Long:  `Generates Dispatcher and Handler code based on packet.proto for the specified languages.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generating code for languages: %v\n", languages)
		fmt.Printf("Output directory: %s\n", outDir)

		// Parse packet.proto
		result, err := parser.Parse("packet.proto")
		if err != nil {
			fmt.Printf("Error parsing packet.proto: %v\n", err)
			return
		}

		fmt.Printf("Found package: %s\n", result.PackageName)
		fmt.Println("Detected payloads:")
		for _, p := range result.Payloads {
			fmt.Printf(" - %s (Field: %s, Type: %s)\n", p.Name, p.FieldName, p.FullName)
		}

		// TODO: Implement code generation logic
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringSliceVar(&languages, "lang", []string{}, "Target languages (go, ts, python, csharp)")
	genCmd.Flags().StringVar(&outDir, "out", "./gen", "Output directory")

	genCmd.MarkFlagRequired("lang")
}
