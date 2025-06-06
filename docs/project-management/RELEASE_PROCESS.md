# 자동 태그 및 릴리스 프로세스

이 문서는 `tenqube-go-http-server` 프로젝트의 자동 태그 및 릴리스 프로세스에 대해 설명합니다.

## 현재 구현

GitHub Actions를 사용하여 자동 버전 관리 및 릴리스 시스템을 구현했습니다. 이 시스템은 네 가지 주요 워크플로우로 구성되어 있습니다:

1. **기본 릴리스 워크플로우** (`release.yml`): 코드가 main 브랜치에 푸시될 때 자동으로 태그와 릴리스 생성
2. **버전 검증 워크플로우** (`version-check.yml`): PR 단계에서 버전 검증
3. **변경 로그 업데이트 워크플로우** (`update-changedlog.yml`): PR 병합 시 docs/CHANGELOG.md 자동 업데이트
4. **Simple 실행 파일 업데이트 워크플로우** (`update-simple-executable.yml`): 버전 변경 시 Simple 실행 파일 자동 빌드 및 릴리스

### 기본 릴리스 워크플로우 작동 방식

1. `main` 브랜치에 코드가 푸시되면 GitHub Actions 워크플로우가 트리거됩니다.
2. 워크플로우는 `version.go` 파일에서 현재 버전을 추출합니다.
3. 추출한 버전으로 GitHub 릴리스를 생성합니다.

### 버전 검증 워크플로우 작동 방식

이 워크플로우는 PR이 `main` 브랜치로 병합되기 전에 실행됩니다:

1. PR 브랜치와 `main` 브랜치에서 각각 버전을 추출합니다.
2. PR의 버전이 `main` 브랜치의 버전과 다른지 확인합니다.
3. PR의 버전이 `main` 브랜치의 버전보다 큰지 확인합니다.
4. 버전이 시맨틱 버전 관리 형식(MAJOR.MINOR.PATCH)을 따르는지 검증합니다.
5. 위 조건 중 하나라도 충족되지 않으면 워크플로우가 실패하고 PR을 병합할 수 없습니다.

이 워크플로우는 버전이 항상 증가하도록 보장하고, 시맨틱 버전 관리 규칙을 따르도록 합니다.

### 변경 로그 업데이트 워크플로우 작동 방식

이 워크플로우는 PR이 `main` 브랜치에 병합될 때 docs/CHANGELOG.md 파일을 자동으로 업데이트합니다:

1. PR이 `main` 브랜치에 병합되면 워크플로우가 트리거됩니다.
2. PR의 제목, 번호, 작성자, 내용, 병합 날짜 등의 정보를 추출합니다.
3. 추출한 정보를 사용하여 docs/CHANGELOG.md 파일에 새 항목을 추가합니다:
   - 파일이 존재하지 않으면 새로 생성합니다.
   - 파일이 이미 존재하면 새 항목을 파일 상단에 추가합니다.
4. 변경 사항을 커밋하고 `main` 브랜치에 푸시합니다.

이 워크플로우는 PR의 내용을 기반으로 변경 로그를 자동으로 생성하므로, PR 작성자는 PR 설명에 변경 사항을 명확하게 기술하는 것이 중요합니다.

### Simple 실행 파일 업데이트 워크플로우 작동 방식

이 워크플로우는 `main` 브랜치에 코드가 푸시되거나 PR이 병합될 때 Simple 실행 파일을 자동으로 빌드하고 릴리스합니다:

1. `main` 브랜치에 코드가 푸시되거나 PR이 병합되면 워크플로우가 트리거됩니다.
2. `version.go` 파일에서 현재 버전을 추출합니다.
3. 다양한 플랫폼(Linux, macOS, Windows)에 대한 Simple 실행 파일을 빌드합니다.
4. 사용 방법이 포함된 README 파일을 생성합니다.
5. 모든 파일을 포함하는 ZIP 아카이브를 생성합니다.
6. 해당 버전의 태그가 이미 존재하는지 확인합니다.
7. 태그가 존재하지 않으면 새 태그를 생성합니다.
8. 빌드된 실행 파일과 ZIP 아카이브를 포함하는 GitHub 릴리스를 생성하거나 업데이트합니다.

이 워크플로우는 버전이 변경될 때마다 자동으로 Simple 실행 파일을 빌드하고 릴리스하므로, 사용자는 항상 최신 버전의 실행 파일을 쉽게 다운로드할 수 있습니다.

## 버전 관리 프로세스

### 버전 업데이트 방법

새 기능이나 버그 수정을 추가할 때는 다음 단계를 따릅니다:

1. 새 브랜치를 생성합니다.
2. 코드 변경을 수행합니다.
3. `version.go` 파일에서 버전을 업데이트합니다:
   - 새 기능 추가: MINOR 버전 증가 (예: 0.1.0 → 0.2.0)
   - 버그 수정: PATCH 버전 증가 (예: 0.1.0 → 0.1.1)
   - 주요 변경: MAJOR 버전 증가 (예: 0.1.0 → 1.0.0)
4. PR을 생성하고 `main` 브랜치로 병합합니다.

버전 검증 워크플로우는 PR이 `main` 브랜치로 병합되기 전에 버전이 올바르게 업데이트되었는지 확인합니다.

### 시맨틱 버전 관리

이 프로젝트는 [시맨틱 버전 관리](https://semver.org/) 규칙을 따릅니다:

- MAJOR 버전: 이전 버전과 호환되지 않는 API 변경
- MINOR 버전: 이전 버전과 호환되는 방식으로 기능 추가
- PATCH 버전: 이전 버전과 호환되는 버그 수정

버전 검증 워크플로우는 버전이 MAJOR.MINOR.PATCH 형식을 따르는지 확인합니다.

## 권장 사항

1. **버전 업데이트 자동화**: 개발 워크플로우에 버전 업데이트를 통합하여 실수를 방지합니다.
2. **PR 체크리스트 사용**: PR 템플릿에 버전 업데이트 확인 항목을 추가합니다.
3. **버전 검증 워크플로우 사용**: `version-check.yml`을 사용하여 PR이 `main`에 병합되기 전에 버전이 올바르게 업데이트되었는지 확인합니다.
4. **변경 로그 유지**: PR 설명에 변경 사항을 명확하게 기술하여 자동 생성된 변경 로그가 유용하도록 합니다.
5. **PR 설명 작성**: PR 설명은 변경 로그에 자동으로 포함되므로, 명확하고 자세한 설명을 작성하는 것이 중요합니다.

이러한 접근 방식을 통해 버전 관리 프로세스와 변경 로그 관리를 더욱 견고하고 자동화할 수 있습니다.
