# SocketGen ⚡️

> **Turn your `.proto` files into type-safe WebSocket dispatchers instantly.**

SocketGen is a CLI tool designed to automate the creation of WebSocket message routing (Dispatcher) and handler interfaces based on Protobuf definitions. It eliminates the boilerplate of writing manual `switch-case` statements and ensures your server and client logic are always in sync.

## Why SocketGen?

Handling WebSocket messages with Protobuf usually involves:

1.  Defining messages in `.proto`.
2.  Parsing the raw bytes.
3.  Reading a header to determine the message type.
4.  **Writing a giant `switch-case` statement to route messages.** 😫
5.  Manually casting payloads.

**SocketGen automates steps 3 through 5.** It parses your `oneof` structure and generates a strongly-typed Dispatcher and Handler Interface for you.

## Features

  * **Multi-Language Support:** Generates code for **Go, TypeScript, Python, C#, Dart, PHP, Ruby, Kotlin, and Java**.
  * **Boilerplate-Free:** No more manual routing logic. Just implement the interface.
  * **Type Safety:** Ensures handlers receive the correct message types at compile time.
  * **Protoc Integration:** Can optionally run `protoc` to generate the underlying Protobuf binding code in one go.

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

## The Protocol Pattern

SocketGen expects a specific structure in your `.proto` file to work its magic. It relies on the **Wrapper Message** pattern using `oneof`.

```protobuf
// packet.proto

// 1. Common Header
message Header {
  int64 timestamp = 1;
  string request_id = 2;
}

// 2. The Wrapper (This is what SocketGen looks for)
message GamePacket {
  Header header = 1; // Always included

  // 3. Payload Oneof (Routes are generated based on this)
  oneof payload {
    LoginReq login_req = 10;
    LoginRes login_res = 11;
    ChatMsg chat_msg = 12;
  }
}
```

## Usage

### 1. Initialize Project

Create a starter template:

```bash
socketgen init
```

This creates a `packet.proto` with the standard structure shown above.

### 2. Generate Code

Run the `gen` command. You can target multiple languages at once.

```bash
socketgen gen --lang=go,ts,csharp --out=./gen --protoc
```

  * `--lang`: Comma-separated list of target languages.
  * `--out`: Output directory (default: `./gen`).
  * `--protoc`: (Optional) Automatically runs `protoc` to generate the base struct/class files.

-----

## 🚀 Generated Code Example

Here is what SocketGen creates for you (e.g., in **Go**):

**1. The Handler Interface** (You implement this)

```go
type PacketHandler interface {
    OnLoginReq(header *Header, msg *LoginReq)
    OnLoginRes(header *Header, msg *LoginRes)
    OnChatMsg(header *Header, msg *ChatMsg)
}
```

**2. The Dispatcher** (Auto-generated)

```go
func Dispatch(data []byte, handler PacketHandler) error {
    pkt := &GamePacket{}
    if err := proto.Unmarshal(data, pkt); err != nil {
        return err
    }

    switch payload := pkt.Payload.(type) {
    case *GamePacket_LoginReq:
        handler.OnLoginReq(pkt.Header, payload.LoginReq)
    case *GamePacket_LoginRes:
        handler.OnLoginRes(pkt.Header, payload.LoginRes)
    case *GamePacket_ChatMsg:
        handler.OnChatMsg(pkt.Header, payload.ChatMsg)
    default:
        return fmt.Errorf("unknown packet type")
    }
    return nil
}
```

-----

## Prerequisites

If you use the `--protoc` flag, you must have the Protocol Buffers compiler and relevant plugins installed:

  * **Protobuf Compiler:** `protoc` ([Install Guide](https://grpc.io/docs/protoc-installation/))
  * **Go:** `protoc-gen-go` (`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`)
  * **TypeScript:** `ts-proto` (`npm install -g ts-proto`)
  * **Dart:** `protoc-gen-dart`
  * **Kotlin/Java:** Standard `protoc` support.

## License

MIT
