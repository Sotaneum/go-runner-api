package handler

import (
	"encoding/json"

	"github.com/Sotaneum/go-logger"
	runner "github.com/Sotaneum/go-runner"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Initialize(options map[string]string, jobControl *JobControlInterface, auth *AuthInterface) *Handler {
	if jobControl == nil || auth == nil {
		return nil
	}

	h.config = options
	h.active = true
	h.jobControl = *jobControl
	h.auth = *auth
	h.runnerChan = make(chan []runner.RunnerInterface)

	go h.fetchJob()

	if h.config["logDisable"] != "true" {
		go logPrint(h.config["path"], runner.NewRunner(h.runnerChan))
	}

	return h
}

func (h *Handler) GetJobList(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	userID, hasUserID := getUserID(c, h.config["path"]+"/user")

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	jobList, fetchJobListErr := h.getJobList(userID)

	responseData := ResponseJobList{}

	for _, job := range jobList {
		if (*job).HasAdminAuthorization(userID) {
			responseData.Owner = append(responseData.Owner, job)
			continue
		}
		if (*job).HasAuthorization(userID) {
			responseData.Editor = append(responseData.Editor, job)
			continue
		}
		responseData.Admin = append(responseData.Admin, job)
	}

	if fetchJobListErr != nil {
		ResposeServerError(c)
		return
	}

	ResponseData(c, responseData)
}

func (h *Handler) GetLogs(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	_, hasUserID := getUserID(c, h.config["path"]+"/user")

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	jobID, hasJobID := getJobID(c)

	if !hasJobID {
		ResposeParamsError(c)
		return
	}

	ResponseData(c, logger.New(h.config["path"]+"/log", jobID+".json").Get())
}

func (h *Handler) GetJob(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	userID, hasUserID := getUserID(c, h.config["path"]+"/user")

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	jobID, hasJobID := getJobID(c)

	if !hasJobID {
		ResposeParamsError(c)
		return
	}

	jobObj, jobObjErr := h.getJobFile(jobID, userID)

	if jobObjErr != nil {
		ResponseNoAuthorization(c)
		return
	}

	ResponseData(c, jobObj)
}

func (h *Handler) UpdateJob(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	userID, hasUserID := getUserID(c, h.config["path"]+"/user")

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	jobObj, jobErr := h.getJob(c, userID)

	if jobErr != nil {
		ResposeParamsError(c)
		return
	}

	(*jobObj).Save(h.config["path"] + "/job")

	go h.fetchJob()

	ResponseData(c, (*jobObj).GetID())
}

func (h *Handler) DeleteJob(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	userID, hasUserID := getUserID(c, h.config["path"]+"/user")

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	jobID, hasJobID := getJobID(c)

	if !hasJobID {
		ResposeParamsError(c)
		return
	}

	jobObj, jobObjErr := h.getJobFile(jobID, userID)

	if jobObjErr != nil || !(*jobObj).HasAdminAuthorization(userID) {
		ResponseNoAuthorization(c)
		return
	}

	err := (*jobObj).Remove(h.config["path"] + "/job")

	if err != nil {
		ResponseCantRemoveJob(c)
		return
	}

	logger.New(h.config["path"]+"/log", jobID+".json").Remove()

	go h.fetchJob()

	ResponseCompleteRemoveJob(c)
}

func (h *Handler) Active(c *gin.Context) {
	userID, hasUserID := getUserID(c, h.config["path"]+"/user")

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	if userID != h.config["adminId"] {
		ResponseNoAuthorization(c)
		return
	}

	active := getActive(c)

	h.active = active

	ResponseData(c, "ok")
}

func (h *Handler) GetHookID(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	userPath := h.config["path"] + "/user"

	userID, hasUserID := getUserID(c, userPath)

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	hook, fetchHookErr := getHook(userID, userPath)

	if fetchHookErr != nil {
		ResposeParamsError(c)
		return
	}

	ResponseData(c, hook)
}

func (h *Handler) ReHookID(c *gin.Context) {
	if !h.active {
		ResponseMaintenance(c)
		return
	}

	userPath := h.config["path"] + "/user"

	userID, hasUserID := getUserID(c, userPath)

	if !hasUserID {
		ResposeNeedLogin(c)
		return
	}

	createHookErr := createHook(userID, userPath)

	if createHookErr != nil {
		ResposeParamsError(c)
		return
	}

	hook, fetchHookErr := getHook(userID, userPath)

	if fetchHookErr != nil {
		ResposeParamsError(c)
		return
	}

	ResponseData(c, hook)
}

func (h *Handler) loadJobList() ([]*JobInterface, error) {
	return h.jobControl.NewList(h.config["path"] + "/job")
}

func (h *Handler) loadJobJSON(data string, owner string) (*JobInterface, error) {
	return h.jobControl.NewByJSON(data, owner)
}

func (h *Handler) loadJobFile(id string) (*JobInterface, error) {
	return h.jobControl.NewByFile(h.config["path"]+"/job", id+".json", "")
}

func (h *Handler) fetchJob() {
	jobList, err := h.loadJobList()

	if err != nil {
		panic(err)
	}

	runnerList := []runner.RunnerInterface{}

	for _, jobObj := range jobList {
		runnerList = append(runnerList, *jobObj)
	}

	h.runnerChan <- runnerList
}

func (h *Handler) getJobFile(id, userID string) (*JobInterface, error) {
	jobObj, err := h.loadJobFile(id)

	if err != nil {
		return nil, err
	}

	if !(*jobObj).HasAuthorization(userID) {
		return nil, ErrorAuthorization
	}

	return jobObj, nil
}

func (h *Handler) getJob(c *gin.Context, userID string) (*JobInterface, error) {

	var data interface{}
	err := c.ShouldBindJSON(&data)

	if err != nil {
		return nil, ErrorJob
	}

	parse, parseErr := json.Marshal(data)

	if parseErr != nil {
		return nil, ErrorJob
	}

	jobObj, jobErr := h.loadJobJSON(string(parse), "")

	if jobErr != nil {
		return nil, ErrorJob
	}

	if !(*jobObj).HasAuthorization(userID) {
		return nil, ErrorAuthorization
	}

	return jobObj, nil
}

func (h *Handler) getJobList(userID string) ([]*JobInterface, error) {
	jobList, err := h.loadJobList()
	if err != nil {
		return nil, err
	}

	authJobList := []*JobInterface{}

	force := userID == h.config["adminId"]

	for _, jobObj := range jobList {
		if force || (*jobObj).HasAuthorization(userID) {
			authJobList = append(authJobList, jobObj)
		}
	}

	return authJobList, nil
}
