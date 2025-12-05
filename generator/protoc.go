package generator

import (
	"fmt"
	"os"
	"os/exec"
)

// GenerateProtoc runs the protoc command for the specified languages
func GenerateProtoc(protoFile string, languages []string, outDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, lang := range languages {
		var args []string

		switch lang {
		case "go":
			// Requires protoc-gen-go installed
			// --go_out=. --go_opt=paths=source_relative is a common pattern
			// We use outDir as the output base
			args = []string{
				"--go_out=" + outDir,
				"--go_opt=paths=source_relative",
				protoFile,
			}
		case "python":
			// Built-in support
			args = []string{
				"--python_out=" + outDir,
				protoFile,
			}
		case "csharp":
			// Built-in support
			args = []string{
				"--csharp_out=" + outDir,
				protoFile,
			}
		case "ts":
			// Uses ts-proto plugin (npm install -g ts-proto)
			// The plugin binary 'protoc-gen-ts_proto' must be in PATH.
			args = []string{
				"--ts_proto_out=" + outDir,
				"--ts_proto_opt=esModuleInterop=true",
				protoFile,
			}
		case "dart":
			// Requires protoc-gen-dart (pub global activate protoc_plugin)
			args = []string{
				"--dart_out=" + outDir,
				protoFile,
			}
		case "php":
			// Built-in support
			args = []string{
				"--php_out=" + outDir,
				protoFile,
			}
		case "ruby":
			// Built-in support
			args = []string{
				"--ruby_out=" + outDir,
				protoFile,
			}
		case "kotlin":
			// Requires protoc-gen-kotlin and usually java_out as well since Kotlin generated code depends on Java
			// We will generate both java and kotlin code in the output directory
			args = []string{
				"--java_out=" + outDir,
				"--kotlin_out=" + outDir,
				protoFile,
			}
		default:
			continue
		}

		cmd := exec.Command("protoc", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Running protoc for %s: %s\n", lang, cmd.String())
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate protobuf code for %s: %w", lang, err)
		}
	}

	return nil
}
