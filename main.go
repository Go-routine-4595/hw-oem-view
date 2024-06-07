package main

import (
	"github.com/Go-routine-4595/how-oem-view/controllers/api"
	"github.com/Go-routine-4595/how-oem-view/controllers/datasource"
	"github.com/Go-routine-4595/how-oem-view/model"
)

func main() {
	var (
		service model.IService
		apiSrv  api.ApiController
		data    *datasource.DataSource
		log     model.IService
	)
	data = datasource.NewDataSource("UAS_external_event_queue.jsonl")
	service = model.NewService(data)
	log = NewLoggingService(service)
	apiSrv = api.NewApiController(log, data.GetKeys())

	apiSrv.Run()
}
