package main

//make init proto update tidy

import (
	"github.com/rs/zerolog"
	"os"
	"time"

	"github.com/Go-routine-4595/how-oem-view/model"
)

type LoggingService struct {
	next model.IService
	log  zerolog.Logger
}

func NewLoggingService(next model.IService) model.IService {

	return &LoggingService{
		next: next,
		log:  zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.InfoLevel).With().Timestamp().Logger(),
	}

}

func (s *LoggingService) DrawOEMChartSeries(val string) (b []byte, err error) {

	defer func(start time.Time) {
		s.log.Info().
			Str("method", "DrawOEMChartSeries").
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.DrawOEMChartSeries(val)
}

func (s *LoggingService) DrawOEMChartBar(val string) (b []byte, err error) {

	defer func(start time.Time) {
		s.log.Info().
			Str("method", "DrawOEMChartBar").
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.DrawOEMChartBar(val)
}

func (s *LoggingService) DrawOEMBarPlot(val string) (b []byte, err error) {

	defer func(start time.Time) {
		s.log.Info().
			Str("method", "DrawOEMBarPlot").
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.DrawOEMBarPlot(val)
}

func (s *LoggingService) DrawMap(val string) (b []byte, err error) {

	defer func(start time.Time) {
		s.log.Info().
			Str("method", "DrawMap").
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.DrawMap(val)
}

func (s *LoggingService) StateGraph(val string) (b []byte, err error) {

	defer func(start time.Time) {
		s.log.Info().
			Str("method", "StateGraph").
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.StateGraph(val)
}
