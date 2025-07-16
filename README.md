# Go HTTP 서버

Go를 위한 간단하고 추상화된 HTTP 서버 라이브러리로, Gin, 표준 net/http 및 AWS Lambda와 같은 인기 있는 프레임워크를 래핑합니다.

## 특징

- 프레임워크에 구애받지 않는 API
- 다양한 HTTP 프레임워크 지원 (Gin, 표준 net/http)
- AWS Lambda 지원
- 간단하고 직관적인 인터페이스
- 사용 및 확장이 쉬움
- 미들웨어 지원
- 라우터 그룹
- TLS 지원
- 정상 종료

## 디렉토리 구조

이 라이브러리는 다음과 같은 디렉토리 구조로 구성되어 있습니다:

- `core/`: 핵심 인터페이스와 타입 정의
  - `context.go`: Context, Server, RouterGroup 등의 인터페이스 정의
- `gin/`: Gin 프레임워크 구현
  - `server.go`: Gin 프레임워크를 사용한 서버 구현
  - `middleware.go`: Gin 프레임워크를 위한 미들웨어 구현
- `std/`: 표준 net/http 패키지 구현
  - `server.go`: 표준 net/http 패키지를 사용한 서버 구현
  - `middleware.go`: 표준 net/http 패키지를 위한 미들웨어 구현
- `middleware/`: 공통 미들웨어 기능
  - `middleware.go`: 로깅 등의 공통 미들웨어 기능 구현
- `server.go`: 루트 패키지에서 서버 생성 함수 제공
- `docs/`: 프로젝트 문서
  - `IMPORTING.md`: 패키지 가져오기 방법에 대한 문서
  - `MIDDLEWARE.md`: 미들웨어 사용 가이드
  - `USAGE.md`: 라이브러리 사용 방법에 대한 자세한 문서
  - `RELEASE_PROCESS.md`: 릴리스 프로세스에 대한 문서

## 설치

```bash
go get github.com/mythofleader/go-http-server
```

## 문서

자세한 문서는 `docs` 디렉토리에서 찾을 수 있습니다:

- [패키지 가져오기 방법](docs/getting-started/IMPORTING.md)
- [미들웨어 사용 가이드](docs/middleware/MIDDLEWARE.md)
- [라이브러리 사용 방법](docs/getting-started/USAGE.md)
- [릴리스 프로세스](docs/project-management/RELEASE_PROCESS.md)

## 사용법

### 기본 예제

```go
package main

import (
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// 새 서버 생성 (기본값은 Gin)
	s, err := server.NewServer(server.FrameworkGin, "8080")
	if err != nil {
		log.Fatalf("서버 생성 실패: %v", err)
	}

	// 라우트 등록
	s.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// 서버 시작
	log.Println("서버가 :8080 포트에서 시작됩니다")
	if err := s.Run(); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
```

### 프레임워크 지정

서버를 생성할 때 사용할 프레임워크를 지정할 수 있습니다:

```go
// Gin 프레임워크 사용
s, err := server.NewServer(server.FrameworkGin, "8080")

// 표준 net/http 패키지 사용
s, err := server.NewServer(server.FrameworkStdHTTP, "8080")
```

### AWS Lambda 지원

AWS Lambda를 사용할 때는 Gin 프레임워크로 서버를 생성한 다음, `Run` 대신 `StartLambda` 메서드를 사용해야 합니다. **중요: Lambda는 Gin 프레임워크에서만 지원되며, 표준 HTTP 서버에서는 지원되지 않습니다.**

```go
// Gin 프레임워크로 서버 생성 (Lambda는 Gin 프레임워크에서만 지원됨)
s, err := server.NewServer(server.FrameworkGin, "8080")
if err != nil {
	log.Fatalf("서버 생성 실패: %v", err)
}

// 라우트 등록 (단일 핸들러)
s.GET("/", func(c server.Context) {
	c.String(http.StatusOK, "Hello, World!")
})

// 라우트 등록 (다중 핸들러)
s.GET("/api/users", 
	logRequestHandler, 
	validateRequestHandler, 
	getUsersHandler)

// Lambda 핸들러 시작 (표준 HTTP 서버에서 호출하면 오류 반환)
if err := s.StartLambda(); err != nil {
	log.Fatalf("Lambda 시작 실패: %v", err)
}
```

