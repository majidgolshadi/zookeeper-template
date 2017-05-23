Zkwatcher
---------
zkwatcher is a zookeeper watcher that try to watch on single namespace and regenerate your configuration file on change
 after that you can run single command sync or async

Usage
-----
```bash
zkwatcher Usage:
  -aSync
        Asyncron command execution
  -cmd string
        Command execute after regenerate config
  -namespace string
        Namespace to watch on (default "/watch")
  -template string
        Template absolut path (default "/etc/zkwatcher/tmp.txt")
  -zookeeper string
        Zookeeper server list example: 192.168.120.1:2181,192.168.120.2:2181,.. (default "192.168.120.81:2181,192.168.120.82:2181")
```