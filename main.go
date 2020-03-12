package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CoronaStats struct {
	Country   string `json: "country"`
	Confirmed string `json: "confirmed"`
	Deaths    string `json: "deaths"`
	Recovered string `json: "recovered"`
}

func stats(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["date"]

	date := "03-11-2020"

	if !ok || len(keys[0]) < 1 {
		fmt.Fprintf(w, "Error: Add date")

	} else {
		date = keys[0]
		log.Print(date, "ok")
	}

	url := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/" + date + ".csv"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		reader := csv.NewReader(response.Body)
		var corona []CoronaStats
		lineNum := 0
		for {
			line, error := reader.Read()
			if lineNum == 0 {
				lineNum++
				continue
			}
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			if line[0] != "" {
				corona = append(corona, CoronaStats{
					Country:   line[0] + ", " + line[1],
					Confirmed: line[3],
					Deaths:    line[4],
					Recovered: line[5],
				})
			} else {
				corona = append(corona, CoronaStats{
					Country:   line[1],
					Confirmed: line[3],
					Deaths:    line[4],
					Recovered: line[5],
				})
			}
			lineNum++
		}

		coronaJSON, _ := json.Marshal(corona)
		fmt.Fprintf(w, string(coronaJSON))
	}

}

func raw(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["date"]

	date := "03-11-2020"

	if !ok || len(keys[0]) < 1 {
		fmt.Fprintf(w, "Error: Add date")

	} else {
		date = keys[0]
		log.Print(date, "ok")
	}

	url := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/" + date + ".csv"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Fprintf(w, bodyString)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/raw", raw)
	router.HandleFunc("/stats", stats)
	log.Print("Listening @ http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
