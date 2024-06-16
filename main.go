package main

import (
	"github.com/Go-routine-4595/how-oem-view/controllers/api"
	"github.com/Go-routine-4595/how-oem-view/controllers/datasource"
	"github.com/Go-routine-4595/how-oem-view/model"
	"os"
)

func main() {
	var (
		service model.IService
		apiSrv  api.ApiController
		data    *datasource.DataSource
		log     model.IService
		args    []string
	)
	args = os.Args

	if len(args) == 1 {
		data = datasource.NewDataSource("")
	} else {
		data = datasource.NewDataSource(args[1])
	}

	service = model.NewService(data)
	log = NewLoggingService(service)
	apiSrv = api.NewApiController(log, data.GetKeys())

	apiSrv.Run()
}
