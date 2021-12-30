package handler

import (
	"github.com/Sotaneum/go-logger"
	"github.com/gin-gonic/gin"
)

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

	jobList, fetchJobListErr := getJobList(h.config["path"], userID, userID == h.config["adminId"])

	responseData := ResponseJobList{}

	for _, job := range jobList {
		if job.HasAdminAuthorization(userID) {
			responseData.Owner = append(responseData.Owner, job)
			continue
		}
		if job.HasAuthorization(userID) {
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

	jobObj, jobObjErr := getJobFile(h.config["path"], jobID, userID)

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

	jobObj, jobErr := getJob(c, userID)

	if jobErr != nil {
		ResposeParamsError(c)
		return
	}

	jobObj.Save(h.config["path"] + "/job")

	go fetchJob(h.config["path"], h.RunnerChan)

	ResponseData(c, jobObj.GetID())
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

	jobObj, jobObjErr := getJobFile(h.config["path"], jobID, userID)

	if jobObjErr != nil || !jobObj.HasAdminAuthorization(userID) {
		ResponseNoAuthorization(c)
		return
	}

	err := jobObj.Remove(h.config["path"] + "/job")

	if err != nil {
		ResponseCantRemoveJob(c)
		return
	}

	logger.New(h.config["path"]+"/log", jobID+".json").Remove()

	go fetchJob(h.config["path"], h.RunnerChan)

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
