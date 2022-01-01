package handler

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

// AuthCallback : 사용자 인증을 해서 세션에 저장합니다.
func (h *Handler) AuthCallback(c *gin.Context) {
	query := c.Request.URL.Query()
	var user User
	h.auth.GetUser(query.Get("code"), &user)
	store := ginsession.FromContext(c)
	store.Set("userID", user.ID)
	err := store.Save()
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.Redirect(302, "/")
}

// Login : 로그인 페이지로 이동시킵니다.
func (h *Handler) Login(c *gin.Context) {
	c.Redirect(302, h.auth.GetLoginRedirectURL())
}
