package main

import (
	"text/template"
	"path"
	"strings"
	"log"
	"os"
)

type Config struct {
	Template string
	Dest string
	templateFile string
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
}

func (c *Config) GenerateConfig(data []Property) {
	t := template.Must(template.New(c.templateFile).Funcs(fns).ParseFiles(c.Template))

	log.Printf("Regenerate config from %s", *srcTemplate)
	err := t.Execute(os.Stdout, data); if err != nil {
		log.Panic(err.Error())
	}
}
