package limiter

import (
	"sync"
	"time"
)

var (
	actions          int
	actionsThreshold int
	mutex            sync.Mutex
	urgent           chan bool
	normal           chan bool
	reset            chan bool
)

//StartRateLimiter starts the rate limiter
func StartRateLimiter(maxActions int) {

	actionsThreshold = maxActions
	reset = make(chan bool)
	urgent = make(chan bool)
	normal = make(chan bool)

	go limitRates()
	go handleRequests()

}

func limitRates() {

	timeToWait := 1 * time.Second
	for {
		time.Sleep(timeToWait)
		mutex.Lock()

		if actions == actionsThreshold {
			reset <- true
		}

		actions = 0
		mutex.Unlock()
	}
}

func increaseActions() {
	mutex.Lock()
	actions++
	mutex.Unlock()
}

func actionsThresholdReached() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return actions == actionsThreshold
}

func handleRequests() {
	for {

		//If we have already reached the maximum
		//amount of actions allowed in our time slice,
		//we need to wait for a signal.
		if actionsThresholdReached() {
			<-reset
		}

		increaseActions()

		//Prioritize urgent actions.
		select {
		case <-urgent:
			continue
		default:
			select {
			case <-urgent:
				continue
			case <-normal:
				continue
			}
		}
	}
}
