// FDLoger
package FDServer

import (
	"fmt"
	"log"
	"os"
)

type FDLog struct {
	Level int32
	Addr  string
	Loger *log.Logger
}

func (this *FDLog) InitFDLog(LogLevel int32, name string) {
	fmt.Println("InitFDLog")
	this.Level = LogLevel
	this.Addr = "./log/" + name
	file, err := os.Create(this.Addr)
	if err != nil {
		fmt.Println("Create LogFile Fail", this.Addr)
		return
	}

	this.Loger = log.New(file, "", log.LstdFlags|log.Llongfile)

}

func (this *FDLog) WriteLog(data string, arg ...interface{}) {
	fmt.Println("WriteLog:", Const_Log_CanWrite, this.Level)
	if Const_Log_CanWrite >= this.Level {
		this.Loger.Printf(data, arg...)
	}
}
