package main

import (
	"fmt"
	"math"
	"math/rand"
)

// Player struct
type Player struct {
	ID    int
	Kills int
}

func main() {

	numPlayers := 5
	maxKills := 40

	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	gameOver := false
	for !gameOver {
		// Calculate weights exponentially with a minimum weight
		weights := make([]float64, numPlayers)
		minWeight := 0.05
		for i, player := range players {
			weights[i] = math.Exp(-0.2 * float64(player.Kills)) // Adjust exponent here
			if weights[i] < minWeight {
				weights[i] = minWeight
			}
		}

		// Select a player based on weights
		randomIndex := weightedRandomSelection(weights)

		// Increment the kill count for the selected player
		players[randomIndex].Kills++

		// Print the kill event
		fmt.Printf("Player %d got a kill! (Total: %d)\n", players[randomIndex].ID, players[randomIndex].Kills)

		// Check if the player has reached maxKills
		if players[randomIndex].Kills >= maxKills {
			gameOver = true
			fmt.Printf("Player %d reached %d kills and won!\n", players[randomIndex].ID, maxKills)
		}
	}

	fmt.Println("Final Kill Counts:")
	for _, player := range players {
		fmt.Printf("Player %d: %d kills\n", player.ID, player.Kills)
	}
}

// weightedRandomSelection selects an index based on the provided weights
func weightedRandomSelection(weights []float64) int {
	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}

	r := rand.Float64() * totalWeight
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return i
		}
	}

	return len(weights) - 1 // Should never reach here
}
