package main

import (
	"fmt"
	"github.com/brenordv/go-request/internal/db"
	"github.com/brenordv/go-request/internal/handlers"
	"github.com/brenordv/go-request/internal/models"
	"github.com/brenordv/go-request/internal/utils"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"net/http"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("go-Request!::GET")
	runtimeConfig := utils.LoadRuntimeConfig("./.configs/get.equip-unavail-qa.config.json")
	handlers.PanicOnError(runtimeConfig.Get.Validate())
	guard := make(chan struct{}, runtimeConfig.Get.MaxParallelRequests)

	u := uuid.New()
	sessionId := u.String()
	fmt.Printf("Your session id is: %s\n", sessionId)

	req, err := runtimeConfig.Get.MakeRequest()

	handlers.PanicOnError(err)

	var wg sync.WaitGroup
	bar := progressbar.Default(int64(runtimeConfig.Get.NumRequests), "Making GET requests")
	dbClient, err := db.NewDatabaseClient()
	handlers.PanicOnError(err)
	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(dbClient)
	dbClient.WaitForPromises = true

	for i := 0; i < runtimeConfig.Get.NumRequests; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(b *progressbar.ProgressBar, r *http.Request, w *sync.WaitGroup, index int, rtCfg *models.RuntimeConfig, g chan struct{}, dbc *db.DatabaseClient, sId string) {
			defer handlers.IgnoreError(b.Add(1))
			defer handlers.ReadIntoVoid(g)

			client := rtCfg.Get.GetHttpClient()
			defer client.CloseIdleConnections()

			if rtCfg.Get.AggressiveMode {
				w.Done()
			} else {
				defer w.Done()
			}

			dbc.PromiseWillAdd(1)

			url := r.URL.String()
			httpRes, err := client.Do(r)
			res := models.NewHttpResponse(url, sessionId, index, httpRes, err)

			key, serialized, err := res.Serialize()
			handlers.PanicOnError(err)

			err = dbc.Add(key, serialized)
			handlers.PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(bar, req, &wg, i, runtimeConfig, guard, dbClient, sessionId)
	}
	wg.Wait()
	fmt.Printf("Done! Elapsed time: %s\n", time.Since(start))
}
