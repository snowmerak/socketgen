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

SocketGen generates idiomatic code for each language, including:
1.  **Dispatcher:** Routes incoming packets to the correct handler method.
2.  **Handler Interface:** Defines the methods you need to implement.
3.  **PacketStream Interface:** Abstraction for reading/writing packets (you implement the network layer).
4.  **Serve Loop:** A helper to continuously read and dispatch packets.
5.  **Send Helpers:** Type-safe functions to wrap and send messages.

<details open>
<summary><strong>Go</strong></summary>

```go
// 1. Handler Interface
type PacketHandler interface {
    OnLoginReq(header *Header, msg *LoginReq)
    OnLoginRes(header *Header, msg *LoginRes)
    OnChatMsg(header *Header, msg *ChatMsg)
}

// 2. PacketStream Interface (Implement this for TCP/WebSocket)
type PacketStream interface {
    ReadPacket() ([]byte, error)
    WritePacket([]byte) error
}

// 3. Serve Loop
func Serve(stream PacketStream, handler PacketHandler) error {
    for {
        data, err := stream.ReadPacket()
        if err != nil {
            return err
        }
        if err := Dispatch(data, handler); err != nil {
            fmt.Println(fmt.Errorf("dispatch error: %w", err))
            continue
        }
    }
}

// 4. Send Helper
func SendLoginReq(stream PacketStream, header *Header, msg *LoginReq) error {
    pkt := &GamePacket{
        Header: header,
        Payload: &GamePacket_LoginReq{LoginReq: msg},
    }
    data, err := proto.Marshal(pkt)
    if err != nil {
        return err
    }
    return stream.WritePacket(data)
}
```
</details>

<details>
<summary><strong>TypeScript</strong></summary>

```typescript
export interface IPacketHandler {
  onLoginReq(header: Header, msg: LoginReq): void;
  // ...
}

export interface IPacketStream {
  readPacket(): Promise<Uint8Array>;
  writePacket(data: Uint8Array): Promise<void>;
}

export async function serve(stream: IPacketStream, handler: IPacketHandler) {
  while (true) {
    const data = await stream.readPacket();
    dispatch(data, handler);
  }
}

export async function sendLoginReq(stream: IPacketStream, header: Header, msg: LoginReq): Promise<void> {
  const pkt = GamePacket.fromPartial({
    header: header,
    loginReq: msg,
  });
  const data = GamePacket.encode(pkt).finish();
  await stream.writePacket(data);
}
```
</details>

<details>
<summary><strong>Python</strong></summary>

```python
class PacketStream(ABC):
    @abstractmethod
    def read_packet(self) -> bytes:
        pass
    @abstractmethod
    def write_packet(self, data: bytes):
        pass

def serve(stream: PacketStream, handler: PacketHandler):
    while True:
        data = stream.read_packet()
        try:
            dispatch(data, handler)
        except Exception as e:
            print(f"Dispatch error: {e}")

def send_login_req(stream: PacketStream, header, msg):
    pkt = GamePacket()
    pkt.header.CopyFrom(header)
    pkt.login_req.CopyFrom(msg)
    stream.write_packet(pkt.SerializeToString())
```
</details>

<details>
<summary><strong>C#</strong></summary>

```csharp
public interface IPacketStream {
    byte[] ReadPacket();
    void WritePacket(byte[] data);
}

public static class PacketDispatcher {
    public static void Serve(IPacketStream stream, IPacketHandler handler) {
        while (true) {
            var data = stream.ReadPacket();
            try {
                Dispatch(data, handler);
            } catch (System.Exception e) {
                System.Console.WriteLine($"Dispatch error: {e}");
            }
        }
    }

    public static void SendLoginReq(IPacketStream stream, Header header, LoginReq msg) {
        var pkt = new GamePacket {
            Header = header,
            LoginReq = msg
        };
        stream.WritePacket(pkt.ToByteArray());
    }
}
```
</details>

<details>
<summary><strong>Java</strong></summary>

```java
interface PacketStream {
    byte[] readPacket() throws java.io.IOException;
    void writePacket(byte[] data) throws java.io.IOException;
}

class PacketDispatcher {
    public static void serve(PacketStream stream, PacketHandler handler) {
        while (true) {
            try {
                byte[] data = stream.readPacket();
                dispatch(data, handler);
            } catch (Exception e) {
                System.err.println("Dispatch error: " + e.getMessage());
            }
        }
    }

    public static void sendLoginReq(PacketStream stream, Header header, LoginReq msg) throws java.io.IOException {
        GamePacket pkt = GamePacket.newBuilder()
            .setHeader(header)
            .setLoginReq(msg)
            .build();
        stream.writePacket(pkt.toByteArray());
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

MPL 2.0
