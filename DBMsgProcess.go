// DataProcess
package FDServer

import (
	"fmt"
	_ "mysql"
	"qbs"
)

func DBInit(dbtype, dsn, dbname string) (*qbs.Qbs, error) {
	qbs.Register(dbtype, dsn, dbname, qbs.NewMysql())
	q, err := qbs.GetQbs()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return q, nil
}
