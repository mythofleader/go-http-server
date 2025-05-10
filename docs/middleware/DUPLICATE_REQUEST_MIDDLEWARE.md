# 중복 요청 방지 미들웨어

중복 요청 방지 미들웨어는 동일한 요청이 중복으로 처리되는 것을 방지하는 기능을 제공합니다. 이 미들웨어는 요청 ID를 생성하고, 이 ID가 이미 처리된 요청인지 확인한 후, 중복된 요청인 경우 409 Conflict 응답을 반환합니다.

## 기능

- 요청 컨텍스트를 기반으로 고유한 요청 ID 생성
- 요청 ID의 중복 여부 확인
- 중복 요청 발견 시 409 Conflict 응답 반환
- 새로운 요청 ID 저장
- 사용자 정의 가능한 ID 생성 및 저장소 인터페이스

## 사용법

### 1. 요청 ID 생성기 구현하기

먼저, 요청 컨텍스트를 기반으로 요청 ID를 생성하는 `RequestIDGenerator` 인터페이스를 구현해야 합니다:

```go
// RequestIDGenerator는 요청 ID를 생성하는 인터페이스입니다.
type RequestIDGenerator interface {
    // GenerateRequestID는 컨텍스트에서 고유한 요청 ID를 생성합니다.
    GenerateRequestID(ctx context.Context) (string, error)
}

// MyRequestIDGenerator는 RequestIDGenerator 인터페이스의 구현체입니다.
type MyRequestIDGenerator struct {}

// GenerateRequestID는 요청에서 고유 ID를 생성합니다.
func (g *MyRequestIDGenerator) GenerateRequestID(ctx context.Context) (string, error) {
    // 요청에서 고유 식별자 추출 (예: 사용자 ID, 요청 경로, 요청 본문 해시 등)
    // 이 예제에서는 간단히 현재 시간과 랜덤 문자열을 조합하여 ID를 생성합니다.
    timestamp := time.Now().UnixNano()
    randomStr := uuid.New().String()

    // 고유 ID 생성
    requestID := fmt.Sprintf("%d-%s", timestamp, randomStr)
    return requestID, nil
}
```

### 2. 요청 ID 저장소 구현하기

다음으로, 요청 ID를 확인하고 저장하는 `RequestIDStorage` 인터페이스를 구현해야 합니다:

```go
// RequestIDStorage는 요청 ID를 확인하고 저장하는 인터페이스입니다.
type RequestIDStorage interface {
    // CheckRequestID는 요청 ID가 저장소에 존재하는지 확인합니다.
    CheckRequestID(requestID string) (bool, error)

    // SaveRequestID는 요청 ID를 저장소에 저장합니다.
    SaveRequestID(requestID string) error
}

// InMemoryRequestIDStorage는 메모리 기반 RequestIDStorage 구현체입니다.
type InMemoryRequestIDStorage struct {
    requestIDs map[string]bool
    mutex      sync.RWMutex
}

// NewInMemoryRequestIDStorage는 새로운 InMemoryRequestIDStorage를 생성합니다.
func NewInMemoryRequestIDStorage() *InMemoryRequestIDStorage {
    return &InMemoryRequestIDStorage{
        requestIDs: make(map[string]bool),
    }
}

// CheckRequestID는 요청 ID가 저장소에 존재하는지 확인합니다.
func (s *InMemoryRequestIDStorage) CheckRequestID(requestID string) (bool, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()

    _, exists := s.requestIDs[requestID]
    return exists, nil
}

// SaveRequestID는 요청 ID를 저장소에 저장합니다.
func (s *InMemoryRequestIDStorage) SaveRequestID(requestID string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.requestIDs[requestID] = true
    return nil
}
```

Redis와 같은 외부 저장소를 사용하는 구현체도 만들 수 있습니다:

```go
// RedisRequestIDStorage는 Redis 기반 RequestIDStorage 구현체입니다.
type RedisRequestIDStorage struct {
    client *redis.Client
    expiry time.Duration
}

// NewRedisRequestIDStorage는 새로운 RedisRequestIDStorage를 생성합니다.
func NewRedisRequestIDStorage(client *redis.Client, expiry time.Duration) *RedisRequestIDStorage {
    return &RedisRequestIDStorage{
        client: client,
        expiry: expiry,
    }
}

// CheckRequestID는 요청 ID가 Redis에 존재하는지 확인합니다.
func (s *RedisRequestIDStorage) CheckRequestID(requestID string) (bool, error) {
    exists, err := s.client.Exists(requestID).Result()
    if err != nil {
        return false, err
    }
    return exists > 0, nil
}

// SaveRequestID는 요청 ID를 Redis에 저장합니다.
func (s *RedisRequestIDStorage) SaveRequestID(requestID string) error {
    _, err := s.client.Set(requestID, "1", s.expiry).Result()
    return err
}
```

