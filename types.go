package handler

import (
	runner "github.com/Sotaneum/go-runner"
	runnerjob "github.com/Sotaneum/go-runner-job"
)

type AuthInterface interface {
	GetLoginRedirectURL() string
	GetUser(code string, user interface{}) error
}

type JobControlInterface interface {
	NewList(path string) ([]runnerjob.BaseJobInterface, error)
	NewByJSON(data, owner string) (runnerjob.BaseJobInterface, error)
	NewByFile(path, name, owner string) (runnerjob.BaseJobInterface, error)
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
	Owner  []runnerjob.BaseJobInterface `json:"owner"`
	Editor []runnerjob.BaseJobInterface `json:"editor"`
	Admin  []runnerjob.BaseJobInterface `json:"admin"`
}
