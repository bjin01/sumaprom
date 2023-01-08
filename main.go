package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
		crypt.PromptUser(create_config)
		step = "create_config"
	} else {
		sumaconf = sumaconf.Decrypt_Sumaconf(conf_file)
		step = "start_collect"
	}

}

func main() {

	if step == "start_collect" {
		MypkgsCollector := new(collector.PkgsCollector)
		MypkgsCollector.Sumainfo = sumaconf

		var port int
		var err error
		if os.Getenv("SUMAPROM_PORT") != "" {
			port, err = strconv.Atoi(os.Getenv("SUMAPROM_PORT"))
			if err != nil {
				log.Fatalf("Failed to get SUMAPROM_PORT number. %s\n", os.Getenv("SUMAPROM_PORT"))
			}
			fmt.Printf("port num: %d\n", port)
			if port < 1024 || port > 30000 {
				log.Fatalf("The port %d you selected is not good. Choose a port between 1024 and 30000 on linux.\n", port)
			}
		} else {
			port = 8888
		}

		prometheus.MustRegister(MypkgsCollector)
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

	}

}
