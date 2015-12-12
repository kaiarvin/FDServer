// DataProcess
package ChatServer

import (
	"fmt"
	"qbs"
)

func DBInit(dbtype, dsn string) (*qbs.Qbs, error) {
	qbs.Register(dbtype, dsn, dbname, qbs.NewMysql())
	q, err := qbs.GetQbs()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return q, nil
}

func FindData()
