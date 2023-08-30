package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lrstanley/girc"
)

const logFilename = "logs/irc_log.txt"

func main() {
	channel := os.Getenv("IRC_CHANNEL")
	server := os.Getenv("IRC_SERVER")
	nickname := os.Getenv("IRC_NICKNAME")
	username := os.Getenv("IRC_USERNAME")
	name := os.Getenv("IRC_NAME")
	timezone := os.Getenv("IRC_TIMEZONE")
	httpUsername := os.Getenv("HTTP_USERNAME")
	httpPassword := os.Getenv("HTTP_PASSWORD")
	httpPort := os.Getenv("HTTP_PORT")

	if channel == "" || server == "" || nickname == "" || username == "" || name == "" || timezone == "" || httpUsername == "" || httpPassword == "" {
		log.Fatalf("Environment variables IRC_CHANNEL, IRC_SERVER, IRC_NICKNAME, IRC_USERNAME, IRC_NAME, IRC_TIMEZONE, HTTP_USERNAME, and HTTP_PASSWORD must be set")
	}

	client := girc.New(girc.Config{
		Server: server,
		Port:   6667,
		Nick:   nickname,
		User:   username,
		Name:   name,
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		message := logMessageWithTimezone(e.Source.Name, e.Last(), timezone)
		appendToFile(logFilename, message+"\n")
	})

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(channel)
	})

	http.HandleFunc("/", basicAuth(displayLogs, httpUsername, httpPassword, "Please enter your username and password"))
	go http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)

	// Handle reconnections
	for {
		if err := client.Connect(); err != nil {
			log.Printf("Failed to connect: %s\n", err)

			var sleepTime time.Duration = 30
			log.Printf("Reconnecting in %d seconds...\n", sleepTime*time.Second)
			time.Sleep(sleepTime * time.Second)
		} else {

			log.Println("Quitting")
			return
		}
	}
}

func basicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

const logTimeFormat = "2006-01-02 15:04:05"

func logMessageWithTimezone(nick, message, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Error loading timezone: %v", err)
		loc = time.UTC
	}

	currentTime := time.Now().In(loc).Format(logTimeFormat)
	return fmt.Sprintf("(%s) %s: %s", currentTime, nick, message)
}

func appendToFile(filename, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		log.Println(err)
	}
}

func displayLogs(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(logFilename)
	if err != nil {
		http.Error(w, "Failed to read logs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, string(content))
}
