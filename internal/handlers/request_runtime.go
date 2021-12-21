package handlers

import (
	"fmt"
	"github.com/brenordv/go-request/internal/core"
	"github.com/brenordv/go-request/internal/db"
	"github.com/brenordv/go-request/internal/models"
	"github.com/brenordv/go-request/internal/parsers"
	"github.com/brenordv/go-request/internal/utils"
	"github.com/schollz/progressbar/v3"
	"os"
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
	var err error
	runtimeConfig := utils.LoadRuntimeConfig(runtimeConfigFile)
	PanicOnError(runtimeConfig.Get.Validate())
	flowControl := models.FlowControl{
		GuardChannel: make(chan struct{}, runtimeConfig.Get.MaxParallelRequests),
	}

	flowControl.GenerateSessionId()
	fmt.Printf("Your session id is: %s\n", flowControl.SessionId)

	flowControl.Request, err = runtimeConfig.Get.MakeRequest()
	PanicOnError(err)

	flowControl.ProgressBar = progressbar.Default(int64(runtimeConfig.Get.NumRequests), "Making GET requests")

	flowControl.Db, err = db.NewDatabaseClient()
	PanicOnError(err)

	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(flowControl.Db)
	flowControl.Db.WaitForPromises = true

	for i := 0; i < runtimeConfig.Get.NumRequests; i++ {
		flowControl.GuardChannel <- struct{}{}
		flowControl.WaitGroup.Add(1)
		go func(fc *models.FlowControl, index int, rtCfg *models.RuntimeConfig) {
			defer IgnoreError(fc.ProgressBar.Add(1))
			defer ReadIntoVoid(fc.GuardChannel)

			client := rtCfg.Get.GetHttpClient()
			defer client.CloseIdleConnections()

			if rtCfg.Get.AggressiveMode {
				fc.WaitGroup.Done()
			} else {
				defer fc.WaitGroup.Done()
			}

			fc.Db.PromiseWillAdd(1)

			url := fc.Request.URL.String()
			httpRes, err := client.Do(fc.Request)
			res := models.NewHttpResponse(url, fc.SessionId, index, httpRes, err)

			key, serialized, err := res.Serialize()
			PanicOnError(err)

			err = fc.Db.Add(key, serialized)
			PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(&flowControl, i, runtimeConfig)
	}
	flowControl.WaitGroup.Wait()
}

func postRequestRoutine(runtimeConfigFile string) {
	var err error
	runtimeConfig := utils.LoadRuntimeConfig(runtimeConfigFile)
	PanicOnError(runtimeConfig.Post.Validate())
	flowControl := models.FlowControl{
		GuardChannel: make(chan struct{}, runtimeConfig.Post.MaxParallelRequests),
	}

	flowControl.GenerateSessionId()
	fmt.Printf("Your session id is: %s\n", flowControl.SessionId)

	flowControl.Request, err = runtimeConfig.Post.MakeRequest()
	PanicOnError(err)

	flowControl.ProgressBar = progressbar.Default(int64(runtimeConfig.Post.NumRequests), "Making POST requests")
	flowControl.Db, err = db.NewDatabaseClient()
	PanicOnError(err)

	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(flowControl.Db)
	flowControl.Db.WaitForPromises = true

	for i := 0; i < runtimeConfig.Post.NumRequests; i++ {
		flowControl.GuardChannel <- struct{}{}
		flowControl.WaitGroup.Add(1)
		go func(fc *models.FlowControl, index int, rtCfg *models.RuntimeConfig) {
			defer IgnoreError(fc.ProgressBar.Add(1))
			defer ReadIntoVoid(fc.GuardChannel)

			client := rtCfg.Post.GetHttpClient()
			defer client.CloseIdleConnections()

			if rtCfg.Post.AggressiveMode {
				fc.WaitGroup.Done()
			} else {
				defer fc.WaitGroup.Done()
			}

			fc.Db.PromiseWillAdd(1)

			url := fc.Request.URL.String()
			httpRes, err := client.Do(fc.Request)
			res := models.NewHttpResponse(url, fc.SessionId, index, httpRes, err)

			key, serialized, err := res.Serialize()
			PanicOnError(err)

			err = fc.Db.Add(key, serialized)
			PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(&flowControl, i, runtimeConfig)
	}
	flowControl.WaitGroup.Wait()
}
