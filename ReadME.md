# Data Alert Storage System

## Overview

Data Alert Storage System in Golang that allows users to send requests for reading and writing data to a
data storage system. APIs provide endpoints for writing and reading alerts.

## Setup

1. Golang Installation : Ensure that you have installed Go programming language installed on your system.
   To download and install Go from official website : https://golang.org/dl/.
2. Project include some dependencies which can be installed through the following commands:

```
go get github.com/gorilla/mux
go get github.com/stretchr/testify/mock
```

3. Build and run :
   - Clone or download the code repository from github repository.
   - Open command prompt or terminal and navigate to the directory where the code is located.
   - Run the code : ``` go run main.go ```
4. Testing : Code includes test cases that ensure functionality of the APIs.
You can run the tests using command : ``` go test ```

5. API Endpoints : 
The API exposes th following endpoints 
   - POST `/alerts`: To write alert data to the database.
   - GET `/alerts` : To read alerts by service_id and time range.

## Result

### Write Alerts

Endpoint : POST - ```localhost:8080/alerts```
Request body :
```
{
    "alert_id": "b950482e9911ec7e41f7ca5e5d9a4241234",
    "service_id": "my_test_service_id",
    "service_name": "my_test_service",
    "model": "my_test_model",
    "alert_type": "anomaly",
    "alert_ts": "1695644166",
    "severity": "warning",
    "team_slack": "slack_ch"
}
```
Response Body :
```
{
   "alert_id":"b950482e9911ec7e41f7ca5e5d9a4241234",
   "error":""
}
```

### Read Alerts
Endpoint : GET - ```localhost:8080/alerts?service_id=my_test_service_id&start_ts=1695644160&end_ts=1695644170```

Response Body :
```
{
    "service_id": "my_test_service_id",
    "service_name": "my_test_service"
    "alerts": [
        {
            "alert_id": "b950482e9911ec7e41f7ca5e5d9a42411",
            "model": "my_test_model",
            "alert_type": "anomaly",
            "alert_ts": "1695644160",
            "severity": "warning",
            "team_slack": "slack_ch"
        },
        {
            "alert_id": "b950482e9911ec7e41f7ca5e5d9a42412",
            "model": "my_test_model",
            "alert_type": "anomaly",
            "alert_ts": "1695644160",
            "severity": "warning",
            "team_slack": "slack_ch"
        },
        {
            "alert_id": "b950482e9911ec7e41f7ca5e5d9a42413",
            "model": "my_test_model",
            "alert_type": "anomaly",
            "alert_ts": "1695644166",
            "severity": "warning",
            "team_slack": "slack_ch"
        },
        {
            "alert_id": "b950482e9911ec7e41f7ca5e5d9a4241234",
            "model": "my_test_model",
            "alert_type": "anomaly",
            "alert_ts": "1695644166",
            "severity": "warning",
            "team_slack": "slack_ch"
        },
        {
            "alert_id": "b950482e9911ec7e41f7ca5e5d9a424123",
            "model": "my_test_model",
            "alert_type": "anomaly",
            "alert_ts": "1695644166",
            "severity": "warning",
            "team_slack": "slack_ch"
        }
    ],
}
```