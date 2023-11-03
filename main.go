package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

// Alert represents the structure of an alert.
type Alert struct {
	AlertID     string `json:"alert_id"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Model       string `json:"model"`
	AlertType   string `json:"alert_type"`
	AlertTS     string `json:"alert_ts"`
	Severity    string `json:"severity"`
	TeamSlack   string `json:"team_slack"`
}

type AlertResponse struct {
	AlertID   string `json:"alert_id"`
	Model     string `json:"model"`
	AlertType string `json:"alert_type"`
	AlertTS   string `json:"alert_ts"`
	Severity  string `json:"severity"`
	TeamSlack string `json:"team_slack"`
}

var db *sql.DB

func main() {
	r := mux.NewRouter()

	// Initialize the SQLite database
	initializeDatabase()

	// POST endpoint for writing alerts
	r.HandleFunc("/alerts", WriteAlert).Methods("POST")

	// GET endpoint for reading alerts by service_id and time range
	r.HandleFunc("/alerts", ReadAlerts).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Initialize the SQLite database.
func initializeDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "myproject.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the alerts table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS alerts (
		alert_id TEXT PRIMARY KEY,
		service_id TEXT NOT NULL,
		service_name TEXT,
		model TEXT,
		alert_type TEXT,
		alert_ts TEXT,
		severity TEXT,
		team_slack TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

// WriteAlert handles the HTTP POST request to write an alert to the database.
func WriteAlert(w http.ResponseWriter, r *http.Request) {
	var alert Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	response := make(map[string]string)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store the alert in the database
	err = storeAlertInDatabase(alert)
	if err != nil {
		response["alert_id"] = alert.AlertID
		response["error"] = err.Error()
		w.WriteHeader(500)
		errorResponse, _ := json.Marshal(response)
		w.Write(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	response["alert_id"] = alert.AlertID
	response["error"] = ""

	json.NewEncoder(w).Encode(response)
}

// Store an alert in the SQLite database.
func storeAlertInDatabase(alert Alert) error {
	_, err := db.Exec(`INSERT INTO alerts (alert_id, service_id, service_name, model, alert_type, alert_ts, severity, team_slack)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, alert.AlertID, alert.ServiceID, alert.ServiceName, alert.Model, alert.AlertType, alert.AlertTS, alert.Severity, alert.TeamSlack)
	return err
}

// ReadAlerts handles the HTTP GET request to read alerts by service_id and time range.
func ReadAlerts(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service_id")
	startTS := r.URL.Query().Get("start_ts")
	endTS := r.URL.Query().Get("end_ts")

	alerts, err, serviceName := getAlertsByServiceAndTime(serviceID, startTS, endTS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(startTS) == 0 || len(endTS) == 0 || len(serviceID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"alert_id": "error",
			"error":    "Please provide service ID, starting and ending Timestamp properly",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(alerts) > 0 {
		w.WriteHeader(http.StatusOK)
		response := make(map[string]interface{})
		response["service_id"] = serviceID
		response["service_name"] = serviceName
		response["alerts"] = alerts

		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"alert_id": "error",
			"error":    "No alerts found for the provided parameters",
		}
		json.NewEncoder(w).Encode(response)
	}
}

// Retrieve alerts by service_id and time range from the SQLite database.
func getAlertsByServiceAndTime(serviceID, startTS, endTS string) ([]AlertResponse, error, string) {
	var alerts []AlertResponse
	serviceName := ""
	rows, err := db.Query(`SELECT alert_id, model, alert_type, alert_ts, severity, team_slack, service_name FROM alerts WHERE service_id = ? AND alert_ts >= ? AND alert_ts <= ?`, serviceID, startTS, endTS)
	if err != nil {
		return nil, err, serviceName
	}
	defer rows.Close()

	for rows.Next() {
		var alert AlertResponse
		err := rows.Scan(&alert.AlertID, &alert.Model, &alert.AlertType, &alert.AlertTS, &alert.Severity, &alert.TeamSlack, &serviceName)
		if err != nil {
			return nil, err, serviceName
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil, serviceName
}
