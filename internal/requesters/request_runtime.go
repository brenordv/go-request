package requesters

import (
	"fmt"
	"github.com/brenordv/go-request/internal/core"
	"github.com/brenordv/go-request/internal/db"
	"github.com/brenordv/go-request/internal/file"
	"github.com/brenordv/go-request/internal/handlers"
	"github.com/brenordv/go-request/internal/models"
	"github.com/brenordv/go-request/internal/parsers"
	"github.com/schollz/progressbar/v3"
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
		runtimeConfig := file.LoadRuntimeConfig(routineConfig)
		if method == core.HttpGet {
			doRequest(&runtimeConfig.Get, core.HttpGet)

		} else if method == core.HttpPost {
			doRequest(&runtimeConfig.Post, core.HttpPost)

		} else {
			fmt.Printf("Method '%s' is not supported. Skipping file...\n", method)
		}
	}
}

func doRequest(runtimeCfg *models.HttpConfig, method string) {
	var err error
	handlers.PanicOnError(runtimeCfg.Validate())
	var wg sync.WaitGroup
	flowControl := models.FlowControl{
		WaitGroup: &wg,
		GuardChannel: make(chan struct{}, runtimeCfg.MaxParallelRequests),
	}

	flowControl.GenerateSessionId()
	fmt.Printf("Your session id is: %s\n", flowControl.SessionId)

	flowControl.Request, err = runtimeCfg.MakeRequest()
	handlers.PanicOnError(err)

	flowControl.ProgressBar = progressbar.Default(int64(runtimeCfg.NumRequests), fmt.Sprintf("Making %s requests", method))

	flowControl.Db, err = db.NewDatabaseClient()
	handlers.PanicOnError(err)

	defer func(dbClient *db.DatabaseClient) {
		_ = dbClient.Close()
	}(flowControl.Db)
	flowControl.Db.WaitForPromises = true

	for i := 0; i < runtimeCfg.NumRequests; i++ {
		flowControl.GuardChannel <- struct{}{}
		flowControl.WaitGroup.Add(1)
		go func(fc *models.FlowControl, index int, rtCfg *models.HttpConfig) {
			defer handlers.IgnoreError(fc.ProgressBar.Add(1))
			defer handlers.ReadIntoVoid(fc.GuardChannel)

			client := rtCfg.GetHttpClient()
			defer client.CloseIdleConnections()

			if rtCfg.AggressiveMode {
				fc.WaitGroup.Done()
			} else {
				defer fc.WaitGroup.Done()
			}

			fc.Db.PromiseWillAdd(1)

			url := fc.Request.URL.String()
			httpRes, err := client.Do(fc.Request)
			res := models.NewHttpResponse(url, fc.SessionId, index, httpRes, err)

			key, serialized, err := res.Serialize()
			handlers.PanicOnError(err)

			err = fc.Db.Add(key, serialized)
			handlers.PanicOnError(err)

			if res.IsInternalServerError {
				fmt.Println(res.String())
			}

		}(&flowControl, i, runtimeCfg)
	}
	flowControl.WaitGroup.Wait()
}
