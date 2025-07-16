# 인증 미들웨어

인증 미들웨어는 HTTP 기본 인증(Basic Authentication) 또는 JWT 베어러 토큰(Bearer Token)을 사용하여 사용자를 인증하는 방법을 제공합니다. 이 미들웨어는 Authorization 헤더에서 자격 증명을 추출하고 사용자 조회 인터페이스를 사용하여 해당 사용자를 찾습니다.

## 기능

- HTTP 기본 인증 지원
- JWT 베어러 토큰 지원
- 사용할 인증 방법 선택 옵션
- 사용자 정의 가능한 사용자 조회 인터페이스
- 인증 실패에 대한 오류 처리
- 요청 컨텍스트에서 인증된 사용자 쉽게 검색

## 사용법

### 1. UserLookupInterface 구현하기

먼저, 자격 증명을 기반으로 사용자를 조회하기 위한 `UserLookupInterface`를 구현해야 합니다:

```go
type UserLookupInterface interface {
    // LookupUserByBasicAuth looks up a user by username and password
    LookupUserByBasicAuth(username, password string) (interface{}, error)

    // LookupUserByJWT looks up a user by JWT claims
    LookupUserByJWT(claims MapClaims) (interface{}, error)
}
```

구현 예시:

```go
type UserService struct {
    // 사용자 저장소 (예: 데이터베이스)
    users map[string]User
}

func (s *UserService) LookupUserByBasicAuth(username, password string) (interface{}, error) {
    // 사용자 이름으로 사용자 조회
    user, exists := s.users[username]
    if !exists {
        return nil, errors.New("user not found")
    }

    // 비밀번호 확인 (실제 앱에서는 적절한 비밀번호 해싱 사용)
    if !verifyPassword(user.PasswordHash, password) {
        return nil, errors.New("invalid password")
    }

    return user, nil
}

func (s *UserService) LookupUserByJWT(claims MapClaims) (interface{}, error) {
    // 클레임에서 사용자 ID 또는 사용자 이름 추출
    sub, ok := claims["sub"].(string)
    if !ok {
        return nil, errors.New("invalid token: missing subject")
    }

    // 사용자 조회
    user, exists := s.users[sub]
    if !exists {
        return nil, errors.New("user not found")
    }

    return user, nil
}
```

### 2. 미들웨어 구성하기

`AuthConfig`를 생성하고 서버에 미들웨어를 추가합니다:

#### 특정 경로 무시하기

인증 미들웨어는 특정 경로에 대한 인증 검사를 건너뛸 수 있습니다. `SkipPaths` 필드에 건너뛸 경로 목록을 설정하여 해당 경로에 대한 인증 검사를 비활성화할 수 있습니다:

```go
authConfig := &server.AuthConfig{
    UserLookup: userService,
    AuthType:   server.AuthTypeJWT,
    JWTSecret:  "your-secret-key",
    SkipPaths: []string{
        "/health",
        "/metrics",
        "/public",
        "/api/*",         // 와일드카드 패턴 - /api로 시작하는 모든 경로
        "/user/:id",      // 파라미터 패턴 - /user/123, /user/abc 등 모든 사용자 ID 경로
    },
}
s.Use(server.AuthMiddleware(authConfig))
```

이렇게 하면 다음 경로에 대한 요청은 인증 검사를 건너뛰게 됩니다:
- 정확한 경로 매칭: `/health`, `/metrics`, `/public`
- 와일드카드 매칭: `/api/users`, `/api/products` 등 `/api`로 시작하는 모든 경로
- 파라미터 패턴 매칭: `/user/123`, `/user/abc` 등 `/user/:id` 형식의 모든 경로

`SkipPaths`는 다음 세 가지 방식으로 경로를 매칭합니다:
1. 정확한 경로 매칭: 지정된 경로와 정확히 일치하는 경우
2. 와일드카드 매칭: `*` 문자를 사용하여 여러 경로를 매칭 (예: `/api/*`)
3. 파라미터 패턴 매칭: `:` 접두사를 사용하여 경로 세그먼트의 파라미터를 매칭 (예: `/user/:id`)

#### 기본 생성자 함수

인증 미들웨어는 인증 유형에 따라 두 가지 기본 생성자 함수를 제공합니다:

##### JWT 인증을 위한 기본 생성자

```go
// JWT 인증을 위한 기본 생성자 (JWTUserLookup과 JWT 비밀 키를 인자로 받음)
s.Use(server.NewDefaultJWTAuthMiddleware(myJWTLookup, "your-jwt-secret"))
```

`NewDefaultJWTAuthMiddleware` 함수는 JWTUserLookup 인터페이스 구현체와 JWT 비밀 키를 인자로 받아 `DefaultAuthConfig()`를 호출하여 기본 구성을 생성하고, AuthType을 AuthTypeJWT로 설정한 다음 이를 `AuthMiddleware` 함수에 전달합니다.

##### 기본 인증을 위한 기본 생성자

```go
// 기본 인증을 위한 기본 생성자 (BasicAuthUserLookup을 인자로 받음)
s.Use(server.NewDefaultBasicAuthMiddleware(myBasicAuthLookup))
```

