package main

import (
	"assignments-my828/servers/gateway/handlers"
	"assignments-my828/servers/gateway/indexes"
	"assignments-my828/servers/gateway/models/users"
	"assignments-my828/servers/gateway/sessions"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

//main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
		the server should listen on. If empty, default to ":80"

	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	tlsKeyPath := os.Getenv("TLSKEY")
	if tlsKeyPath == "" {
		fmt.Printf("No private key path!")
	}
	tlsCertPath := os.Getenv("TLSCERT")
	if tlsCertPath == "" {
		fmt.Printf("No certificate found!")
	}

	sessionKey := os.Getenv("SESSIONKEY")
	if sessionKey == "" {
		sessionKey = "sessionkey"
	}

	//the address of your redis session store server
	redisAddr := os.Getenv("REDISADDR")
	if redisAddr == "" {
		fmt.Printf("No redis address found!")
	}

	// the full data source name to
	// pass as the second parameter to sql.Open()
	// rootPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
	// dsn := fmt.Sprintf("root:%s@tcp(users:3306)/userinfo", rootPassword)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	sessionStore := sessions.NewRedisStore(client, time.Hour)

	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("No sql found!")
		os.Exit(1)
	}
	trie := indexes.NewTrie()
	defer db.Close()
	userStore := users.NewSQLStore(db)
	mux := http.NewServeMux()

	context := &handlers.Context{
		Key:          sessionKey,
		SessionStore: sessionStore,
		UsersStore:   userStore,
		SearchIndex:  trie,
	}

	// connect to redis
	if _, err := client.Ping().Result(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
	}
	// connect to mysql
	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
	}

	messageAddr := strings.Split(os.Getenv("MESSAGEADDR"), ",")
	summaryAddr := strings.Split(os.Getenv("SUMMARYADDR"), ",")
	messageAddrs := []*url.URL{}
	summaryAddrs := []*url.URL{}
	for _, addr := range messageAddr {
		parseAddr, err := url.Parse(addr)
		if err != nil {
			fmt.Printf("error parsing message address: %v\n", err)
		}
		messageAddrs = append(messageAddrs, parseAddr)
	}

	for _, addr := range summaryAddr {
		parseAddr, err := url.Parse(addr)
		if err != nil {
			fmt.Printf("error parsing message address: %v\n", err)
		}
		summaryAddrs = append(summaryAddrs, parseAddr)
	}

	messageProxy := &httputil.ReverseProxy{Director: handlers.CustomDirector(messageAddrs)}
	summaryProxy := &httputil.ReverseProxy{Director: handlers.CustomDirector(summaryAddrs)}

	mux.Handle("/v1/messages/:messageID", messageProxy)
	mux.Handle("/v1/channels/:channelID/members", messageProxy)
	mux.Handle("/v1/channels/:channelID", messageProxy)
	mux.Handle("/v1/channel", messageProxy)

	mux.Handle("/v1/summary", summaryProxy)

	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)

	//wrap new mux with CORS middleware handler
	wrappedMux := handlers.NewCorsHandler(mux)

	log.Printf("!!!!!!server is listening at http://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
