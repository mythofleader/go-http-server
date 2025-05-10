# API 키 미들웨어

API 키 미들웨어는 요청 헤더에서 API 키를 확인하여 API에 대한 접근을 제어하는 기능을 제공합니다.

## 기본 사용법

가장 기본적인 사용법은 다음과 같습니다:

```go
package main

import (
    server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
    s, _ := server.NewServer(server.FrameworkGin, "8080")

    // API 키 미들웨어 구성
    apiKeyConfig := &server.APIKeyConfig{
        APIKey: "your-api-key-here", // 예상되는 API 키 값
    }

    // 보호된 라우트 그룹 생성
    protected := s.Group("/api")
    protected.Use(server.APIKeyMiddleware(apiKeyConfig))

    // 보호된 라우트 추가
    protected.GET("/data", func(c server.Context) {
        c.JSON(200, map[string]string{"message": "Protected data"})
    })

    s.Run()
}
```

## 사용자 정의 오류 메시지

오류 메시지를 사용자 정의하려면 다음과 같이 `UnauthorizedMessage` 필드를 설정합니다:

```go
apiKeyConfig := &server.APIKeyConfig{
    APIKey:              "your-api-key-here",
    UnauthorizedMessage: "유효한 API 키가 필요합니다",
}
```

## 기본 생성자 함수

API 키 미들웨어는 기본 구성을 사용하는 생성자 함수를 제공합니다:

```go
// 기본 API 키 미들웨어 생성자 (API 키를 인자로 받음)
s.Use(server.NewDefaultAPIKeyMiddleware("your-api-key-here"))
```

`NewDefaultAPIKeyMiddleware` 함수는 API 키를 인자로 받아 `DefaultAPIKeyConfig()`를 호출하여 기본 구성을 생성하고, API 키를 설정한 다음 이를 `APIKeyMiddleware` 함수에 전달합니다.

더 많은 설정이 필요한 경우 다음과 같이 사용할 수 있습니다:

```go
// 기본 구성을 가져와서 API 키 설정
apiKeyConfig := server.DefaultAPIKeyConfig()
apiKeyConfig.APIKey = "your-api-key-here"
apiKeyConfig.UnauthorizedMessage = "사용자 정의 오류 메시지"
s.Use(server.APIKeyMiddleware(apiKeyConfig))
```

## 작동 방식

API 키 미들웨어는 다음과 같이 작동합니다:

1. 요청 헤더에서 `x-api-key` 값을 확인합니다.
2. API 키가 없거나 구성에 지정된 값과 일치하지 않으면 401 Unauthorized 응답을 반환합니다.
3. API 키가 유효하면 요청이 다음 미들웨어나 핸들러로 전달됩니다.

## 클라이언트 사용법

클라이언트는 다음과 같이 `x-api-key` 헤더를 포함해야 합니다:

```
GET /api/data HTTP/1.1
Host: example.com
x-api-key: your-api-key-here
```

## 주의사항

- API 키는 HTTPS를 통해 전송되어야 합니다. HTTP를 통해 API 키를 전송하면 중간자 공격에 취약할 수 있습니다.
- API 키 미들웨어는 간단한 인증 메커니즘을 제공하지만, 더 복잡한 인증 요구 사항에는 인증 미들웨어를 사용하는 것이 좋습니다.
