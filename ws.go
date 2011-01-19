package main;

import (
	"./statmsg"
	"./nullstore"
	"runtime"
	"os"
	"io/ioutil"
	"strings"
	"fmt"
	"http"
)

// Channel size 
const chanBacklog = 1000

func loadRedirects(filename string) (redirects map[string] string) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Can't read the map file!")
	}

	redirects = make(map[string] string)
	lines := strings.Split(string(contents), "\n", -1)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		s := strings.Split(line, " ", 2)
		redirects[s[0]] = s[1]
	}

	return redirects
}

func statsUpdate(statChan chan *statmsg.Statmsg) {
	for {
		stat := <- statChan
		go nullstore.Update(stat)
	}
}

func makeRedirectServer(redirects map[string] string,
	statChan chan *statmsg.Statmsg) (http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request) {
		key := req.URL.Path[1:]
		url, exists := redirects[key]

		if !exists {
			http.NotFound(w, req)
			return
		}

		w.SetHeader("Location", url)
		w.WriteHeader(http.StatusMovedPermanently)

		var stat statmsg.Statmsg
		stat.Key = key
		stat.IP = w.RemoteAddr()
		stat.Referer = req.Referer
		stat.UA = req.UserAgent
		statChan <- &stat
	}
}

func main() {
	runtime.GOMAXPROCS(8)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s redirect_map_file\n",
			os.Args[0])
		return
	}

	fmt.Printf("Loading redirects map...\n")
	redirects := loadRedirects(os.Args[1])

	statChan := make(chan *statmsg.Statmsg, chanBacklog)
	go statsUpdate(statChan)

	fmt.Printf("Starting web server...\n")
	http.HandleFunc("/", makeRedirectServer(redirects, statChan))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't start: %s\n", err.String())
		return
	}
}
