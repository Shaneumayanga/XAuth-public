package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Shaneumayanga/XAuth/api"
)

//website :~ shaneumayanga.com
//email :~ me@shaneumayanga.com
//hosted on :~ xauth.shaneumayanga.com

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Addr:    ":" + port,
		Handler: api.Start(),
	}

	go func() {
		fmt.Println("Server running on port :8080")
		log.Fatal(server.ListenAndServe())
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Server shutdown started")
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	server.Shutdown(ctx)
	fmt.Println("Server gracefully shutted down")

}
