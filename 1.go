package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "net/http"
	_ "strconv"
)

type msgs struct {
	Id    int
	Name  string
	Msg   string
	Pid   int
	Child []msgs
}

var db *sql.DB

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, _ = db.Exec("create table accounts (name text, password text, isAdmin int)")
	_, _ = db.Exec("create table msgs (id int primary key auto_increment, name text, msg text, pid int)")

	r := gin.Default()
	r.GET("/register", register)
	r.GET("/login", login)
}

func initDB() (err error) {
	dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var pwd string
	row := db.QueryRow("select password from accounts where username = ?", username)
	row.Scan(&pwd)
	if username == "" || password == "" {
		c.JSON(200, gin.H{"result": "账号和密码不能为空"})
		return
	}
	if pwd != password {
		c.JSON(200, gin.H{"result": "该账号已被注册"})
		_, _ = db.Exec("insert into accounts values ( ? , ? , ? )", username, password, 0)
	}
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var pwd string
	row := db.QueryRow("select password from accounts where username = ?", username)
	err := row.Scan(&pwd)
	if username == "" || password == "" {
		c.JSON(200, gin.H{"result": "账号和密码不能为空"})
		return
	}
	if err == nil {
		c.JSON(200, gin.H{"result": "账号不存在"})
		return
	}
	if pwd != password {
		c.JSON(200, gin.H{"result": "密码错误"})
	}
	c.SetCookie("username", username, 10, "/", "localhost", false, true)
	c.JSON(200, gin.H{"result": "登录成功"})

}
func showMsgs(c *gin.Context) {

	var ret msgs
	ret.Id = 0
	ret.Name = "账号名"
	ret.Msg = "留言"
	ret.Pid = 0
	getAllChild(&ret)
	c.JSON(200, ret)

}
func getAllChild(msg *msgs) {

	msg.Child = make([]msgs, 0, 0)
	row, _ := db.Query("select * from msgs where pid = ?", msg.Id)
	for row.Next() {
		var m msgs
		_ = row.Scan(&m.Id, &m.Name, &m.Msg, &m.Pid)
		getAllChild(&m)
		msg.Child = append(msg.Child, m)
	}
}
