package resource

import (
	"github.com/relaunch-cot/service-post/handler"
	"github.com/relaunch-cot/service-post/repositories"
	"github.com/relaunch-cot/service-post/server"
)

var Repositories repositories.Repositories
var Handler handler.Handlers
var Server server.Servers

func Inject() {
	mysqlClient := OpenMysqlConn()

	Repositories.Inject(mysqlClient)
	Handler.Inject(&Repositories)
	Server.Inject(&Handler)
}
