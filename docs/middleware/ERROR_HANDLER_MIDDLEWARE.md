# 에러 핸들러 미들웨어

에러 핸들러 미들웨어는 API 요청 처리 중 발생하는 에러를 캐치하고 적절한 HTTP 응답을 반환하는 기능을 제공합니다. 이 미들웨어는 특히 400, 401, 403, 500 등의 HTTP 상태 코드에 대한 에러 클래스를 제공하여 일관된 에러 응답을 생성할 수 있게 합니다.

## 기본 사용법

가장 기본적인 사용법은 다음과 같습니다:

```go
package main

import (
    "fmt"
    server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
    s, _ := server.NewServer(server.FrameworkGin, "8080")

    // 기본 에러 핸들러 미들웨어 추가
    errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
    s.Use(errorHandlerMiddleware.Middleware(nil))

    // 라우트 등록
    s.GET("/", func(c server.Context) {
        c.String(200, "Hello, World!")
    })

    s.GET("/error", func(c server.Context) {
        // 표준 에러 생성
        err := fmt.Errorf("잘못된 요청입니다")

        // 에러를 BadRequestHttpError로 래핑하고 panic으로 발생시키면 에러 핸들러 미들웨어가 처리
        panic(server.NewBadRequestHttpError(err))
    })

    s.Run()
}
```

## 에러 핸들러 설정

에러 핸들러 설정을 변경하려면 다음과 같이 `ErrorHandlerConfig`를 사용합니다:

```go
errorHandlerConfig := &server.ErrorHandlerConfig{
    DefaultErrorMessage: "서버 오류가 발생했습니다", // 기본 에러 메시지 변경
    DefaultStatusCode:   500,                  // 기본 상태 코드 변경
}
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
s.Use(errorHandlerMiddleware.Middleware(errorHandlerConfig))
```

## 프레임워크별 에러 핸들러 미들웨어

v1.1.0부터 tenqube-go-http-server는 프레임워크별 에러 핸들러 미들웨어를 제공하여 각 프레임워크에 최적화된 에러 처리를 할 수 있습니다. 이제 이러한 미들웨어는 각 프레임워크별 디렉토리에 구현되어 있습니다:

- `gin/errorhandler.go`: Gin 프레임워크를 위한 에러 핸들러 미들웨어
- `std/errorhandler.go`: 표준 HTTP 패키지를 위한 에러 핸들러 미들웨어
- `middleware/errorhandler.go`: 공통 에러 핸들러 미들웨어 기능

### 권장 사용법

가능하면 다음과 같이 프레임워크별 에러 핸들러 미들웨어를 직접 사용하는 것이 좋습니다:

```go
s, _ := server.NewServer(server.FrameworkGin, "8080")

// 에러 핸들러 미들웨어 추가
s.Use(s.GetErrorHandlerMiddleware().Middleware(nil))
```

또는

```go
s, _ := server.NewServer(server.FrameworkStdHTTP, "8080")

// 에러 핸들러 미들웨어 추가
s.Use(s.GetErrorHandlerMiddleware().Middleware(nil))
```

이 방식은 항상 올바른 프레임워크별 에러 핸들러 미들웨어를 사용하므로 에러를 더 정확하게 처리할 수 있습니다.

### 기존 방식과의 호환성

기존 방식은 더 이상 지원되지 않습니다. 항상 프레임워크별 에러 핸들러 미들웨어를 직접 사용해야 합니다:

```go
s, _ := server.NewServer(server.FrameworkGin, "8080")
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
s.Use(errorHandlerMiddleware.Middleware(nil))
```

이 방식은 항상 올바른 프레임워크별 에러 핸들러 미들웨어를 사용하므로 에러를 더 정확하게 처리할 수 있습니다.

## 기본 구성 사용하기

에러 핸들러 미들웨어는 기본 구성을 사용할 수 있습니다:

```go
// 프레임워크별 에러 핸들러 미들웨어 가져오기
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()

// 기본 구성으로 미들웨어 추가
s.Use(errorHandlerMiddleware.Middleware(nil))
```

