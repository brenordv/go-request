package models

type RuntimeConfig struct {
	Get HttpConfig `json:"get"`
	Post HttpConfig `json:"post"`
}
