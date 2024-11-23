package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		UserName string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	// 요청 본문에서 데이터를 읽어와서 구조체에 디코딩한다.
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 유효성 검사 수행
	err := l.Validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// 로그인 서비스 호출
	jwt, err := l.Service.Login(ctx, body.UserName, body.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt,
	}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
