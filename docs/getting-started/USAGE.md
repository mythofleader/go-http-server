# Tenqube Go HTTP 서버 사용 방법

이 문서는 Go 프로젝트에서 Tenqube HTTP 서버 라이브러리를 사용하는 방법에 대한 자세한 지침을 제공합니다.

## 설치

라이브러리를 설치하려면 다음 명령을 사용하세요:

```bash
go get github.com/tenqube/tenqube-go-http-server
```

## 디렉토리 구조

이 라이브러리는 다음과 같은 디렉토리 구조로 구성되어 있습니다:

- `core/`: 핵심 인터페이스와 타입 정의
  - `context.go`: Context, Server, RouterGroup 등의 인터페이스 정의
  - `gin/`: Gin 프레임워크 구현
    - `server.go`: Gin 프레임워크를 사용한 서버 구현
    - `logging.go`: Gin 프레임워크를 위한 로깅 미들웨어 구현
    - `errorhandler.go`: Gin 프레임워크를 위한 에러 핸들러 미들웨어 구현
  - `std/`: 표준 net/http 패키지 구현
    - `server.go`: 표준 net/http 패키지를 사용한 서버 구현
    - `logging.go`: 표준 net/http 패키지를 위한 로깅 미들웨어 구현
    - `errorhandler.go`: 표준 net/http 패키지를 위한 에러 핸들러 미들웨어 구현
  - `middleware/`: 공통 미들웨어 기능
    - `logging.go`: 로깅 미들웨어 구현
    - `timeout.go`: 타임아웃 미들웨어 구현
    - `errorhandler.go`: 에러 핸들러 미들웨어 구현
    - `auth.go`: 인증 미들웨어 구현
    - `apikey.go`: API 키 미들웨어 구현
    - `cors.go`: CORS 미들웨어 구현
    - `errors/`: 에러 클래스 정의
- `server.go`: 루트 패키지에서 서버 생성 함수 제공

이 구조는 라이브러리를 확장하거나 커스터마이징하려는 사용자에게 유용합니다. 예를 들어, 새로운 HTTP 프레임워크를 추가하려면 해당 프레임워크에 대한 새 디렉토리를 만들고 core 인터페이스를 구현하면 됩니다.

## 기본 사용법

다음은 라이브러리를 사용하여 HTTP 서버를 생성하는 간단한 예제입니다:

```go
package main

import (
	"log"
	"net/http"

	server "github.com/tenqube/tenqube-go-http-server"
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

## 프레임워크 선택

이 라이브러리는 여러 HTTP 프레임워크를 지원합니다. 서버를 생성할 때 사용할 프레임워크를 지정할 수 있습니다:

```go
// Gin 프레임워크 사용 (기본값)
s, err := server.NewServer(server.FrameworkGin, "8080")

// 표준 net/http 패키지 사용
s, err := server.NewServer(server.FrameworkStdHTTP, "8080")
```

## AWS Lambda 지원

AWS Lambda를 사용할 때는 Gin 프레임워크로 서버를 생성한 다음, `Run` 대신 `StartLambda` 메서드를 사용해야 합니다. **중요: Lambda는 Gin 프레임워크에서만 지원되며, 표준 HTTP 서버에서는 지원되지 않습니다.**

## 404 Not Found 및 405 Method Not Allowed 핸들러

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

### Gin 프레임워크와 Lambda 사용하기

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

// Lambda 핸들러 시작
if err := s.StartLambda(); err != nil {
	log.Fatalf("Lambda 시작 실패: %v", err)
}
```

### 표준 HTTP 서버와 Lambda 사용 시 주의사항

표준 HTTP 서버에서는 Lambda를 지원하지 않습니다. 표준 HTTP 서버에서 `StartLambda` 메서드를 호출하면 다음과 같은 오류가 반환됩니다:

```go
// 표준 HTTP 서버 생성
s, err := server.NewServer(server.FrameworkStdHTTP, "8080")
if err != nil {
	log.Fatalf("서버 생성 실패: %v", err)
}

// Lambda 핸들러 시작 시도
err = s.StartLambda()
if err != nil {
	// 다음 오류가 반환됨: "Lambda is only supported with the Gin framework"
	log.Fatalf("Lambda 시작 실패: %v", err)
}
```

### AWS Lambda와 API 프록시 어댑터 사용하기

이 라이브러리는 내부적으로 `github.com/awslabs/aws-lambda-go-api-proxy` 패키지를 사용하여 Gin 서버를 Lambda 핸들러로 변환합니다:

- Gin 프레임워크: `github.com/awslabs/aws-lambda-go-api-proxy/gin`

이 패키지를 사용하려면 다음 명령으로 설치해야 합니다:

```bash
go get github.com/awslabs/aws-lambda-go-api-proxy
```

그런 다음 `StartLambda` 메서드를 사용하여 Lambda 핸들러를 시작할 수 있습니다. 이 메서드는 내부적으로 다음과 같은 코드를 사용합니다:

