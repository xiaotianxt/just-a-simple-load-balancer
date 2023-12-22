package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"

	"gopkg.in/yaml.v2"
)

type ServerPool struct {
	authList []string
	current  uint64
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.authList)))
}

func (s *ServerPool) GetNextAuth() string {
	return s.authList[s.NextIndex()]
}

func (s *ServerPool) AddAuth(auth string) {
	s.authList = append(s.authList, auth)
}

type Config struct {
	AuthList []string `yaml:"authList"`
	Host     string   `yaml:"host"`
}

func main() {
	serverPool := &ServerPool{}

	// Read config
	data, err := os.ReadFile("/app/config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	println(config.AuthList)

	// Add auths
	for _, auth := range config.AuthList {
		serverPool.AddAuth(auth)
		println("Add auth: " + auth)
	}

	backend := &url.URL{
		Scheme: "https",
		Host:   config.Host,
	}
	println("Connecting to: " + config.Host)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := serverPool.GetNextAuth()

		println("Using auth: " + auth)

		director := func(req *http.Request) {
			req.URL.Scheme = backend.Scheme
			req.URL.Host = backend.Host
			req.Header["X-Forwarded-Host"] = req.Header["Host"]
			req.Host = backend.Host
			req.Header["Authorization"] = []string{auth}
		}

		proxy := &httputil.ReverseProxy{Director: director,
			ModifyResponse: func(resp *http.Response) error {
				if resp.StatusCode != http.StatusOK {
					// print resp body json
					body, err := httputil.DumpResponse(resp, true)
					if err != nil {
						log.Fatalln("Fatal error ", err.Error())
					} else {
						log.Println("Error: ", resp.StatusCode)
						log.Println(string(body))
					}
				}
				return nil
			},
		}
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8088", nil)
}
