package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body) // body를 JSON으로 인코딩
	if err != nil {
		fmt.Printf("encode response error: %v", err)
		w.WriteHeader(http.StatusInternalServerError) // 인코딩에 실패하면 500 에러를 반환
		rsp := ErrResponse{                           // ErrResponse 구조체를 사용하여 에러 응답을 생성
			Message: http.StatusText(http.StatusInternalServerError),
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil { // JSON으로 인코딩한 rsp를 응답으로 반환
			fmt.Printf("write error response error: %v", err)
		}
		return
	}

	w.WriteHeader(status) // 상태 코드를 설정
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
