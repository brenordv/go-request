package handlers

import (
	"fmt"
	"github.com/brenordv/go-request/internal/core"
	"github.com/brenordv/go-request/internal/db"
	"github.com/brenordv/go-request/internal/models"
	"github.com/brenordv/go-request/internal/parsers"
	"github.com/brenordv/go-request/internal/utils"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"net/http"
	"os"
	"sync"
)

func ExecRequests(appName string, method string) {
	runtimeConfigs := parsers.GetConfigFilenames()
	runtimeConfigsCount := len(runtimeConfigs)

	if runtimeConfigsCount == 0 && len(os.Args) == 1 {
		fmt.Print("You must pass at least one config file in the command line.")
		fmt.Printf("Example: %s ./my-get-requests.json\n", appName)
		return
	}
	if runtimeConfigsCount == 0 {
		fmt.Print("No valid runtime config files found.")
		fmt.Print("Try passing one in the command line.")
		fmt.Printf("Example: %s ./my-get-requests.json\n", appName)
		return
	}

	for _, routineConfig := range runtimeConfigs {
		if method == core.HttpGet {
			getRequestRoutine(routineConfig)

		} else if method == core.HttpPost {
			postRequestRoutine(routineConfig)

		} else {
			fmt.Printf("Method '%s' is not supported. Skipping file...\n", method)
		}
	}
}

func getRequestRoutine(runtimeConfigFile string) {
	runtimeConfig := utils.LoadRuntimeConfig(runtimeConfigFile)
	PanicOnError(runtimeConfig.Get.Validate())
	guard := make(chan struct{}, runtimeConfig.Get.MaxParallelRequests)

	u := uuid.New()
	sessionId := u.String()
	fmt.Printf("Your session id is: %s\n", sessionId)

	req, err := runtimeConfig.Get.MakeRequest()

	PanicOnError(err)

	var wg sync.WaitGroup
	bar := progressbar.Default(int64(runtimeConfig.Get.NumRequests), "Making GET requests")
	dbClient, err := db.NewDatabaseClient()
	PanicOnError(err)
	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(dbClient)
	dbClient.WaitForPromises = true

	for i := 0; i < runtimeConfig.Get.NumRequests; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(b *progressbar.ProgressBar, r *http.Request, w *sync.WaitGroup, index int, rtCfg *models.RuntimeConfig, g chan struct{}, dbc *db.DatabaseClient, sId string) {
			defer IgnoreError(b.Add(1))
			defer ReadIntoVoid(g)

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
			PanicOnError(err)

			err = dbc.Add(key, serialized)
			PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(bar, req, &wg, i, runtimeConfig, guard, dbClient, sessionId)
	}
	wg.Wait()
}

func postRequestRoutine(runtimeConfigFile string) {
	runtimeConfig := utils.LoadRuntimeConfig(runtimeConfigFile)
	PanicOnError(runtimeConfig.Post.Validate())
	guard := make(chan struct{}, runtimeConfig.Post.MaxParallelRequests)

	u := uuid.New()
	sessionId := u.String()
	fmt.Printf("Your session id is: %s\n", sessionId)

	var wg sync.WaitGroup
	bar := progressbar.Default(int64(runtimeConfig.Post.NumRequests), "Making POST requests")
	dbClient, err := db.NewDatabaseClient()
	PanicOnError(err)
	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(dbClient)
	dbClient.WaitForPromises = true

	for i := 0; i < runtimeConfig.Post.NumRequests; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(b *progressbar.ProgressBar, w *sync.WaitGroup, index int, rtCfg *models.RuntimeConfig, g chan struct{}, dbc *db.DatabaseClient, sId string) {
			defer IgnoreError(b.Add(1))
			defer ReadIntoVoid(g)

			client := rtCfg.Post.GetHttpClient()
			defer client.CloseIdleConnections()

			req, err := runtimeConfig.Post.MakeRequest()
			PanicOnError(err)

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
			PanicOnError(err)

			err = dbc.Add(key, serialized)
			PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(bar, &wg, i, runtimeConfig, guard, dbClient, sessionId)
	}
	wg.Wait()
}