표준 HTTP 서버에서 `StartLambda`를 호출하면 "Lambda is only supported with the Gin framework" 오류가 반환됩니다.

### 404 Not Found 및 405 Method Not Allowed 핸들러

존재하지 않는 경로(404 Not Found)나 허용되지 않는 메서드(405 Method Not Allowed)에 대한 요청을 처리하기 위한 핸들러를 등록할 수 있습니다:

```go
// 404 Not Found 핸들러 등록
s.NoRoute(func(c server.Context) {
	c.JSON(http.StatusNotFound, server.NewErrorResponse(http.StatusNotFound, "페이지를 찾을 수 없습니다"))
})

// 405 Method Not Allowed 핸들러 등록
s.NoMethod(func(c server.Context) {
	c.JSON(http.StatusMethodNotAllowed, server.NewErrorResponse(http.StatusMethodNotAllowed, "허용되지 않는 메서드입니다"))
})
```

핸들러를 등록하지 않으면 기본 핸들러가 자동으로 적용됩니다. 기본 핸들러는 404 Not Found 및 405 Method Not Allowed 오류를 적절한 상태 코드와 메시지로 반환합니다:

```go
// 기본 404 Not Found 핸들러는 다음과 같이 동작합니다
func(c server.Context) {
	path := c.Request().URL.Path
	err := fmt.Errorf("Route not found: %s", path)
	c.Error(server.NewNotFoundHttpError(err))
}

// 기본 405 Method Not Allowed 핸들러는 다음과 같이 동작합니다
func(c server.Context) {
	method := c.Request().Method
	path := c.Request().URL.Path
	err := fmt.Errorf("Method %s not allowed for path %s", method, path)
	c.Error(server.NewMethodNotAllowedHttpError(err))
}
```

이러한 핸들러는 에러 핸들러 미들웨어와 함께 사용하여 일관된 에러 응답을 제공할 수 있습니다:

```go
// 404 Not Found 핸들러 등록 (에러 핸들러 미들웨어 사용)
s.NoRoute(func(c server.Context) {
	err := fmt.Errorf("페이지를 찾을 수 없습니다")
	c.Error(server.NewNotFoundHttpError(err))
})

// 405 Method Not Allowed 핸들러 등록 (에러 핸들러 미들웨어 사용)
s.NoMethod(func(c server.Context) {
	err := fmt.Errorf("허용되지 않는 메서드입니다")
	c.Error(server.NewMethodNotAllowedHttpError(err))
})
```

### 미들웨어

```go
// 미들웨어 추가
s.Use(func(c server.Context) {
	log.Printf("요청: %s %s", c.Request().Method, c.Request().URL.Path)
	// 요청 처리 계속
})
```

#### 미들웨어 등록 순서

미들웨어 등록 순서는 애플리케이션의 동작에 중요한 영향을 미칩니다. 올바른 순서로 미들웨어를 등록하지 않으면 예상치 못한 동작이 발생할 수 있습니다. 다음은 권장되는 미들웨어 등록 순서입니다:

1. **에러 핸들러 미들웨어** (반드시 첫 번째로 등록)
   - 이후의 모든 미들웨어와 핸들러에서 발생하는 에러와 패닉을 캐치합니다.
   - 다른 미들웨어에서 발생하는 에러를 적절히 처리하기 위해 반드시 첫 번째로 등록해야 합니다.

2. **타임아웃 미들웨어** (사용하는 경우)
   - 요청 타임아웃을 제어하고 장시간 실행되는 요청을 방지합니다.

3. **CORS 미들웨어** (사용하는 경우)
   - Cross-Origin Resource Sharing 헤더를 처리합니다.

4. **로깅 미들웨어** (에러 핸들러 미들웨어 이후에 등록)
   - 상태 코드와 에러를 포함한 요청 세부 정보를 로깅합니다.
   - 에러 핸들러 미들웨어 이후에 등록하여 에러를 올바르게 캡처해야 합니다.

5. **커스텀 미들웨어**
   - 애플리케이션에서 제공하는 추가 미들웨어입니다.

자세한 내용은 [미들웨어 사용 가이드](docs/middleware/MIDDLEWARE.md)를 참조하세요.