또는 `DefaultErrorHandlerConfig()`를 명시적으로 사용할 수 있습니다:

```go
// 프레임워크별 에러 핸들러 미들웨어 가져오기
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()

// 기본 구성으로 미들웨어 추가
s.Use(errorHandlerMiddleware.Middleware(middleware.DefaultErrorHandlerConfig()))
```

기본 구성은 다음과 같은 값을 사용합니다:
- DefaultErrorMessage: "Internal Server Error"
- DefaultStatusCode: 500

## 에러 구조체 사용하기

라이브러리는 다음과 같은 HTTP 에러 구조체를 제공합니다:

1. `BadRequestHttpError` (400): 잘못된 요청 에러
2. `UnauthorizedHttpError` (401): 인증 실패 에러
3. `ForbiddenHttpError` (403): 권한 없음 에러
4. `NotFoundHttpError` (404): 리소스를 찾을 수 없음 에러
5. `InternalServerHttpError` (500): 서버 내부 에러
6. `ServiceUnavailableHttpError` (503): 서비스 사용 불가 에러

이러한 에러 구조체는 다음과 같이 사용할 수 있습니다:

```go
// 표준 에러 생성
err := fmt.Errorf("잘못된 요청 파라미터")

// 400 Bad Request 에러로 래핑
badRequestErr := server.NewBadRequestHttpError(err)

// 401 Unauthorized 에러로 래핑
unauthorizedErr := server.NewUnauthorizedHttpError(fmt.Errorf("인증이 필요합니다"))

// 403 Forbidden 에러로 래핑
forbiddenErr := server.NewForbiddenHttpError(fmt.Errorf("접근 권한이 없습니다"))

// 500 Internal Server Error로 래핑
internalErr := server.NewInternalServerHttpError(fmt.Errorf("서버 오류가 발생했습니다"))

// 핸들러에서 에러 반환
s.GET("/bad-request", func(c server.Context) {
    // 표준 에러 생성
    err := fmt.Errorf("잘못된 요청입니다")

    // 에러를 BadRequestHttpError로 래핑하고 panic으로 발생시키면 에러 핸들러 미들웨어가 처리
    panic(server.NewBadRequestHttpError(err))
})

// Context 인터페이스의 Error 메서드를 사용할 수 있습니다
s.GET("/unauthorized", func(c server.Context) {
    // 표준 에러 생성
    err := fmt.Errorf("인증이 필요합니다")

    // 에러를 UnauthorizedHttpError로 래핑하고 Error 메서드를 사용하여 추가
    // 이 에러는 에러 핸들러 미들웨어에 의해 자동으로 처리됩니다
    c.Error(server.NewUnauthorizedHttpError(err))
})

// 컨텍스트에서 모든 에러를 가져올 수 있습니다
s.GET("/errors", func(c server.Context) {
    // 여러 에러 추가
    c.Error(fmt.Errorf("첫 번째 에러"))
    c.Error(fmt.Errorf("두 번째 에러"))

    // 모든 에러 가져오기
    errors := c.Errors()
    for i, err := range errors {
        fmt.Printf("에러 %d: %v\n", i+1, err)
    }
})
```

## 컨텍스트에서 에러 가져오기

Context 인터페이스는 `Errors()` 메서드를 제공하여 컨텍스트에 추가된 모든 에러를 가져올 수 있습니다. 이 메서드는 `Error()` 메서드로 추가된 모든 에러를 포함하는 슬라이스를 반환합니다.

```go
// 여러 에러 추가
c.Error(fmt.Errorf("첫 번째 에러"))
c.Error(fmt.Errorf("두 번째 에러"))

// 모든 에러 가져오기
errors := c.Errors()
for i, err := range errors {
    fmt.Printf("에러 %d: %v\n", i+1, err)
}
```

