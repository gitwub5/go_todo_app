package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gitwub5/go_todo_app/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminated server: %v", err)
		os.Exit(1)
	}
}

// (-) func run (ctx context.Context, l net.Listenr) error {
func run(ctx context.Context) error {
	cfg, err := config.New() // 환경 변수를 파싱하여 Config 구조체를 생성한다.
	if err != nil {
		return err
	}
	// Config 구조체의 포트 번호로 네트워크 리스너를 생성한다.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}
