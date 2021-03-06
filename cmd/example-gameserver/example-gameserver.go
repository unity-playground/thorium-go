package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/jaybennett89/thorium-go/client"
	"github.com/jaybennett89/thorium-go/model"
	"github.com/jaybennett89/thorium-go/requests"
	"time"

	"github.com/go-martini/martini"
)

var machineKey string
var listenPort int
var servicePort int
var game model.Game
var players map[string]*model.Character

func main() {
	log.Print("running a mock game server")

	flag.StringVar(&machineKey, "key", "", "machine key")
	flag.IntVar(&game.GameId, "id", 0, "identifies this game within the cluster")
	flag.IntVar(&listenPort, "listen", 0, "game server listen port")
	flag.IntVar(&servicePort, "service", 0, "machine local service port")
	flag.StringVar(&game.Map, "map", "mp_sandbox", "game map: mp_sandbox, mp_openworld")
	flag.StringVar(&game.Mode, "mode", "tutorial", "game mode: basic, tutorial, freeforall")
	flag.IntVar(&game.MinimumLevel, "minlvl", 0, "minimum level of player")
	flag.IntVar(&game.MaximumPlayers, "maxplayers", 16, "maximum player count")

	flag.Parse()

	if game.GameId == 0 || listenPort == 0 || servicePort == 0 || game.Map == "" || game.Mode == "" {
		log.Fatal("bad arguments")
	}

	var data request.RegisterGameServer
	data.MachineKey = machineKey
	data.Port = listenPort
	data.GameId = game.GameId
	jsonBytes, err := json.Marshal(&data)

	time.Sleep(180 * time.Millisecond)

	endpoint := fmt.Sprintf("http://localhost:%d/games/register_server", servicePort)
	var req *http.Request
	req, err = http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Die - failed to register")
	}

	players = make(map[string]*model.Character)

	m := martini.Classic()
	m.Get("/status", handleStatusRequest)
	m.Post("/connect", handleConnectRequest)
	m.Post("/move", handleMoveRequest)
	m.Post("/disconnect", handleDisconnect)
	m.RunOnAddr(fmt.Sprintf(":%d", listenPort))
}

func handleStatusRequest(httpReq *http.Request) (int, string) {
	return 200, "OK"
}

type ConnectToken struct {
	SessionKey  string `json:"sessionKey"`
	CharacterId int    `json:"characterId"`
}

func handleConnectRequest(httpReq *http.Request) (int, string) {

	var req ConnectToken
	decoder := json.NewDecoder(httpReq.Body)
	err := decoder.Decode(&req)
	if err != nil {

		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	serviceEndpoint := fmt.Sprintf("localhost:%d", servicePort)

	rc, body, err := client.PlayerConnect(serviceEndpoint, game.GameId, machineKey, req.SessionKey, req.CharacterId)
	if err != nil {

		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	if rc != 200 {

		fmt.Println("status: ", rc, " body: ", body)
		return 500, "Internal Server Error"
	}

	var resp request.PlayerConnectResponse
	err = json.Unmarshal([]byte(body), &resp)

	players[req.SessionKey] = resp.Character

	fmt.Println("instantiate player: ", body)
	return 200, "OK"
}

type Move struct {
	SessionKey string        `json:"sessionKey"`
	MoveDir    model.Vector3 `json:"movedir"`
}

// this method is for testing only! not a real game server!!
func handleMoveRequest(httpReq *http.Request) (int, string) {

	var req Move
	decoder := json.NewDecoder(httpReq.Body)
	err := decoder.Decode(&req)
	if err != nil {

		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	players[req.SessionKey].Position.X += req.MoveDir.X
	players[req.SessionKey].Position.Y += req.MoveDir.Y
	players[req.SessionKey].Position.Z += req.MoveDir.Z

	// for testing purposes we will update the character in the database after every move

	serviceEndpoint := fmt.Sprintf("localhost:%d", servicePort)

	rc, body, err := client.UpdateCharacter(serviceEndpoint, machineKey, players[req.SessionKey])
	if err != nil {

		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	if rc != 200 {

		fmt.Println("status: ", rc, " body: ", body)
		return 500, "Internal Server Error"
	}

	return 200, "OK"
}

type Disconnect struct {
	SessionKey string `json:"sessionKey"`
}

func handleDisconnect(httpReq *http.Request) (int, string) {

	var req Disconnect
	decoder := json.NewDecoder(httpReq.Body)
	err := decoder.Decode(&req)
	if err != nil {

		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	serviceEndpoint := fmt.Sprintf("localhost:%d", servicePort)

	rc, body, err := client.PlayerDisconnect(serviceEndpoint, machineKey, game.GameId, players[req.SessionKey])
	if err != nil {

		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	if rc != 200 {

		fmt.Println("status: ", rc, " body: ", body)
		return 500, "Internal Server Error"
	}

	return 200, "OK"
}
