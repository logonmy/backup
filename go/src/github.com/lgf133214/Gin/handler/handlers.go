package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lgf133214/Gin/config"
	"github.com/lgf133214/Gin/db"
	"github.com/lgf133214/Gin/model"
	"github.com/lgf133214/Gin/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Maps map[string]interface{}

func (m *Maps) Set(key string, value interface{}) *Maps {
	(*m)[key] = value
	return m
}

func Param(c *gin.Context) *Maps {
	m := make(Maps)
	m["Login"] = c.GetBool("Login")
	m["Superuser"] = c.GetBool("Superuser")
	m["TodayViews"] = c.GetInt("TodayViews")
	return &m
}

func Index(c *gin.Context) {
	ok, page := GetPage(c)
	page--
	if ok {
		c.HTML(http.StatusOK, "index.html", Param(c).Set("Posts", db.GetPosts(bson.D{}, options.Find().SetSort(
			bson.D{{"pub_time", -1}}).SetLimit(10).SetSkip(page*10))).Set("Title", "Lazy Li"))
	}
}

func GetPage(c *gin.Context) (bool, int64) {
	page := c.DefaultQuery("page", "1")
	parseInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || parseInt < 1 {
		DefaultIndex(c)
		return false, 0
	}
	return true, parseInt
}

func DefaultIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", Param(c).Set("Posts", db.GetPosts(bson.D{}, options.Find().SetSort(
		bson.D{{"pub_time", -1}}).SetLimit(10))).Set("Title", "Lazy Li"))
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", Param(c).Set("Title", "你要的页面找不到啦！"))
	c.Abort()
}

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", Param(c).Set("Title", "关于"))
}

func Contact(c *gin.Context) {
	if c.GetBool("Login") {
		id := c.GetString("UserId")
		obj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusOK, "contact.html", Param(c).Set("Title", "联系作者"))
			return
		}

		ok, user := db.GetUser(bson.D{{"_id", obj}})
		if ok {
			c.HTML(http.StatusOK, "contact.html", Param(c).Set("Title", "联系作者").Set("Email", user.Mail))
			return
		}
	}
	c.HTML(http.StatusOK, "contact.html", Param(c).Set("Title", "联系作者"))
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Logout(c *gin.Context) {
	if c.GetBool("Login") {
		session, err := c.Request.Cookie("Session")
		if err != nil {
			if err != http.ErrNoCookie {
				log.Println(err)
			}
			return
		}
		if db.DelSession(session.Value) {
			c.Redirect(http.StatusTemporaryRedirect, "/")
		}
	}
}

func ResetPw(c *gin.Context) {
	c.HTML(http.StatusOK, "forget.html", Param(c).Set("Title", "修改密码"))
}

func View(c *gin.Context) {
	id := c.Param("postId")
	if id == "" {
		NotFound(c)
		return
	}
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		NotFound(c)
		return
	}

	posts := db.GetPosts(bson.D{{"_id", obj}})
	if len(posts) != 1 {
		NotFound(c)
	} else {
		c.HTML(http.StatusOK, "blog-detail.html", Param(c).Set("Post", posts[0]))
	}
}

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func LoginPost(c *gin.Context) {
	uname := c.PostForm("username")
	passwd := c.PostForm("password")
	ok, user := db.GetUser(bson.D{{"mail", uname}, {"password", passwd}})
	if !ok {
		ok, user = db.GetUser(bson.D{{"name", uname}, {"password", passwd}})
	}
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "用户名或密码错误", "tip": "是不是不知道账号密码！？？"})
		return
	}
	if !user.Validated {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "还没激活呢, 快去邮箱看看吧", "tip": "不知道要激活吗？！！！"})
		return
	}

	session := util.GetSession(user.Mail)
	c.SetCookie("Session", session, int(config.SessionExpireTime*3600*24), "/", config.Domain, false, true)
	db.UpdateUser(user.Id, bson.D{{"$set",
		bson.D{{"session", session}, {"expire_time",
			time.Now().Unix() + config.SessionExpireTime*3600*24}}}})
	c.JSON(http.StatusOK, map[string]string{"status": "ok", "msg": "登录成功", "tip": "欢迎，" + user.Name})
}

func Activation(c *gin.Context) {
	code := c.Param("verifyCode")
	ok, ret := db.GetVerifyCode(code, true)
	if !ok {
		NotFound(c)
		return
	}
	if ret.ExpireTime > time.Now().Unix() {
		db.UpdateUser(ret.UserId, bson.D{{"$set", bson.D{{
			"validated", true,
		}, {
			"join_time", time.Now().Unix(),
		}}}})
		c.HTML(http.StatusOK, "login.html", gin.H{"Msg": "激活成功"})
	} else {
		db.DelUser(ret.UserId)
		c.HTML(http.StatusOK, "login.html", gin.H{"Msg": "验证码过期，请重新申请注册"})
	}
	db.DelVerifyCode(ret.Id)
}

