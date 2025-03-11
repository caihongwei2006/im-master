package service

import (
	"fmt"
	"im-master/models"
	"strconv"

	"net/http"

	"im-master/utils"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList godoc
// @Summary      获取用户列表
// @Description  返回系统中所有用户的列表
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{} "返回用户列表和消息"
// @Router       /user/userlist [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": "service.go  Ready!!!!!!",
		"data":    data,
	})
}

// CreateUser godoc
// @Summary      创建新用户
// @Description  根据提供的用户名和密码创建新用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name       query    string  true  "用户名"
// @Param        password   query    string  true  "密码"
// @Param        repassword query    string  true  "确认密码"
// @Success      200  {object}  map[string]string  "创建成功响应"
// @Failure      400  {object}  map[string]string  "密码不匹配错误"
// @Router       /user/createuser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword") //本来repassword依靠前端处理，后端不用管
	if password != repassword {
		c.JSON(400, gin.H{
			"message": "Passwords do not match",
		})
		return
	}
	user.Password = password
	models.CreateUser(&user)
	c.JSON(200, gin.H{
		"message": "Succeeded in creating usser!!!!!!!!!!!!!",
	})
}

// DeleteUser godoc
// @Summary      删除用户
// @Description  根据用户名删除指定用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name  query  string  true  "要删除的用户名"
// @Success      200  {object}  map[string]string  "删除成功响应"
// @Router       /user/deleteuser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	models.DeleteUser(&user)
	c.JSON(200, gin.H{
		"message": "删除先辈成功",
	})
}

// UpdateUser godoc
// @Summary      更新用户信息
// @Description  根据用户ID更新用户名和密码
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id        formData  integer  true  "用户ID"
// @Param        name      formData  string   true  "新用户名"
// @Param        password  formData  string   true  "新密码"
// @Param        email  formData  string   true  "新密码"
// @Param        phone  formData  string   true  "新密码"
// @Success      200  {object}  map[string]string  "更新成功响应"
// @Router       /user/updateuser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	fmt.Println("update:    ", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println("err: ", err)
		c.JSON(400, gin.H{
			"message": "参数格式不匹配，需要正确的电话号码",
		})
	}

	models.UpdateUser(&user)
	c.JSON(200, gin.H{
		"message": "更新先辈成功",
	})
}

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
