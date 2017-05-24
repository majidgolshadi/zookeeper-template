package main

import (
	"log"
	"os/exec"
)

type Command struct {
	Cmd   string
	Args  []string
	Async bool
}

var err error

func (c *Command) Execute() {
	if c.Cmd == "" {
		return
	}

	cmd := exec.Command(c.Cmd, c.Args...)

	if c.Async {
		err = cmd.Start()
	} else {
		err = cmd.Run()
	}

	if err != nil {
		log.Fatal(cmd.Stderr)
	}

	log.Print(cmd.Stdout)
}
