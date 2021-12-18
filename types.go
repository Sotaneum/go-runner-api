package handler

import (
	runner "github.com/Sotaneum/go-runner"
)

type FileSystem interface {
	Pull() error
	Push() error
}

type SSO interface {
	GetLoginRedirectURL() string
	GetUser(code string, user interface{}) error
}

type Handler struct {
	FileSystem FileSystem
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
