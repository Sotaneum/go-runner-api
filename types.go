package handler

import (
	requestjob "github.com/Sotaneum/go-request-job"
	runner "github.com/Sotaneum/go-runner"
)

type SSO interface {
	GetLoginRedirectURL() string
	GetUser(code string, user interface{}) error
}

type Handler struct {
	SSO        SSO
	config     map[string]string
	active     bool
	RunnerChan chan []runner.RunnerInterface
}

// User : 사용자 정보
type User struct {
	ID string `json:"accountId"`
}

type UserConfig struct {
	Hook string `json:"hook"`
}

type ResponseJobList struct {
	Owner  []*requestjob.RequestJob `json:"owner"`
	Editor []*requestjob.RequestJob `json:"editor"`
	Admin  []*requestjob.RequestJob `json:"admin"`
}
