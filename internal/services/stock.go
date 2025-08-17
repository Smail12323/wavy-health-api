package services

import (
	"math/rand"
	"sync"
	"time"
)

var (
	// StockPrice holds the current price of the stock
	StockPrice float64

	// mu is a read/write mutex to safely access StockPrice concurrently
	mu sync.RWMutex
)

// InitStockPrice initializes the stock price with a given starting value
// and starts a background goroutine to simulate price changes over time
func InitStockPrice(initial float64) {
	StockPrice = initial
	go simulate()
}

// simulate runs in a goroutine and randomly changes the stock price every second
func simulate() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	for {
		delta := float64(rand.Intn(4)+2) * 0.1 // Random change between 0.2 and 0.5

		mu.Lock() // Lock for writing
		if rand.Float64() < 0.6 {
			StockPrice += delta // 60% chance to increase price
		} else {
			StockPrice -= delta // 40% chance to decrease price
			if StockPrice < 0 {
				StockPrice = 0 // Prevent negative stock price
			}
		}
		mu.Unlock()

		time.Sleep(1 * time.Second) // Update every second
	}
}

// GetPrice returns the current stock price in a thread-safe manner
func GetPrice() float64 {
	mu.RLock()         // Lock for reading
	defer mu.RUnlock() // Unlock after reading
	return StockPrice
}
