# 타임아웃 미들웨어

타임아웃 미들웨어는 API 요청이 지정된 시간 내에 응답하지 않을 경우 자동으로 타임아웃 응답을 반환하는 기능을 제공합니다. 기본적으로 2초 후에 타임아웃 응답을 반환합니다.

## 기본 사용법

가장 기본적인 사용법은 다음과 같습니다:

```go
package main

import (
    server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
    s, _ := server.NewServer(server.FrameworkGin, "8080")

    // 기본 타임아웃 미들웨어 추가 (2초 후 타임아웃)
    s.Use(server.TimeoutMiddleware(nil))

    // 라우트 등록
    s.GET("/", func(c server.Context) {
        c.String(200, "Hello, World!")
    })

    s.Run()
}
```

## 타임아웃 시간 설정

타임아웃 시간을 설정하려면 다음과 같이 `TimeoutConfig`를 사용합니다:

```go
timeoutConfig := &server.TimeoutConfig{
    Timeout: 5 * time.Second, // 5초로 타임아웃 설정
}
s.Use(server.TimeoutMiddleware(timeoutConfig))
```

## 기본 생성자 함수

타임아웃 미들웨어는 기본 구성을 사용하는 생성자 함수를 제공합니다:

```go
// 기본 타임아웃 미들웨어 생성자 (2초 타임아웃)
s.Use(server.NewDefaultTimeoutMiddleware())
```

`NewDefaultTimeoutMiddleware` 함수는 `DefaultTimeoutConfig()`를 호출하여 기본 구성을 생성하고, 이를 `TimeoutMiddleware` 함수에 전달합니다. 이는 다음과 동일합니다:

```go
s.Use(server.TimeoutMiddleware(nil))
```

또는

```go
s.Use(server.TimeoutMiddleware(server.DefaultTimeoutConfig()))
```

기본 구성은 2초의 타임아웃 시간을 사용합니다.

## 타임아웃 동작 방식

타임아웃 미들웨어는 다음과 같이 동작합니다:

1. 요청이 들어오면 타임아웃 타이머를 시작합니다.
2. 지정된 시간 내에 응답이 완료되면 정상적으로 응답을 반환합니다.
3. 지정된 시간 내에 응답이 완료되지 않으면 503 Service Unavailable 상태 코드와 함께 타임아웃 메시지를 반환합니다.

타임아웃 미들웨어는 장시간 실행되는 API 요청으로 인한 서버 리소스 고갈을 방지하고, 클라이언트에게 적절한 응답 시간을 보장하는 데 유용합니다.
