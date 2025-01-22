package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// // Player struct
// type Player struct {
// 	ID    int
// 	Kills int
// }

func main1() {
	// rand.Seed(time.Now().UnixNano())

	numPlayers := 5
	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect shared access to players slice
	maxKills := 40    // Target number of kills to end the simulation
	gameOver := false // Flag to indicate if a player has reached maxKills

	wg.Add(numPlayers)
	for i := 0; i < numPlayers; i++ {
		go simulateKills(i, &players, &wg, &mu, &gameOver, maxKills)
	}

	wg.Wait()

	fmt.Println("Final Kill Counts:")
	for _, player := range players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
}

func simulateKills(playerIndex int, players *[]Player, wg *sync.WaitGroup, mu *sync.Mutex, gameOver *bool, maxKills int) {
	defer wg.Done()

	for {
		// Check if game is over
		mu.Lock()
		if *gameOver {
			mu.Unlock()
			return
		}
		mu.Unlock()

		duration := time.Duration(rand.Intn(751)+250) * time.Millisecond
		time.Sleep(duration)

		mu.Lock()
		(*players)[playerIndex].Kills++

		// Check if current player reached maxKills
		if (*players)[playerIndex].Kills >= maxKills {
			*gameOver = true // Signal that the game is over
			fmt.Printf("Player %d reached %d kills and won!\n", (*players)[playerIndex].ID, maxKills)
			mu.Unlock()
			return // End this goroutine
		}

		fmt.Printf("Player %d got a kill! (Total: %d)\n", (*players)[playerIndex].ID, (*players)[playerIndex].Kills)
		mu.Unlock()
	}
}
