package cmd

import (
	"fmt"

	"github.com/snowmerak/socketgen/generator"
	"github.com/snowmerak/socketgen/parser"
	"github.com/spf13/cobra"
)

var (
	languages  []string
	outDir     string
	withProtoc bool
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code for selected languages",
	Long:  `Generates Dispatcher and Handler code based on packet.proto for the specified languages.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generating code for languages: %v\n", languages)
		fmt.Printf("Output directory: %s\n", outDir)

		// Run protoc if requested
		if withProtoc {
			fmt.Println("Running protoc...")
			if err := generator.GenerateProtoc("packet.proto", languages, outDir); err != nil {
				fmt.Printf("Warning: Failed to run protoc: %v\n", err)
				fmt.Println("Make sure you have 'protoc' and necessary plugins installed.")
			} else {
				fmt.Println("Successfully generated protobuf bindings.")
			}
		}

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

		for _, lang := range languages {
			var err error
			switch lang {
			case "go":
				fmt.Println("Generating Go code...")
				err = generator.GenerateGo(result, outDir)
			case "ts":
				fmt.Println("Generating TypeScript code...")
				err = generator.GenerateTS(result, outDir)
			case "python":
				fmt.Println("Generating Python code...")
				err = generator.GeneratePython(result, outDir)
			case "csharp":
				fmt.Println("Generating C# code...")
				err = generator.GenerateCSharp(result, outDir)
			case "dart":
				fmt.Println("Generating Dart code...")
				err = generator.GenerateDart(result, outDir)
			case "php":
				fmt.Println("Generating PHP code...")
				err = generator.GeneratePHP(result, outDir)
			case "ruby":
				fmt.Println("Generating Ruby code...")
				err = generator.GenerateRuby(result, outDir)
			case "kotlin":
				fmt.Println("Generating Kotlin code...")
				err = generator.GenerateKotlin(result, outDir)
			case "java":
				fmt.Println("Generating Java code...")
				err = generator.GenerateJava(result, outDir)
			default:
				fmt.Printf("Warning: Language '%s' is not supported yet.\n", lang)
				continue
			}

			if err != nil {
				fmt.Printf("Error generating %s code: %v\n", lang, err)
			} else {
				fmt.Printf("Successfully generated %s code.\n", lang)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringSliceVar(&languages, "lang", []string{}, "Target languages (go, ts, python, csharp, dart, php, ruby, kotlin, java)")
	genCmd.Flags().StringVar(&outDir, "out", "./gen", "Output directory")
	genCmd.Flags().BoolVar(&withProtoc, "protoc", false, "Generate protobuf bindings using protoc")

	genCmd.MarkFlagRequired("lang")
}
