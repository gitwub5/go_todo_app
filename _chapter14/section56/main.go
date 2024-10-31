// 포트 번호 변경할 수 있도록 만들기
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 { // os.Args는 명령행 인수를 저장하는 슬라이스이다. 명령행 인수가 2개가 아니면 종료한다.
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]                    // 명령행 인수로 받은 포트 번호를 사용한다.
	l, err := net.Listen("tcp", ":"+p) // net.Listen() 함수로 네트워크 연결을 생성한다.
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

// 동적으로 포트 번호를 할당할 수 있도록 변경한다.
func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// 인수로 받은 net.Listener를 이용하므로 Addr 필드는 지정하지 않는다.
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// ListenAndServe 메서드가 아닌 Serve 메서드로 변경한다.
		if err := s.Serve(l); err != nil && // Serve 메서드로 서버를 기동한다.
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
