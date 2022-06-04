package es

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	EsDump *esDump
)

type esDump struct {
	Client  *elasticsearch.Client
	index   string
	size    int
	query   string
	complex bool
}

func NewEsDump(urlStr, query string, size int, complexFlag bool) (EsDump *esDump) {
	EsDump = new(esDump)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}
	if query != "" {
		EsDump.query = query
	}

	EsDump.size = size * 10

	var user, password string
	passwordStr, hasPwd := u.User.Password()
	if hasPwd {
		if passwordStr == "" {
			log.Println("password can not be empty.")
			os.Exit(1)
		}
		password = passwordStr
		user = u.User.Username()
	}

	path := strings.Split(u.Path, "/")
	if path[1] == "" {
		log.Println("index can not be empty.")
		os.Exit(1)
	}
	EsDump.index = path[1]
	EsDump.complex = complexFlag

	var host string
	switch {
	case strings.ToUpper(u.Scheme) == "HTTPS":
		host = u.Scheme + "://" + u.Host
	case strings.ToUpper(u.Scheme) == "HTTP":
		host = u.Scheme + "://" + u.Host
	default:
		log.Println("only support http or https protocol")
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	config := elasticsearch.Config{}
	if user != "" && password != "" {
		config.Addresses = []string{host}
		config.Username = user
		config.Password = password
		config.Transport = tr
		Es, err := elasticsearch.NewClient(config)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		EsDump.Client = Es
	} else if user == "" && password == "" {
		config.Addresses = []string{host}
		Es, err := elasticsearch.NewClient(config)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		EsDump.Client = Es
	}
	return EsDump
}
