# 로깅 미들웨어

로깅 미들웨어는 API 요청에 대한 로그를 생성하고 콘솔에 출력하거나 원격 URL로 전송하는 기능을 제공합니다.

## 기본 사용법

가장 기본적인 사용법은 다음과 같습니다:

```go
package main

import (
    server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
    s, _ := server.NewServer(server.FrameworkGin, "8080")

    // 기본 로깅 미들웨어 추가 (콘솔에만 로그 출력)
    loggingMiddleware := s.GetLoggingMiddleware()
    s.Use(loggingMiddleware.Middleware(nil))

    // 라우트 등록
    s.GET("/", func(c server.Context) {
        c.String(200, "Hello, World!")
    })

    s.Run()
}
```

## 커스텀 필드 추가

로그에 커스텀 필드를 추가하려면 다음과 같이 `LoggingConfig`를 사용합니다:

```go
loggingConfig := &server.LoggingConfig{
    CustomFields: map[string]string{
        "environment": "development",
        "version":     "0.0.1",
        "app_name":    "my-awesome-app",
    },
}
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(loggingConfig))
```

## 원격 로깅

로그를 원격 URL로 전송하려면 다음과 같이 `RemoteURL`을 설정합니다:

```go
remoteLoggingConfig := &server.LoggingConfig{
    RemoteURL: "https://your-logging-service.com/api/logs",
    LoggingToRemote: true,
    CustomFields: map[string]string{
        "environment": "production",
        "version":     "0.0.1",
    },
}
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(remoteLoggingConfig))
```

원격 로깅은 비동기적으로 처리되므로 API 응답 시간에 영향을 주지 않습니다.

## 특정 경로 무시하기

로깅 미들웨어는 특정 경로에 대한 로깅을 건너뛸 수 있습니다. `SkipPaths` 필드에 건너뛸 경로 목록을 설정하여 해당 경로에 대한 로깅을 비활성화할 수 있습니다:

```go
loggingConfig := &server.LoggingConfig{
    SkipPaths: []string{
        "/health",
        "/metrics",
        "/favicon.ico",
    },
    CustomFields: map[string]string{
        "version": "0.0.1",
    },
}
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(loggingConfig))
```

이렇게 하면 `/health`, `/metrics`, `/favicon.ico` 경로에 대한 요청은 로깅되지 않습니다.

## 콘솔 및 원격 로깅 설정

로깅 미들웨어는 `LoggingToConsole` 및 `LoggingToRemote` 필드를 통해 로깅 동작을 제어할 수 있습니다:

```go
// 콘솔에만 로그 출력 (기본값)
consoleOnlyConfig := &server.LoggingConfig{
    LoggingToConsole: true,
    LoggingToRemote: false,
    CustomFields: map[string]string{
        "version": "0.0.1",
    },
}
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(consoleOnlyConfig))

// 원격 URL에만 로그 출력
remoteOnlyConfig := &server.LoggingConfig{
    LoggingToConsole: false,
    LoggingToRemote: true,
    RemoteURL: "https://your-logging-service.com/api/logs",
    CustomFields: map[string]string{
        "version": "0.0.1",
    },
}
loggingMiddleware = s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(remoteOnlyConfig))

// 콘솔과 원격 URL 모두에 로그 출력
bothConfig := &server.LoggingConfig{
    LoggingToConsole: true,
    LoggingToRemote: true,
    RemoteURL: "https://your-logging-service.com/api/logs",
    CustomFields: map[string]string{
        "version": "0.0.1",
    },
}
loggingMiddleware = s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(bothConfig))
```

`LoggingToConsole`이 `true`로 설정되면 콘솔에 로그가 출력됩니다. `LoggingToRemote`이 `true`로 설정되고 `RemoteURL`이 제공되면 원격 URL로 로그가 전송됩니다. 두 필드 모두 `true`로 설정하면 콘솔과 원격 URL 모두에 로그가 출력됩니다.

기본적으로 `LoggingToConsole`은 `true`이고 `LoggingToRemote`는 `false`입니다.

## 기본 구성 사용하기

로깅 미들웨어는 기본 구성을 사용할 수 있습니다:

```go
// 프레임워크별 로깅 미들웨어 가져오기
loggingMiddleware := s.GetLoggingMiddleware()

// 기본 구성으로 미들웨어 추가
s.Use(loggingMiddleware.Middleware(nil))
```

또는 `DefaultLoggingConfig()`를 명시적으로 사용할 수 있습니다:

```go
// 프레임워크별 로깅 미들웨어 가져오기
loggingMiddleware := s.GetLoggingMiddleware()

// 기본 구성으로 미들웨어 추가
s.Use(loggingMiddleware.Middleware(middleware.DefaultLoggingConfig()))
```

