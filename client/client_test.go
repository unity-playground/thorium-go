package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"thorium-go/model"
	"thorium-go/requests"
	"time"
)

var masterEndpoint = "thorium-sky.net:6960"
var sessionKey string
var characterIds []int
var gameList []model.Game

var user string = "test"
var password string = "test"

// Test 1A: Service Status
// HTTP GET /status

func Test1A_Ping(t *testing.T) {
	statusCode, _, err := GetStatus(masterEndpoint)

	if err != nil {
		t.Fail()
	} else if statusCode != 200 {
		log.Printf("recieved bad status code %d", statusCode)
		t.Fail()
	} else {
		log.Print("Test 1A Success")
	}
}

// Test 2A: Register
// HTTP POST /clients/register

func Test2A_Register(t *testing.T) {

	fmt.Println("Test 2A: Register")

	tries := 3

	for tries > 0 {

		// pick a random number between 10k-1M
		rand.Seed(time.Now().UTC().UnixNano())
		randNum := rand.Intn(990000) + 10000
		// create user string
		user = fmt.Sprintf("test%d", randNum)

		fmt.Println(user)
		// execute request
		responseCode, body, err := Register(masterEndpoint, user, password)
		if err != nil {
			log.Print(err)
			t.FailNow()
		}

		fmt.Printf("register response: status %d, body %s\n", responseCode, body)
		if responseCode == 200 {

			var resp request.LoginResponse
			json.Unmarshal([]byte(body), &resp)

			if resp.SessionKey != "" {
				// success
				sessionKey = resp.SessionKey
				characterIds = resp.CharacterIDs
				return
			}
		}

		tries = tries - 1
	}

	// fail if all tries are used up
	log.Print("all register tries used")
	t.Fail()
}

// Test 2B: Disconnect
// HTTP POST /clients/disconnect
func Test2B_Disconnect(t *testing.T) {

	fmt.Println("Test 2B: Disconnect")

	rc, _, err := Disconnect(masterEndpoint, sessionKey)
	if err != nil {
		log.Print(err)
		t.FailNow()
	}

	log.Print("disconnect response code = ", rc)
	if rc != 200 {
		t.FailNow()
	}
}

// Test 2C: Login
// HTTP POST /clients/login

func Test2C_Login(t *testing.T) {
	fmt.Println("Test 2C: Login")

	// execute request
	responseCode, body, err := Login(masterEndpoint, user, password)
	if err != nil {
		log.Print(err)
		t.FailNow()
	}

	fmt.Printf("login response: status %d, %s\n", responseCode, body)
	if responseCode != 200 {
		t.Fail()
	} else {
		var resp request.LoginResponse
		json.Unmarshal([]byte(body), &resp)

		if resp.SessionKey == "" {
			t.FailNow()
		}

		sessionKey = resp.SessionKey
		characterIds = resp.CharacterIDs
	}
}

// Test 3A: Create Character
// HTTP POST /characters/new

func Test3A_CreateCharacter(t *testing.T) {
	fmt.Println("Test 3A: Create Character")

	// execute request
	rc, body, err := CreateCharacter(masterEndpoint, sessionKey, user, 1)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("create character response: status %d, %s\n", rc, body)
	if rc != 200 {
		t.FailNow()
	}

	var resp request.NewCharacterResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println("failed json unmarshal")
		t.FailNow()
	}

	if resp.CharacterId == 0 {
		fmt.Println("response returned invalid character id")
		t.FailNow()
	}
	characterIds = append(characterIds, resp.CharacterId)
	fmt.Println("characterIds: ", characterIds)
}

// Test 3B: Get Character Info

func Test3B_GetCharacter(t *testing.T) {

	fmt.Println("Test3B: Get Character Info")

	count := len(characterIds)
	if count == 0 {
		t.FailNow()
	}

	for i := 0; i < count; i++ {

		rc, body, err := GetCharacter(masterEndpoint, characterIds[i])
		if err != nil {
			t.FailNow()
		}

		fmt.Printf("get character response: status %d, %s\n", rc, body)
		if rc != 200 {
			t.FailNow()
		}
	}
}

// Test 4A: Query For Game List
func Test4A_GameGameList(t *testing.T) {

	fmt.Println("Test4A: Get Game List")

	rc, body, err := GetGameList(masterEndpoint)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("get game list response: status %d, %s\n", rc, body)

	if rc != 200 {
		t.FailNow()
	}

	err = json.Unmarshal([]byte(body), &gameList)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

// Test 4B: Create New Game
func Test4B_CreateNewGame(t *testing.T) {

	fmt.Println("Test4B: Create New Game")

	mapName := "Map_Sandbox"
	mode := "Tutorial"
	minimumLevel := 1
	maxPlayers := 16

	rc, body, err := CreateNewGame(masterEndpoint, sessionKey, mapName, mode, minimumLevel, maxPlayers)
	if err != nil {

		log.Print(err)
		t.FailNow()
	}

	fmt.Printf("new game response: status %d, %s\n", rc, body)

	if rc != 201 { // HTTP 201 Created

		log.Print("failed to create game")
		t.FailNow()
	}

	var data request.CreateNewGameResponse
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {

		log.Print(err)
		t.FailNow()
	}

	if data.GameId == 0 {

		log.Print("bad response - gameId is 0")
		t.FailNow()
	}

	tries := 10

	for tries > 0 {

		rc, body, err := GetServerInfo(masterEndpoint, data.GameId)
		if err != nil {

			log.Print(err)
			t.FailNow()
		}

		fmt.Printf("get server info response: status %d, %s\n", rc, body)

		if rc == 200 {

			log.Print("game server ready")
			break

		} else if rc == 202 {

			log.Print("game server loading...")

		} else {

			log.Print("game server unavailable")
			t.FailNow()
		}

		tries = tries - 1
		time.Sleep(50 * time.Millisecond)
	}

	if tries == 0 {

		log.Print("server is not ready after 10 tries")
		t.FailNow()

	}
}

// Test 4C: Join Existing Game
