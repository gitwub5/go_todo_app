package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitwub5/go_todo_app/entity"
	"github.com/go-playground/validator/v10"
)

/*
	vaildator 패키지는 Unmarshal한 타입에 validate라는 태그를 사용해 해당 필드의 검증 조건을 설정할 수 있다.
	설정한 조건은 *validator.Validate 메서드를 사용해 검증할 수 있다. 정의 완료된 조건도 있지만, 커스텀 조건도 정의할 수 있다.
*/

// AddTask 구조체는 HTTP 요청을 처리하고 새로운 Task를 저장하는 핸들러이다.
type AddTask struct {
	Service   AddTaskService // 비즈니스 로직을 처리하는 서비스
	Validator *validator.Validate
}

// ServeHTTP는 HTTP 요청을 처리하는 메서드로, AddTask 핸들러의 엔트리 포인트이다. (POST /task)
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 요청 본문에서 데이터를 읽어와서 구조체에 디코딩한다.
	var b struct {
		Title string `json:"title" validate:"required"` // Title 필드는 JSON에서 가져오며, 필수 값임을 검증합니다.
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		// 요청 본문 디코딩에 실패하면 에러 응답을 반환한다.
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 유효성 검사를 수행. Title 필드가 비어 있는 경우 에러가 반환된다.
	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t, err := at.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 성공 시 생성된 Task의 ID를 JSON 응답으로 반환한다.
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
