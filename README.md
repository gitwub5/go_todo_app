# go_todo_app
GDG server session: todo app (Golang) 

refer to: https://github.com/budougumi0617/go_todo_app 

## 만들고자 하는 애플리케이션에 대해
이 리포지토리에서 만들고자 하는 웹 애플리케이션은 인증 기능이 포함된 TODO 작업을 관리하는 API 서버입니다.

최종적으로 다음의 엔드포인트를 구현할 예정입니다.

| HTTP 메서드 | 경로         | 설명                       |
|-------------|--------------|----------------------------|
| POST        | `/register`  | 새로운 사용자를 등록         |
| POST        | `/login`     | 등록된 사용자 정보로 액세스 토큰을 획득 |
| POST        | `/tasks`     | 액세스 토큰을 사용하여 작업을 등록 |
| GET         | `/tasks`     | 액세스 토큰을 사용하여 작업을 조회 |
| GET         | `/admin`     | 관리자 권한의 사용자만 접근 가능 |

`Docker Compose`를 이용하여 API 서버, MySQL, Redis를 시작합니다.   
주로 실행할 명령어는 `Makefile`에 사전에 정의되어 있습니다.
