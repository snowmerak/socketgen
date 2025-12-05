# SocketGen

SocketGen is a CLI tool designed to automate the creation of WebSocket message routing (Dispatcher) and handler interfaces based on Protobuf definitions. It aims to maximize development productivity by generating boilerplate code for multiple languages.

## Features

*   **Multi-Language Support:** Generates code for Go, TypeScript, Python, C#, Dart, PHP, Ruby, Kotlin, and Java.
*   **Protobuf Integration:** Parses `.proto` files to understand message structures and `oneof` payloads.
*   **Dispatcher Generation:** Automatically creates switch-case logic to route messages to appropriate handler methods.
*   **Handler Interface:** Defines clear interfaces for handling each message type, ensuring type safety and consistency.
*   **Protoc Execution:** Can optionally run `protoc` to generate the underlying Protobuf binding code.

## Installation

```bash
go install github.com/snowmerak/socketgen@latest
```

Or build from source:

```bash
git clone https://github.com/snowmerak/socketgen.git
cd socketgen
go build -o socketgen
```

## Usage

### 1. Initialize Project

Create a basic `packet.proto` template:

```bash
socketgen init
```

This will create a `packet.proto` file with a standard structure including a `Header` and a `GamePacket` wrapper with a `oneof` payload.

### 2. Define Messages

Edit `packet.proto` to add your own messages.

```protobuf
message MyMessage {
  string content = 1;
}

message GamePacket {
  Header header = 1;
  oneof payload {
    // ... existing ...
    MyMessage my_message = 13;
  }
}
```

### 3. Generate Code

Run the `gen` command to generate Dispatcher and Handler code. You can specify multiple languages.

```bash
socketgen gen --lang=go,ts,python,csharp,java --out=./gen
```

**Options:**

*   `--lang`: Comma-separated list of target languages.
    *   Supported: `go`, `ts`, `python`, `csharp`, `dart`, `php`, `ruby`, `kotlin`, `java`
*   `--out`: Output directory (default: `./gen`).
*   `--protoc`: If set, also runs `protoc` to generate the Protobuf binding code (requires `protoc` installed).

### Example

```bash
# Generate Go and TypeScript code, including Protobuf bindings
socketgen gen --lang=go,ts --out=./output --protoc
```

## Prerequisites

*   **Protobuf Compiler (`protoc`):** Required if you use the `--protoc` flag.
*   **Language Plugins:** Some languages require specific `protoc` plugins:
    *   **Go:** `protoc-gen-go`
    *   **TypeScript:** `ts-proto` (`npm install -g ts-proto`)
    *   **Dart:** `protoc-gen-dart` (`pub global activate protoc_plugin`)
    *   **Kotlin:** `protoc-gen-kotlin`

## License

MIT
