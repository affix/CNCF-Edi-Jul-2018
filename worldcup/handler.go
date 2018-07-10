package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Match []struct {
	Venue           string      `json:"venue"`
	Location        string      `json:"location"`
	Status          string      `json:"status"`
	Time            interface{} `json:"time"`
	FifaID          string      `json:"fifa_id"`
	Weather         interface{} `json:"weather"`
	Attendance      interface{} `json:"attendance"`
	Officials       interface{} `json:"officials"`
	StageName       string      `json:"stage_name"`
	HomeTeamCountry string      `json:"home_team_country"`
	AwayTeamCountry string      `json:"away_team_country"`
	Datetime        time.Time   `json:"datetime"`
	Winner          interface{} `json:"winner"`
	WinnerCode      interface{} `json:"winner_code"`
	HomeTeam        struct {
		Country   string      `json:"country"`
		Code      string      `json:"code"`
		Goals     int         `json:"goals"`
		Penalties interface{} `json:"penalties"`
	} `json:"home_team"`
	AwayTeam struct {
		Country   string      `json:"country"`
		Code      string      `json:"code"`
		Goals     int         `json:"goals"`
		Penalties interface{} `json:"penalties"`
	} `json:"away_team"`
	HomeTeamEvents    []interface{} `json:"home_team_events"`
	AwayTeamEvents    []interface{} `json:"away_team_events"`
	LastEventUpdateAt interface{}   `json:"last_event_update_at"`
	LastScoreUpdateAt time.Time     `json:"last_score_update_at"`
}

// Handle a serverless request
func Handle(req []byte) string {
	var match Match

	response, err := http.Get("https://worldcup.sfg.io/matches/today")
	if err != nil {
		return fmt.Sprintf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		json.Unmarshal(contents, &match)
		return CallFunction("http://figlet:8080", fmt.Sprintf("%s [ %d - %d ] %s", match[0].HomeTeam.Code, match[0].HomeTeam.Goals, match[0].AwayTeam.Goals, match[0].AwayTeam.Code))
	}
}

func CallFunction(function string, text string) string {
	textBytes := []byte(text)
	req, err := http.NewRequest("POST", function, bytes.NewBuffer(textBytes))
	if err != nil {
		return fmt.Sprintf("http.NewRequest() error: %v\n", err)
	}

	c := &http.Client{}
	resp, _ := c.Do(req)
	if err != nil {
		return fmt.Sprintf("http.Do() error: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("ioutil.ReadAll() error: %v\n", err)
	}

	return fmt.Sprintf(string(data))
}
