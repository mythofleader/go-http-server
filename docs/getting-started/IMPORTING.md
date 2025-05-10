# 다른 Go 프로젝트에서 이 패키지를 가져오는 방법

이 문서는 다른 Go 프로젝트에서 `go-http-server` 패키지를 가져올 수 있도록 하는 방법을 설명합니다.

## 현재 상태

이 패키지는 이미 `go.mod` 파일에 모듈 경로 `github.com/mythofleader/go-http-server`로 Go 모듈로 설정되어 있습니다. 이는 패키지가 게시 및 가져오기가 가능한 상태임을 의미합니다.

## GitHub에 게시하기

다른 Go 프로젝트에서 이 패키지를 가져올 수 있도록 하려면 GitHub에 게시해야 합니다:

1. `github.com/mythofleader/go-http-server`에 GitHub 저장소 생성
2. 코드를 이 저장소에 푸시:

```bash
# git을 아직 초기화하지 않았다면
git init
git add .
git commit -m "Initial commit"

# GitHub 저장소를 원격으로 추가
git remote add origin https://github.com/mythofleader/go-http-server.git

# GitHub에 푸시
git push -u origin main
```

## 패키지 버전 관리

Go 모듈은 시맨틱 버전 관리를 사용합니다. 패키지의 버전을 생성하려면:

1. 시맨틱 버전으로 릴리스에 태그 지정:

```bash
git tag v0.1.0
git push origin v0.1.0
```

주요 버전 업데이트(v2 이상)의 경우 `go.mod`의 모듈 경로를 업데이트해야 합니다:

```
module github.com/mythofleader/go-http-server/v2
```

## 다른 프로젝트에서 이 패키지를 가져오는 방법

패키지가 GitHub에 게시되고 버전이 태그되면 다른 프로젝트에서 다음과 같이 가져올 수 있습니다:

```bash
go get github.com/mythofleader/go-http-server
```

Go 코드에서는 다음과 같이 가져올 수 있습니다:

```go
import "github.com/mythofleader/go-http-server"
```

또는 별칭을 사용하여:

```go
import server "github.com/mythofleader/go-http-server"
```

## 빌드나 압축이 필요 없음

Go 모듈은 소스 저장소에서 직접 가져옵니다. 패키지를 빌드하거나 압축할 필요가 없습니다. 누군가 패키지를 가져오면 Go는 자동으로 GitHub에서 소스 코드를 다운로드하고 해당 프로젝트의 일부로 컴파일합니다.

## 비공개 저장소

저장소가 비공개인 경우, 사용자는 Go가 패키지를 다운로드할 때 GitHub로 인증하도록 Git을 구성해야 합니다. 다음을 사용할 수 있습니다:

1. GOPRIVATE 환경변수 설정
```aiignore
# 단일 저장소의 경우
go env -w GOPRIVATE=github.com/your-org/your-repo

# 조직의 모든 저장소의 경우
go env -w GOPRIVATE=github.com/your-org/*

# 여러 도메인/조직의 경우 쉼표로 구분
go env -w GOPRIVATE=github.com/your-org/*,gitlab.com/another-org/*
```

2. SSH 키 생성 (없는 경우):
```aiignore
ssh-keygen -t ed25519 -C "your-email@example.com"
```

3. GitHub에 SSH 공개키 등록:
    - GitHub → Settings → SSH and GPG keys → New SSH key
    - `~/.ssh/id_ed25519.pub` 내용을 복사하여 등록

4. `~/.gitconfig` 설정
```aiignore
[url "ssh://git@github.com/"]
    insteadOf = https://github.com/
```

5. Git 연결 테스트
```bash
ssh -T git@github.com
```

## 패키지 업데이트

패키지를 변경할 때:

1. 코드 변경
2. 변경 사항 커밋
3. 새 버전 태그 지정
4. 변경 사항과 새 태그를 GitHub에 푸시

사용자는 다음과 같이 최신 버전으로 업데이트할 수 있습니다:

```bash
go get -u github.com/mythofleader/go-http-server
```

또는 특정 버전으로:

```bash
go get github.com/mythofleader/go-http-server@v0.2.0
```

## 패키지 테스트

게시하기 전에 모든 테스트가 통과하는지 확인하세요:

```bash
go test ./...
```

이렇게 하면 패키지를 가져오는 사용자가 작동하는 버전을 얻을 수 있습니다.
