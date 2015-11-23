// ChatServer project main.go
package main

import (
	"database/sql"
	"fmt"
)

func main() {
	db := sql.Open("mysql", "root/arvin@tcp(localhost:3306)/test?cha-utf-8")
	fmt.Println("Hello World!")
	fmt.Println("This's ok!!!")
}
