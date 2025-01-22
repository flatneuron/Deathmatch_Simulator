package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Player struct
type Player struct {
	ID    int
	Kills int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numPlayers := 50
	maxKills := 40

	players := make([]Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = Player{ID: i + 1, Kills: 0}
	}

	gameOver := false
	for !gameOver {
		// Select a random player
		randomIndex := rand.Intn(numPlayers)

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
