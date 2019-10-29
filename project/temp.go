package project

var baseControllerTemp = `package controller

import (
	"fmt"
	"github.com/will110/ginhelper/ctl"
	"{{baseDir}}/pkg/myerror"
	"{{baseDir}}/pkg/utils"
	"time"
)

type BaseController struct {
	RunStartTime time.Time
	ctl.Controller
}

//子类初始化
type InitialiseInterface interface {
	Initialise()
}

func (m *BaseController) Prepare() {
	m.RunStartTime = time.Now()
	if app, ok := m.AppController.(InitialiseInterface); ok {
		app.Initialise()
	}
	//这里可以做一些初始化的工作
}

func (m *BaseController) Finish() {
	fmt.Println(5, "Finish")
	//在这个方法可以，可以做一些收尾的工作，比如记录访问接口的日志
}

/**
params1 data
params2 message
params3 code   负数类(即小于0的)的表示失败，1表示成功.
code可以在调用setResponse时设置，这个是全局设置，如：
m.SetResponse(nil, "error message", -110122)

也可以在每个错误信息时，设置code,表示局部的错误信息, 如：
a := utils.NewErrorCode(-110123, "error message")
m.SetResponse(nil, a)

或者重置局部code
a := utils.NewErrorCode(-110123, "error message")
m.SetResponse(nil, a, -110122) 这样的话会打印-110122, 不会打印-110123

params4 other  其他的一些参数，用于方法中要记的日志
*/
func (m *BaseController) SetResponse(params ...interface{}) {
	var code = 1
	var message string
	var data interface{}

	lenBuf := len(params)
	//如果没有传参数，data就设置空
	if lenBuf == 0 {
		data = struct{}{}
	}

	//如果有一个，就把第一个设置到data中
	if lenBuf > 0 {
		if params[0] != nil {
			data = params[0]
		} else {
			data = struct{}{}
		}
	} else {
		data = struct{}{}
	}

	if lenBuf > 1 {
		isCommomError := true
		switch params[1].(type) {
		case string:
			message = params[1].(string)

		case *utils.ErrorCode:
			b := params[1].(*utils.ErrorCode)
			if len(b.Content) == 0 {
				if v, ok := myerror.CodeList[b.CodeId]; ok {
					message = v
				}
			} else {
				message = b.Content
			}

			code = b.CodeId
			isCommomError = false

		default:
			message = "conversion is error"
		}

		if isCommomError {
			code = myerror.CommonErrorOfErrorCode
		}
	}

	if lenBuf > 2 {
		code = params[2].(int)
	}

	responseList := utils.ResponseList{
		RunTime: time.Since(m.RunStartTime).Seconds(),
		Code:    code,
		Message: message,
		Data:    data,
	}

	m.C.JSON(200, responseList)
}
`

var userControllerTemp = `package controller

import (
	"{{baseDir}}/filter"
)

var userFilter *filter.UserFilter
type UserController struct {
	BaseController
}

func (m *UserController) Initialise()  {
	userFilter = filter.NewUserFilter(m.C)
}

func (m *UserController) GetListForm() {
	str, err := userFilter.GetListForm()
	if err != nil {
		m.SetResponse(str, err)
		return
	}

	m.SetResponse(str)
}

func (m *UserController) GetListQuery() {
	str, err := userFilter.GetListQuery()
	if err != nil {
		//m.SetResponse(str, err, -110010) 也可以重写接受过来的错误码，自己定义
		m.SetResponse(str, err)
		return
	}

	m.SetResponse(str)
}

func (m *UserController) GetListJson() {
	str, err := userFilter.GetListJson()
	if err != nil {
		m.SetResponse(str, err)
		return
	}

	m.SetResponse(str)
}

func (m *UserController) GetListXml() {
	str, err := userFilter.GetListXml()
	if err != nil {
		m.SetResponse(str, err)
		return
	}

	m.SetResponse(str)
}
`

var appConf = `name = "hello"
version = 1.0
`

var devAppConf = `include "app.conf"

[db]
host = 127.0.0.1:3306
user = root
password = 123456
database = test

[redis]
host = 127.0.0.1:6379
password =
database = 10

[mongodb]
host = 127.0.0.1:27017
name = "hello"

#请保持在最后一行
include "local_app.conf"

`

var localAppConf = `name = "hello-local"
`

var previewAppConf = `include "app.conf"
#当前环境的配置可以配置在这里

#请保持在最后一行
include "local_app.conf"
`

var prodAppConf = `include "app.conf"
#当前环境的配置可以配置在这里

#请保持在最后一行
include "local_app.conf"
`

