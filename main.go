package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/StephanUllmann/group-randomizer-api/config"
	"github.com/StephanUllmann/group-randomizer-api/controllers"
	"github.com/StephanUllmann/group-randomizer-api/db"
)

func main() {
	var app config.AppConfig
	var err error

	app.DB, err = db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	// 'IF NOT EXISTS' doesn't work
	// err = db.InitDB(app.DB)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	repo := controllers.NewRepository(&app)
	controllers.NewHandler(repo)

	server := http.NewServeMux()

	server.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.Method)
		w.Write([]byte("Hello"))
	})
	server.HandleFunc("/api", controllers.Repo.RouteGroups)



	port := ":"+ os.Getenv("PORT")
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	log.Fatalf("hostname not acquired, %v", err)
	// }

	fmt.Printf("listening on http://localhost%s\n", port)
	err = http.ListenAndServe(port, server)
	if err != nil {
		log.Fatal(err)
	}
}