기본 구성은 다음과 같은 값을 사용합니다:
- RemoteURL: "" (원격 로깅 비활성화)
- CustomFields: 빈 맵
- Environment: "dev" (콘솔 로깅 활성화)
- SkipPaths: 빈 슬라이스

## 로그 데이터 구조

로그 데이터는 다음과 같은 구조로 생성됩니다:

```go
type ApiLog struct {
    ClientIp      string            `json:"client_ip"`
    Timestamp     string            `json:"timestamp"`
    Method        string            `json:"method"`
    Path          string            `json:"path"`
    Protocol      string            `json:"protocol"`
    StatusCode    int               `json:"status_code"`
    Latency       int64             `json:"latency"`
    UserAgent     string            `json:"user_agent"`
    Error         string            `json:"error"`
    RequestId     string            `json:"request_id"`
    Authorization string            `json:"authorization"`
    CustomFields  map[string]string `json:"custom_fields,omitempty"`
}
```

각 필드의 의미는 다음과 같습니다:

- `ClientIp`: 클라이언트 IP 주소
- `Timestamp`: 요청 시간 (RFC3339 형식)
- `Method`: HTTP 메소드 (GET, POST 등)
- `Path`: 요청 경로
- `Protocol`: HTTP 프로토콜 버전
- `StatusCode`: HTTP 상태 코드 (기본값: 200)
- `Latency`: 요청 처리 시간 (밀리초)
- `UserAgent`: 사용자 에이전트 문자열
- `Error`: 오류 메시지 (오류가 없는 경우 "none"으로 설정됨)
- `RequestId`: 요청 ID (X-Request-ID 헤더에서 추출, 없으면 생성)
- `Authorization`: 인증 정보 (개발 환경에서는 전체 토큰이 로깅되고, 프로덕션 환경에서는 토큰이 마스킹 처리됨)
- `CustomFields`: 사용자 정의 필드

## 로그 출력 예시

콘솔에 출력되는 로그의 예시는 다음과 같습니다:

```json
{
  "client_ip": "127.0.0.1",
  "timestamp": "2023-06-01T12:34:56Z",
  "method": "GET",
  "path": "/api/users",
  "protocol": "HTTP/1.1",
  "status_code": 200,
  "latency": 42,
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
  "error": "none",
  "request_id": "1685622896123456789",
  "authorization": "Bearer [MASKED]",
  "custom_fields": {
    "environment": "development",
    "version": "0.0.1"
  }
}
```

## 프레임워크별 로깅 미들웨어

v1.1.0부터 tenqube-go-http-server는 프레임워크별 로깅 미들웨어를 제공하여 실제 HTTP 상태 코드와 오류 정보를 캡처할 수 있습니다. 이제 이러한 미들웨어는 각 프레임워크별 디렉토리에 구현되어 있습니다:

- `gin/middleware.go`: Gin 프레임워크를 위한 로깅 미들웨어
- `std/middleware.go`: 표준 HTTP 패키지를 위한 로깅 미들웨어
- `middleware/middleware.go`: 공통 로깅 미들웨어 기능

### Gin 프레임워크

Gin 프레임워크에서는 `gin.Context.Writer.Status()`를 사용하여 상태 코드를 캡처하고, `gin.Context.Errors`를 사용하여 오류 정보를 캡처합니다.

```go
s, _ := server.NewServer(server.FrameworkGin, "8080")
s.Use(s.GetLoggingMiddleware().Middleware(nil))
```

### 표준 HTTP

표준 HTTP에서는 `http.ResponseWriter`를 래핑하여 상태 코드를 캡처합니다.

```go
s, _ := server.NewServer(server.FrameworkStdHTTP, "8080")
s.Use(s.GetLoggingMiddleware().Middleware(nil))
```

### 기존 방식과의 호환성

기존 방식은 더 이상 지원되지 않습니다. 항상 프레임워크별 로깅 미들웨어를 직접 사용해야 합니다:

```go
s, _ := server.NewServer(server.FrameworkGin, "8080")
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(nil))
```

이 방식은 항상 올바른 프레임워크별 로깅 미들웨어를 사용하므로 상태 코드와 오류 정보를 정확하게 캡처할 수 있습니다.

## 미들웨어 순서

로깅 미들웨어는 응답 상태 코드를 정확하게 캡처하기 위해 **에러 핸들러 미들웨어 이후에** 추가해야 합니다. 이는 에러 핸들러 미들웨어가 상태 코드를 설정하고, 로깅 미들웨어가 이 상태 코드를 캡처해야 하기 때문입니다.

