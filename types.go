package handler

import (
	"time"

	runner "github.com/Sotaneum/go-runner"
)

type SSO interface {
	GetLoginRedirectURL() string
	GetUser(code string, user interface{}) error
}

type JobInterface interface {
	HasAuthorization(userID string) bool
	HasAdminAuthorization(userID string) bool
	IsRun(t time.Time) bool
	GetID() string
	Run() interface{}
	Remove(path string) error
	Save(path string)
}

type JobControlInterface interface {
	NewList(path string) ([]*JobInterface, error)
	NewByJSON(data, owner string) (*JobInterface, error)
	NewByFile(path, name, owner string) (*JobInterface, error)
}

type Handler struct {
	SSO        SSO
	config     map[string]string
	active     bool
	runnerChan chan []runner.RunnerInterface
	jobControl JobControlInterface
}

// User : 사용자 정보
type User struct {
	ID string `json:"accountId"`
}

type UserConfig struct {
	Hook string `json:"hook"`
}

type ResponseJobList struct {
	Owner  []*JobInterface `json:"owner"`
	Editor []*JobInterface `json:"editor"`
	Admin  []*JobInterface `json:"admin"`
}
