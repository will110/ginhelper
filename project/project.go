package project

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var dirList = []string{
	"command",
	"conf",
	"controller",
	"filter",
	"model",
	"model/user",
	"model/userdetail",
	"pkg",
	"pkg/utils",
	"pkg/db",
	"pkg/myerror",
	"pkg/param",
	"router",
	"runtime",
	"static",
	"servicelogic",
	"test",
}

var modName string
var appName string

func GenerateProject() {
	modName = getModName()
	if len(modName) == 0 {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			log.Fatal("GOPATH environment variable is not set or empty")
		}

		appDir, err := getAppDir()
		if err != nil {
			log.Fatalln(err)
		}

		modName = appDir

		buf := strings.Split(appDir, "/")
		appName = buf[len(buf) - 1]
	} else {
		appName = modName
	}

	fmt.Println("Project creating ...")
	generateDir()
	generateAllFile()
	fmt.Println("Congratulations on the completion of the project. Take a look at it.")
}

func generateDir() {
	currentPath, _ := os.Getwd()
	for _, v := range dirList {
		fp := path.Join(currentPath, v)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if er := os.MkdirAll(fp, 0777); er != nil {
				log.Fatalf("Could not create the "+ v +" directory: %s", er)
			}
		}
	}
}

func generateAllFile() {
	//controller
	generateFile( "controller", "BaseController", baseControllerTemp)
	generateFile( "controller", "UserController", userControllerTemp)

	//config
	generateFile( "conf", "app", appConf, "conf")
	generateFile( "conf", "test_app", testAppConf, "conf")
	generateFile( "conf", "prod_app", prodAppConf, "conf")
	generateFile( "conf", "preview_app", previewAppConf, "conf")
	generateFile( "conf", "local_app", localAppConf, "conf")
	generateFile( "conf", "dev_app", devAppConf, "conf")

	//filter
	generateFile( "filter", "UserFilter", userFilterTemp)

	//router
	generateFile( "router", "router", routerTemp)

	//pkg db
	generateFile( "pkg/db", "db", pkgDbdbTemp)
	generateFile( "pkg/db", "mongodb", pkgDbMongodbTemp)
	generateFile( "pkg/db", "redis", pkgDbRedisTemp)
	//pkg myerror
	generateFile( "pkg/myerror", "errorCode", pkgMyerrorTemp)

	//pkg param
	generateFile( "pkg/param", "user", pkgParamTemp)

	//pkg utils
	generateFile( "pkg/utils", "config", pkgUtilsConfigTemp)
	generateFile( "pkg/utils", "engine", pkgUtilsEngineTemp)
	generateFile( "pkg/utils", "error", pkgUtilsErrorTemp)
	generateFile( "pkg/utils", "file", pkgUtilsFileTemp)
	generateFile( "pkg/utils", "response", pkgUtilsResponseTemp)

	//model
	generateFile( "model/user", "User", modelUserTemp)
	generateFile( "model/userdetail", "UserDetail", modelUserDetailTemp)

	//servicelogic
	generateFile( "servicelogic", "UserLogic", servicelogicUserLogicTemp)

	//main
	generateFile( "", "config", mainConfigTemp)
	generateFile( "", "main", mainTemp)

	//Gitignore
	generateGitignoreFile()

}

func generateGitignoreFile() {
	currentPath, _ := os.Getwd()
	fpath := path.Join(currentPath, ".gitignore")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create gitignore file: %s", err)
	}

	gitignoreTmep += appName + ".exe\n" + appName + ".exe*\n" + appName+"\n"
	_, _ = f.WriteString(gitignoreTmep)
	_ = f.Close()
}

func getAppDir() (string, error) {
	currentPath, _ := os.Getwd()
	if strings.Index(currentPath, "src") == -1 {
		return "", errors.New("You must create file in the src directory or go mod project")
	}

	fileList := strings.Split(currentPath, "src")
	if len(fileList[1]) == 0 {
		return "", errors.New("You must create file in the src directory or go mod project")
	}

	dir := strings.Replace(fileList[1], string(filepath.Separator), "/", -1)[1:]

	return dir, nil
}

func getModName() string {
	var modName string
	_ = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if path != "go.mod" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		fread := bufio.NewReader(f)
		for {
			strBuf, _, err := fread.ReadLine()
			if err != nil {
				continue
			}
			str := string(strBuf)

			if strings.Contains(str, "module") {
				str = strings.ReplaceAll(str, "module", "")
				str = strings.ReplaceAll(str, " ", "")
				str = strings.ReplaceAll(str, "	", "")
				modName = str
				break
			}
		}

		return nil
	})

	return modName
}

func generateFile(dir, fileName, str string, ext ...string) {
	currentPath, _ := os.Getwd()
	extStr := "go"
	if len(ext) > 0 {
		extStr = ext[0]
	}
	fpath := path.Join(currentPath, dir, fileName + "." + extStr)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create %s/%s.go" +" file: %s", dir, fileName, err)
	}

	str = strings.Replace(str, "{{baseDir}}", modName, -1)
	_, _ = f.WriteString(str)
	_ = f.Close()
	cmd := exec.Command("gofmt", "-w", fpath)
	_ = cmd.Run()
}