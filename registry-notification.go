package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"
)

// Configuration represents configuration structure of the server
type Configuration struct {
	Port      int    `json:"port"`
	ServerKey string `json:"serverkey"`
	ServerCrt string `json:"servercrt"`
}

// Config is instance of Configruation
var Config Configuration

// version represents version of the server
var version string

// RequestHandler handles incoming HTTP request
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, r.Proto, r.Host, r.RemoteAddr, r.Header)
	if r.Method == "GET" {
		// print out all request headers
		fmt.Fprintf(w, "%s %s %s \n", r.Method, r.URL, r.Proto)
		for k, v := range r.Header {
			h := strings.ToLower(k)
			if strings.Contains(h, "hmac") || strings.Contains(h, "cookie") {
				continue
			}
			fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
		}
		fmt.Fprintf(w, "Host = %q\n", r.Host)
		fmt.Fprintf(w, "RemoteAddr= %q\n", r.RemoteAddr)
		fmt.Fprintf(w, "\n\nFinding value of \"Accept\" %q\n", r.Header["Accept"])

		page := "Hello from Go\n"
		w.Write([]byte(page))
	} else {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprint(w, string(requestDump))
		}
	}
}

// helper function to parse the config
func parseConfig(configFile string) error {
	if configFile == "" {
		Config.Port = 9215
		return nil
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// helper function to return version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("httpgo git=%s go=%s date=%s", version, goVersion, tstamp)
}

// main function
func main() {
	var config string
	flag.StringVar(&config, "config", "", "configuration file")
	var version bool
	flag.BoolVar(&version, "version", false, "print version information about the server")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)
	}
	// log time, filename, and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := parseConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", RequestHandler)
	if Config.ServerKey != "" && Config.ServerCrt != "" {
		server := &http.Server{
			Addr: fmt.Sprintf(":%d", Config.Port),
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
				//             ClientAuth: tls.RequestClientCert,
			},
		}
		err = server.ListenAndServeTLS(Config.ServerCrt, Config.ServerKey)
		if err != nil {
			fmt.Println("Unable to start the server", err)
		}
	} else {
		http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil)
	}
}
