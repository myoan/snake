package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func doHealth(ctx context.Context) {
	t := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done")
			return
		case <-t:
			fmt.Println("tick")
		}
	}
}

func main() {
	fmt.Println("Start main")
	ctx := context.Background()
	go doHealth(ctx)

	http.HandleFunc("/hoge", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	// <-ctx.Done()
}