이 기능은 미들웨어나 핸들러에서 발생한 여러 에러를 수집하고 처리하는 데 유용합니다. 에러 핸들러 미들웨어는 로깅이 활성화된 경우 컨텍스트의 모든 에러를 로그에 기록합니다.

## 에러 응답 형식

에러 핸들러 미들웨어는 다음과 같은 JSON 형식으로 에러 응답을 반환합니다:

```json
{
  "error": {
    "code": 400,
    "message": "잘못된 요청 파라미터"
  }
}
```

## 표준화된 에러 응답 구조체 사용하기

라이브러리는 에러 응답을 생성하기 위한 표준화된 구조체와 헬퍼 함수를 제공합니다:

```go
// ErrorDetail은 응답에서 에러 세부 정보의 구조를 나타냅니다.
type ErrorDetail struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// ErrorResponse는 에러 응답의 구조를 나타냅니다.
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}
```

이러한 구조체를 사용하여 일관된 에러 응답을 생성할 수 있습니다:

```go
// 기본 에러 응답 생성
response := server.NewErrorResponse(http.StatusBadRequest, "잘못된 요청 파라미터")

// 특정 HTTP 상태 코드에 대한 헬퍼 함수 사용
badRequestResponse := server.NewBadRequestResponse("잘못된 요청 파라미터")
unauthorizedResponse := server.NewUnauthorizedResponse("인증이 필요합니다")
forbiddenResponse := server.NewForbiddenResponse("접근 권한이 없습니다")
notFoundResponse := server.NewNotFoundResponse("리소스를 찾을 수 없습니다")
conflictResponse := server.NewConflictResponse("리소스가 이미 존재합니다")
internalErrorResponse := server.NewInternalServerErrorResponse("서버 오류가 발생했습니다")
serviceUnavailableResponse := server.NewServiceUnavailableResponse("서비스를 사용할 수 없습니다")

// HTTPError에서 ErrorResponse 생성
httpErr := server.NewBadRequestError("잘못된 요청")
errorResponse := server.FromHTTPError(httpErr)

// 핸들러에서 에러 응답 반환
s.GET("/bad-request", func(c server.Context) {
    c.JSON(http.StatusBadRequest, server.NewBadRequestResponse("잘못된 요청 파라미터"))
})
```

이 표준화된 구조체를 사용하면 에러 응답의 형식이 일관되게 유지되며, 키 값을 빠뜨리거나 오타가 발생할 가능성이 줄어듭니다.

## 새로운 에러 구조체 사용하기

v0.5.0부터 라이브러리는 error 인터페이스를 임베딩하는 새로운 HTTP 에러 구조체를 제공합니다. 이 구조체들은 기존 에러 객체를 래핑하여 HTTP 상태 코드를 부여하는 간단한 방법을 제공합니다:

1. `BadRequestHttpError` (400): 잘못된 요청 에러
2. `UnauthorizedHttpError` (401): 인증 실패 에러
3. `ForbiddenHttpError` (403): 권한 없음 에러
4. `NotFoundHttpError` (404): 리소스를 찾을 수 없음 에러
5. `InternalServerHttpError` (500): 서버 내부 에러
6. `ServiceUnavailableHttpError` (503): 서비스 사용 불가 에러

이러한 새로운 에러 구조체는 다음과 같이 사용할 수 있습니다:

```go
// 표준 에러 생성
err := fmt.Errorf("잘못된 요청 파라미터")

// 400 Bad Request 에러로 래핑
badRequestErr := server.NewBadRequestHttpError(err)

// 401 Unauthorized 에러로 래핑
unauthorizedErr := server.NewUnauthorizedHttpError(fmt.Errorf("인증이 필요합니다"))

// 403 Forbidden 에러로 래핑
forbiddenErr := server.NewForbiddenHttpError(fmt.Errorf("접근 권한이 없습니다"))

// 404 Not Found 에러로 래핑
notFoundErr := server.NewNotFoundHttpError(fmt.Errorf("리소스를 찾을 수 없습니다"))

// 500 Internal Server Error로 래핑
internalErr := server.NewInternalServerHttpError(fmt.Errorf("서버 오류가 발생했습니다"))

// 503 Service Unavailable 에러로 래핑
serviceUnavailableErr := server.NewServiceUnavailableHttpError(fmt.Errorf("서비스를 사용할 수 없습니다"))

// 핸들러에서 에러 반환
s.GET("/new-bad-request", func(c server.Context) {
    // 표준 에러 생성
    err := fmt.Errorf("잘못된 요청 파라미터")

    // BadRequestHttpError로 래핑
    httpErr := server.NewBadRequestHttpError(err)

    // 에러를 panic으로 발생시키면 에러 핸들러 미들웨어가 처리
    panic(httpErr)
})
```