### 서버 초기화 로깅

서버가 시작될 때 서버 정보, 미들웨어 구성, 라우트 정보 등이 자동으로 로깅됩니다:

```
[GIN] Creating new Gin server on port 8080
[MIDDLEWARE] ErrorHandler middleware configured:
[MIDDLEWARE]   - Default error message: Internal Server Error
[MIDDLEWARE]   - Default status code: 500
[MIDDLEWARE]   - Log errors: true
[MIDDLEWARE] Logging middleware configured:
[MIDDLEWARE]   - Logging to console: true
[MIDDLEWARE]   - Logging to remote: false
[MIDDLEWARE]   - Custom fields:
[MIDDLEWARE]     version: 1.0.0
[MIDDLEWARE]     environment: development
[GIN] Adding middleware: github.com/mythofleader/go-http-server/core/gin.ErrorHandlerMiddleware.Middleware
[GIN] Adding middleware: github.com/mythofleader/go-http-server/core/gin.LoggingMiddleware.Middleware
[GIN] Server starting on :8080
[GIN] Using Gin framework version: 1.9.1
[GIN] Middleware registered:
[GIN]   1. github.com/mythofleader/go-http-server/core/gin.ErrorHandlerMiddleware.Middleware
[GIN]   2. github.com/mythofleader/go-http-server/core/gin.LoggingMiddleware.Middleware
[GIN] Routes registered:
[GIN]   1. GET /
[GIN]   2. GET /api/users
[GIN] Server is ready to handle requests
```

표준 HTTP 서버를 사용할 경우 로그 접두사는 `[STD]`로 표시됩니다.

#### 로깅 미들웨어에서 특정 경로 무시하기

로깅 미들웨어는 특정 경로에 대한 로깅을 건너뛸 수 있습니다:

```go
loggingConfig := &server.LoggingConfig{
	SkipPaths: []string{
		"/health",
		"/metrics",
		"/favicon.ico",
	},
	CustomFields: map[string]string{
		"version": "1.0.0",
	},
}
// 서버에서 프레임워크별 로깅 미들웨어 가져오기
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(loggingConfig))
```

#### 인증 미들웨어에서 특정 경로 무시하기

인증 미들웨어는 특정 경로에 대한 인증 검사를 건너뛸 수 있습니다:

```go
authConfig := &server.AuthConfig{
	UserLookup: userService,
	AuthType:   server.AuthTypeJWT,
	JWTSecret:  "your-secret-key",
	SkipPaths: []string{
		"/health",
		"/metrics",
		"/public",
	},
}
s.Use(server.AuthMiddleware(authConfig))
```

#### 기본 미들웨어 생성자 사용하기

각 미들웨어에는 기본 구성을 사용하는 생성자 함수가 있습니다. 이 함수들은 미들웨어 이름 앞에 `NewDefault`를 붙여서 명명됩니다:

```go
// 기본 타임아웃 미들웨어 추가 (2초 타임아웃)
s.Use(server.NewDefaultTimeoutMiddleware())

// 기본 CORS 미들웨어 추가 (모든 도메인 허용)
s.Use(server.NewDefaultCORSMiddleware())

// 프레임워크별 로깅 미들웨어 가져오기 및 기본 구성으로 사용
loggingMiddleware := s.GetLoggingMiddleware()
s.Use(loggingMiddleware.Middleware(middleware.DefaultLoggingConfig()))

// 프레임워크별 에러 핸들러 미들웨어 가져오기 및 기본 구성으로 사용
errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
s.Use(errorHandlerMiddleware.Middleware(middleware.DefaultErrorHandlerConfig()))
```

일부 미들웨어는 추가 구성이 필요하지만, 기본 생성자를 사용할 수 있습니다:

```go
// API 키 미들웨어 추가 (API 키를 인자로 받음)
s.Use(server.NewDefaultAPIKeyMiddleware("your-api-key"))

// JWT 인증 미들웨어 추가 (JWTUserLookup과 JWT 비밀 키를 인자로 받음)
s.Use(server.NewDefaultJWTAuthMiddleware(myJWTLookup, "your-jwt-secret"))

// 기본 인증 미들웨어 추가 (BasicAuthUserLookup을 인자로 받음)
s.Use(server.NewDefaultBasicAuthMiddleware(myBasicAuthLookup))
```

