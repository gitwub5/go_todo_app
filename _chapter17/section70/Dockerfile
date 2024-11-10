# 다단계 빌드 방식을 활용하여 build 단계와 deploy 단계로 나누어 구성한다.
# build 단계에서는 애플리케이션을 빌드하고, deploy 단계에서는 최종 배포에 필요한 파일만 포함시켜 이미지를 경량화한다.

# <배포용 컨테이너에 포함할 바이너리를 생성하는 빌드 단계>
# Go는 빌드 후 단일 바이너리 파일만으로 배포할 수 있으므로, 최종 컨테이너에는 이 바이너리만 포함시키면 된다.
FROM golang:1.18.2-bullseye AS deploy-builder

# 애플리케이션 소스 코드 작업을 위한 작업 디렉터리 설정
WORKDIR /app

# go.mod 및 go.sum 파일을 복사하여 종속성을 먼저 설치
COPY go.mod go.sum ./
# Go 모듈 종속성 다운로드
RUN go mod download  

# 전체 소스 코드 복사
COPY . .

# 애플리케이션을 빌드하여 최종 바이너리 생성
# -trimpath 옵션은 경로 정보를 제거하여 바이너리 크기를 줄이고, -ldflags "-w -s" 옵션은 디버깅 정보를 제거해 경량화한다.
RUN go build -trimpath -ldflags "-w -s" -o app

# ----------------------------------------------------------------

# <배포용 컨테이너 구성 단계>
# 이 단계에서는 deploy-builder에서 생성된 바이너리 파일만 복사하여 배포용으로 최적화된 경량 컨테이너를 만든다.
FROM debian:bullseye-slim AS deploy

# 최신 패키지 정보를 가져와 배포 환경을 업데이트
RUN apt-get update

# build 단계에서 생성한 바이너리 파일을 최종 컨테이너로 복사
COPY --from=deploy-builder /app/app .

# 애플리케이션을 실행하도록 기본 명령어 설정
CMD ["./app"]


# ----------------------------------------------------------------

# <로컬 개발 환경에서 사용하는 자동 새로고침 환경을 위한 빌드 단계>
# 개발 시 변경 사항을 자동으로 반영하기 위한 환경을 설정한다.
FROM golang:1.23 AS dev 

# 애플리케이션 소스 코드 작업을 위한 작업 디렉터리 설정
WORKDIR /app

# air 패키지 설치 (로컬 개발 중 파일 변경 시 애플리케이션을 자동으로 다시 빌드하고 새로고침)
#RUN go install github.com/cosmtrek/air@latest
# air 패키지 최신 버전에서는 Go 1.23 이상의 버전을 요구하므로, Go 1.23 버전을 사용한다.
RUN go install github.com/air-verse/air@latest

# 개발 환경에서 air 명령어를 실행하여 자동 새로고침 기능 활성화
CMD ["air"]