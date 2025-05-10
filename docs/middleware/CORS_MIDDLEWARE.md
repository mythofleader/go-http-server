# CORS 미들웨어

CORS(Cross-Origin Resource Sharing) 미들웨어는 다른 도메인에서의 API 요청을 제어하는 기능을 제공합니다. 특정 도메인 목록을 설정하면 해당 도메인만 허용하고, 도메인 목록을 설정하지 않으면 모든 도메인을 허용합니다.

## 기본 사용법

가장 기본적인 사용법은 다음과 같습니다:

```go
package main

import (
    server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
    s, _ := server.NewServer(server.FrameworkGin, "8080")

    // 기본 CORS 미들웨어 추가 (모든 도메인 허용)
    s.Use(server.CORSMiddleware(nil))

    // 라우트 등록
    s.GET("/", func(c server.Context) {
        c.String(200, "Hello, World!")
    })

    s.Run()
}
```

## 특정 도메인만 허용하기

특정 도메인만 허용하려면 다음과 같이 `CORSConfig`를 사용합니다:

```go
corsConfig := &server.CORSConfig{
    AllowedDomains: []string{
        "https://example.com",
        "https://api.example.com",
    },
}
s.Use(server.CORSMiddleware(corsConfig))
```

## 고급 설정

CORS 미들웨어는 다양한 설정을 제공합니다:

```go
corsConfig := &server.CORSConfig{
    // 허용할 도메인 목록 (비어있으면 모든 도메인 허용)
    AllowedDomains: []string{
        "https://example.com",
        "https://api.example.com",
    },

    // 허용할 HTTP 메서드
    AllowedMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",

    // 허용할 HTTP 헤더
    AllowedHeaders: "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, X-Requested-With",

    // 자격 증명(쿠키 등) 포함 여부
    AllowCredentials: true,

    // 프리플라이트 요청 캐시 시간(초)
    MaxAge: 86400, // 24시간
}
s.Use(server.CORSMiddleware(corsConfig))
```

## 기본 생성자 함수

CORS 미들웨어는 기본 구성을 사용하는 생성자 함수를 제공합니다:

```go
// 기본 CORS 미들웨어 생성자 (모든 도메인 허용)
s.Use(server.NewDefaultCORSMiddleware())
```

`NewDefaultCORSMiddleware` 함수는 `DefaultCORSConfig()`를 호출하여 기본 구성을 생성하고, 이를 `CORSMiddleware` 함수에 전달합니다. 기본 구성은 모든 도메인을 허용하고, 일반적인 HTTP 메서드와 헤더를 허용합니다.

이는 다음과 동일합니다:

```go
s.Use(server.CORSMiddleware(nil))
```

또는

```go
s.Use(server.CORSMiddleware(server.DefaultCORSConfig()))
```

## 작동 방식

CORS 미들웨어는 다음과 같이 작동합니다:

1. 요청 헤더에서 `Origin` 값을 확인합니다.
2. `AllowedDomains`가 비어있으면 모든 도메인을 허용합니다.
3. `AllowedDomains`에 도메인이 지정되어 있으면, 요청의 `Origin`이 허용된 도메인 목록에 있는지 확인합니다.
4. 허용된 도메인이면 적절한 CORS 헤더를 설정합니다.
5. OPTIONS 메서드(프리플라이트 요청)인 경우 200 OK 응답을 반환합니다.
6. 그렇지 않으면 다음 미들웨어나 핸들러로 요청을 전달합니다.

## CORS 헤더

미들웨어는 다음과 같은 CORS 헤더를 설정합니다:

- `Access-Control-Allow-Origin`: 허용된 도메인 또는 `*`(모든 도메인)
- `Access-Control-Allow-Methods`: 허용된 HTTP 메서드
- `Access-Control-Allow-Headers`: 허용된 HTTP 헤더
- `Access-Control-Allow-Credentials`: 자격 증명 포함 여부
- `Access-Control-Max-Age`: 프리플라이트 요청 캐시 시간

## 주의사항

- 보안을 위해 가능하면 `AllowedDomains`를 설정하여 특정 도메인만 허용하는 것이 좋습니다.
- `AllowCredentials`가 `true`인 경우, `Access-Control-Allow-Origin`은 `*`가 될 수 없으며 특정 도메인이어야 합니다.
- 프로덕션 환경에서는 필요한 메서드와 헤더만 허용하는 것이 좋습니다.
