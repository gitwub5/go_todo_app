// 최소한의 코드만 사용한 웹 서버
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting server on :18080")

	err := http.ListenAndServe( // 첫 번째 인자로 포트 번호를 지정하고, 두 번째 인자로 핸들러를 지정한다.
		":18080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	)
	if err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

/*
동작 확인 방법:
1. $ go run main.go
Starting server on :18080
2. $ curl http://localhost:18080/World
Hello, World!
*/
