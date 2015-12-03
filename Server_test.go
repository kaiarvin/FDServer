// ChatServer project main.go
//package main

//import (
//	"database/sql"
//	"fmt"
//)

//func main() {
//	db := sql.Open("mysql", "root/arvin@tcp(localhost:3306)/test?cha-utf-8")
//	fmt.Println("Hello World!  123123123")
//	fmt.Println("This's ok!!!")
//}
package ChatServer

import "testing"

func TestServer(t *testing.T) {
	var s ChatServer
	s.ReadConf("./ChatServer.conf")

	//	if port, ok := s.ChatConfigData["ChatPort"]; ok {
	//		if a, err := strconv.Atoi(port); err != nil {
	//			fmt.Println("Conve ChatPort to int Error!")
	//			return
	//		} else {
	//			s.Port = a
	//			fmt.Println(a)
	//		}
	//	} else {
	//		fmt.Println("Cant find ChatPort")
	//	}

	//	fmt.Println(s.Port)

	//	fmt.Println("ChatServer .stop")

	s.InitServer()
}
