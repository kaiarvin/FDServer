// Timer
package FDServer

import (
	"fmt"
	"time"
)

type FDTimer struct {
	TimerName string
	StartTime int64
	EndTime   int64
	Second    int64
	RunFunc   func()
}

type TestTimer struct {
}

func (this *Server) RunTimer() {

	testTime := new(TestTimer)
	fmt.Println("SetTime 0")
	this.MinuteTimers = append(this.MinuteTimers, new(FDTimer))
	this.MinuteTimers[0].SetTime("Test", 1, 0, 60, testTime.CallTime)

	this.MinuteTimers = append(this.MinuteTimers, new(FDTimer))
	this.MinuteTimers[1].SetTime("Test", 2, 0, 60, testTime.CallTime)

	SecondTime := time.Now()
	MinuteTime := time.Now()
	for {
		if len(this.SecondTimers) == 0 && len(this.MinuteTimers) == 0 && len(this.OneTimesTimers) == 0 {

			return
		}
		if time.Now().After(SecondTime.Add(1 * time.Second)) {
			nowUnixTime := time.Now().Unix()
			for i, fdTime := range this.SecondTimers {
				if fdTime.StartTime <= nowUnixTime {
					fdTime.RunFunc()
					fdTime.StartTime = nowUnixTime + fdTime.Second
				}
				if fdTime.EndTime <= nowUnixTime && fdTime.EndTime != 0 {
					this.SecondTimers[i] = new(FDTimer)
				}
			}
			SecondTime = time.Now()
		}

		if time.Now().After(MinuteTime.Add(1 * time.Minute)) {
			nowUnixTime := time.Now().Unix()
			for i, fdTime := range this.MinuteTimers {
				if fdTime.StartTime <= nowUnixTime {
					fdTime.RunFunc()
					fdTime.StartTime = nowUnixTime + fdTime.Second
				}
				if fdTime.EndTime <= nowUnixTime && fdTime.EndTime != 0 {
					this.MinuteTimers[i] = new(FDTimer)
				}
			}
			MinuteTime = time.Now()
		}
	}
}

func (this *FDTimer) SetTime(name string, start, end, second int64, run func()) {
	if end < start && end != 0 {
		return
	}
	this.TimerName = name
	this.StartTime = time.Now().Unix() + start*second
	if end != 0 {
		this.EndTime = time.Now().Unix() + end*second
	} else {
		this.EndTime = 0
	}

	this.Second = second
	this.RunFunc = run
}

func (this *TestTimer) CallTime() {
	fmt.Println("Call Timer 0!")
}
