package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "root"        // Replace with your MySQL username
	dbPassword = "dragon"      // Replace with your MySQL password
	dbName     = "time_log_db" // Replace with your database name

)

var db *sql.DB
var dbHost = os.Getenv("DB_HOST")

// Initialize a log file
func initLogger() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Logging initialized")
}
func initDB() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	log.Println("Connected to MySQL database")
}

func getCurrentTime(w http.ResponseWriter, r *http.Request) {
	// Get current time in UTC
	//currentTime := time.Now().UTC()
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Loading time zone - failed", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now().In(location)

	// Define the time format (e.g., "2006-01-02 15:04:05")
	timeFormat := "2006-01-02 15:04:05"
	torontoTime := currentTime.Format(timeFormat)
	// Convert to Toronto time
	//loc, err := time.LoadLocation("America/Toronto")
	// if err != nil {
	// 	http.Error(w, "Failed to load Toronto timezone", http.StatusInternalServerError)
	// 	log.Printf("Error loading timezone: %v", err)
	// 	return
	// }
	//torontoTime := currentTime.In(loc)

	// Insert into MySQL
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", torontoTime)
	if err != nil {
		http.Error(w, "Failed to log time to database", http.StatusInternalServerError)
		log.Printf("Error inserting timestamp: %v", err)
		return
	}

	// Respond with JSON
	response := map[string]string{
		"current_time": torontoTime,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAllTimes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, timestamp FROM time_log")
	if err != nil {
		http.Error(w, "Failed to retrieve times from database", http.StatusInternalServerError)
		log.Printf("Error fetching times: %v", err)
		return
	}
	defer rows.Close()

	var times []map[string]string
	for rows.Next() {
		var id int
		var timestamp string
		if err := rows.Scan(&id, &timestamp); err != nil {
			http.Error(w, "Failed to parse database results", http.StatusInternalServerError)
			log.Printf("Error scanning row: %v", err)
			return
		}
		times = append(times, map[string]string{
			"id":        fmt.Sprintf("%d", id),
			"timestamp": timestamp,
		})
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(times)
	log.Printf("Fetched %d time entries from database", len(times))
}

func main() {
	// Initialize database and logger
	initLogger()
	initDB()
	defer db.Close()

	// Set up HTTP server
	http.HandleFunc("/current-time", getCurrentTime)
	http.HandleFunc("/all-times", getAllTimes)

	port := ":80"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
