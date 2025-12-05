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
			// Requires a plugin. We'll assume 'protoc-gen-ts' or similar is available via built-in or npm.
			// Common plugins: ts-proto, protoc-gen-ts
			// For now, let's try to use a generic plugin flag if the user has one,
			// but since we can't guess, we might skip or try a common one.
			// Let's try 'ts_out' which implies protoc-gen-ts is in PATH.
			// If it fails, we'll return an error.
			args = []string{
				"--ts_out=" + outDir,
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
