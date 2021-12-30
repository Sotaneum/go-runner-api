package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	ginsession "github.com/go-session/gin-session"
	"github.com/google/uuid"

	gocreatefolder "github.com/Sotaneum/go-create-folder"
	file "github.com/Sotaneum/go-json-file"
	"github.com/Sotaneum/go-logger"
	"github.com/Sotaneum/go-runner"
	"github.com/gin-gonic/gin"
)

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

func logPrint(path string, run *runner.Runner) {
	for {
		res := <-run.ResultCh
		if len(res) > 0 {
			for id, resData := range res {
				logger.New(path+"/log", id+".json").Add(resData)
			}
		}
	}
}
