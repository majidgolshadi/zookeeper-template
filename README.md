Zkwatcher
---------
zkwatcher is a zookeeper watcher that watch on a namespace. On receive create or delete event it can regenerate your configuration file (based on defined [template](https://golang.org/pkg/text/template/)) and/or run command sync/async

Usage
-----
```bash
zkwatcher Usage:
    -aSync
          Asyncron command execution (default false)
    -cmd string
          Command execute after regenerate config
    -destConf string
          Generated config absolut path
    -namespace string
          Namespace to watch on (default "/watch")
    -template string
          srcTemplate absolut path
    -zookeeper string
          Zookeeper server list example: 192.168.120.1:2181,192.168.120.2:2181,.. (default "192.168.120.81:2181,192.168.120.82:2181")
```

Additional functions
--------------------
**method:** `plus1 $i`
**return:** $i+1
**usage:**
```bash
[{{$n := len .}}{{range  $i, $e := .}}{{if $i}}, {{end}}{{if eq (plus1 $i) $n}}{{end}}"{{.Value}}"{{end}}].
```