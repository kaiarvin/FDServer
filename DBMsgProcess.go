// DataProcess
package FDServer

import (
	"fmt"
	_ "mysql"
	"qbs"
)

type Account struct {
	Id       int `qbs:"pk"`
	UserName string
	UserPw   string
}

func DBInit(dbtype, dsn, dbname string) (*qbs.Qbs, error) {
	qbs.Register(dbtype, dsn, dbname, qbs.NewMysql())
	q, err := qbs.GetQbs()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return q, nil
}

func ProcessRegistAccountId(q *qbs.Qbs, data *Account) {
	fmt.Println("ProcessRegistAccountId:", data)
	accountTB := &Account{}
	err := q.WhereEqual("user_name", data.UserName).Find(accountTB)
	if err == nil {
		fmt.Println("Have This UserName:", accountTB)
		return
	}

	n, err := q.Save(data)
	if err != nil {
		fmt.Println("Save:", err)
		return
	}
	if n != 1 {
		fmt.Print("Save n:", n, data)
	} else {
		fmt.Println("Save Success!!!")
	}
}
