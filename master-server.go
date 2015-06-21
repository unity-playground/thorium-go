package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"thorium-go/database"
)
import "github.com/go-martini/martini"
import "thorium-go/requests"

func main() {
	fmt.Println("hello world")

	m := martini.Classic()

	// status
	m.Get("/status", handleGetStatusRequest)

	// client
	m.Post("/clients/login", handleClientLogin)
	m.Post("/clients/register", handleClientRegister)
	m.Post("/clients/disconnect", handleClientDisconnect)

	// characters
	m.Post("/new", handleCreateCharacter)
	m.Get("/:id", handleGetCharacter)
	m.Get("/:id/simple", handleGetCharSimple)

	// games
	m.Post("/games/:id/register_server", handleRegisterServer)
	m.Post("/games/:id/server_status", handleGameServerStatus)

	m.Post("/games/new_request", handleGameRequest) // deprecate
	m.Post("/games/new", handleGameRequest)

	m.Get("/games", handleGetServerList)
	m.Get("/games/:id", handleGetGameInfo)
	m.Post("/games/join", handleClientJoinGame)
	m.Post("/games/join_queue", handleClientJoinQueue)

	// machines
	m.Post("/machines/register", handleRegisterMachine)
	m.Post("/machines/register_new", handleRegisterMachine) // deprecate

	m.Post("/machines/:id/unregister", handleUnregisterMachine)
	m.Delete("/machines/:id", handleUnregisterMachine)

	m.Run()
}

func handleGetStatusRequest(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleClientLogin(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleClientRegister(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleClientDisconnect(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleCreateCharacter(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleGetCharacter(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleGetCharSimple(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleClientJoinGame(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleClientJoinQueue(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleGameServerStatus(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleGetServerList(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleGetGameInfo(httpReq *http.Request) (int, string) {
	return 500, "Not Implemented"
}

func handleRegisterMachine(httpReq *http.Request) (int, string) {

	decoder := json.NewDecoder(httpReq.Body)
	var req request.RegisterMachine
	err := decoder.Decode(&req)
	if err != nil {
		logerr("Error decoding machine register request", err)
		return 500, "Internal Server Error"
	}

	if req.Port == 0 {
		fmt.Println("No Port Given")
		return 500, "No Port Given"
	} else {
		fmt.Println("register port = ", req.Port)
	}

	machineIp := strings.Split(httpReq.RemoteAddr, ":")[0]

	var machineId int
	machineId, err = thordb.RegisterMachine(machineIp, req.Port)
	if err != nil {
		logerr("error registering machine", err)
		return 500, "Internal Server Error"
	}
	fmt.Println("machine registered, ip=", machineIp)
	return 200, strconv.Itoa(machineId)
}

func handleUnregisterMachine(params martini.Params) (int, string) {

	machineId, err := strconv.Atoi(params["id"])
	if err != nil {
		logerr(fmt.Sprint("unable to convert request parameter, id=", params["id"]), err)
		return 400, "Bad Request"
	}

	success, err := thordb.UnregisterMachine(machineId)
	if err != nil || success {
		logerr("unable to remove machine registry", err)
	}

	fmt.Println("machine unregistered, id=", machineId)
	return 200, "OK"
}

func handleGameRequest(httpReq *http.Request) (int, string) {
	fmt.Println("[ThoriumNET] master-server.handleGameRequest")

	decoder := json.NewDecoder(httpReq.Body)
	var req request.NewGame
	err := decoder.Decode(&req)
	if err != nil {
		logerr("unable to decode body data", err)
		return 500, "Internal Server Error"
	}

	if req.Map == "" {
		fmt.Println("No Map Name Given")
		return 400, "Missing Parameters"
	}

	var gameId int
	gameId, err = thordb.RegisterNewGame(req.Map, req.MaxPlayers)
	if err != nil {
		fmt.Println("[ThoriumNET] unable to insert new game record")
		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	fmt.Println("[ThoriumNET] new game, id=", strconv.Itoa(gameId))
	return 200, "OK"
}

func handleRegisterServer(httpReq *http.Request, params martini.Params) (int, string) {
	decoder := json.NewDecoder(httpReq.Body)
	var req request.RegisterGame
	err := decoder.Decode(&req)
	if err != nil {
		logerr("Error decoding machine register request", err)
		return 500, "Internal Server Error"
	}

	if req.Port == 0 {
		fmt.Println("No Port Given")
		return 400, "Missing Parameters"
	}

	var gameId int
	gameId, err = strconv.Atoi(params["id"])
	if err != nil {
		logerr(fmt.Sprintf("unable to convert parameter id=%s to integer", params["id"]), err)
		return 400, "Bad Request"
	}

	fmt.Println("[ThoriumNET] master-server.handleRegisterServer ID=", gameId)

	exists, err := thordb.CheckExists(gameId)
	if err != nil {
		logerr("unable to connect to DB", err)
		return 500, "Internal Server Error"
	}

	if !exists {
		fmt.Println("game id ", strconv.Itoa(gameId), " does not exist")
		return 400, "Bad Request"
	}

	registered, err := thordb.RegisterActiveGame(gameId, req.MachineId, req.Port)
	if err != nil || !registered {
		logerr("unable to register game", err)
		return 500, "Internal Server Error"
	}

	fmt.Println("Found game ", gameId)
	return 200, "OK"

}

// TODO: Refactor into logging package
func logerr(msg string, err error) {
	fmt.Println("[ThoriumNET] ", msg)
	fmt.Println(err)
}