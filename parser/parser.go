package parser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// PayloadMessage represents a message type that can be carried in the GamePacket payload
type PayloadMessage struct {
	Name      string // The type name (e.g., "LoginReq")
	FieldName string // The field name in the oneof (e.g., "login_req")
	FullName  string // The full proto name (e.g., "packet.LoginReq")
}

// ParseResult holds the extracted information from the proto file
type ParseResult struct {
	PackageName string
	Payloads    []PayloadMessage
}

// Parse runs protoc to generate a descriptor set and then parses it to extract GamePacket info
func Parse(protoFile string) (*ParseResult, error) {
	// 1. Check if protoc is installed
	_, err := exec.LookPath("protoc")
	if err != nil {
		return nil, fmt.Errorf("protoc is not installed or not in PATH. Please install Protocol Buffers compiler")
	}

	// 2. Generate FileDescriptorSet using protoc
	// We output to a temporary file
	tmpFile := "temp_descriptor.pb"
	defer os.Remove(tmpFile)

	cmd := exec.Command("protoc",
		"--descriptor_set_out="+tmpFile,
		"--include_imports",
		protoFile,
	)

	// Capture stderr to show protoc errors if any
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run protoc: %w", err)
	}

	// 3. Read the generated descriptor file
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read descriptor file: %w", err)
	}

	// 4. Unmarshal into FileDescriptorSet
	var fileDescSet descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fileDescSet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal descriptor set: %w", err)
	}

	// 5. Analyze the descriptor to find GamePacket and its payload
	return analyzeDescriptor(&fileDescSet, protoFile)
}

func analyzeDescriptor(fds *descriptorpb.FileDescriptorSet, targetFile string) (*ParseResult, error) {
	var targetFileDesc *descriptorpb.FileDescriptorProto

	// Find the descriptor for the target file
	// Note: protoFile path might need normalization to match what's in the descriptor set
	// For simplicity, we'll look for the file that matches the input filename (base name)

	// Simple strategy: Look for the file that matches the input filename (base name)
	targetBase := filepath.Base(targetFile)
	for _, fd := range fds.File {
		if strings.HasSuffix(fd.GetName(), targetBase) {
			targetFileDesc = fd
			break
		}
	}

	if targetFileDesc == nil {
		// Fallback: use the last one (often the main file in simple cases)
		if len(fds.File) > 0 {
			targetFileDesc = fds.File[len(fds.File)-1]
		} else {
			return nil, fmt.Errorf("no file descriptors found")
		}
	}

	result := &ParseResult{
		PackageName: targetFileDesc.GetPackage(),
		Payloads:    []PayloadMessage{},
	}

	// Find "GamePacket" message
	var gamePacketMsg *descriptorpb.DescriptorProto
	for _, msg := range targetFileDesc.MessageType {
		if msg.GetName() == "GamePacket" {
			gamePacketMsg = msg
			break
		}
	}

	if gamePacketMsg == nil {
		return nil, fmt.Errorf("message 'GamePacket' not found in %s", targetFile)
	}

	// Find "payload" oneof field
	// In DescriptorProto, OneofDecl contains the names of oneofs.
	// Field contains the fields, which refer to OneofIndex.

	oneofIndex := -1
	for i, oneof := range gamePacketMsg.OneofDecl {
		if oneof.GetName() == "payload" {
			oneofIndex = i
			break
		}
	}

	if oneofIndex == -1 {
		return nil, fmt.Errorf("'payload' oneof field not found in GamePacket")
	}

	// Collect fields belonging to this oneof
	for _, field := range gamePacketMsg.Field {
		if field.OneofIndex != nil && int(*field.OneofIndex) == oneofIndex {
			// This field is part of the payload oneof

			// TypeName usually returns ".package.MessageName"
			fullType := field.GetTypeName()
			typeName := fullType
			if lastDot := strings.LastIndex(fullType, "."); lastDot != -1 {
				typeName = fullType[lastDot+1:]
			}

			result.Payloads = append(result.Payloads, PayloadMessage{
				Name:      typeName,
				FieldName: field.GetName(),
				FullName:  strings.TrimPrefix(fullType, "."),
			})
		}
	}

	return result, nil
}
