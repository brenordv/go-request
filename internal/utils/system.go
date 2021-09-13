package utils

import (
	"encoding/json"
	"github.com/brenordv/go-request/internal/core"
	"github.com/brenordv/go-request/internal/models"
	"os"
	"path"
)

func LoadRuntimeConfig(file string) *models.RuntimeConfig {
	var runtimeConfig models.RuntimeConfig
	f, _ := os.Open(file)
	jsonParser := json.NewDecoder(f)
	_ = jsonParser.Decode(&runtimeConfig)
	runtimeConfig.Get.RequestType = core.HttpGet
	runtimeConfig.Post.RequestType = core.HttpPost

	return &runtimeConfig
}

func GetAppDir() (string, error) {
	execPath, err := os.Executable()
	return path.Dir(execPath), err
}