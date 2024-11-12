package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gitwub5/go_todo_app/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminated server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New() // config.New 함수를 사용하여 설정을 초기화한다.
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port)) // net.Listen 함수를 사용하여 TCP 네트워크에서 주소를 사용하여 네트워크 리스너를 생성한다.
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	mux := NewMux()        // NewMux 함수를 사용하여 ServeHTTP 메서드를 구현한 http.Handler 인터페이스를 반환한다.
	s := NewServer(l, mux) // NewServer 함수를 사용하여 서버를 생성한다.
	return s.Run(ctx)      // Run 메서드를 사용하여 서버를 실행한다.
}
