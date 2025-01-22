package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Player struct
type Player struct {
	ID    int
	Kills int32
}

// GameState holds the overall game information
type GameState struct {
	Players  []Player
	MaxKills int32
	Mu       sync.Mutex // Mutex to protect access to GameState
}

// KillRequest is a request to register a kill
type KillRequest struct {
	PlayerID int
}

func main() {

	numPlayers := 5
	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	gameState := &GameState{
		Players:  players,
		MaxKills: 40,
	}

	killChan := make(chan KillRequest)
	doneChan := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(numPlayers + 1)

	go gameManager(gameState, killChan, doneChan, &wg)

	for i := 0; i < numPlayers; i++ {
		go playerRoutine(killChan, doneChan, &wg, numPlayers)
	}

	wg.Wait()

	gameState.Mu.Lock()
	fmt.Println("Final Kill Counts:")
	for _, player := range gameState.Players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
	gameState.Mu.Unlock()
}

func gameManager(gameState *GameState, killChan chan KillRequest, doneChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return
		case killReq := <-killChan:
			gameState.Mu.Lock()

			playerIndex := killReq.PlayerID - 1

			if playerIndex >= 0 && playerIndex < len(gameState.Players) {
				gameState.Players[playerIndex].Kills++

				fmt.Printf("Player %d got a kill! (Total: %d)\n",
					gameState.Players[playerIndex].ID,
					gameState.Players[playerIndex].Kills)

				if gameState.Players[playerIndex].Kills >= gameState.MaxKills {
					fmt.Printf("Player %d reached %d kills and won!\n",
						gameState.Players[playerIndex].ID,
						gameState.MaxKills)

					// Signal game over to all player routines
					for i := 0; i < len(gameState.Players); i++ {
						doneChan <- true
					}
					gameState.Mu.Unlock()
					return
				}
			}

			gameState.Mu.Unlock()
		}
	}
}

func playerRoutine(killChan chan KillRequest, doneChan chan bool, wg *sync.WaitGroup, numPlayers int) {
	defer wg.Done()

	for {
		time.Sleep(time.Duration(rand.Intn(76)+25) * time.Millisecond)
		select {
		case <-doneChan:
			return
		default:
			randomIndex := rand.Intn(numPlayers)
			killChan <- KillRequest{PlayerID: randomIndex + 1}
		}
	}
}