#### Gin 프레임워크

```go
var ginLambda *ginadapter.GinLambdaALB
ginLambda = ginadapter.NewALB(ginEngine)

func handler(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
    return ginLambda.ProxyWithContext(ctx, req)
}

lambda.Start(handler)
```

## 고급 사용법

### 미들웨어

로깅, 인증 등과 같은 횡단 관심사를 처리하기 위해 서버에 미들웨어를 추가할 수 있습니다.

```go
// 미들웨어 추가
s.Use(func(c server.Context) {
	log.Printf("요청: %s %s", c.Request().Method, c.Request().URL.Path)
	// 요청 처리 계속
})
```

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
[GIN] Adding middleware: github.com/tenqube/tenqube-go-http-server/core/middleware.ErrorHandlerMiddleware.func1
[GIN] Adding middleware: github.com/tenqube/tenqube-go-http-server/core/middleware.LoggingMiddleware.func1
[GIN] Server starting on :8080
[GIN] Using Gin framework version: 1.9.1
[GIN] Middleware registered:
[GIN]   1. github.com/tenqube/tenqube-go-http-server/core/middleware.ErrorHandlerMiddleware.func1
[GIN]   2. github.com/tenqube/tenqube-go-http-server/core/middleware.LoggingMiddleware.func1
[GIN] Routes registered:
[GIN]   1. GET /
[GIN]   2. GET /api/users
[GIN] Server is ready to handle requests
```

표준 HTTP 서버를 사용할 경우 로그 접두사는 `[STD]`로 표시됩니다.

### 라우터 그룹

라우터 그룹을 사용하면 관련 라우트를 함께 그룹화하고 미들웨어를 적용할 수 있습니다.

```go
// 라우터 그룹 생성
api := s.Group("/api")
{
	// 그룹에 미들웨어 추가
	api.Use(authMiddleware)

	// 라우트 등록
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

`JSON` 메서드를 사용하여 쉽게 JSON 응답을 반환할 수 있습니다.

```go
func jsonHandler(c server.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello, JSON!",
		"status":  "success",
	})
}
```

### 요청 바인딩

`Bind` 또는 `BindJSON` 메서드를 사용하여 요청 데이터를 구조체에 바인딩할 수 있습니다.

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

`RunTLS` 메서드를 사용하여 TLS 지원으로 서버를 시작할 수 있습니다.

```go
// TLS로 서버 시작
if err := s.RunTLS(":8443", "cert.pem", "key.pem"); err != nil {
	log.Fatalf("서버 시작 실패: %v", err)
}
```

### 서버 종료

서버를 종료하는 방법에는 두 가지가 있습니다:

#### 즉시 종료 (Stop)

`Stop` 메서드를 사용하여 서버를 즉시 종료할 수 있습니다. 이 메서드는 현재 진행 중인 연결을 기다리지 않고 서버를 즉시 중지합니다.

```go
// 서버 즉시 종료
if err := s.Stop(); err != nil {
	log.Fatalf("서버 종료 실패: %v", err)
}
log.Println("서버가 즉시 중지되었습니다")
```

#### 정상 종료 (Shutdown)

`Shutdown` 메서드를 사용하여 서버를 정상적으로 종료할 수 있습니다. 이 메서드는 현재 진행 중인 연결이 완료될 때까지 기다린 후 서버를 중지합니다.

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

	// 서버 정상 종료
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("서버 종료 실패: %v", err)
	}
	log.Println("서버가 정상적으로 중지되었습니다")
}()
```

## 명령줄 예제

라이브러리에는 명령줄 플래그를 사용하여 프레임워크를 선택하고 Lambda 모드를 활성화하는 방법을 보여주는 예제가 포함되어 있습니다:

```go
package main

import (
	"flag"
	"log"
	"net/http"

	server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
	// 명령줄 플래그 파싱
	framework := flag.String("framework", "gin", "사용할 HTTP 프레임워크 (gin, std)")
	lambdaMode := flag.Bool("lambda", false, "AWS Lambda 모드로 실행")
	flag.Parse()

	// 지정된 프레임워크를 기반으로 새 서버 생성
	var s server.Server
	var err error

	switch *framework {
	case "gin":
		s, err = server.NewServer(server.FrameworkGin, "8080")
	case "std":
		s, err = server.NewServer(server.FrameworkStdHTTP, "8080")
	default:
		// 기본값은 Gin
		s, err = server.NewServer(server.FrameworkGin, "8080")
	}

	if err != nil {
		log.Fatalf("서버 생성 실패: %v", err)
	}

	// 라우트 등록...

	// 서버 시작
	if *lambdaMode {
		// Lambda 모드인 경우 StartLambda 사용
		log.Println("Lambda 서버 시작")
		// Lambda 모드 시작
		if err := s.StartLambda(); err != nil {
			log.Fatalf("Lambda 시작 실패: %v", err)
		}
	} else {
		// 일반 모드인 경우 Run 사용
		log.Printf("%s 프레임워크로 서버 시작 (:8080)", *framework)
		if err := s.Run(); err != nil {
			log.Fatalf("서버 시작 실패: %v", err)
		}
	}
}
```

## 서버 빌더 사용하기

서버 빌더를 사용하면 컨트롤러와 미들웨어를 쉽게 구성하고, 컨트롤러에서 로깅 무시 경로와 인증 검사 무시 경로를 자동으로 수집할 수 있습니다.

```go
// 서버 빌더 생성 (방법 1: 프레임워크와 포트 지정)
builder := server.NewServerBuilder(server.FrameworkGin, "8080")

