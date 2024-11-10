package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gitwub5/go_todo_app/entity"
	"github.com/gitwub5/go_todo_app/store"
	"github.com/go-playground/validator/v10"
)

/*
	vaildator 패키지는 Unmarshal한 타입에 validate라는 태그를 사용해 해당 필드의 검증 조건을 설정할 수 있다.
	설정한 조건은 *validator.Validate 메서드를 사용해 검증할 수 있다. 정의 완료된 조건도 있지만, 커스텀 조건도 정의할 수 있다.
*/

// AddTask 구조체는 HTTP 요청을 처리하고 새로운 Task를 저장하는 핸들러이다.
type AddTask struct {
	Store     *store.TaskStore    // TaskStore는 Task 엔티티를 저장하기 위한 인터페이스이다.
	Validator *validator.Validate // validator.Validate 구조체를 사용하여 유효성 검사를 수행한다.
}

// ServeHTTP는 HTTP 요청을 처리하는 메서드로, AddTask 핸들러의 엔트리 포인트이다.
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
	err := validator.New().Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Task를 생성하고 저장소에 추가한다.
	t := &entity.Task{
		Title:   b.Title,
		Status:  "todo",
		Created: time.Now(),
	}
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 성공 시 생성된 Task의 ID를 JSON 응답으로 반환한다.
	rsp := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
