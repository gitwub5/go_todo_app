package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// main_test.go를 server_test.go로 변경
func TestServer_Run(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0") // net.Listen 함수를 사용하여 TCP 서버를 생성한다.
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // mux 변수를 http.HandlerFunc로 변경
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:]) // fmt.Fprintf 함수를 사용하여 응답을 작성한다.
	})

	eg.Go(func() error {
		s := NewServer(l, mux) // NewServer 함수를 사용하여 서버를 생성한다.
		return s.Run(ctx)      // Run 메서드를 사용하여 서버를 실행한다.
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)

	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