일부 미들웨어는 여전히 추가 구성이 필요하므로 기본 생성자를 직접 사용하면 패닉이 발생합니다:

```go
// 주의: 이 함수는 추가 구성이 필요하므로 직접 호출하면 패닉이 발생합니다
// server.NewDefaultDuplicateRequestMiddleware() // RequestIDGenerator와 RequestIDStorage가 필요합니다
```

더 많은 설정이 필요한 경우 다음과 같이 구성을 제공할 수 있습니다:

```go
// API 키 미들웨어 구성
apiKeyConfig := server.DefaultAPIKeyConfig()
apiKeyConfig.APIKey = "your-api-key"
apiKeyConfig.UnauthorizedMessage = "사용자 정의 오류 메시지"
s.Use(server.APIKeyMiddleware(apiKeyConfig))

// JWT 인증 미들웨어 구성
authConfig := server.DefaultAuthConfig()
authConfig.AuthType = server.AuthTypeJWT
authConfig.JWTLookup = myJWTLookup
authConfig.JWTSecret = "your-jwt-secret"
authConfig.UnauthorizedMessage = "사용자 정의 오류 메시지"
s.Use(server.AuthMiddleware(authConfig))

// 기본 인증 미들웨어 구성
authConfig := server.DefaultAuthConfig()
authConfig.AuthType = server.AuthTypeBasic
authConfig.BasicAuthLookup = myBasicAuthLookup
authConfig.UnauthorizedMessage = "사용자 정의 오류 메시지"
s.Use(server.AuthMiddleware(authConfig))

// 중복 요청 방지 미들웨어 구성
dupReqConfig := server.DefaultDuplicateRequestConfig()
dupReqConfig.RequestIDGenerator = myRequestIDGenerator
dupReqConfig.RequestIDStorage = myRequestIDStorage
s.Use(server.DuplicateRequestMiddleware(dupReqConfig))
```

### 라우터 그룹

```go
// 라우터 그룹 생성
api := s.Group("/api")
{
	api.GET("/users", getUsersHandler)
	api.POST("/users", createUserHandler)
}
```

### 컨트롤러 인터페이스

컨트롤러 인터페이스를 사용하면 관련 라우트를 그룹화하고 재사용 가능한 컨트롤러 컴포넌트를 만들 수 있습니다:

```go
// 컨트롤러 구현
type UserController struct {
	// 필요한 의존성
	userService UserService
}

// HTTP 메서드 반환
func (r *UserController) GetHttpMethod() server.HttpMethod {
	return server.GET
}

// 경로 반환
func (r *UserController) GetPath() string {
	return "/api/users"
}

// 로깅 무시 경로 반환
func (r *UserController) GetLogIgnorePath() string {
	return "/api/users/health"
}

// 인증 검사 무시 경로 반환
func (r *UserController) GetAuthCheckIgnorePath() string {
	return "/api/users/public"
}

// 핸들러 함수 반환 (여러 핸들러 지원)
func (r *UserController) Handler() []server.HandlerFunc {
	// 여러 핸들러를 반환할 수 있음
	return []server.HandlerFunc{
		r.logRequest,    // 로깅 핸들러
		r.validateInput, // 입력 검증 핸들러
		r.getUsers,      // 실제 비즈니스 로직 핸들러
	}
}

// 로깅 핸들러
func (r *UserController) logRequest(c server.Context) {
	log.Printf("요청: %s %s", c.Request().Method, c.Request().URL.Path)
}

// 입력 검증 핸들러
func (r *UserController) validateInput(c server.Context) {
	// 입력 검증 로직
}

// 실제 비즈니스 로직 핸들러
func (r *UserController) getUsers(c server.Context) {
	// 사용자 목록 조회 로직
}

// 컨트롤러 등록
userController := &UserController{userService: myUserService}
s.RegisterRouter(userController)

// 또는 라우터 그룹에 등록
api := s.Group("/api")
api.RegisterRouter(userController)
```

### JSON 응답

```go
func jsonHandler(c server.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello, JSON!",
		"status":  "success",
	})
}
```

### 요청 바인딩

