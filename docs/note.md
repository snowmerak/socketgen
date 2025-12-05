
# 📄 [설계 문서 v2] SocketGen

## 1\. 개요 (Overview)

  * **목표:** WebSocket 통신을 위한 Protobuf 메시지 정의부터, 각 언어별 **메시지 라우팅(Dispatcher) 및 핸들러(Handler)** 코드를 자동으로 생성하여 개발 생산성을 극대화한다.
  * **형태:** 독립 실행 가능한 CLI (Command Line Interface) 도구
  * **구현 언어 (Base):** **Go (Golang)** - 빠르고, 멀티 플랫폼 빌드가 용이하며, Protobuf 파싱 라이브러리가 강력함.
  * **지원 언어 (First-party):** **Go, TypeScript, Python, C\#**

## 2\. 워크플로우 (User Experience)

개발자는 단 두 개의 명령어만 알면 됩니다: `init`과 `gen`.

```bash
# 1. 프로젝트 초기화 (기본 템플릿 생성)
$ ps-cli init
> Created 'packet.proto' with basic structure.

# ... 사용자가 packet.proto에 메시지 추가 ...

# 2. 코드 생성 (언어 선택)
$ ps-cli gen --lang=go,ts,python,csharp --out=./gen
> Generating Go code... Done.
> Generating TypeScript code... Done.
> Generating Python code... Done.
> Generating C# code... Done.
```

## 3\. 프로토콜 구조 (Standard Template)

`ps-cli init` 명령어가 생성해주는 기본 `.proto` 파일 구조입니다.

```protobuf
syntax = "proto3";
package packet;

// [헤더]: 모든 패킷에 포함될 메타데이터
message Header {
  int64 timestamp = 1;
  string request_id = 2;
}

// [페이로드]: 실제 전송할 데이터들 (사용자가 추가하는 부분)
message LoginReq { string id = 1; string pw = 2; }
message LoginRes { bool success = 1; }
message ChatMsg  { string text = 1; }

// [패킷 래퍼]: 네트워크 전송 단위
message GamePacket {
  Header header = 1;

  // 도구는 이 'oneof'를 파싱하여 분기문을 작성합니다.
  oneof payload {
    LoginReq login_req = 10;
    LoginRes login_res = 11;
    ChatMsg chat_msg = 12;
  }
}
```

## 4\. 언어별 생성 전략 (Generation Strategy)

CLI는 `.proto` 파일을 파싱한 뒤, 각 언어의 문법에 맞는 \*\*Dispatcher(분배기)\*\*와 **Handler Interface**를 생성합니다.

### A. Go (Server/Client)

  * **특징:** `interface`를 활용한 핸들러 정의.
  * **생성 파일:** `packet_dispatcher.go`
  * **구조:**
    ```go
    // Handler Interface
    type PacketHandler interface {
        OnLoginReq(header *Header, msg *LoginReq)
        OnChatMsg(header *Header, msg *ChatMsg)
        // ...
    }

    // Dispatcher
    func Dispatch(data []byte, handler PacketHandler) error {
        pkt := &GamePacket{}
        proto.Unmarshal(data, pkt)
        
        switch payload := pkt.Payload.(type) {
        case *GamePacket_LoginReq:
            handler.OnLoginReq(pkt.Header, payload.LoginReq)
        case *GamePacket_ChatMsg:
            handler.OnChatMsg(pkt.Header, payload.ChatMsg)
        }
        return nil
    }
    ```

### B. TypeScript (Web/Node.js)

  * **특징:** Strict Typing 지원. `protobuf.js` 또는 `ts-proto` 기반 객체 활용.
  * **생성 파일:** `PacketDispatcher.ts`
  * **구조:**
    ```typescript
    export interface IPacketHandler {
      onLoginReq(header: Header, msg: LoginReq): void;
      onChatMsg(header: Header, msg: ChatMsg): void;
    }

    export function dispatch(data: Uint8Array, handler: IPacketHandler) {
      const pkt = GamePacket.decode(data);
      
      if (pkt.loginReq) handler.onLoginReq(pkt.header, pkt.loginReq);
      else if (pkt.chatMsg) handler.onChatMsg(pkt.header, pkt.chatMsg);
    }
    ```

### C. Python (AI/Server/Tool)

  * **특징:** Type Hinting 적극 활용. `ABC`(Abstract Base Class)로 인터페이스 강제.
  * **생성 파일:** `packet_dispatcher.py`
  * **구조:**
    ```python
    from abc import ABC, abstractmethod
    from .packet_pb2 import GamePacket

    class PacketHandler(ABC):
        @abstractmethod
        def on_login_req(self, header, msg): pass
        @abstractmethod
        def on_chat_msg(self, header, msg): pass

    def dispatch(data: bytes, handler: PacketHandler):
        pkt = GamePacket()
        pkt.ParseFromString(data)
        
        type_str = pkt.WhichOneof('payload')
        if type_str == 'login_req':
            handler.on_login_req(pkt.header, pkt.login_req)
        elif type_str == 'chat_msg':
            handler.on_chat_msg(pkt.header, pkt.chat_msg)
    ```

### D. C\# (Unity/Server)

  * **특징:** `partial class`나 인터페이스 활용. Unity 호환성 고려.
  * **생성 파일:** `PacketDispatcher.cs`
  * **구조:**
    ```csharp
    public interface IPacketHandler {
        void OnLoginReq(Header header, LoginReq msg);
        void OnChatMsg(Header header, ChatMsg msg);
    }

    public static class PacketDispatcher {
        public static void Dispatch(byte[] data, IPacketHandler handler) {
            var pkt = GamePacket.Parser.ParseFrom(data);
            
            switch (pkt.PayloadCase) {
                case GamePacket.PayloadOneofCase.LoginReq:
                    handler.OnLoginReq(pkt.Header, pkt.LoginReq);
                    break;
                case GamePacket.PayloadOneofCase.ChatMsg:
                    handler.OnChatMsg(pkt.Header, pkt.ChatMsg);
                    break;
            }
        }
    }
    ```

## 5\. 기술 스택 및 라이브러리 (Tech Stack)

  * **Main Language:** Go 1.21+
  * **CLI Framework:** `github.com/spf13/cobra` (명령어 관리의 표준)
  * **Protobuf Parser:** `google.golang.org/protobuf/reflect/protoreflect` (Proto 파일의 구조를 동적으로 읽기 위함)
      * *전략:* `.proto` 텍스트 자체를 파싱하기보다, `protoc`를 실행시켜 나오는 `FileDescriptorSet` 바이너리를 Go에서 읽어들이는 방식이 가장 정확합니다.
  * **Template Engine:** Go 내장 `text/template` (강력하고 외부 의존성 없음)

## 6\. 개발 마일스톤

1.  **Phase 1 (Skeleton):**
      * Go 프로젝트 세팅 및 Cobra로 `init`, `gen` 명령어 껍데기 구현.
      * `init` 실행 시 기본 `packet.proto` 파일 생성 기능 구현.
2.  **Phase 2 (Parser):**
      * `protoc`를 서브 프로세스로 호출하여 `.proto` 파일 정보를 읽어오는 로직 구현.
      * `GamePacket` 내의 `oneof` 필드 리스트 추출 로직 구현.
3.  **Phase 3 (Generator - Go/TS):**
      * 가장 수요가 많은 Go와 TypeScript용 템플릿 작성 및 코드 생성 기능 구현.
4.  **Phase 4 (Generator - C\#/Py):**
      * C\#과 Python용 템플릿 추가.
      * 최종 테스트 및 문서화.
