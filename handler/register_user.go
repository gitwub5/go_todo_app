package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitwub5/go_todo_app/entity"
	"github.com/go-playground/validator/v10"
)

type RegisterUser struct {
	Service   RegisterUserService
	Validator *validator.Validate
}

func (ru *RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}

	// 요청 본문에서 데이터를 읽어와서 구조체에 디코딩한다.
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 유효성 검사를 수행. 필수 값이 비어 있는 경우 에러가 반환된다.
	if err := ru.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// 사용자 등록
	u, err := ru.Service.RegisterUser(ctx, b.Name, b.Password, b.Role)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID entity.UserID `json:"id"`
	}{ID: u.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
