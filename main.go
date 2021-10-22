package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

//const frameSize int64 = int64(15 * time.Second)

func main() {

	var (
		signals     = []os.Signal{os.Interrupt, os.Kill}
		dataJSON    = make(chan DataItem, 1)
		dataSSE     = make(chan DataItem, 1)
		timeRowJSON = &TimeRow{}
		timeRowSSE  = &TimeRow{}

		counterJSON = int64(0)
		sumJSON     = float64(0)

		counterSSE = int64(0)
		sumSSE     = float64(0)
		conf       = NewConfig()
	)

	ctx, cf := context.WithCancel(context.WithValue(context.Background(), "CONFIG", conf))

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, signals...)
	outPrintTicker := time.NewTicker(time.Duration(conf.UpdatePeriod()))
	if conf.Mode() == "ALL" || conf.Mode() == "JSON" {
		NewJSLProvider().Do(ctx, dataJSON)
	}
	if conf.Mode() == "ALL" || conf.Mode() == "SSE" {
		NewEventProvider().Do(ctx, dataSSE)
	}
	go func() {
		s := <-osSignal
		log.Printf("Shudown : %v", s)
		cf()
	}()

	for {
		select {
		case <-outPrintTicker.C:
			if conf.Mode() == "ALL" || conf.Mode() == "JSON" {
				log.Printf("JSON sma : %v  avg:%v", timeRowJSON.Average(), sumJSON/float64(counterJSON))
			}

			if conf.Mode() == "ALL" || conf.Mode() == "SSE" {
				log.Printf("SSE sma : %v  avg:%v", timeRowSSE.Average(), sumSSE/float64(counterSSE))
			}

		case <-ctx.Done():
			outPrintTicker.Stop()
			if conf.Mode() == "ALL" || conf.Mode() == "JSON" {
				log.Printf("JSON sma : %v  avg:%v", timeRowJSON.Average(), sumJSON/float64(counterJSON))
			}
			if conf.Mode() == "ALL" || conf.Mode() == "SSE" {
				log.Printf("SSE sma : %v  avg:%v ", timeRowSSE.Average(), sumSSE/float64(counterSSE))
			}
			return

		case v := <-dataJSON:
			timeRowJSON.Add(v)
			timeRowJSON.UpdateDataFrame(conf.TimeFrameWidth())
			sumJSON += v.Value
			counterJSON++

		case v := <-dataSSE:
			timeRowSSE.Add(v)
			timeRowSSE.UpdateDataFrame(int64(conf.TimeFrameWidth()) /*conf.TimeFrameWidth()*/)
			sumSSE += v.Value
			counterSSE++

		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}
