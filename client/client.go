package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/jaybennett89/thorium-go/requests"
)

import "bytes"
import "io/ioutil"

func GetStatus(masterEndpoint string) (int, string, error) {

	url := fmt.Sprintf("http://%s/status", masterEndpoint)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		return 0, "", err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Print("ping master - error:\n", err)
		return 0, "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(bodyBytes), nil
}

func Register(masterEndpoint string, username string, password string) (int, string, error) {

	// create request data struct in memory
	var loginReq request.Authentication
	loginReq.Username = username
	loginReq.Password = password

	// marshal request data into json byte array
	jsonBytes, err := json.Marshal(&loginReq)
	if err != nil {
		return 0, "", err
	}

	// create the http request struct
	url := fmt.Sprintf("http://%s/clients/register", masterEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Print("error with request: ", err)
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// create the http client struct and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return 0, "", err
	}

	// read the body into a byte array and return the results
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(body), nil
}

func Login(masterEndpoint string, username string, password string) (int, string, error) {

	// create request data struct in memory
	var loginReq request.Authentication
	loginReq.Username = username
	loginReq.Password = password

	// marshal request data into json byte array
	jsonBytes, err := json.Marshal(&loginReq)
	if err != nil {
		return 0, "", err
	}

	// create the http request struct
	url := fmt.Sprintf("http://%s/clients/login", masterEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Print("error with request: ", err)
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// create the http client struct and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return 0, "", err
	}

	// read the body into a byte array and return the results
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(body), nil
}

func Disconnect(masterEndpoint string, token string) (int, string, error) {

	var disconnectReq request.Disconnect
	disconnectReq.SessionKey = token

	// marshal request data into json byte array
	jsonBytes, err := json.Marshal(&disconnectReq)
	if err != nil {
		return 0, "", err
	}

	url := fmt.Sprintf("http://%s/clients/disconnect", masterEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Print("error with request: ", err)
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return 0, "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(body), nil
}

func CreateCharacter(masterEndpoint string, sessionKey string, name string, classId int) (int, string, error) {

	var charCreateReq request.CreateCharacter
	charCreateReq.SessionKey = sessionKey
	charCreateReq.Name = name
	charCreateReq.ClassId = classId
	jsonBytes, err := json.Marshal(&charCreateReq)
	if err != nil {
		return 0, "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/characters/new", masterEndpoint), bytes.NewBuffer(jsonBytes))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error with request: ", err)
		return 0, "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

func SelectCharacter(masterEndpoint string, sessionKey string, characterId int) (int, string, error) {

	selectCharacter := request.SelectCharacter{

		SessionKey:  sessionKey,
		CharacterId: characterId,
	}

	json, err := json.Marshal(&selectCharacter)
	if err != nil {
		return 0, "", err
	}

	url := fmt.Sprintf("http://%s/characters/select", masterEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return 0, "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(body), nil
}

func GetGameList(masterEndpoint string) (int, string, error) {

	url := fmt.Sprintf("http://%s/games", masterEndpoint)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Print(err)
		return 0, "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return 0, "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(body), nil
}

func CreateNewGame(masterEndpoint string, sessionKey string, mapName string, gameMode string, minimumLevel int, maxPlayers int) (int, string, error) {

	data := request.CreateNewGame{
		SessionKey:   sessionKey,
		Map:          mapName,
		GameMode:     gameMode,
		MinimumLevel: minimumLevel,
		MaxPlayers:   maxPlayers}

	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return 0, "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/games", masterEndpoint), bytes.NewBuffer(jsonBytes))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error with request: ", err)
		return 0, "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

func GetServerInfo(masterEndpoint string, gameId int) (int, string, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/games/%d/server_info", masterEndpoint, gameId), bytes.NewBuffer([]byte("")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error with request: ", err)
		return 0, "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

func JoinGame(masterEndpoint string, gameId int, sessionKey string) (int, string, error) {

	data := request.JoinGame{
		GameId:     gameId,
		SessionKey: sessionKey}

	json, err := json.Marshal(&data)
	if err != nil {

		return 0, "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/games/join", masterEndpoint), bytes.NewBuffer(json))
	if err != nil {

		return 0, "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return 0, "", err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}
