package handler

import (
	"log"
	"net/http"

	"github.com/gitwub5/go_todo_app/auth"
)

/*
인증 및 권한이 필요한 엔드포인트를 만들어도 단순히 HTTP 헤더에 JWT로 저장되어 있으면 다른 service 패키지 등에서 참조할 수 없다.
따라서 미들웨어를 사용해 JWT에 포함된 인증 및 권한 정보를 추출하고, 이를 context에 저장해 다른 패키지에서 사용할 수 있도록 한다.
*/

// 액세스 토큰을 검증하고, 사용자 ID와 권한을 포함시킨 context를 반환하는 미들웨어를 반환한다.
func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// JWTer의 FillContext 메서드를 사용해 context에 사용자 ID와 권한을 저장한다.
			req, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

// 인증 정보로부터 admin 권한을 가진 사용자인지 확인하는 미들웨어를 반환한다.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("AdminMiddleware")
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
