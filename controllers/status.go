package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RespMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func CheckStatus(c *gin.Context) {
	resp := &RespMsg{}
	resp.Code = 0
	resp.Msg = "OK aigcd 7: " + time.Now().String()
	c.JSON(200, resp)
	return

}
