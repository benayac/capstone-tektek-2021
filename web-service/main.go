package main

import (
	"capstone/middleware"
	"capstone/repository"
	"capstone/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	host := "database"
	port := 5432
	user := "user"
	password := "password"
	dbName := "db"
	db, err := repository.GetConnection(host, port, user, password, dbName)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error Load DB: %v", err.Error()))
	}
	reportService := services.ReportService{DB: db}
	scoreService := services.ScoreService{DB: db}
	r := mux.NewRouter()

	r.HandleFunc("/report", reportService.InsertReportHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/report/new", reportService.InsertNewReportHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/score", scoreService.GetScoreFromIdHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/score/update", scoreService.UpdateScoreHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/score/new", scoreService.InsertScoreHandler).Methods("GET", "OPTIONS")

	r.Use(middleware.DefaultHeader)
	portServer := os.Getenv("PORT")
	if portServer == "" {
		portServer = "8080" // Default port if not specified
	}
	fmt.Printf("Running Web Service on Port %v\n", portServer)
	log.Fatal(http.ListenAndServe(":"+portServer, r))
}
