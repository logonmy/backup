package middlerware

import (
	"github.com/gin-gonic/gin"
	"github.com/lgf133214/Gin/config"
	"github.com/lgf133214/Gin/db"
	"github.com/lgf133214/Gin/util"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

var (
	limiter = rate.NewLimiter(200, 800)
)

func AuthMiddleware(c *gin.Context) {
	c.Set("Login", false)
	c.Set("Superuser", false)
	cookie, err := c.Request.Cookie("Session")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println(err)
		}
		return
	}

	session := cookie.Value
	ok, user := db.GetUser(bson.D{{"session", session}})
	if !ok {
		return
	}

	if time.Now().Unix() < user.ExpireTime {
		c.Set("UserId", user.Id)
		c.Set("Login", true)
		if user.SuperUser {
			c.Set("Superuser", true)
		}
		db.UpdateUser(user.Id, bson.D{{"$set", bson.D{{"last_login_time", time.Now().Unix()}}}})
	} else {
		db.UpdateUser(user.Id, bson.D{{"$set", bson.D{{"session", nil}}}})
		c.SetCookie("Session", "", -1, "/", config.Domain, false, true)
	}
}

func TodayViewsIncr(c *gin.Context) {
	addr := util.GetAddr(c.Request.RemoteAddr)
	ok, s := db.GetDailyView()
	c.Set("TodayViews", 0)
	if ok {
		for _, i := range s.Address {
			if i == addr {
				c.Set("TodayViews", len(s.Address))
				return
			}
		}
		c.Set("TodayViews", len(s.Address)+1)
		db.DailyViewsPushAddr(addr)
	}
}

func Limiter(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, map[string]string{"msg": "Too Many Requests"})
		c.Abort()
	}
}

func EveryDayWork() {
	// todo daily work
}
