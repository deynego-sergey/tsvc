package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const frameSize int64 = int64(15 * time.Second)

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
			log.Printf("JSON sma : %v  avg:%v   len:%v", timeRowJSON.Average(), sumJSON/float64(counterJSON), len(*timeRowJSON))
			log.Printf("SSE sma : %v  avg:%v   len:%v", timeRowSSE.Average(), sumSSE/float64(counterSSE), len(*timeRowSSE))

		case <-ctx.Done():
			outPrintTicker.Stop()
			log.Printf("JSON sma : %v  avg:%v   len:%v", timeRowJSON.Average(), sumJSON/float64(counterJSON), len(*timeRowJSON))
			log.Printf("SSE sma : %v  avg:%v   len:%v", timeRowSSE.Average(), sumSSE/float64(counterSSE), len(*timeRowSSE))
			return

		case v := <-dataJSON:
			timeRowJSON.Add(v)
			timeRowJSON.UpdateWindow(frameSize)
			sumJSON += v.Value
			counterJSON++
		//	log.Printf( "JSON sma : %v  avg:%v   len:%v", timeRowJSON.Average(), sumJSON/float64(counterJSON) , len(*timeRowJSON))

		case v := <-dataSSE:
			timeRowSSE.Add(v)
			timeRowSSE.UpdateWindow(frameSize)
			sumSSE += v.Value
			counterSSE++
		//	log.Printf( "SSE sma : %v  avg:%v   len:%v", timeRowSSE.Average(), sumSSE/float64(counterSSE) , len(*timeRowSSE))

		default:
			runtime.Gosched()
		}
	}
}
