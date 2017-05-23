package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	srcTemplate *string
	destConf     *string
	command      *string
	aSync	     *bool

	zookeeperProperties []Property
	eventChannel        <-chan zk.Event
	child               []string

	commandExistence bool
)

func flags() {
	zookeeper = flag.String("zookeeper", "192.168.120.81:2181,192.168.120.82:2181", "Zookeeper server list example: 192.168.120.1:2181,192.168.120.2:2181,..")
	namespace = flag.String("namespace", "/watch", "Namespace to watch on")
	srcTemplate = flag.String("srcTemplate", "", "srcTemplate absolut path")
	destConf = flag.String("destConf", "", "Generated config absolut path")
	aSync = flag.Bool("aSync", false, "Asyncron command execution")
	command = flag.String("cmd", "", "Command execute after regenerate config")

	flag.Parse()
}

func main() {
	flags()

	if *srcTemplate != "" && *destConf == "" {
		println("destConf option must be set")
		flag.Usage()
		os.Exit(2)
	}

	if *srcTemplate == "" && *destConf != "" {
		println("srcTemplate option must be set")
		flag.Usage()
		os.Exit(2)
	}

	if *command == "" && *srcTemplate == "" {
		println("At least you have to set srcTemplate path or a command to execute after change event")
		flag.Usage()
		os.Exit(2)
	}

	splittedCommand := strings.Split(*command, " ")
	cmd := &Command{
		Cmd: splittedCommand[0],
		Args: splittedCommand[1:],
		Async: *aSync,
	}


	conf := &Config{
		Template: *srcTemplate,
	}
	conf.Init()

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
			value, _, _ := conn.Get(fmt.Sprintf("%s/%s", *namespace, key))
			if string(value) != "" { // bypass directory
				zookeeperProperties = append(zookeeperProperties, Property{
					Key:   key,
					Value: string(value),
				})
			}
		}

		if *srcTemplate != "" {
			conf.GenerateConfig(zookeeperProperties)
		}
		
		if commandExistence {
			cmd.Execute()
		}
	}
}
