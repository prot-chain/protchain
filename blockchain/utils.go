package blockchain

import "fmt"

func LogPerformance(action string, duration int64) {
	fmt.Printf("Action: %s, Duration: %dms\n", action, duration)
}
