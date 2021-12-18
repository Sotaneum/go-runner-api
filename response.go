package handler

import "github.com/gin-gonic/gin"

// ResposeNeedLogin : 로그인이 필요한 경우 발생합니다.
func ResposeNeedLogin(c *gin.Context) {
	c.JSON(401, gin.H{"code": 401, "message": "로그인이 필요한 서비스입니다."})
}

// ResposeServerError : 서버에 에러가 발생했을 경우 발생합니다.
func ResposeServerError(c *gin.Context) {
	c.JSON(500, gin.H{"code": 500, "message": "서버 초기화에 문제가 발생했습니다."})
}

// ResposeParamsError : 파라미터에 오류가 있을 경우 발생합니다.
func ResposeParamsError(c *gin.Context) {
	c.JSON(400, gin.H{"code": 400, "message": "잘못된 파라미터 입니다."})
}

// ResponseNotFoundPage : 찾을 수 없는 페이지일 경우 발생합니다.
func ResponseNotFoundPage(c *gin.Context) {
	c.JSON(404, gin.H{"code": 404, "message": "접근 할 수 없는 페이지입니다!"})
}

// ResponseNoAuthorization : 권한이 없는 사용자가 요청했을 경우 발생합니다.
func ResponseNoAuthorization(c *gin.Context) {
	c.JSON(401, gin.H{"code": 401, "message": "권한이 없거나 존재하지 않은 Job입니다."})
}

// ResponseCantRemoveJob : OS 문제로 인한 삭제가 불가능할 경우 발생합니다.
func ResponseCantRemoveJob(c *gin.Context) {
	c.JSON(500, gin.H{"code": 500, "message": "Job를 삭제할 수 없었습니다. 다시 시도해주세요."})
}

// ResponseCompleteRemoveJob : Job이 성공적으로 삭제되었을 경우 발생합니다.
func ResponseCompleteRemoveJob(c *gin.Context) {
	c.JSON(200, gin.H{"code": 200, "message": "Job이 삭제되었습니다."})
}

// ResponseData : 데이터를 전송해야하는 경우 발생합니다.
func ResponseData(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 200, "data": data})
}

// ResponseMaintenance : 시스템 점검 중일 경우 메시지를 반환합니다.
func ResponseMaintenance(c *gin.Context) {
	c.JSON(500, gin.H{"code": 500, "message": "시스템 점검 중으로 작업이 중단됩니다. \n시스템 점검이 길어지면 문의 부탁드립니다."})
}
