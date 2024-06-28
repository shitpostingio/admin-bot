package antiflood

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/config/structs"
	"github.com/shitpostingio/admin-bot/repository"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														STRUCTS														   *
 *																													   *
 ***********************************************************************************************************************
 */

//FloodState represent infos about reports generated in a group
type FloodState struct {
	handled bool
	mutex   sync.Mutex
	counter int
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												CONSTS AND VARS														   *
 *																													   *
 ***********************************************************************************************************************
 */

const (
	floodKey = "flood"
)

var (
	cfg   *structs.AntiFloodConfiguration
	state FloodState
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														START														   *
 *																													   *
 ***********************************************************************************************************************
 */

//Start starts flood monitoring
func Start() {
	cfg = repository.GetAntiFloodConfiguration()
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													FLOOD CONTROL													   *
 *																													   *
 ***********************************************************************************************************************
 */

//IncreaseFloodCounter increments by `amount` the flood counter
func IncreaseFloodCounter(amount int) {

	state.mutex.Lock()
	state.counter += amount

	if state.counter > cfg.Threshold {
		if !state.handled {
			state.handled = true
			reportFlood()
		}
	}

	state.mutex.Unlock()
	go decreaseFloodCounterAfterTime(amount)

}

//decreaseFloodCounterAfterTime decreases by `amount` the flood counter after `floodRoutineLifespan`
func decreaseFloodCounterAfterTime(amount int) {

	time.Sleep(time.Duration(cfg.RoutineLifeSpan) * time.Second)

	if IsFlood() {
		if repository.GetTestingStatus() {
			time.Sleep(15 * time.Second)
		} else {
			time.Sleep(5 * time.Minute)
		}
	}

	state.mutex.Lock()
	state.counter -= amount

	if state.counter < cfg.Threshold {
		state.handled = false
	}

	state.mutex.Unlock()

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														REPORTING													   *
 *																													   *
 ***********************************************************************************************************************
 */

//reportFlood reports a flood to the report channel
func reportFlood() {
	report := "⚠️⚠️⚠️WE ARE BEING FLOODED⚠️⚠️⚠️"
	log.Warn(report)
	_ = adminbot.SendPlainTextMessage(repository.GetTelegramConfiguration().ReportChannelID, report, true)
}

/*
 ***********************************************************************************************************************
 *																													   *
 *														ACCESSORS													   *
 *																													   *
 ***********************************************************************************************************************
 */

//IsFlood returns true if the chat is being flooded
func IsFlood() bool {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	return state.counter > cfg.Threshold
}
