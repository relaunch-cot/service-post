package mysql

import "github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"

type mysqlResource struct {
	mysqlClient *mysql.Client
}

type IMySqlPost interface {
}

func NewMysqlRepository(mysqlClient *mysql.Client) IMySqlPost {
	return &mysqlResource{
		mysqlClient: mysqlClient,
	}
}
