package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // context.WithCancel() 함수로 취소 처리를 위한 context.Context 타입값을 생성한다.
	eg, ctx := errgroup.WithContext(ctx)                    // errgroup.WithContext() 함수로 에러 그룹을 생성하고, 고루틴을 생성한다.
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTP 서버의 반환값을 검증한다.
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	// run 함수에 종료 알림을 전송한다.
	cancel()
	// run 함수의 반환값을 검증한다.
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}

/*
Go 테스트 방법:
1. 파일명이 _test.go로 끝나는 파일에 테스트 함수를 작성한다.
2. 테스트 함수는 Test로 시작하고, *testing.T 타입의 인수를 받는다.
3. 테스트 함수 내에서 *testing.T 타입의 메서드를 호출해서 테스트 결과를 출력한다.
4. 테스트 파일이 있는 디렉터리에서 go test 명령을 실행한다.

동작 확인 방법:
1. $ go test -v
=== RUN   TestRun
--- PASS: TestRun (0.00s)
PASS
ok      github.com/gitwub5/go_todo_app  0.145s
*/
