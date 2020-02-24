package common

import "time"

func GetTimeSlot(genesisTime int64, pointInTime int64, slotTime int64) uint64 {
	slotTimeDur := time.Duration(slotTime)
	blockTime := time.Unix(pointInTime, 0)
	timePassed := blockTime.Sub(time.Unix(genesisTime, 0)).Round(slotTimeDur)
	timeSlot := uint64(int64(timePassed.Seconds()) / slotTime)
	return timeSlot
}
