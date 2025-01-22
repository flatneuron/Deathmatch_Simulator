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
	GameOver bool
	Mu       sync.Mutex // Mutex to protect access to GameState
}

// KillRequest is a request to register a kill
type KillRequest struct {
	PlayerID int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numPlayers := 5
	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	gameState := &GameState{
		Players:  players,
		MaxKills: 40,
		GameOver: false,
	}

	killChan := make(chan KillRequest) // Channel for kill requests
	doneChan := make(chan bool)        // Channel to signal game over

	var wg sync.WaitGroup
	wg.Add(numPlayers + 1) // Add 1 for the game manager goroutine

	// Start the game manager goroutine
	go gameManager(gameState, killChan, doneChan, &wg)

	// Start player goroutines
	for i := 0; i < numPlayers; i++ {
		go playerRoutine(killChan, doneChan, &wg, gameState)
	}

	wg.Wait()

	fmt.Println("Final Kill Counts:")
	for _, player := range gameState.Players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
}

// gameManager is the central goroutine that manages the game state and enforces order
func gameManager(gameState *GameState, killChan chan KillRequest, doneChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-doneChan: // Game over signal
			return
		case killReq := <-killChan: // Kill request received
			gameState.Mu.Lock()

			playerID := killReq.PlayerID
			playerIndex := -1
			for i, p := range gameState.Players {
				if p.ID == playerID {
					playerIndex = i
					break
				}
			}

			if playerIndex != -1 {
				gameState.Players[playerIndex].Kills++

				fmt.Printf("Player %d got a kill! (Total: %d)\n",
					gameState.Players[playerIndex].ID,
					gameState.Players[playerIndex].Kills)

				if gameState.Players[playerIndex].Kills >= gameState.MaxKills {
					gameState.GameOver = true
					fmt.Printf("Player %d reached %d kills and won!\n",
						gameState.Players[playerIndex].ID,
						gameState.MaxKills)

					close(doneChan) // Signal game over
				}
			}

			gameState.Mu.Unlock()
		}
	}
}

// playerRoutine simulates a player's actions
func playerRoutine(killChan chan KillRequest, doneChan chan bool, wg *sync.WaitGroup, gameState *GameState) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return
		default:
			// Simulate some activity and then send a kill request
			time.Sleep(time.Duration(rand.Intn(751)+250) * time.Millisecond)

			// Select a random player for the kill (needs access to gameState.Players)
			gameState.Mu.Lock()
			randomIndex := rand.Intn(len(gameState.Players))
			killChan <- KillRequest{PlayerID: gameState.Players[randomIndex].ID}
			gameState.Mu.Unlock()
		}
	}
}
