package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/tedpearson/ecobeedash/util"
)

// Tokens to access ecobee
type Tokens struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

// Client that delegates to http.Client but tries to get new tokens on auth errors
type Client struct {
	client http.Client
	tokens Tokens
	apiKey string
}

// NewClient returns a new http client that can get new tokens
func NewClient(apiKey string) *Client {
	return &Client{tokens: getTokens(apiKey), apiKey: apiKey}
}

// Get Performs an HTTP GET, getting new tokens on auth error
func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.tokens.Access)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// retry via recursion
	if checkRetryStatus(data) == true {
		c.refreshTokens()
		return c.Get(url)
	}
	return data, nil
}

type status struct {
	Code    int
	Message string
}

type statusResponse struct {
	Status status
}

func checkRetryStatus(data []byte) bool {
	var s statusResponse
	err := json.Unmarshal(data, &s)
	util.CheckError(err, "json parse error")
	if s.Status.Code == 14 {
		return true
	}
	if s.Status.Code != 0 {
		log.Fatal(s.Status.Message)
	}
	return false
}

type authResponse struct {
	Code      string
	EcobeePin string
}

func parseAndWriteTokens(resp *http.Response, destination interface{}) {
	jsonTxt, err := ioutil.ReadAll(resp.Body)
	util.CheckError(err, "reading from response failed")
	err = json.Unmarshal(jsonTxt, &destination)
	util.CheckError(err, "Couldn't parse json")
	ioutil.WriteFile("tokens.json", jsonTxt, 0755)
}

func (c *Client) refreshTokens() {
	theURL := "https://api.ecobee.com/token"
	resp, err := http.PostForm(theURL, url.Values{"grant_type": {"refresh_token"}, "code": {c.tokens.Refresh}, "client_id": {c.apiKey}})
	util.CheckError(err, "Couldn't get tokens from ecobee")
	var t Tokens
	parseAndWriteTokens(resp, &t)
	c.tokens = t
}

func getTokens(apiKey string) Tokens {
	data, err := ioutil.ReadFile("tokens.json")
	if err == nil {
		var t Tokens
		err = json.Unmarshal(data, &t)
		util.CheckError(err, "Couldn't read tokens")
		return t
	}
	// ask user to get pin.
	theURL := "https://api.ecobee.com/authorize?response_type=ecobeePin&scope=smartRead&client_id=" + apiKey
	resp, err := http.Get(theURL)
	util.CheckError(err, "Couldn't get pin from ecobee")
	defer resp.Body.Close()
	jsonTxt, err := ioutil.ReadAll(resp.Body)
	util.CheckError(err, "reading from response failed")
	var jsonResponse authResponse
	err = json.Unmarshal(jsonTxt, &jsonResponse)
	util.CheckError(err, "Couldn't decode json")
	// code and ecobeePin
	fmt.Println("Go to https://www.ecobee.com/consumerportal/index.html#/my-apps/add/new and enter this PIN:")
	fmt.Println(jsonResponse.EcobeePin)
	fmt.Println("Press enter when done.")
	fmt.Scanln()
	// get tokens
	theURL = "https://api.ecobee.com/token"
	resp, err = http.PostForm(theURL, url.Values{"grant_type": {"ecobeePin"}, "code": {jsonResponse.Code}, "client_id": {apiKey}})
	util.CheckError(err, "Couldn't get tokens from ecobee")
	var t Tokens
	parseAndWriteTokens(resp, &t)
	return t
}
