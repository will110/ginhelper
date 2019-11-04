# ginhelper
ginhelper 是一个gin web框架的快速搭建项目的助手，帮助你快速开发项目、简单、易用。


## 一、使用方法

### 1、下载
go get -u github.com/will110/ginhelper


### 2、生成项目

__(1)、如果你使用go mod管理项目，操作如下，这里以windows为例，其他系统操作一样，这里的操作都在d:/goMod/test目录下进行。__

    1、生成 go.mod文件
    
        set GO111MODULE=on
        go mod init hello


    2、生成项目
    
        ginhelper create


    3、更新依赖包
    
        go mod tidy
    
    
    4、编译
    
        go build
    

    5、运行
    
        ./hello
        
        
    6、测试
    
        curl http://127.0.0.1:8099/user/get-list-query?user=a3c&password=23
        或者直接在浏览器上访问，返回内容如下
        {
          "run_time": 0,
          "code": -110020,
          "message": "Key: 'Login.User' Error:Field validation for 'User' failed on the 'eq' tag",
          "data": ""
        }


    7、review code

        你可以通过某个开发工具，如：goland, 来打开你的项目，然后慢慢、静静的看一下代码



__(2)、如果你只是在src目录下开发项目，你可以这样操作，这里以windows为例，其他系统操作一样，这里的操作都在d:/goProject/src/hiworld目录下进行。__
        
    1、创建项目
    
        ginhelper create
       
    2、下载其他的依赖
    
       go get -u github.com/gin-gonic/gin
       go get github.com/astaxie/beego/config  #这里我们可以使用beego的config, 如果你还想使用beego的其他模块请自行下载，如：httplib、logs、cache等模块
       go get github.com/beego/bee  #可选下载
       github.com/jinzhu/gorm
       github.com/go-redis/redis
       
       
    3、编译
    
       go build 或者bee run
       
       
    4、运行
    
       ./hiworld
       
       
    5、测试
       
       curl http://127.0.0.1:8099/user/get-list-query?user=a3c&password=23
       或者直接在浏览器上访问，返回内容如下
       {
         "run_time": 0,
         "code": -110020,
         "message": "Key: 'Login.User' Error:Field validation for 'User' failed on the 'eq' tag",
         "data": ""
       }
    


## 目录介绍

    command  
        脚本目录，这个目录下可以存放一些脚本的业务逻辑，你根据自己业务的不同，然后创建自己的目录
        
    conf
        配置文件目录，这个目录下存放了各个环境的配置文件，然后我们只需要在环境变量中配置WebRunMode对应的值即可.
        比如: WebRunMode=test, 那么使用的配置文件就是test_app.conf, 而local_app.conf是本地配置文件优先级高于
        任何配置文件，如果你在test_app.conf中设置name = "hi", 然后在local_app.conf中设置name = "hello", 那么最终
        name的值是hello。
        
    controller
        控制器目录，这里是接口地址的入口点
    
    filter
        数据过滤层， 对每个接口请求过来的数据，进行严格的清洗，然后组装相应的数据，以便后面的程序使用
    
    model
        model层，可以存放一些关于表对应的结构体、字段或者简单操作表数据的方法
    
    pkg
        pkg存放自己封装的公共方法或公共变量、常量等
        
    router
        接口地址的路由配置
        
    runtime
        存放运行时的一些缓存文件、日志等
        
    servicelogic
        业务逻辑层，可以把大量的业务逻辑写在这里
        
    static
        静态目录，可以存放图片、文件等
        
    test
        接口测试目录



当然了大家也可以根据自己的喜好设计自己喜欢的目录结构，项目框架整合了mysql、redis、mongodb等相关配置，
默认代码是注释的，你需要把mysql、redis、mongodb的配置文件进行配置，然后打开注释，试着运行一下。