### 3. 미들웨어 구성하기

`DuplicateRequestConfig`를 생성하고 서버에 미들웨어를 추가합니다:

```go
// 메인 함수에서
func main() {
    // 서버 생성
    srv, _ := server.NewServer(server.FrameworkGin, "8080")

    // 요청 ID 생성기 및 저장소 생성
    idGenerator := &MyRequestIDGenerator{}
    idStorage := NewInMemoryRequestIDStorage()

    // 중복 요청 방지 미들웨어 구성
    dupReqConfig := &server.DuplicateRequestConfig{
        RequestIDGenerator: idGenerator,
        RequestIDStorage:   idStorage,
        ConflictMessage:    "중복 요청이 감지되었습니다", // 선택 사항: 사용자 정의 오류 메시지
    }

    // 서버에 미들웨어 추가
    srv.Use(server.DuplicateRequestMiddleware(dupReqConfig))

    // 또는 특정 라우트 그룹에 추가
    api := srv.Group("/api")
    api.Use(server.DuplicateRequestMiddleware(dupReqConfig))

    // 라우트 등록
    api.POST("/orders", createOrderHandler)

    // 서버 시작
    srv.Run()
}
```

### 기본 생성자 함수

중복 요청 방지 미들웨어는 기본 구성을 사용하는 생성자 함수를 제공합니다:

```go
// 기본 중복 요청 방지 미들웨어 생성자
// 주의: 이 함수는 RequestIDGenerator와 RequestIDStorage가 필요하므로 패닉이 발생합니다
// server.NewDefaultDuplicateRequestMiddleware()
```

`NewDefaultDuplicateRequestMiddleware` 함수는 `DefaultDuplicateRequestConfig()`를 호출하여 기본 구성을 생성하고, 이를 `DuplicateRequestMiddleware` 함수에 전달합니다. 그러나 `DefaultDuplicateRequestConfig()`는 `RequestIDGenerator`와 `RequestIDStorage` 필드를 설정하지 않으므로, 이 함수를 직접 호출하면 패닉이 발생합니다.

대신 다음과 같이 사용해야 합니다:

```go
// 기본 구성을 가져와서 필요한 필드 설정
dupReqConfig := server.DefaultDuplicateRequestConfig()
dupReqConfig.RequestIDGenerator = idGenerator
dupReqConfig.RequestIDStorage = idStorage
srv.Use(server.DuplicateRequestMiddleware(dupReqConfig))
```

## 오류 처리

미들웨어는 다음과 같은 경우에 오류를 반환합니다:

- 요청 ID를 생성할 수 없는 경우: 500 Internal Server Error
- 요청 ID를 확인할 수 없는 경우: 500 Internal Server Error
- 요청 ID가 이미 존재하는 경우: 409 Conflict
- 요청 ID를 저장할 수 없는 경우: 500 Internal Server Error

`ConflictMessage` 필드를 설정하여 중복 요청 오류 메시지를 사용자 정의할 수 있습니다.

## 주의사항

- 요청 ID 생성 로직은 애플리케이션의 요구사항에 맞게 신중하게 설계해야 합니다. 너무 광범위하면 중복 요청 방지 기능이 제대로 작동하지 않을 수 있고, 너무 구체적이면 정상적인 요청이 차단될 수 있습니다.
- 분산 시스템에서는 Redis와 같은 공유 저장소를 사용하여 모든 서버 인스턴스에서 요청 ID를 확인할 수 있도록 해야 합니다.
- 요청 ID 저장소에 적절한 만료 시간을 설정하여 오래된 요청 ID가 자동으로 제거되도록 하는 것이 좋습니다.
- 이 미들웨어는 멱등성이 필요한 API 엔드포인트(예: 결제 처리, 주문 생성 등)에 특히 유용합니다.

## 전체 예제

전체 작동 예제는 [중복 요청 방지 예제](../../examples/duprequest/main.go)를 참조하세요.