// 또는 포트를 나중에 설정하는 방법 (자동으로 8000-9000 사이의 사용 가능한 포트 할당)
// builder := server.NewServerBuilder(server.FrameworkGin)
// builder.WithDefaultPort()

// 또는 Gin 프레임워크를 사용하는 서버 빌더 생성 (방법 3: 더 간단한 방법)
// builder := server.NewGinServerBuilder()
// builder.WithDefaultPort() // 자동으로 8000-9000 사이의 사용 가능한 포트 할당

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
    AllowedDomains: []string{"*"},
    AllowedMethods: "GET, POST, PUT, DELETE, PATCH",
})

// 에러 핸들러 구성 추가
builder.WithErrorHandler(server.ErrorHandlerConfig{
    DefaultErrorMessage: "서버 오류가 발생했습니다",
    DefaultStatusCode:   500,
    LogErrors:           true,
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
3. 포트 구성: `WithDefaultPort`: 8000-9000 사이의 사용 가능한 포트를 자동으로 할당 (NewServerBuilder에서 포트를 지정하지 않은 경우 필수)
4. 로깅 구성: `WithLogging`, `WithRemoteLogging`
5. 타임아웃 구성: `WithTimeout`
6. CORS 구성: `WithCORS`
7. 에러 핸들러 구성: `WithErrorHandler`
8. 기본 미들웨어 활성화:
   - `WithDefaultLogging(console ...bool)`: 기본 로깅 미들웨어 활성화 (console 파라미터로 콘솔 로깅 활성화/비활성화 가능, 파라미터가 없으면 기본값은 true)
   - `WithDefaultTimeout`: 기본 타임아웃 미들웨어 활성화
   - `WithDefaultCORS`: 기본 CORS 미들웨어 활성화
   - `WithDefaultErrorHandling`: 기본 에러 핸들러 미들웨어 활성화

서버 빌더는 컨트롤러에서 로깅 무시 경로(`GetLogIgnorePath`)와 인증 검사 무시 경로(`GetAuthCheckIgnorePath`)를 자동으로 수집하여 미들웨어에 전달합니다. 이를 통해 각 컨트롤러에서 무시할 경로를 지정하고, 서버 빌더가 이를 자동으로 처리하도록 할 수 있습니다.

### 서버 빌더 사용 예시

다음은 서버 빌더를 사용하여 컨트롤러와 미들웨어를 구성하는 예시입니다:

```go
package main

import (
    "log"
    "time"

    server "github.com/tenqube/tenqube-go-http-server"
)

func main() {
    // 서버 빌더 생성 (방법 1: 프레임워크와 포트 지정)
    builder := server.NewServerBuilder(server.FrameworkGin, "8080")

    // 또는 Gin 프레임워크와 8080 포트를 사용하는 서버 빌더 생성 (방법 2: 더 간단한 방법)
    // builder := server.NewGinServerBuilder()

    // 컨트롤러 추가
    userController := &UserController{userService: myUserService}
    productController := &ProductController{productService: myProductService}
    builder.AddControllers(userController, productController)

    // 로깅 구성 추가
    builder.WithLogging(map[string]string{
        "environment": "development",
        "version":     "1.0.0",
    })

    // 타임아웃 구성 추가
    builder.WithTimeout(server.TimeoutConfig{
        Timeout: 5 * time.Second,
    })

    // 서버 빌드 및 시작
    s, err := builder.Build()
    if err != nil {
        log.Fatalf("서버 빌드 실패: %v", err)
    }

    log.Println("서버가 :8080 포트에서 시작됩니다")
    if err := s.Run(); err != nil {
        log.Fatalf("서버 시작 실패: %v", err)
    }
}
```

이 예시에서는 서버 빌더를 사용하여 컨트롤러와 미들웨어를 구성하고, 서버를 빌드하여 시작합니다. 서버 빌더는 컨트롤러에서 로깅 무시 경로와 인증 검사 무시 경로를 자동으로 수집하여 미들웨어에 전달합니다.

## 테스트

라이브러리에는 기능을 검증하는 테스트가 포함되어 있습니다. 다음 명령을 사용하여 테스트를 실행할 수 있습니다:

```bash
go test github.com/tenqube/tenqube-go-http-server
```

## 기여

기여는 환영합니다! Pull Request를 자유롭게 제출해 주세요.
