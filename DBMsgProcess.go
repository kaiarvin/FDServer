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

type Character struct {
	Id    int `qbs:"pk"`
	Level int
	Name  string
	Sex   int8
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

func ProcessRegistAccountId(q *qbs.Qbs, data *Account, Gname string) int {
	accountTB := &Account{}
	err := q.WhereEqual("user_name", data.UserName).Find(accountTB)
	if err == nil {
		fmt.Println("Have This UserName:", accountTB)
		return 1
	}

	n, err := q.Save(data)
	if err != nil {
		fmt.Println("Save Error:", err)
		return 2
	}
	if n != 1 {
		fmt.Print("Save n:", n, data)
		return 3
	} else {
		err := q.WhereEqual("user_name", data.UserName).Find(accountTB)
		if err != nil {
			fmt.Println("Cannt Find UserName:", data.UserName)
			return 4
		}

		//rel, err1 := q.Exec("INSERT INTO `character` (id,level,name,sex) VALUE (?,0,?,0)", accountTB.Id, Gname)
		character := &Character{Id: accountTB.Id, Name: accountTB.UserName}
		n, err1 := q.Save(character)
		if err1 != nil || n != 1 {
			fmt.Println("Save character Fail. ERROR: ", err1)
			return 5
		}

		fmt.Println("Save Success.")
		return 0
	}
}

func CheckUserEnterWorld(q *qbs.Qbs, data *ReqUserEnterServer) *Account {
	account := new(Account)
	err := q.WhereEqual("user_name", data.Uname).Find(account)
	if err != nil {
		fmt.Println("EnterWorld Error[", "Uname: ", data.Uname, "Error: ", err)
		return nil
	}

	return account
}

func GetDBCharacter(q *qbs.Qbs, id int) *Character {
	character := new(Character)
	err := q.WhereEqual("id", id).Find(character)
	if err != nil {
		fmt.Println("GetDBCharacter Error[", "id: ", id, "Error: ", err, "]")
		return nil
	}

	return character
}
