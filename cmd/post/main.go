package main

import (
	"fmt"
	"github.com/brenordv/go-request/internal/core"
	"github.com/brenordv/go-request/internal/handlers"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("go-Request!::POST")
	defer fmt.Printf("Done! Elapsed time: %s\n", time.Since(start))

	handlers.ExecRequests("go-post.exe", core.HttpPost)
}