var testAppConf = `include "app.conf"
#当前环境的配置可以配置在这里

#请保持在最后一行
include "local_app.conf"
`

var userFilterTemp = `package filter

import (
	"github.com/gin-gonic/gin"
	"{{baseDir}}/pkg/param"
	"{{baseDir}}/pkg/utils"
	"{{baseDir}}/servicelogic"
)

type UserFilter struct {
	c *gin.Context
}

func NewUserFilter(c *gin.Context) *UserFilter {
	return &UserFilter{c: c}
}

type LoginJson struct {
	User     string `+"`json:\"user\" binding:\"required\"`"+`
	Password string `+"`json:\"password\" binding:\"required\"`"+`
}

type LoginXml struct {
	User     string `+"`xml:\"user\" binding:\"required\"`"+`
	Password string `+"`xml:\"password\" binding:\"required\"`"+`
}

type Login struct {
	User     string `+"`form:\"user\" binding:\"required,len=3,eq=abc\"`"+`
	Password string `+"`form:\"password\" binding:\"required\"`"+`
}

//数据清洗，组装相关操作, 这里只作简单的介绍，根据业务场景自己组装相关数量
func (m *UserFilter) GetListQuery() (string, *utils.ErrorCode) {
	var person Login
	if err := m.c.ShouldBindQuery(&person); err != nil{
		return "", utils.NewErrorCode(-110020, err.Error())
	}

	//数据组装，为后面程序准备
	p := &param.UserForm{
		User:     person.User,
		Password: person.Password,
		BusinessId: 1,
	}
	list := servicelogic.NewUserLogic(m.c).GetListQuery(p)

	return list, nil
}

//数据清洗，组装相关操作, 这里只作简单的介绍，根据业务场景自己组装相关数量
func (m *UserFilter) GetListJson() (map[string]interface{}, *utils.ErrorCode) {
	var json LoginJson
	if err := m.c.ShouldBindJSON(&json); err != nil {
		//NewErrorCode也可以直接填写一个错误码，然后错误码信息从myerror.CodeList列表中获取，所以错误信息可以维护在这个公共的地方。
		//也可以在这里自己定义错误码和内容，如：utils.NewErrorCode(-150101, "hello")
		return nil, utils.NewErrorCode(-140101)
	}

	p := &param.UserForm{
		User:     json.User,
		Password: json.Password,
		BusinessId: 1,
	}
	list, err := servicelogic.NewUserLogic(m.c).GetListJson(p)
	if err != nil {
		return nil, utils.NewErrorCode(-110021, err.Error())
	}

	return list, nil
}

//数据清洗，组装相关操作, 这里只作简单的介绍，根据业务场景自己组装相关数量
func (m *UserFilter) GetListXml() (string, *utils.ErrorCode) {
	var xml LoginXml
	if err := m.c.ShouldBindXML(&xml); err != nil {
		return "", utils.NewErrorCode(-110021, err.Error())
	}

	p := &param.UserForm{
		User:     xml.User,
		Password: xml.Password,
		BusinessId: 1,
	}
	list := servicelogic.NewUserLogic(m.c).GetListXml(p)

	return list, nil
}

//数据清洗，组装相关操作, 这里只作简单的介绍，根据业务场景自己组装相关数量
func (m *UserFilter) GetListForm() (string, *utils.ErrorCode) {
	var form Login
	if err := m.c.ShouldBind(&form); err != nil {
		return "", utils.NewErrorCode(-110010, err.Error())
	}

	p := &param.UserForm{
		User:     form.User,
		Password: form.Password,
		BusinessId: 1,
	}
	list := servicelogic.NewUserLogic(m.c).GetListForm(p)

	return list, nil
}
`

var routerTemp = `package router

import (
	"github.com/will110/ginhelper/routerhelper"
	"{{baseDir}}/controller"
	"{{baseDir}}/pkg/utils"
)

func InitRouter() {
	r := utils.NewGinDefault()
	//第一种路由方式
	r.POST("/user/get-list-json", routerhelper.BindRouter(new(controller.UserController), "GetListJson"))
	r.POST("/user/get-list-form", routerhelper.BindRouter(new(controller.UserController), "GetListForm"))

	//第二种路由方式
	r1 := r.Group("user")
	{
		r1.GET("/get-list-query", routerhelper.BindRouter(new(controller.UserController), "GetListQuery"))
		r1.GET("/get-list-xml", routerhelper.BindRouter(new(controller.UserController), "GetListXml"))
	}
}
`

