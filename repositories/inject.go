package repositories

import (
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"

	MysqlRepository "github.com/relaunch-cot/service-post/repositories/mysql"
)

type Repositories struct {
	Mysql MysqlRepository.IMySqlPost
}

func (r *Repositories) Inject(mysqlClient *mysql.Client) {
	r.Mysql = MysqlRepository.NewMysqlRepository(mysqlClient)
}
