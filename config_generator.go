package main

import (
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

type Config struct {
	Template     string
	Dest         string
	templateFile string

	file *os.File
}

var fns = template.FuncMap{
	"plus1": func(x int) int {
		return x + 1
	},
}

func (c *Config) Init() {
	c.templateFile = strings.Replace(
		path.Base(c.Template),
		c.Template,
		"", 1)

	c.file, err = os.OpenFile(c.Dest, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (c *Config) GenerateConfig(data []Property) {
	if c.Template == "" {
		return
	}

	t := template.Must(template.New(c.templateFile).Funcs(fns).ParseFiles(c.Template))

	log.Printf("Regenerate config from %s", *srcTemplate)
	err := t.Execute(c.file, data)
	if err != nil {
		log.Panic(err.Error())
	}
}

func (c *Config) Close() {
	c.file.Close()
}
