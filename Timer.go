// Timer
package FDServer

import (
	"time"
)

type FDTimer struct {
	StartTime int64
	EndTime   int64
	RunFunc   func()
}

func (this *Server) RunTimer() {

	SecondTime := time.Now()
	MinuteTime := time.Now()
	for {
		if len(this.SecondTimeers) == 0 && len(this.MinuteTimers) == 0 && len(this.OneTimesTimers) == 0 {
			return
		}

		if time.Now().After(SecondTime.Add(1 * time.Second)) {
			nowUnixTime := time.Now().Unix()
			for i, fdTime := range this.SecondTimeers {
				if fdTime.StartTime <= nowUnixTime {
					fdTime.RunFunc()
				}
				if fdTime.EndTime <= nowUnixTime && fdTime.EndTime != 0 {
					this.SecondTimeers[i] = new(FDTimer)
				}
			}
			SecondTime = time.Now()
		}

		if time.Now().After(MinuteTime.Add(1 * time.Minute)) {
			nowUnixTime := time.Now().Unix()
			for i, fdTime := range this.MinuteTimers {
				if fdTime.StartTime <= nowUnixTime {
					fdTime.RunFunc()
				}
				if fdTime.EndTime <= nowUnixTime && fdTime.EndTime != 0 {
					this.MinuteTimers[i] = new(FDTimer)
				}
			}
			MinuteTime = time.Now()
		}
	}
}

func (this *FDTimer) SetTime(start, end int64, run func()) {
	this.StartTime = start
	this.EndTime = end
	this.RunFunc = run
}