var pkgDbdbTemp = `package db

/*
import (
	"flag"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/astaxie/beego"
)

var (
	Conn *gorm.DB
)

func GetDbConnect() (*gorm.DB, error) {
	dbConStr := fmt.Sprintf("%s:%s@tcp(%s)/%s%s?charset=utf8&parseTime=True&loc=Local", beego.AppConfig.String("db::user"), beego.AppConfig.String("db::password"), beego.AppConfig.String("db::host"), beego.AppConfig.String("db::database"), getConnectDbName())
	fmt.Println(dbConStr, "dbConStr")
	db, err := gorm.Open("mysql", dbConStr)
	Conn = db

	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" || runMode == "test" {
		db.LogMode(true)
	}

	return db, err
}

func getConnectDbName() string {
	dbName := flag.String("db", "", "")
	flag.Parse()

	return *dbName
}
*/
`

var pkgDbMongodbTemp = `package db

/*
import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	Mongodb *mongo.Client
)

//连接mongodb
func GetMongodbClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	host := fmt.Sprintf("mongodb://%s", beego.AppConfig.String("mongodb::host"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))
	cancel()
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2 * time.Second)
	err = client.Ping(ctx, readpref.Primary())
	cancel()
	if err != nil {
		return nil, err
	}

	Mongodb = client

	return client, nil
}
*/
`

var pkgDbRedisTemp = `package db

/*
import (
	"github.com/go-redis/redis"
	"github.com/astaxie/beego"
)

var (
	Redis *redis.Client
)

//连接redis
func GetRedisClient() (*redis.Client, error){
	db ,_ := beego.AppConfig.Int("redis::database")
	redisOption := &redis.Options{
		Addr:     beego.AppConfig.String("redis::host"),
		DB:       db,
	}

	runMode := beego.AppConfig.String("runmode")
	if runMode != "prod" {
		redisOption.Password = beego.AppConfig.String("redis::password")
	}
	client := redis.NewClient(redisOption)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	Redis = client

	return client, nil
}
*/
`

var pkgMyerrorTemp = `package myerror

const (
	//通用的错误提示
	CommonErrorOfErrorCode = -101

	//没有权限
	NoPermissionOfErrorCode = -120
	//没有找到
	NotFoundOfErrorCode = -404
	//服务器错误
	ServerInternalErrorOfErrorCode = -500

)

//错误码规则定义：一共6位，比如：110101， 第一第二位表示项目，第三第四位表示哪个controller, 最后二位表示这个控制器下面的错误码
//比如如下
var CodeList = map[int]string{
	-140101:"你好呀~~~~~~",
	-140102:"你好呀1",
	-140103:"你好呀2",
	-140104:"你好呀3",
	-140105:"你好呀4",
}
`

var pkgParamTemp = `package param

type UserForm struct {
	User string
	Password string
	BusinessId uint64
}
`

var pkgUtilsConfigTemp = `package utils

var (
	MyConfig *WebConfig
)

type WebConfig struct {
	DbHost string
	DbName string
	DbPassword string
	Version float64
}
`

var pkgUtilsEngineTemp = `package utils

import "github.com/gin-gonic/gin"

var R *gin.Engine

func NewGinDefault() *gin.Engine {
	r := gin.Default()
	R = r

	return r
}
`

var pkgUtilsErrorTemp = `package utils

import "bytes"

type ErrorCode struct {
	CodeId int
	Content string
}

func NewErrorCode(codeId int, content ...string) *ErrorCode {
	buf := bytes.NewBuffer(nil)
	for _, v := range content {
		buf.WriteString(v)
	}

	return &ErrorCode{
		CodeId: codeId,
		Content: buf.String(),
	}
}
`

var pkgUtilsFileTemp = `package utils

import "os"

func DirIsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
`

var pkgUtilsResponseTemp = `package utils

type ResponseList struct {
	RunTime	float64		`+"`json:\"run_time\"`"+`
	Code    int         `+"`json:\"code\"`"+`
	Message string      `+"`json:\"message\"`"+`
	Data    interface{} `+"`json:\"data\"`"+`
}
`

var modelUserTemp = `package user

import (
	"fmt"
	//"{{baseDir}}/pkg/db"
)

type Users struct {
	UserId uint64
	BusinessId uint64
	UserName string
}

//获取表名
func GetTableName() string {
	return "users"
}

func GetField() []string {
	return []string{
		"user_id","business_id", "user_name",
	}
}

/*
根据商户id和用户id，获取用户信息
注意：如果是简单的数据库操作或者是公共的方法，也可以封装在model中，
当然了也可以封装在servicelogic中，这个根据业务场景来决定。
 */
/*
func GetUserListByUserId(userName string, businessId uint64) (*Users, error) {
	userList := &Users{}
	err := db.Conn.Table(GetTableName()).Select(GetField()).
		Where("user_id = ? AND business_id = ?", userName, businessId).
		Find(userList).Error
	if err != nil {
		return nil, fmt.Errorf("select memeber info is fail, err: %v", err)
	}

	return userList, nil
}
*/
`

