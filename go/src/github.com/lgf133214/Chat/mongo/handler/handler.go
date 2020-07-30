package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lgf133214/Chat/mongo/db"
	"github.com/lgf133214/Chat/mongo/model"
	"net/http"
	"strconv"
	"time"
)

func GetMsg(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	t, ok := c.GetQuery("last_time")
	if !ok {
		return
	}
	lastTime, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return
	}

	ret, ok := db.GetMsgs(lastTime)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": ret})
}

func StoreMsg(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	name, ok := c.GetPostForm("name")
	if !ok {
		return
	}
	msg, ok := c.GetPostForm("msg")
	if !ok {
		return
	}

	ret := model.Msg{Content: msg, UserName: name, Time: time.Now().Unix()}
	ok = db.StoreMsg(ret)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    "发送成功",
		})
	}
}