이 새로운 접근 방식의 장점은 다음과 같습니다:

1. 기존 에러 객체를 그대로 사용할 수 있습니다.
2. 에러 메시지를 별도로 지정할 필요 없이 원래 에러의 메시지를 사용합니다.
3. 코드가 더 간결하고 직관적입니다.
4. 에러 핸들러 미들웨어가 자동으로 적절한 HTTP 상태 코드를 설정합니다.

## 커스텀 에러 구조체 만들기

라이브러리에서 제공하는 HTTP 에러 구조체를 상속하여 커스텀 에러 구조체를 만들 수 있습니다. 이를 통해 특정 유형의 에러에 대한 더 구체적인 에러 타입을 정의할 수 있습니다.

예를 들어, 잘못된 요청 파라미터에 대한 커스텀 에러 구조체를 만들려면 다음과 같이 할 수 있습니다:

```go
// InvalidRequestParamError는 잘못된 요청 파라미터에 대한 400 Bad Request 에러를 나타냅니다.
type InvalidRequestParamError struct {
    BadRequestHttpError
}

// NewInvalidRequestParamError는 새로운 InvalidRequestParamError를 생성합니다.
func NewInvalidRequestParamError(err error) *InvalidRequestParamError {
    return &InvalidRequestParamError{
        BadRequestHttpError: *NewBadRequestHttpError(err),
    }
}
```

이 커스텀 에러 구조체는 다음과 같이 사용할 수 있습니다:

```go
// 표준 에러 생성
err := fmt.Errorf("잘못된 요청 파라미터: 'id'는 양의 정수여야 합니다")

// InvalidRequestParamError로 래핑
httpErr := NewInvalidRequestParamError(err)

// Context.Error 메서드를 사용하여 에러 설정
c.Error(httpErr)
```

이 에러는 BadRequestHttpError를 상속하므로, 에러 핸들러 미들웨어에 의해 400 Bad Request 응답으로 처리됩니다. 이를 통해 특정 유형의 에러에 대한 더 구체적인 에러 타입을 정의하면서도 기존 에러 처리 메커니즘을 그대로 활용할 수 있습니다.

라이브러리는 이미 InvalidRequestParamError 구조체를 제공하므로, 다음과 같이 직접 사용할 수 있습니다:

```go
// 표준 에러 생성
err := fmt.Errorf("잘못된 요청 파라미터: 'id'는 양의 정수여야 합니다")

// InvalidRequestParamError로 래핑
httpErr := server.NewInvalidRequestParamError(err)

// Context.Error 메서드를 사용하여 에러 설정
c.Error(httpErr)
```

## 에러 핸들러 동작 방식

에러 핸들러 미들웨어는 다음과 같이 동작합니다:

1. 요청 처리 중 panic이 발생하면 이를 캐치하여 에러로 변환합니다.
2. 에러가 HTTP 에러 클래스인 경우 해당 상태 코드와 메시지로 응답합니다.
3. 에러가 HTTP 에러 클래스가 아닌 경우 기본 상태 코드와 메시지로 응답합니다.
4. 설정에 따라 에러를 로깅합니다.

에러 핸들러 미들웨어는 API 응답의 일관성을 유지하고, 클라이언트에게 적절한 에러 정보를 제공하는 데 유용합니다.
