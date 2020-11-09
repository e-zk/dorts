package functions

import (
	"fmt"
	"os/exec"
	"regexp"
	"text/template"
)

func xrdbGrep(reg string) string {
	xrdbCmd := exec.Command("xrdb", "-query")

	output, err := xrdbCmd.Output()
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(reg)
	match := r.FindStringSubmatch(fmt.Sprintf("%s", output))

	return match[1]
}

func GetFuncs() template.FuncMap {
	return template.FuncMap{
		"XrdbGrep": xrdbGrep,
	}
}
