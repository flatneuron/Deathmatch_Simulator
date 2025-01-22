package main

import (
	"fmt"
	"math/rand"
	"sync"
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

	wg.Add(numPlayers)
	for i := 0; i < numPlayers; i++ {
		go func(playerIndex int) {
			defer wg.Done()

			for {
				mu.Lock()
				if gameOver {
					mu.Unlock()
					return
				}
				mu.Unlock()

				// Removed time.Sleep() for maximum speed

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
		}(i)
	}

	wg.Wait()

	fmt.Println("Final Kill Counts:")
	for _, player := range players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
}
