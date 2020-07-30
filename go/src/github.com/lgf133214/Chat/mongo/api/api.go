package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lgf133214/Chat/mongo/handler"
	"log"
)

func Run() {
	r:=gin.Default()
	{
		r.GET("/getMsg", handler.GetMsg)
		r.POST("/sendMsg", handler.StoreMsg)
	}
	err:=r.Run(":8848")
	if err != nil {
		log.Fatal(err)
	}
}
