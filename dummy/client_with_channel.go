package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Player struct
type Player struct {
	ID    int
	Kills int32 // Use int32 for atomic operations
}

func main() {
	// rand.Seed(time.Now().UnixNano())

	numPlayers := 50
	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	var wg sync.WaitGroup
	maxKills := int32(40)   // Use int32 for atomic operations
	done := make(chan bool) // Channel to signal game over

	wg.Add(numPlayers)
	for i := 0; i < numPlayers; i++ {
		go func(playerIndex int) {
			defer wg.Done()

			for {
				select {
				case <-done: // Check if done channel is closed
					return
				default: // Non-blocking default case
					// Select a random player
					randomIndex := rand.Intn(numPlayers)

					// Atomically increment kills
					atomic.AddInt32(&players[randomIndex].Kills, 1)

					// Check if current player reached maxKills
					if atomic.LoadInt32(&players[randomIndex].Kills) >= maxKills {
						fmt.Printf("Player %d reached %d kills and won!\n", players[randomIndex].ID, maxKills)
						close(done) // Signal game over by closing the channel
						return      // End this goroutine
					}

					fmt.Printf("Player %d got a kill! (Total: %d)\n", players[randomIndex].ID, atomic.LoadInt32(&players[randomIndex].Kills))
				}
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("Final Kill Counts:")
	for _, player := range players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
}
