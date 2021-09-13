package main

import (
	"fmt"
	"github.com/brenordv/go-request/internal/db"
	"github.com/brenordv/go-request/internal/handlers"
	"github.com/brenordv/go-request/internal/models"
	"github.com/brenordv/go-request/internal/utils"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("go-Request!::POST")
	runtimeConfig := utils.LoadRuntimeConfig("./.configs/post.metrics-local.config.json")
	handlers.PanicOnError(runtimeConfig.Post.Validate())
	guard := make(chan struct{}, runtimeConfig.Post.MaxParallelRequests)

	u := uuid.New()
	sessionId := u.String()
	fmt.Printf("Your session id is: %s\n", sessionId)

	var wg sync.WaitGroup
	bar := progressbar.Default(int64(runtimeConfig.Post.NumRequests), "Making POST requests")
	dbClient, err := db.NewDatabaseClient()
	handlers.PanicOnError(err)
	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(dbClient)
	dbClient.WaitForPromises = true

	for i := 0; i < runtimeConfig.Post.NumRequests; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(b *progressbar.ProgressBar, w *sync.WaitGroup, index int, rtCfg *models.RuntimeConfig, g chan struct{}, dbc *db.DatabaseClient, sId string) {
			defer handlers.IgnoreError(b.Add(1))
			defer handlers.ReadIntoVoid(g)

			client := rtCfg.Post.GetHttpClient()
			defer client.CloseIdleConnections()

			req, err := runtimeConfig.Post.MakeRequest()
			handlers.PanicOnError(err)

			if rtCfg.Post.AggressiveMode {
				w.Done()
			} else {
				defer w.Done()
			}

			dbc.PromiseWillAdd(1)

			url := req.URL.String()
			httpRes, err := client.Do(req)
			res := models.NewHttpResponse(url, sessionId, index, httpRes, err)

			key, serialized, err := res.Serialize()
			handlers.PanicOnError(err)

			err = dbc.Add(key, serialized)
			handlers.PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(bar, &wg, i, runtimeConfig, guard, dbClient, sessionId)
	}
	wg.Wait()
	fmt.Printf("Done! Elapsed time: %s\n", time.Since(start))
}
