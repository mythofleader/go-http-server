# 미들웨어 사용 가이드

이 문서는 tenqube-go-http-server 라이브러리에서 제공하는 미들웨어 기능에 대한 개요입니다. 각 미들웨어에 대한 자세한 내용은 해당 문서를 참조하세요.

## 제공되는 미들웨어

tenqube-go-http-server 라이브러리는 다음과 같은 미들웨어를 제공합니다:

### [로깅 미들웨어](LOGGING_MIDDLEWARE.md)

로깅 미들웨어는 API 요청에 대한 로그를 생성하고 콘솔에 출력하거나 원격 URL로 전송하는 기능을 제공합니다. 개발 및 프로덕션 환경에 따라 로깅 동작을 제어할 수 있으며, 커스텀 필드를 추가할 수 있습니다.

[자세히 보기](LOGGING_MIDDLEWARE.md)

### [타임아웃 미들웨어](TIMEOUT_MIDDLEWARE.md)

타임아웃 미들웨어는 API 요청이 지정된 시간 내에 응답하지 않을 경우 자동으로 타임아웃 응답을 반환하는 기능을 제공합니다. 기본적으로 2초 후에 타임아웃 응답을 반환하며, 타임아웃 시간을 설정할 수 있습니다.

[자세히 보기](TIMEOUT_MIDDLEWARE.md)

### [에러 핸들러 미들웨어](ERROR_HANDLER_MIDDLEWARE.md)

에러 핸들러 미들웨어는 API 요청 처리 중 발생하는 에러를 캐치하고 적절한 HTTP 응답을 반환하는 기능을 제공합니다. 400, 401, 403, 500 등의 HTTP 상태 코드에 대한 에러 클래스를 제공하여 일관된 에러 응답을 생성할 수 있습니다.

[자세히 보기](ERROR_HANDLER_MIDDLEWARE.md)

### [인증 미들웨어](AUTH_MIDDLEWARE.md)

인증 미들웨어는 HTTP 기본 인증(Basic Authentication) 또는 JWT 베어러 토큰(Bearer Token)을 사용하여 사용자를 인증하는 기능을 제공합니다. 사용자 정의 사용자 조회 인터페이스를 구현하여 다양한 인증 방식을 지원할 수 있습니다.

[자세히 보기](AUTH_MIDDLEWARE.md)

### [API 키 미들웨어](API_KEY_MIDDLEWARE.md)

API 키 미들웨어는 요청 헤더에서 API 키를 확인하여 API에 대한 접근을 제어하는 기능을 제공합니다. 간단한 인증 메커니즘으로, 헤더의 x-api-key 값을 검증하여 접근을 허용하거나 거부합니다.

[자세히 보기](API_KEY_MIDDLEWARE.md)

### [CORS 미들웨어](CORS_MIDDLEWARE.md)

CORS 미들웨어는 다른 도메인에서의 API 요청을 제어하는 기능을 제공합니다. 특정 도메인 목록을 설정하면 해당 도메인만 허용하고, 도메인 목록을 설정하지 않으면 모든 도메인을 허용합니다.

[자세히 보기](CORS_MIDDLEWARE.md)

### [중복 요청 방지 미들웨어](DUPLICATE_REQUEST_MIDDLEWARE.md)

중복 요청 방지 미들웨어는 동일한 요청이 중복으로 처리되는 것을 방지하는 기능을 제공합니다. 요청 ID를 생성하고, 이 ID가 이미 처리된 요청인지 확인한 후, 중복된 요청인 경우 409 Conflict 응답을 반환합니다. 사용자 정의 ID 생성 및 저장소 인터페이스를 구현하여 다양한 방식으로 중복 요청을 감지할 수 있습니다.

[자세히 보기](DUPLICATE_REQUEST_MIDDLEWARE.md)

## 미들웨어 등록 순서

미들웨어 등록 순서는 애플리케이션의 동작에 중요한 영향을 미칩니다. 올바른 순서로 미들웨어를 등록하지 않으면 예상치 못한 동작이 발생할 수 있습니다. 다음은 권장되는 미들웨어 등록 순서입니다:

