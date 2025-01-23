# Deathmatch Simulation in Go

This program simulates a simple deathmatch scenario where multiple players randomly kill each other until one player reaches a specified kill limit. It's built using Go's concurrency features, including goroutines and channels, to demonstrate an idiomatic approach to managing game state and event ordering.

**Features:**

*   **Concurrent Players:** Each player is a goroutine.
*   **Centralized Game State:** A "game manager" goroutine controls the game state and ensures events are processed in order.
*   **Random Kill Events:** Players randomly select another player to kill.
*   **Channel-Based Termination:** Uses a `doneChan` channel for clean goroutine termination, avoiding shared memory flags and potential deadlocks.
*   **Configurable:** Easily adjust the number of players and the winning kill count.

**How to Run:**

1. Make sure you have Go installed ([https://go.dev/]).
2. Save the code as a `.go` file (e.g., `main.go`).
3. Run from your terminal:

    ```bash
    cd client
    go run client.go
    ```

**Example Output:**