var modelUserDetailTemp = `package userdetail

type UserDetails struct {
	UserId uint64
	UserName string
	Sex uint8
	Age uint8
}

//获取表名
func GetTableName() string {
	return "user_details"
}

func GetField() []string {
	return []string{
		"user_id","user_name","sex","age",
	}
}
`

var servicelogicUserLogicTemp = `
package servicelogic

import (
	"github.com/gin-gonic/gin"
	//"{{baseDir}}/model/user"
	"{{baseDir}}/model/userdetail"
	//"{{baseDir}}/pkg/db"
	"{{baseDir}}/pkg/param"
)

type UserLogic struct {
	c *gin.Context
}

func NewUserLogic(c *gin.Context) *UserLogic {
	return &UserLogic{c: c}
}

func (m *UserLogic) GetListQuery(p *param.UserForm) string {
	//相关业务逻辑
	name := "hello" + p.User

	return name
}

func (m *UserLogic) GetListJson(p *param.UserForm) (map[string]interface{}, error) {
	//相关业务逻辑 比如：操作数据库
	//userList, err := user.GetUserListByUserId(p.User, 101)
	//if err != nil {
	//	return nil, err
	//}

	detailInfo := &userdetail.UserDetails{}
	//err = db.Conn.Table(userdetail.GetTableName()).Select(userdetail.GetField()).
	//	Where("user_id = ?", userList.UserId).Find(detailInfo).Error
	//if err != nil {
	//	return nil, err
	//}

	list := make(map[string]interface{})
	list["name"] = detailInfo.UserName
	list["age"] = detailInfo.Age
	list["password"] = p.Password

	return list, nil
}

func (m *UserLogic) GetListXml(p *param.UserForm) string {
	//相关业务逻辑
	name := "hello" + p.User

	return name
}


func (m *UserLogic) GetListForm(p *param.UserForm) string {
	//相关业务逻辑
	name := "hello" + p.User

	return name
}
`
var mainTemp = `package main

import (
	"log"
	"{{baseDir}}/pkg/db"
	"{{baseDir}}/pkg/utils"
	"{{baseDir}}/router"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func main() {
	InitConfig()
	router.InitRouter()

	/*
	//连接mysql数据库
	conn, err := db.GetDbConnect()
	if err != nil {
		log.Fatal("mysql 连接失败， %v", err)
	}
	//把表名的复数形式去掉
	conn.SingularTable(true)
	//设置mysql的空闲连接数为10个
	conn.DB().SetMaxIdleConns(20)
	conn.DB().SetMaxOpenConns(1000)
	conn.DB().SetConnMaxLifetime(2 * time.Minute)

	//连接redis
	redisClient, err := db.GetRedisClient()
	if err != nil {
		log.Fatal("redis 连接失败, %v", err)
	}

	defer func() {
		err = conn.Close()
		spew.Dump(err, "mysql")

		err = redisClient.Close()
		spew.Dump(err, "redis")
	}()
	*/
	log.Fatal(utils.R.Run(":8099"))

}
`

var gitignoreTmep = `runtime/log/
debug
.idea
local_app.conf
`

var mainConfigTemp = `package main

import (
	"fmt"
	"log"
	"os"
	"{{baseDir}}/pkg/utils"
	"github.com/astaxie/beego"
)

var envMode = "WebRunMode"
var env = "dev"

//初始化配置文件
func InitConfig() {
	//环境配置
	bufEnv := os.Getenv(envMode)
	if len(bufEnv) > 0 {
		env = bufEnv
	}

	configName := fmt.Sprintf("conf/%s_app.conf", env)
	err := beego.LoadAppConfig("ini", configName)
	if err != nil {
		log.Fatal(err)
	}

	configList, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	utils.MyConfig = configList
}

/*
加载配置文件到缓存中
如果你对配置的性能要求比较高，你可以提前把一些参数添加到缓存中，在后面的调用中，直接取缓存即可。
比如：如下
*/
func LoadConfig() (*utils.WebConfig, error) {
	config := &utils.WebConfig{}
	config.DbHost = beego.AppConfig.String("db::host")
	config.DbName = beego.AppConfig.String("db::user")
	config.DbPassword = beego.AppConfig.String("db::password")
	version, err := beego.AppConfig.Float("version")
	if err != nil {
		config.Version = version
	}

	return config, err
}
`