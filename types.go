package handler

import (
	runner "github.com/Sotaneum/go-runner"
)

type AuthInterface interface {
	GetLoginRedirectURL() string
	GetUser(code string, user interface{}) error
}

type JobControlInterface interface {
	NewList(path string) (interface{}, error)
	NewByJSON(data, owner string) (interface{}, error)
	NewByFile(path, name, owner string) (interface{}, error)
}

type Handler struct {
	auth       AuthInterface
	jobControl JobControlInterface
	config     map[string]string
	active     bool
	runnerChan chan []runner.JobInterface
}

// User : 사용자 정보
type User struct {
	ID string `json:"accountId"`
}

type UserConfig struct {
	Hook string `json:"hook"`
}

type ResponseJobList struct {
	Owner  []interface{} `json:"owner"`
	Editor []interface{} `json:"editor"`
	Admin  []interface{} `json:"admin"`
}
