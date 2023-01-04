package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sumaprom/collector"
	"sumaprom/crypt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	usage = `This is SUMA prometheus exporter running on port :8888
	usage: %s -sumaconf <sumaconf.yaml>

	You need to provide a yaml file with SUMA api login data.
	server: <YOUR-SUMA-SERVER>
	user: admin
	password: <encrypted-PASSWD>

	Use "-create-sumaconf <path-to-sumaconf-file>" to create a config file with encrypted password.

Options:
`
)

var (
	sumaconf crypt.Sumaconf
	step     string
)

func init() {

	var conf_file = flag.String("sumaconf", "", "provide the suma conf file with login data.")
	var create_config = flag.String("create-sumaconf", "", "Create a config file with login data.")
	flag.Usage = func() { // [1]
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(*conf_file) == 0 && len(*create_config) == 0 {
		log.Fatal("sumaconf not provided or create a new sumaconf. Exit")
		step = "exit"
	} else if len(*conf_file) == 0 && len(*create_config) > 0 {
		crypt.PromptUser()
		step = "create_config"
	} else {
		sumaconf = sumaconf.Decrypt_Sumaconf(conf_file)
		fmt.Printf("server: %s\nuserid: %s\npassword: %s\n", sumaconf.Server, sumaconf.Userid, sumaconf.Password)
		step = "start_collect"
	}

}

func main() {

	if step == "start_collect" {
		fmt.Printf("aaa %s %s %s\n", sumaconf.Server, sumaconf.Userid, sumaconf.Password)
		MypkgsCollector := new(collector.PkgsCollector)
		MypkgsCollector.Sumainfo = sumaconf
		prometheus.MustRegister(MypkgsCollector)

		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8888", nil))
	}

}
