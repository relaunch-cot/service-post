package resource

import (
	"context"
	"log"

	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
	"github.com/relaunch-cot/service-post/config"
)

func OpenMysqlConn() *mysql.Client {
	ctx := context.Background()

	client, err := mysql.InitMySQL(ctx, config.MYSQL_USER, config.MYSQL_PASS, config.MYSQL_HOST, config.MYSQL_PORT, config.MYSQL_DBNAME)
	if err != nil {
		log.Fatal("failed to open myslq connection: ", err)
	}

	return client
}
