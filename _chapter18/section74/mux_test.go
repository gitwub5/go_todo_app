package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()                              // 테스트용 레코드 생성
	r := httptest.NewRequest(http.MethodGet, "/health", nil) // 테스트용 요청 생성
	sut := NewMux()                                          // NewMux 함수를 사용하여 테스트 대상을 생성한다.
	sut.ServeHTTP(w, r)                                      // ServeHTTP 메서드를 사용하여 요청을 처리한다.
	resp := w.Result()                                       // 테스트용 레코드의 결과를 가져온다.
	t.Cleanup(func() { _ = resp.Body.Close() })              // 테스트 종료 시 레코드의 바디를 닫는다.
	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
