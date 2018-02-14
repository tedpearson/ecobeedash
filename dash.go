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

func main() {
	config := configuration.LoadConfig("default.conf")
	// get api key from config

	// check if tokens exist.
	apiKey = config.GetString("api-key")
	client := http.NewClient(apiKey)
	// currentTemp(client)
	// summary(client)
	runtimeReport(client)

	// access token expires after 1 hour, refresh token after 1 year
	// return http client delegate, which will retry requests on auth errors, after refreshing tokens

}
