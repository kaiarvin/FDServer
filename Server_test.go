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

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	var s ChatServer
	s.ReadConf("./ChatServer.conf", processLine)

	r := bufio.NewReader(os.Stdin)

	//for {
	b, _, _ := r.ReadLine()
	line := string(b)
	//tokens := strings.Split(line, " ")
	fmt.Println(line)
	//}
	fmt.Println("中文")
}
