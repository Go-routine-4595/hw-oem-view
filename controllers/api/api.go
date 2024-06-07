package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Go-routine-4595/how-oem-view/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type ApiController struct {
	srv  model.IService
	keys map[string]string
}

func NewApiController(s model.IService, k map[string]string) ApiController {
	return ApiController{
		srv:  s,
		keys: k,
	}
}

func (a ApiController) Run() {

	e := echo.New()
	e.Debug = true
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Dur("Latency", v.Latency).
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	e.File("/favicon.ico", "images/favicon.ico")
	//e.GET("/favicon.ico", func(c echo.Context) error {
	//	return c.Blob(http.StatusOK, "image/png", []byte{})
	//})
	e.GET("/menu", a.testMenuEcho)
	e.GET("/oem-serie", a.drawOEMChartSeries)
	//e.GET("/custom", drawCustomChart)
	e.GET("/oem-bar", a.drawOEMChartBar)
	e.GET("/test", a.drawOEMBarPlot)
	e.GET("/map", a.drawMap)
	e.GET("/state", a.StateGraph)

	e.Logger.Fatal(e.Start(":8080"))

	//http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
	//	res.Write([]byte{})
	//})

}

func (a ApiController) drawOEMChartSeries(c echo.Context) error {

	// Query parameters are returned as a map[string][]string
	// params := req.URL.Query()
	// To get the value of a specific key:
	// (assuming the parameter exists)
	// val := params.Get("key")
	// //res.Header().Set("Content-Type", "image/png")

	val := c.QueryParam("key")

	res, _ := a.srv.DrawOEMChartSeries(val)
	c.Blob(http.StatusOK, "image/png", res)
	return nil
}

func (a ApiController) drawOEMChartBar(c echo.Context) error {

	val := c.QueryParam("key")

	res, _ := a.srv.DrawOEMChartBar(val)
	c.Blob(http.StatusOK, "image/png", res)
	return nil
}

func (a ApiController) drawOEMBarPlot(c echo.Context) error {

	val := c.QueryParam("key")

	res, _ := a.srv.DrawOEMChartBar(val)
	c.Blob(http.StatusOK, "image/png", res)
	return nil
}

func (a ApiController) drawMap(c echo.Context) error {

	val := c.QueryParam("key")

	res, err := a.srv.DrawMap(val)
	if err != nil {
		var message *bytes.Buffer

		message = new(bytes.Buffer)
		if strings.Contains(err.Error(), "No key found") {
			fmt.Fprintln(message, "<html><body>")
			fmt.Fprintf(message, "<h1>Missing equipement to map </h1>")
			c.HTML(http.StatusNotFound, string(message.String()))
			return err
		}
		if strings.Contains(err.Error(), "no value found") {
			fmt.Fprintln(message, "<html><body>")
			fmt.Fprintf(message, "<h1>No value for equipement: %s</h1>", val)
			c.HTML(http.StatusNotFound, string(message.String()))
			return err
		}
	}
	c.Blob(http.StatusOK, "image/png", res)
	return nil
}

func (a ApiController) StateGraph(c echo.Context) error {

	val := c.QueryParam("key")

	res, err := a.srv.StateGraph(val)
	if err != nil {
		var message *bytes.Buffer

		message = new(bytes.Buffer)
		if strings.Contains(err.Error(), "key not found") {
			fmt.Fprintln(message, "<html><body>")
			fmt.Fprintf(message, "<h1>Missing equipement to map </h1>")
			c.HTML(http.StatusNotFound, string(message.String()))
			return err
		}
		if strings.Contains(err.Error(), "no value found") {
			fmt.Fprintln(message, "<html><body>")
			fmt.Fprintf(message, "<h1>No value for equipement: %s</h1>", val)
			c.HTML(http.StatusNotFound, string(message.String()))
			return err
		}
	}
	c.Blob(http.StatusOK, "image/png", res)
	return nil
}
