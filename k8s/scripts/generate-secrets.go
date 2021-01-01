package main

import (
	"flag"
	"github.com/thanhpk/randstr"
	"io/ioutil"
	"os"
	"text/template"
	"encoding/base64"
)

type Replacements struct {
	Password1 string
	Password2 string
}

func main() {
	sourceFile := flag.String("source", "", "source file")
	destFile := flag.String("dest", "", "destination file")
	flag.Parse()
	if sourceFile == nil || destFile == nil || len(*sourceFile) == 0 || len(*destFile) == 0 {
		panic("Both --source and --dest arguments required")
	}
	content, err := ioutil.ReadFile(*sourceFile)
	if err != nil {
		panic(err)
	}
	// a-zA-Z0-9 plus special characters that are easy to copy and paste in bash.
	letters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@_~#+-$%^."
	replacements := Replacements{
		base64.StdEncoding.EncodeToString([]byte(randstr.String(32, letters))),
		base64.StdEncoding.EncodeToString([]byte(randstr.String(32, letters))),
	}
	tpl, err := template.New("secrets-template").Parse(string(content))
	if err != nil {
		panic(err)
	}
	out, err := os.Create(*destFile)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(out, replacements)
	if err != nil {
		panic(err)
	}
	err = out.Close()
	if err != nil {
		panic(err)
	}
}