```go
// 에러 핸들러 미들웨어 추가 (먼저)
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
s.Use(errorHandlerMiddleware.Middleware(errorHandlerConfig))

// 로깅 미들웨어 추가 (나중에)
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(loggingConfig))
```

잘못된 순서로 미들웨어를 추가하면 오류 응답의 상태 코드가 로그에 정확하게 반영되지 않을 수 있습니다.

## 주의사항

### 미들웨어 순서의 중요성

미들웨어 순서가 중요합니다. 에러 핸들러 미들웨어는 로깅 미들웨어보다 먼저 추가되어야 합니다. 이렇게 해야 에러 핸들러 미들웨어가 설정한 상태 코드가 로깅 미들웨어에 의해 정확하게 캡처됩니다.

### 권장 사용법

가능하면 다음과 같이 프레임워크별 로깅 미들웨어를 직접 사용하는 것이 좋습니다:

```go
s, _ := server.NewServer(server.FrameworkGin, "8080")

// 에러 핸들러 미들웨어 추가 (먼저)
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
s.Use(errorHandlerMiddleware.Middleware(errorHandlerConfig))

// 로깅 미들웨어 추가 (나중에)
s.Use(s.GetLoggingMiddleware().Middleware(nil))
```

또는

```go
s, _ := server.NewServer(server.FrameworkStdHTTP, "8080")

// 에러 핸들러 미들웨어 추가 (먼저)
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
s.Use(errorHandlerMiddleware.Middleware(errorHandlerConfig))

// 로깅 미들웨어 추가 (나중에)
s.Use(s.GetLoggingMiddleware().Middleware(nil))
```

이 방식은 항상 올바른 프레임워크별 로깅 미들웨어를 사용하므로 상태 코드와 오류 정보를 정확하게 캡처할 수 있습니다.

### 기타 주의사항

- 원격 로깅 기능을 사용할 때는 로그 서버의 가용성을 고려해야 합니다. 로그 서버에 문제가 있어도 API 서버의 동작에는 영향을 주지 않습니다.
- 인증 토큰은 보안을 위해 마스킹 처리됩니다.

### 상태 코드 및 지연 시간 로깅 문제

v0.4.7 이전 버전에서는 일반 미들웨어 구현(`server.LoggingMiddleware`)에서 상태 코드가 실제로는 500 오류가 발생했는데도 200으로 로깅되고, 지연 시간(latency)이 0으로 로깅되는 문제가 있었습니다. 이 문제는 다음과 같은 이유로 발생했습니다:

1. 미들웨어가 `c.Next()`를 호출하여 다음 핸들러를 실행한 후, 즉시 지연 시간을 계산하고 상태 코드를 캡처했습니다.
2. 그러나 `c.Next()`는 다음 핸들러가 완료될 때까지 블록하지 않고, 단순히 제어를 다음 핸들러로 넘기고 계속 실행됩니다.
3. 이로 인해 핸들러가 완료되기 전에 지연 시간 계산과 상태 코드 캡처가 이루어져, 지연 시간이 0으로 기록되고 기본 상태 코드인 200이 로깅되었습니다.

v0.4.7부터는 이 문제를 해결하기 위해 다음과 같은 변경이 이루어졌습니다:

1. 지연 시간 계산과 상태 코드 캡처를 `defer` 함수로 이동하여 핸들러가 완료된 후에 실행되도록 했습니다.
2. `c.Next()` 호출 후 작은 지연(1 나노초)을 추가하여 `defer` 함수가 핸들러 완료 후에 실행되도록 보장했습니다.

이 변경으로 인해 상태 코드와 지연 시간이 정확하게 로깅됩니다. 특히 오류가 발생한 경우 실제 오류 상태 코드(예: 500)가 로그에 정확히 반영됩니다.

### 프레임워크별 로깅 미들웨어 사용 권장

v0.4.8부터는 상태 코드를 정확하게 캡처하기 위해 프레임워크별 로깅 미들웨어를 사용하는 것이 권장됩니다. 일반 미들웨어 구현은 더 이상 지원되지 않으므로, 항상 프레임워크별 로깅 미들웨어를 사용해야 합니다.

프레임워크별 로깅 미들웨어를 사용하려면 다음과 같이 `srv.GetLoggingMiddleware().Middleware`를 사용하세요:

```go
// 권장 방식 (항상 상태 코드를 정확하게 캡처)
loggingMiddleware := srv.GetLoggingMiddleware()
srv.Use(loggingMiddleware.Middleware(loggingConfig))
```

프레임워크별 로깅 미들웨어는 해당 프레임워크의 응답 작성자를 올바르게 래핑하여 상태 코드를 정확하게 캡처합니다. 이는 특히 400, 500 등의 오류 상태 코드를 로깅할 때 중요합니다.