```go
func createUserHandler(c server.Context) {
	var user struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// 사용자 처리...
	c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      1,
		"name":    user.Name,
		"message": "사용자가 성공적으로 생성되었습니다",
	})
}
```

### TLS 지원

```go
// TLS로 서버 시작
if err := s.RunTLS(":8443", "cert.pem", "key.pem"); err != nil {
	log.Fatalf("서버 시작 실패: %v", err)
}
```

### 정상 종료

```go
// 별도의 고루틴에서
go func() {
	// 인터럽트 신호 대기
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// 종료를 위한 데드라인 생성
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 서버 종료
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("서버 종료 실패: %v", err)
	}
	log.Println("서버가 정상적으로 중지되었습니다")
}()
```

### 서버 빌더 사용하기

서버 빌더를 사용하면 컨트롤러와 미들웨어를 쉽게 구성하고, 컨트롤러에서 로깅 무시 경로와 인증 검사 무시 경로를 자동으로 수집할 수 있습니다.

```go
// 서버 빌더 생성 (방법 1: 프레임워크와 포트 지정)
builder := server.NewServerBuilder(server.FrameworkGin, "8080")

// 또는 Gin 프레임워크와 8080 포트를 사용하는 서버 빌더 생성 (방법 2: 더 간단한 방법)
// builder := server.NewGinServerBuilder()

// 컨트롤러 추가
userController := &UserController{userService: myUserService}
productController := &ProductController{productService: myProductService}
builder.AddControllers(userController, productController)

// 로깅 구성 추가 (컨트롤러의 GetLogIgnorePath 값이 자동으로 수집됨)
builder.WithLogging(map[string]string{
    "environment": "development",
    "version":     "1.0.0",
})

// 타임아웃 구성 추가
builder.WithTimeout(server.TimeoutConfig{
    Timeout: 5 * time.Second,
})

// CORS 구성 추가
builder.WithCORS(server.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
})

// 에러 핸들러 구성 추가
builder.WithErrorHandler(server.ErrorHandlerConfig{
    DefaultErrorMessage: "서버 오류가 발생했습니다",
    DefaultStatusCode:   500,
})

// 커스텀 미들웨어 추가
builder.AddMiddleware(func(c server.Context) {
    log.Printf("요청: %s %s", c.Request().Method, c.Request().URL.Path)
})

// 서버 빌드 및 시작
s, err := builder.Build()
if err != nil {
    log.Fatalf("서버 빌드 실패: %v", err)
}

// 서버 시작
if err := s.Run(); err != nil {
    log.Fatalf("서버 시작 실패: %v", err)
}
```

서버 빌더는 다음과 같은 기능을 제공합니다:

1. 컨트롤러 추가: `AddController`, `AddControllers`
2. 미들웨어 추가: `AddMiddleware`, `AddMiddlewares`
3. 로깅 구성: `WithLogging`, `WithRemoteLogging`
4. 타임아웃 구성: `WithTimeout`
5. CORS 구성: `WithCORS`
6. 에러 핸들러 구성: `WithErrorHandler`
7. 기본 미들웨어 활성화:
   - `WithDefaultLogging`: 기본 로깅 미들웨어 활성화
   - `WithDefaultTimeout`: 기본 타임아웃 미들웨어 활성화
   - `WithDefaultCORS`: 기본 CORS 미들웨어 활성화
   - `WithDefaultErrorHandling`: 기본 에러 핸들러 미들웨어 활성화
8. 404 Not Found 및 405 Method Not Allowed 핸들러 구성:
   - `WithNoRoute`: 커스텀 404 Not Found 핸들러 설정
   - `WithNoMethod`: 커스텀 405 Method Not Allowed 핸들러 설정

서버 빌더는 컨트롤러에서 로깅 무시 경로(`GetLogIgnorePath`)와 인증 검사 무시 경로(`GetAuthCheckIgnorePath`)를 자동으로 수집하여 미들웨어에 전달합니다. 이를 통해 각 컨트롤러에서 무시할 경로를 지정하고, 서버 빌더가 이를 자동으로 처리하도록 할 수 있습니다.

## 라이센스

MIT 라이센스 - 이 소프트웨어는 MIT 라이센스 하에 배포됩니다. 자세한 내용은 LICENSE 파일을 참조하세요.
