package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Property struct {
	Key   string
	Value string
}

var (
	zookeeper    *string
	namespace    *string
	templatePath *string
	command      *string
	aSync	     *bool

	zookeeperProperties []Property
	eventChannel        <-chan zk.Event
	child               []string
)

func flags() {
	zookeeper = flag.String("zookeeper", "192.168.120.81:2181,192.168.120.82:2181", "Zookeeper server list example: 192.168.120.1:2181,192.168.120.2:2181,..")
	namespace = flag.String("namespace", "/producer", "Namespace to watch on")
	templatePath = flag.String("template", "/etc/zkwatcher/tmp.txt", "Template absolut path")
	command = flag.String("cmd", "", "Command execute after regenerate config")
	aSync = flag.Bool("aSync", false, "Asyncron command execution")

	flag.Parse()
}

func main() {
	flags()

	splittedCommand := strings.Split(*command, " ")
	commandExistence := len(splittedCommand) > 1

	templateFile := template.Must(template.New("tmp.txt").ParseFiles(*templatePath))

	conn, _, _ := zk.Connect(strings.Split(*zookeeper, ","), time.Second)
	defer conn.Close()

	log.Printf("connected to zookeepers: %v", *zookeeper)
	log.Printf("Watch on: %s", *namespace)

	for true {
		zookeeperProperties = []Property{}
		_, _, eventChannel, _ = conn.ChildrenW(*namespace)

		<-eventChannel
		child, _, eventChannel, _ = conn.ChildrenW(*namespace)
		for _, key := range child {
			value, _, _ := conn.Get(fmt.Sprintf("%s/%s", "/producer", key))
			if string(value) != "" { // bypass directory
				zookeeperProperties = append(zookeeperProperties, Property{
					Key:   key,
					Value: string(value),
				})
			}
		}

		log.Printf("Regenerate config from %s", *templatePath)
		err := templateFile.Execute(os.Stdout, zookeeperProperties)
		if err != nil {
			log.Panic(err.Error())
		}

		if commandExistence {
			log.Printf("Execute command %s", *command)
			cmd := exec.Command(splittedCommand[0], splittedCommand[1:]...)

			if *aSync {
				err = cmd.Start()
			} else {
				err = cmd.Run()
			}

			if err != nil {
				log.Fatal(cmd.Stderr)
			}

			log.Print(cmd.Stdout)
		}
	}
}
