package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
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

	// GET endpoint for Ariana Grande UI page
	r.HandleFunc("/ariana", ArianaPageHandler).Methods("GET")

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Initialize the SQLite database.
func initializeDatabase() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "myproject.db"
	}
	db, err = sql.Open("sqlite3", dbPath)
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

// ArianaPageHandler serves a simple HTML UI page for Ariana Grande and her albums.
func ArianaPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Ariana Grande — Albums</title>
  <style>
    :root { --bg: #0b0b10; --card: #151520; --text: #f4f4f8; --muted: #b9b9c6; --accent: #9b81ff; --chip: #222232; }
    * { box-sizing: border-box; }
    body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Inter, Helvetica, Arial, sans-serif; background: radial-gradient(1200px 800px at 20% -10%, #1d1d2b 0%, #0b0b10 45%), #0b0b10; color: var(--text); }
    a { color: inherit; text-decoration: none; }
    .container { max-width: 1100px; margin: 0 auto; padding: 32px 20px 64px; }
    .header { display: flex; align-items: center; gap: 20px; }
    .avatar { width: 84px; height: 84px; border-radius: 18px; background: linear-gradient(145deg, #1b1b29, #0e0e16); box-shadow: 0 12px 24px rgba(0,0,0,0.4), inset 0 0 0 1px #2a2a3d; display: grid; place-items: center; font-weight: 800; letter-spacing: 0.5px; color: var(--accent); }
    .title { display: flex; flex-direction: column; gap: 6px; }
    .title h1 { margin: 0; font-size: 28px; }
    .title .meta { color: var(--muted); font-size: 14px; }
    .chips { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 10px; }
    .chip { background: var(--chip); color: var(--muted); font-size: 12px; padding: 6px 10px; border-radius: 999px; border: 1px solid #2a2a3d; }
    .grid { display: grid; grid-template-columns: repeat(1, minmax(0, 1fr)); gap: 18px; margin-top: 26px; }
    @media (min-width: 560px) { .grid { grid-template-columns: repeat(2, 1fr); } }
    @media (min-width: 900px) { .grid { grid-template-columns: repeat(3, 1fr); } }
    .card { background: linear-gradient(145deg, #171726, #10101a); border: 1px solid #25253a; border-radius: 16px; padding: 16px; display: flex; flex-direction: column; gap: 10px; transition: transform .2s ease, box-shadow .2s ease; }
    .card:hover { transform: translateY(-2px); box-shadow: 0 16px 32px rgba(0,0,0,0.35); }
    .album-title { font-weight: 700; font-size: 16px; }
    .album-year { color: var(--muted); font-size: 13px; }
    .footer { margin-top: 36px; color: var(--muted); font-size: 13px; text-align: center; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="avatar">AG</div>
      <div class="title">
        <h1>Ariana Grande</h1>
        <div class="meta">American singer, songwriter, and actress</div>
        <div class="chips">
          <span class="chip">Pop</span>
          <span class="chip">R&B</span>
          <span class="chip">Vocal</span>
        </div>
      </div>
    </div>

    <div style="margin-top: 24px; font-size: 14px; color: var(--muted);">Studio albums</div>
    <div class="grid">
      <div class="card">
        <div class="album-title">Yours Truly</div>
        <div class="album-year">2013</div>
      </div>
      <div class="card">
        <div class="album-title">My Everything</div>
        <div class="album-year">2014</div>
      </div>
      <div class="card">
        <div class="album-title">Dangerous Woman</div>
        <div class="album-year">2016</div>
      </div>
      <div class="card">
        <div class="album-title">Sweetener</div>
        <div class="album-year">2018</div>
      </div>
      <div class="card">
        <div class="album-title">Thank U, Next</div>
        <div class="album-year">2019</div>
      </div>
      <div class="card">
        <div class="album-title">Positions</div>
        <div class="album-year">2020</div>
      </div>
      <div class="card">
        <div class="album-title">Eternal Sunshine</div>
        <div class="album-year">2024</div>
      </div>
    </div>

    <div class="footer">This is a static demo page served by the Go server at /ariana.</div>
  </div>
</body>
</html>`
	_, _ = w.Write([]byte(html))
}
