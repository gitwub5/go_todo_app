package handler

import (
	"testing"
)

/*
	Go에서는 여러 개의 입력 및 기대값을 조합해서 공통화된 실행 순서로 테스트하는 패턴을 TDD(Test-Driven Development)라고 한다.

	테스트 입력값이나 기댓값을 파일에 저장하고 불러와서 테스트하는 방법을 Golden 테스트라고 한다.
	*json.golden이라는 파일명을 사용하여 테스트 입력값이나 기댓값을 저장한다.
*/

func TestAddTask(t *testing.T) {
	// type want struct {
	// 	status  int
	// 	rspFile string
	// }
	// // 테스트 케이스 정의
	// tests := map[string]struct {
	// 	reqFile string
	// 	want    want
	// }{
	// 	"ok": {
	// 		reqFile: "testdata/add_task/ok_req.json.golden",
	// 		want: want{
	// 			status:  http.StatusOK,
	// 			rspFile: "testdata/add_task/ok_rsp.json.golden",
	// 		},
	// 	},
	// 	"badRequest": {
	// 		reqFile: "testdata/add_task/bad_req.json.golden",
	// 		want: want{
	// 			status:  http.StatusBadRequest,
	// 			rspFile: "testdata/add_task/bad_rsp.json.golden",
	// 		},
	// 	},
	// }
	// for n, tt := range tests {
	// 	tt := tt
	// 	t.Run(n, func(t *testing.T) {
	// 		t.Parallel() // 테스트를 병렬로 실행한다.

	// 		w := httptest.NewRecorder()
	// 		r := httptest.NewRequest(
	// 			http.MethodPost,
	// 			"/tasks",
	// 			bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
	// 		)

	// 		sut := AddTask{
	// 			Store: &store.TaskStore{
	// 				Tasks: map[entity.TaskID]*entity.Task{},
	// 			},
	// 			Validator: validator.New(),
	// 		}
	// 		sut.ServeHTTP(w, r)

	// 		resp := w.Result()
	// 		testutil.AssertResponse(t,
	// 			resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
	// 		)
	// 	})
	// }
}