func ToResetPw(c *gin.Context) {
	code := c.Param("verifyCode")
	ok, ret := db.GetVerifyCode(code, false)
	if !ok {
		NotFound(c)
		return
	}
	if ret.ExpireTime > time.Now().Unix() {
		session := util.GetSession(ret.Code)
		db.UpdateUser(ret.UserId, bson.D{{"$set", bson.D{{"session", session}, {"expire_time", -1}}}})
		c.HTML(http.StatusOK, "reset.html", gin.H{"Session": session})
	} else {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "验证码过期，请重新申请", "tip": ""})
	}
	db.DelVerifyCode(ret.Id)
}

func SendMail(c *gin.Context) {
	mail, ok := c.GetPostForm("email")
	if !ok || !util.VerifyEmailFormat(mail) {
		JsonNotAccept(c)
		return
	}
	name, ok := c.GetPostForm("name")
	if !ok {
		JsonNotAccept(c)
		return
	}
	msg, ok := c.GetPostForm("message")
	if !ok {
		JsonNotAccept(c)
		return
	}
	ok, _ = db.GetUser(bson.D{{"mail", mail}})
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "账号不存在, 请注册后使用", "tip": ""})
		return
	}
	ok = util.SendMessage(name, msg, mail)
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "发送失败，邮件可能被外星人劫持了。手动发试试？", "tip": ""})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "ok", "msg": "发送成功，感谢您的每一次反馈", "tip": ""})
}

func ResetPwPost(c *gin.Context) {
	s, ok := c.GetPostForm("session")
	if !ok {
		JsonNotAccept(c)
		return
	}
	ok, user := db.GetUser(bson.D{{"session", s}})
	if !ok {
		JsonNotAccept(c)
		return
	}
	pw, ok := c.GetPostForm("password")
	if !ok {
		JsonNotAccept(c)
		return
	}
	ok = db.UpdateUser(user.Id, bson.D{{"$set", bson.D{{"password", pw}, {"session", nil}}}})
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"msg": "修改失败", "status": "no", "tip": ""})
	} else {
		c.JSON(http.StatusOK, map[string]string{"msg": "修改成功", "status": "ok", "tip": ""})
	}
}

func Search(c *gin.Context) {

}

func RegisterPost(c *gin.Context) {
	mail, ok := c.GetPostForm("email")
	if !ok || !util.VerifyEmailFormat(mail) {
		JsonNotAccept(c)
		return
	}
	uname, ok := c.GetPostForm("username")
	if !ok {
		JsonNotAccept(c)
		return
	}
	pw, ok := c.GetPostForm("password")
	if !ok {
		JsonNotAccept(c)
		return
	}

	ok, _ = db.GetUser(bson.D{{"mail", mail}})
	if ok {
		c.JSON(http.StatusOK, map[string]string{"msg": "邮箱已被使用", "status": "no", "tip": ""})
		return
	}
	ok, _ = db.GetUser(bson.D{{"name", uname}})
	if ok {
		c.JSON(http.StatusOK, map[string]string{"msg": "用户名已被使用", "status": "no", "tip": ""})
		return
	}

	code := util.GetSession(mail)
	ok = db.AddUser(model.User{Name: uname, PassWord: pw, Mail: mail, Icon: "/assets/images/blog/user.jpg"})
	if !ok {
		JsonError(c)
		return
	}
	ok, user := db.GetUser(bson.D{{"mail", mail}})
	if !ok {
		JsonError(c)
		return
	}

	ok = db.AddVerifyCode(model.VerifyCode{ExpireTime: time.Now().Add(time.Minute * 10).Unix(), Code: code, UserId: user.Id, Register: true})
	if !ok {
		JsonError(c)
		return
	}

	ok = util.SendVerifyCode(code, true, mail)
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "邮件发送失败，请检查邮箱正确性或更换邮箱", "tip": ""})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "ok", "msg": "发送成功，请点击邮件内的激活链接激活", "tip": "十分钟内有效，如未收到请检查是否在垃圾箱内"})
}

func JsonError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, map[string]string{"status": "no", "msg": "出错啦！等会试试吧", "tip": ""})
}

func JsonNotAccept(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, map[string]string{"status": "no", "msg": "请求不符合符合规则", "tip": ""})
}

func ResetSendCode(c *gin.Context) {
	mail, ok := c.GetPostForm("email")
	if !ok || !util.VerifyEmailFormat(mail) {
		JsonNotAccept(c)
		return
	}
	ok, user := db.GetUser(bson.D{{"mail", mail}})
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "未注册的邮箱", "tip": ""})
		return
	}
	code := util.GetSession(user.Mail)
	ok = db.AddVerifyCode(model.VerifyCode{Code: code, UserId: user.Id, ExpireTime: time.Now().Add(time.Minute * 10).Unix()})
	if !ok {
		JsonError(c)
		return
	}
	ok = util.SendVerifyCode(code, false, mail)
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "发送失败，等会再来试试吧？", "tip": ""})
		return
	}
	db.UpdateUser(user.Id, bson.D{{"$set", bson.D{{"session", nil}, {"expire_time", nil}}}})
	c.JSON(http.StatusOK, map[string]string{"status": "ok", "msg": "发送成功，十分钟内有效，请尽快完成修改", "tip": ""})
}
