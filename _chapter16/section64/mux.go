package main

import "net/http"

// 어떤 처리를 어떤 URL 패스로 공개할지 라우팅하는 NewMux 함수 구현
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { // HTTP 서버가 실행중인지 확인하는 /health 엔드포인트
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 정적 분석 오류를 회피하기 위해 명시적으로 반환값을 버린다.
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