1. **에러 핸들러 미들웨어** (반드시 첫 번째로 등록)
   - 이 미들웨어는 이후의 모든 미들웨어와 핸들러에서 발생하는 에러와 패닉을 캐치합니다.
   - 다른 미들웨어에서 발생하는 에러를 적절히 처리하기 위해 반드시 첫 번째로 등록해야 합니다.

2. **타임아웃 미들웨어** (사용하는 경우)
   - 요청 타임아웃을 제어하고 장시간 실행되는 요청을 방지합니다.

3. **CORS 미들웨어** (사용하는 경우)
   - Cross-Origin Resource Sharing 헤더를 처리합니다.

4. **로깅 미들웨어** (에러 핸들러 미들웨어 이후에 등록)
   - 이 미들웨어는 상태 코드와 에러를 포함한 요청 세부 정보를 로깅합니다.
   - 에러 핸들러 미들웨어 이후에 등록하여 에러를 올바르게 캡처해야 합니다.

5. **커스텀 미들웨어**
   - 애플리케이션에서 제공하는 추가 미들웨어입니다.

### 예시 코드

다음은 미들웨어를 올바른 순서로 등록하는 예시입니다:

```
// 1. 에러 핸들러 미들웨어 (반드시 첫 번째로 등록)
errorHandler := s.GetErrorHandlerMiddleware()
s.Use(errorHandler.Middleware(nil))

// 2. 타임아웃 미들웨어 (선택 사항)
timeoutConfig := &server.TimeoutConfig{
    Timeout: 2 * time.Second,
}
s.Use(server.TimeoutMiddleware(timeoutConfig))

// 3. CORS 미들웨어 (선택 사항)
corsConfig := &server.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
}
s.Use(server.CORSMiddleware(corsConfig))

// 4. 로깅 미들웨어 (에러 핸들러 이후에 등록)
loggingConfig := &server.LoggingConfig{
    LoggingToConsole: true,
}
s.Use(s.GetLoggingMiddleware().Middleware(loggingConfig))

// 5. 커스텀 미들웨어
s.Use(CustomMiddleware())
```

전체 예제는 `examples/middleware/main.go` 파일을 참조하세요.

## 미들웨어 흐름 제어

tenqube-go-http-server 라이브러리는 미들웨어 내에서 흐름을 제어하기 위한 `Next()` 메서드를 제공합니다. 이 메서드를 사용하면 미들웨어에서 다음 핸들러를 호출하고, 그 후에 추가 작업을 수행할 수 있습니다.

### Next() 메서드 사용 예시

```
func CustomMiddleware() server.HandlerFunc {
    return func(c server.Context) {
        // 핸들러 실행 전 작업
        log.Printf("핸들러 실행 전: %s %s", c.Request().Method, c.Request().URL.Path)

        // 시작 시간 기록
        start := time.Now()

        // 다음 핸들러 호출
        c.Next()

        // 핸들러 실행 후 작업
        duration := time.Since(start)
        log.Printf("핸들러 실행 후: %s %s (소요 시간: %v)", c.Request().Method, c.Request().URL.Path, duration)
    }
}
```

### 조건부 흐름 제어

`Next()` 메서드를 사용하면 특정 조건에 따라 다음 핸들러를 호출할지 여부를 결정할 수 있습니다.

```
func ConditionalMiddleware() server.HandlerFunc {
    return func(c server.Context) {
        path := c.Request().URL.Path

        if path == "/skip" {
            // 특정 경로에 대해 다음 핸들러를 호출하지 않고 직접 응답
            log.Printf("핸들러 건너뛰기: %s", path)
            c.String(http.StatusOK, "핸들러가 건너뛰어졌습니다!")
            return
        }

        // 다른 경로에 대해서는 다음 핸들러 호출
        log.Printf("다음 핸들러로 계속: %s", path)
        c.Next()
    }
}
```

### 예제 코드

미들웨어 흐름 제어에 대한 전체 예제는 `examples/middleware/main.go`를 참조하세요.
