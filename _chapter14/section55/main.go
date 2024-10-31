// 리펙토링과 테스트 코드
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
	// golang.org/x/sync에 있는 errgroup 패키지를 임포트한다.
	// errgroup.Group 타입을 사용하면, 반환값에 오류가 포함되는 고루틴의 병렬 처리를 간단히 구현할 수 있다.
)

func main() {
	// run 함수로 처리를 분리한다.
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

// run 함수는 HTTP 서버를 기동하고 종료하는 역할을 수행한다.
// context.Context 타입값을 인수로 받아서 외부에서 취소 처리를 받으면 서버를 종료한다.
func run(ctx context.Context) error {
	// http.Server 타입값을 생성하고, http.ListenAndServe() 메서드로 서버를 기동한다.
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	// errgroup.WithContext() 함수로 에러 그룹을 생성하고, 고루틴을 생성한다.
	eg, ctx := errgroup.WithContext(ctx)

	// 다른 고루틴에서 HTTP 서버를 기동한다.
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil &&
			// http.ErrServerClosed 에러가 발생하면 정상 종료로 간주한다.
			// http.Server.Shutdown()가 정상 종료된 것을 나타내므로 이상 처리가 아니다.
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})
	// 채널로부터의 알림(종료 알림)을 기다렸다가 요청이 오면 서버를 종료한다.
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil { // http.Server.Shutdown() 메서드로 서버를 종료할 수 있다.
		log.Printf("failed to shutdown: %+v", err)
	}
	// errgroup.Wait() 메서드로 모든 고루틴의 종료를 기다린다.
	return eg.Wait()
}
