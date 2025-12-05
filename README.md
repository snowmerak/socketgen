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

## 🚀 Generated Code Examples

SocketGen generates idiomatic code for each language.

<details open>
<summary><strong>Go</strong></summary>

```go
// 1. Handler Interface
type PacketHandler interface {
    OnLoginReq(header *Header, msg *LoginReq)
    OnLoginRes(header *Header, msg *LoginRes)
    OnChatMsg(header *Header, msg *ChatMsg)
}

// 2. Dispatcher
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
</details>

<details>
<summary><strong>TypeScript</strong></summary>

```typescript
export interface IPacketHandler {
  onLoginReq(header: Header, msg: LoginReq): void;
  onLoginRes(header: Header, msg: LoginRes): void;
  onChatMsg(header: Header, msg: ChatMsg): void;
}

export function dispatch(data: Uint8Array, handler: IPacketHandler) {
  const pkt = GamePacket.decode(data);
  if (pkt.loginReq) {
    handler.onLoginReq(pkt.header!, pkt.loginReq!);
  }
  else if (pkt.loginRes) {
    handler.onLoginRes(pkt.header!, pkt.loginRes!);
  }
  else if (pkt.chatMsg) {
    handler.onChatMsg(pkt.header!, pkt.chatMsg!);
  }
}
```
</details>

<details>
<summary><strong>Python</strong></summary>

```python
class PacketHandler(ABC):
    @abstractmethod
    def on_login_req(self, header, msg):
        pass
    @abstractmethod
    def on_login_res(self, header, msg):
        pass
    @abstractmethod
    def on_chat_msg(self, header, msg):
        pass

def dispatch(data: bytes, handler: PacketHandler):
    pkt = GamePacket()
    pkt.ParseFromString(data)
    
    type_str = pkt.WhichOneof('payload')
    if type_str == 'login_req':
        handler.on_login_req(pkt.header, pkt.login_req)
    elif type_str == 'login_res':
        handler.on_login_res(pkt.header, pkt.login_res)
    elif type_str == 'chat_msg':
        handler.on_chat_msg(pkt.header, pkt.chat_msg)
```
</details>

<details>
<summary><strong>C#</strong></summary>

```csharp
public interface IPacketHandler {
    void OnLoginReq(Header header, LoginReq msg);
    void OnLoginRes(Header header, LoginRes msg);
    void OnChatMsg(Header header, ChatMsg msg);
}

public static class PacketDispatcher {
    public static void Dispatch(byte[] data, IPacketHandler handler) {
        var pkt = GamePacket.Parser.ParseFrom(data);
        
        switch (pkt.PayloadCase) {
            case GamePacket.PayloadOneofCase.LoginReq:
                handler.OnLoginReq(pkt.Header, pkt.LoginReq);
                break;
            case GamePacket.PayloadOneofCase.LoginRes:
                handler.OnLoginRes(pkt.Header, pkt.LoginRes);
                break;
            case GamePacket.PayloadOneofCase.ChatMsg:
                handler.OnChatMsg(pkt.Header, pkt.ChatMsg);
                break;
        }
    }
}
```
</details>

<details>
<summary><strong>Java</strong></summary>

```java
public interface PacketHandler {
    void onLoginReq(Header header, LoginReq msg);
    void onLoginRes(Header header, LoginRes msg);
    void onChatMsg(Header header, ChatMsg msg);
}

class PacketDispatcher {
    public static void dispatch(byte[] data, PacketHandler handler) throws InvalidProtocolBufferException {
        GamePacket pkt = GamePacket.parseFrom(data);
        
        switch (pkt.getPayloadCase()) {
            case LOGIN_REQ:
                handler.onLoginReq(pkt.getHeader(), pkt.getLoginReq());
                break;
            case LOGIN_RES:
                handler.onLoginRes(pkt.getHeader(), pkt.getLoginRes());
                break;
            case CHAT_MSG:
                handler.onChatMsg(pkt.getHeader(), pkt.getChatMsg());
                break;
        }
    }
}
```
</details>

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
