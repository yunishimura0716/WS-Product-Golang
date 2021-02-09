package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type counters struct {
	sync.Mutex
	view  int
	click int
}

var (
	c = counters{}

	content = []string{"sports", "entertainment", "business", "education"}

	counterLogs = map[string]map[string]int{
		"sports": {"views": 0, "clicks": 0},
		"entertainment": {"views": 0, "clicks": 0},
		"education": {"views": 0, "clicks": 0},
		"business": {"views": 0, "clicks": 0},
	}

	store []map[string]map[string]int

	timeToSleep = 5
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]

	c.Lock()
	c.view++
	c.Unlock()

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data)
	}

	// save counter
	// counter := map[string]int{
	// 	"views": c.view,
	// 	"clicks": c.click,
	// }
	counterLogs[data]["views"] = c.view
	counterLogs[data]["clicks"] = c.click
	fmt.Printf("%s views:%d, clicks:%d\n", data, counterLogs[data]["views"], counterLogs[data]["clicks"])
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click++
	c.Unlock()

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}
}

func isAllowed() bool {
	return true
}

func counterStoring() {
	for true {
		sleep()
		uploadCounters()
	}
}

func uploadCounters() error {
	t := time.Now()
	currCounters := make(map[string]map[string]int)
	for i:= 0; i < len(content); i++ {
		data := content[i]
		key := data + ":" + t.Format("2006-01-02 15:04")
		fmt.Printf("%s views:%d, clicks:%d\n", key, counterLogs[data]["views"], counterLogs[data]["clicks"])

		// fmt.Printf("!")
		var m = map[string]int{
			"views": counterLogs[data]["views"],
			"clicks": counterLogs[data]["clicks"],
		}
		currCounters[key] = m
	}
	store = append(store, currCounters)
	return nil
}

// sleep 5 seconds
func sleep() {
	time.Sleep(time.Duration(timeToSleep) * time.Second)
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)
	go counterStoring()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
