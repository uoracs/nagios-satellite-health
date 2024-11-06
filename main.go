package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SatelliteAPIPingResponse struct {
	Results struct {
		Foreman struct {
			Database struct {
				Active     bool   `json:"active"`
				DurationMs string `json:"duration_ms"`
			} `json:"database"`
			Cache struct {
				Servers []struct {
					Status     string `json:"status"`
					DurationMs string `json:"duration_ms"`
				} `json:"servers"`
			} `json:"cache"`
		} `json:"foreman"`
		Katello struct {
			Services struct {
				Candlepin struct {
					Status     string `json:"status"`
					DurationMs string `json:"duration_ms"`
				} `json:"candlepin"`
				CandlepinAuth struct {
					Status     string `json:"status"`
					DurationMs string `json:"duration_ms"`
				} `json:"candlepin_auth"`
				ForemanTasks struct {
					Status     string `json:"status"`
					DurationMs string `json:"duration_ms"`
				} `json:"foreman_tasks"`
				KatelloEvents struct {
					Status     string `json:"status"`
					Message    string `json:"message"`
					DurationMs string `json:"duration_ms"`
				} `json:"katello_events"`
				CandlepinEvents struct {
					Status     string `json:"status"`
					Message    string `json:"message"`
					DurationMs string `json:"duration_ms"`
				} `json:"candlepin_events"`
				Pulp3 struct {
					Status     string `json:"status"`
					DurationMs string `json:"duration_ms"`
				} `json:"pulp3"`
				Pulp3Content struct {
					Status     string `json:"status"`
					DurationMs string `json:"duration_ms"`
				} `json:"pulp3_content"`
			} `json:"services"`
			Status string `json:"status"`
		} `json:"katello"`
	} `json:"results"`
}

func ExitErrorWithMessage(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	server, found := os.LookupEnv("SATELLITE_SERVER")
	if !found {
		ExitErrorWithMessage("SATELLITE_SERVER environment variable is not set")
	}
	url := "https://" + server + "/api/ping"

	resp, err := http.Get(url)
	if err != nil {
		ExitErrorWithMessage(fmt.Sprintf("failed to get response from Satellite server: %v", err))
	}
	if resp.StatusCode != 200 {
		ExitErrorWithMessage("response status code is not 200")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ExitErrorWithMessage(fmt.Sprintf("failed to read response body: %v", err))
	}

	var satResp SatelliteAPIPingResponse
	err = json.Unmarshal(body, &satResp)
	if err != nil {
		ExitErrorWithMessage(fmt.Sprintf("failed to unmarshal response body: %v", err))
	}

	if satResp.Results.Foreman.Database.Active != true {
		ExitErrorWithMessage("database is not active")
	}

	if satResp.Results.Foreman.Cache.Servers[0].Status != "ok" {
		ExitErrorWithMessage("foreman cache not ok")
	}

	if satResp.Results.Katello.Services.Candlepin.Status != "ok" {
		ExitErrorWithMessage("candlepin service not ok")
	}

	if satResp.Results.Katello.Services.CandlepinAuth.Status != "ok" {
		ExitErrorWithMessage("candlepin_auth service not ok")
	}

	if satResp.Results.Katello.Services.ForemanTasks.Status != "ok" {
		ExitErrorWithMessage("foreman tasks service not ok")
	}

	if satResp.Results.Katello.Services.KatelloEvents.Status != "ok" {
		ExitErrorWithMessage("katello events service not ok")
	}

	if satResp.Results.Katello.Services.CandlepinEvents.Status != "ok" {
		ExitErrorWithMessage("candlepin events service not ok")
	}

	if satResp.Results.Katello.Services.Pulp3.Status != "ok" {
		ExitErrorWithMessage("pulp3 service not ok")
	}

	if satResp.Results.Katello.Services.Pulp3Content.Status != "ok" {
		ExitErrorWithMessage("pulp3 content not ok")
	}

	if satResp.Results.Katello.Status != "ok" {
		ExitErrorWithMessage("katello not ok")
	}

	fmt.Println("OK")
	os.Exit(0)
}