`NewDefaultBasicAuthMiddleware` 함수는 BasicAuthUserLookup 인터페이스 구현체를 인자로 받아 `DefaultAuthConfig()`를 호출하여 기본 구성을 생성하고, AuthType을 AuthTypeBasic으로 설정한 다음 이를 `AuthMiddleware` 함수에 전달합니다.

##### 더 많은 설정이 필요한 경우

더 많은 설정이 필요한 경우 다음과 같이 사용할 수 있습니다:

```go
// JWT 인증을 위한 구성
authConfig := server.DefaultAuthConfig()
authConfig.AuthType = server.AuthTypeJWT
authConfig.JWTLookup = myJWTLookup
authConfig.JWTSecret = "your-jwt-secret"
authConfig.UnauthorizedMessage = "사용자 정의 오류 메시지"
s.Use(server.AuthMiddleware(authConfig))

// 또는 기본 인증을 위한 구성
authConfig := server.DefaultAuthConfig()
authConfig.AuthType = server.AuthTypeBasic
authConfig.BasicAuthLookup = myBasicAuthLookup
authConfig.UnauthorizedMessage = "사용자 정의 오류 메시지"
s.Use(server.AuthMiddleware(authConfig))
```

##### 이전 버전과의 호환성

이전 버전과의 호환성을 위해 `NewDefaultAuthMiddleware` 함수도 제공되지만, 이 함수는 더 이상 사용되지 않으며 `NewDefaultJWTAuthMiddleware` 또는 `NewDefaultBasicAuthMiddleware` 함수를 사용하는 것이 좋습니다.

```go
// 이전 버전과의 호환성을 위한 함수 (더 이상 사용되지 않음)
// 주의: 이 함수는 UserLookup과 JWTSecret이 필요하므로 패닉이 발생합니다
// server.NewDefaultAuthMiddleware()
```

```go
// 사용자 서비스 생성
userService := NewUserService()

// 기본 인증을 위한 인증 미들웨어 구성
basicAuthConfig := &server.AuthConfig{
    UserLookup: userService,
    AuthType:   server.AuthTypeBasic,
    
    // 선택 사항: 사용자 정의 오류 메시지
    UnauthorizedMessage: "Authentication required",
    ForbiddenMessage:    "Access denied",
}
    
// 또는 JWT 인증을 위한 인증 미들웨어 구성
jwtAuthConfig := &server.AuthConfig{
    UserLookup: userService,
    AuthType:   server.AuthTypeJWT,
    JWTSecret:  "your-secret-key",
    
    // 선택 사항: 사용자 정의 오류 메시지
    UnauthorizedMessage: "Authentication required",
    ForbiddenMessage:    "Access denied",
}

// 라우트 그룹에 미들웨어 추가
protectedBasic := server.Group("/api/basic")
protectedBasic.Use(server.AuthMiddleware(basicAuthConfig))

protectedJWT := server.Group("/api/jwt")
protectedJWT.Use(server.AuthMiddleware(jwtAuthConfig))

// 또는 특정 라우트에 추가
server.GET("/protected-basic", server.AuthMiddleware(basicAuthConfig), handleProtectedRoute)
server.GET("/protected-jwt", server.AuthMiddleware(jwtAuthConfig), handleProtectedRoute)
```

### 3. 인증된 사용자 접근하기

라우트 핸들러에서 요청 컨텍스트에서 인증된 사용자에 접근할 수 있습니다:

```go
func handleProtectedRoute(c server.Context) {
    // 컨텍스트에서 인증된 사용자 가져오기
    user, ok := server.GetUserFromContext(c.Request().Context())
    if !ok {
        c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
        return
    }

    // 사용자 타입으로 타입 변환
    u, ok := user.(User)
    if !ok {
        c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
        return
    }

    // 사용자 데이터 사용
    c.JSON(http.StatusOK, map[string]interface{}{
        "id":       u.ID,
        "username": u.Username,
        "role":     u.Role,
    })
}
```

## 인증 방법

### 기본 인증

기본 인증의 경우, 클라이언트는 다음 형식의 Authorization 헤더를 보내야 합니다:

```
Authorization: Basic base64(username:password)
```

### JWT 베어러 토큰

JWT 베어러 토큰 인증의 경우, 클라이언트는 다음 형식의 Authorization 헤더를 보내야 합니다:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

JWT 토큰은 `AuthConfig`에 지정된 비밀 키를 사용하여 HMAC-SHA256(HS256)으로 서명되어야 합니다.

## 오류 처리

미들웨어는 인증 실패에 대해 적절한 HTTP 상태 코드를 반환합니다:

- 401 Unauthorized: 자격 증명이 없거나 유효하지 않음
- 403 Forbidden: 유효한 자격 증명이지만 권한 부족

`AuthConfig`의 `UnauthorizedMessage` 및 `ForbiddenMessage` 필드를 설정하여 오류 메시지를 사용자 정의할 수 있습니다.

## 전체 예제

전체 작동 예제는 [인증 예제](../../examples/auth/main.go)를 참조하세요.
