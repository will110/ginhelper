package main

import (
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/will110/ginhelper/project"
)

const _version = "v1.1.1"

func main() {
	flag.Parse()
	arg := flag.Args()
	if len(arg) == 0 {
		getHelp()
		return
	}

	switch arg[0] {
	case "create":
		project.GenerateProject()
	case "version":
		fmt.Println("ginhelper " + _version)
	case "help":
		getHelp()
	default:
		getHelp()
	}
}

func getHelp() {
	title :=aurora.White("Ginhelper is a gin web framework helper, Help you develop fast.")
	usage := aurora.Magenta("USAGE")
	t1 := aurora.BrightWhite("ginhelper [arguments]")
	t2 := aurora.Underline(aurora.Magenta("AVAILABLE ARGUMENTS"))
	t3 := aurora.BrightWhite("create       ")
	t4 := aurora.Underline(aurora.BrightWhite("Create single gin web framework"))
	t5 := aurora.BrightWhite("version      ")
	t6 := aurora.Underline(aurora.BrightWhite("Prints the current GinHelper version"))
	t7 := aurora.Underline(aurora.BrightWhite("Use ginhelper help for more information."))

	buf := fmt.Sprintf("%v\n\n%v\n    %v\n\n%v\n\n    %v%v\n    %v%v\n\n\n\n%v", title, usage, t1, t2, t3, t4, t5, t6, t7)
	fmt.Println(buf)
}