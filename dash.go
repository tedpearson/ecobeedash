package main

import (
	"fmt"

	"github.com/tedpearson/ecobeedash/http"

	"github.com/go-akka/configuration"
)

var (
	apiKey string
)

func currentTemp(client *http.Client) {
	data, err := client.Get("https://api.ecobee.com/1/thermostat?format=json&body=" +
		"{\"selection\":{\"selectionType\":\"registered\",\"selectionMatch\":\"\",\"includeRuntime\":true}}")
	if err == nil {
		fmt.Println(string(data))
	} else {
		fmt.Println(err)
	}
}

func summary(client *http.Client) {
	data, err := client.Get("https://api.ecobee.com/1/thermostatSummary?format=json&body=" +
		"{\"selection\":{\"selectionType\":\"registered\",\"selectionMatch\":\"\",\"includeEquipmentStatus\":true}}")
	if err == nil {
		fmt.Println(string(data))
	} else {
		fmt.Println(err)
	}
}

func runtimeReport(client *http.Client) {
	data, err := client.Get("https://api.ecobee.com/1/runtimeReport?format=json&body=" +
		`{"startDate":"2018-02-09","endDate":"2018-02-12","columns":"` +
		`auxHeat1,compHeat1,compHeat2,fan,outdoorTemp,outdoorHumidity,zoneAveTemp,zoneHeatTemp",` +
		`"selection":{"selectionType":"thermostats","selectionMatch":"511810481935"}}`)
	if err == nil {
		fmt.Println(string(data))
	} else {
		fmt.Println(err)
	}
}

/*
TODO
figure out what information to collect, either via discovery, or by user input.
on first run, collect all historical data possible, with all information from above
on subsequent runs, figure out what data exists and get new data since then
while running, periodically poll for more information
while running, poll thermostat every 5 mins to get an idea of what the next 1 hour of data will look like. store it separately though.
make sure to stay under an api call quota
all of this data needs to be stored in a sqlite db: https://github.com/mattn/go-sqlite3; store in 5 min intervals
add apis to get the data! maybe change their formats, maybe keep a similar format to ecobee

add an api to calculate cost based on tables for efficiency and configured cost for stage per time period
*/

func main() {
	config := configuration.LoadConfig("default.conf")
	// get api key from config

	// check if tokens exist.
	apiKey = config.GetString("api-key")
	client := http.NewClient(apiKey)
	// currentTemp(client)
	// summary(client)
	runtimeReport(client)

}
