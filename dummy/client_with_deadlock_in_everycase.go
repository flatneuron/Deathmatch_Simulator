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
	Kills int
}

func main() {

	numPlayers := 5
	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	maxKills := 40
	gameOver := false

	wg.Add(numPlayers) // Still need to add the correct number of goroutines to the wait group
	for i := 0; i < numPlayers; i++ {
		go func(playerIndex int) { // Capture playerIndex correctly in the closure
			defer wg.Done()

			for {
				mu.Lock()
				if gameOver {
					mu.Unlock()
					return
				}
				mu.Unlock()

				duration := time.Duration(rand.Intn(751)+250) * time.Millisecond
				time.Sleep(duration)

				mu.Lock()
				// Select a random player
				randomIndex := rand.Intn(numPlayers)
				players[randomIndex].Kills++

				if players[randomIndex].Kills >= maxKills {
					gameOver = true
					fmt.Printf("Player %d reached %d kills and won!\n", players[randomIndex].ID, maxKills)
					mu.Unlock()
					return
				}

				fmt.Printf("Player %d got a kill! (Total: %d)\n", players[randomIndex].ID, players[randomIndex].Kills)
				mu.Unlock()
			}
		}(i) // Pass i as an argument to the goroutine function
	}

	wg.Wait()

	fmt.Println("Final Kill Counts:")
	for _, player := range players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
}
