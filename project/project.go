package project

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var dirList = []string{
	"commands",
	"conf",
	"controllers",
	"filters",
	"models",
	"pkg",
	"routers",
	"runtime",
	"serviceLogics",
	"tests",
}

func GenerateProject() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("GOPATH environment variable is not set or empty")
	}

	generateDir()
	generateControllerFile()
	generateRouterFile()
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

func generateControllerFile() {
	currentPath, _ := os.Getwd()
	fpath := path.Join(currentPath, "controllers", "BaseController.go")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create BaseController.go file: %s", err)
	}

	_, _ = f.WriteString(baseControllerTemp)
	_ = f.Close()
	cmd := exec.Command("gofmt", "-w", fpath)
	_ = cmd.Run()


	fpath = path.Join(currentPath, "controllers", "UserController.go")
	f, err = os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create UserController.go file: %s", err)
	}

	_, _ = f.WriteString(userControllerTemp)
	_ = f.Close()
	cmd = exec.Command("gofmt", "-w", fpath)
	_ = cmd.Run()
}

func generateRouterFile() {
	currentPath, _ := os.Getwd()
	if strings.Index(currentPath, "src") == -1 {
		log.Fatalln("you create file in src directory")
	}

	fpath := path.Join(currentPath, "routers", "router.go")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create router.go file: %s", err)
	}

	fileList := strings.Split(currentPath, "src")
	if len(fileList[1]) == 0 {
		log.Fatalln("you create file in src directory")
	}

	fileList[1] = strings.Replace(fileList[1], string(filepath.Separator), "/", -1)
	fileList[1] = fileList[1][1:]
	routerTemp = strings.Replace(routerTemp, "{{controllers}}", fileList[1], -1)
	routerTemp = strings.Replace(routerTemp, "{{utils}}", fileList[1], -1)
	_, _ = f.WriteString(routerTemp)
	_ = f.Close()
	cmd := exec.Command("gofmt", "-w", fpath)
	_ = cmd.Run()
}