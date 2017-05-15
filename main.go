package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"fmt"
	"flag"
	"strings"
)

var (
	zookeeper *string
	namespace *string
	templatePath *string

	child []string
)

func flags()  {
	zookeeper = flag.String("zookeeper", "127.0.0.1", "Zookeeper server list example: 192.168.120.1:2181,192.168.120.2:2181,..")
	namespace = flag.String("namespace", "/", "Namespace to watch on")
	namespace = flag.String("template", "/etc/zkreloader/tmp.txt", "Template absolut path")

	flag.Parse()
}

func main() {
	flags()

	conn, _, err := zk.Connect(strings.Split(*zookeeper, ","), time.Second)
	defer conn.Close()
	_, _, ech, err := conn.ChildrenW(*namespace)
	if err != nil {
		panic(err.Error())
	}

	for true {
		for event := range ech {
			child, _, ech, _ = conn.ChildrenW(event.Path)
			for _, key := range child {
				value, _ ,_ := conn.Get(fmt.Sprintf("%s/%s",event.Path, key))
			}
		}
	}
}
