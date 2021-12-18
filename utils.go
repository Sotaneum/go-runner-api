package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ginsession "github.com/go-session/gin-session"
	"github.com/google/uuid"

	gocreatefolder "github.com/Sotaneum/go-create-folder"
	file "github.com/Sotaneum/go-json-file"
	"github.com/Sotaneum/go-logger"
	"github.com/Sotaneum/go-runner"
	"github.com/gin-gonic/gin"

	requestjob "github.com/Sotaneum/go-request-job"
)

func New(options map[string]string) *Handler {
	handler := new(Handler)
	handler.config = options
	handler.RunnerChan = make(chan []runner.RunnerInterface)
	handler.active = true

	return handler
}

func decodeAuthorization(c *gin.Context) (string, string, error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return "", "", ErrorAuthorization
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return "", "", ErrorBearer
	}
	decodeAuth, _ := base64.RawStdEncoding.DecodeString(token)
	parse := strings.Split(string(decodeAuth), ":")
	if len(parse) != 2 {
		return "", "", ErrorBearer
	}
	return parse[0], parse[1], nil
}

func getUserID(c *gin.Context, path string) (string, bool) {
	store := ginsession.FromContext(c)
	data, ok := store.Get("userID")
	if ok {
		return fmt.Sprintf("%s", data), ok
	}

	userId, hook, parseErr := decodeAuthorization(c)
	if parseErr != nil {
		return "", false
	}

	uHook, fetchErr := getHook(userId, path)
	if fetchErr != nil {
		return "", false
	}

	if hook != uHook {
		return "", false
	}

	return userId, true
}

func getJobID(c *gin.Context) (string, bool) {
	ID := c.Request.URL.Query().Get("id")
	return ID, ID != ""
}

func getJobFile(path, id, userID string) (*requestjob.RequestJob, error) {
	jobObj, err := requestjob.NewByFile(path+"/job", id+".json", "")

	if err != nil {
		return nil, err
	}

	if !jobObj.HasAuthorization(userID) {
		return nil, requestjob.ErrorNoAuthorization
	}

	return jobObj, nil
}

func getJob(c *gin.Context, userID string) (*requestjob.RequestJob, error) {

	var data interface{}
	err := c.ShouldBindJSON(&data)

	if err != nil {
		return nil, requestjob.ErrorCantCreateJob
	}

	parse, parseErr := json.Marshal(data)

	if parseErr != nil {
		return nil, requestjob.ErrorCantCreateJob
	}

	jobObj, jobErr := requestjob.NewByJSON(string(parse), "")

	if jobErr != nil {
		return nil, requestjob.ErrorCantCreateJob
	}

	if !jobObj.HasAuthorization(userID) {
		return nil, requestjob.ErrorNoAuthorization
	}

	return jobObj, nil
}

func getJobList(path, userID string) ([]*requestjob.RequestJob, error) {
	jobList, err := requestjob.NewList(path + "/job")
	if err != nil {
		return nil, err
	}

	authJobList := []*requestjob.RequestJob{}

	for _, jobObj := range jobList {
		if jobObj.HasAuthorization(userID) {
			authJobList = append(authJobList, jobObj)
		}
	}

	return authJobList, nil
}

func getActive(c *gin.Context) bool {
	ID := c.Request.URL.Query().Get("value")
	return ID == "true"
}

func createHook(id, path string) error {
	createFolderErr := gocreatefolder.CreateFolder(path, 0755)

	if createFolderErr != nil {
		return createFolderErr
	}

	f := file.File{Path: path, Name: id + ".json"}
	data := f.Load()

	userConfig := UserConfig{}

	if data == "" {
		json.Unmarshal([]byte("{}"), &userConfig)
	} else {
		json.Unmarshal([]byte(data), &userConfig)
	}

	userConfig.Hook = uuid.New().String()

	f.SaveObject(userConfig)

	return nil
}

func getHook(id, path string) (string, error) {
	createFolderErr := gocreatefolder.CreateFolder(path, 0755)

	if createFolderErr != nil {
		return "", createFolderErr
	}

	f := file.File{Path: path, Name: id + ".json"}

	data := f.Load()
	if data == "" {
		createHook(id, path)
		return getHook(id, path)
	}

	userConfig := UserConfig{}
	json.Unmarshal([]byte(data), &userConfig)

	if userConfig.Hook == "" {
		createHook(id, path)
		return getHook(id, path)
	}

	return userConfig.Hook, nil
}

func FetchJob(path string, runnerChan chan []runner.RunnerInterface) {
	jobList, err := requestjob.NewList(path + "/job")

	if err != nil {
		panic(err)
	}

	runnerList := []runner.RunnerInterface{}

	for _, jobObj := range jobList {
		runnerList = append(runnerList, jobObj)
	}

	runnerChan <- runnerList
}

func LogPrint(path string, run *runner.Runner) {
	for {
		res := <-run.ResultCh
		if len(res) > 0 {
			for id, resData := range res {
				logger.New(path+"/log", id+".json").Add(resData)
			}
		}
	}
}

func Update(fileSystem FileSystem) {
	target := 10 * 60
	term := 1
	for {
		if term >= target {
			fileSystem.Push()
			fmt.Println("저장되었습니다.", term, target)
			term = 0
		}
		term = term + 1
		time.Sleep(time.Second * 1)
	}
}
