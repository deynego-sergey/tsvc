package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	JsonUrl string = "http://78.47.241.61:24433/stream"
	EvtUrl  string = "http://78.47.241.61:24433/sse"
)

type (
	//
	SourceDate struct {
		Data float64 `json:"data"`
	}
	//
	DataItem struct {
		Value float64
		Time  int64
	}

	IDataProvider interface {
		Do(ctx context.Context, data chan<- DataItem) error
	}

	EventDataProvider struct{}

	JSLDataProvider struct{}
)

func NewEventProvider() IDataProvider {
	return &EventDataProvider{}
}

func NewJSLProvider() IDataProvider {
	return &JSLDataProvider{}
}

func (js *JSLDataProvider) Do(ctx context.Context, ch chan<- DataItem) (e error) {
	if request, e := http.NewRequest(http.MethodGet, JsonUrl, nil); e == nil {
		client := &http.Client{}

		if response, e := client.Do(request); e == nil {
			reader := bufio.NewReader(response.Body)
			decoder := json.NewDecoder(reader)
			c, _ := context.WithCancel(ctx)
			go func(ctx context.Context, ch chan<- DataItem) {
				for {
					select {
					case <-ctx.Done():
						response.Body.Close()
						return

					default:
						if decoder.More() {
							var src SourceDate
							if e = decoder.Decode(&src); e == nil {
								data := DataItem{
									Value: src.Data,
									Time:  time.Now().UnixNano(),
								}
								ch <- data
							}
						}
					}
				}
			}(c, ch)
		}
	}
	return
}

func (ev *EventDataProvider) Do(ctx context.Context, ch chan<- DataItem) (e error) {
	// смотрим в конфиге, нужно ли запускать этот провадер

	if request, e := http.NewRequest(http.MethodGet, EvtUrl, nil); e == nil {
		client := &http.Client{}
		if response, e := client.Do(request); e == nil {
			reader := bufio.NewReader(response.Body)
			c, _ := context.WithCancel(ctx)
			go func(tx context.Context, ch chan<- DataItem) {
				for {
					select {

					case <-ctx.Done():
						response.Body.Close()
						return

					default:
						if srcString, e := reader.ReadString('\n'); e == nil {
							//log.Println(srcString)
							re := regexp.MustCompile("^data:")
							if re.MatchString(srcString) {
								var src SourceDate
								if e = json.Unmarshal([]byte(re.ReplaceAllString(srcString, "")), &src); e == nil {
									ch <- DataItem{
										Value: src.Data,
										Time:  time.Now().UnixNano(),
									}
								} else {
									log.Println(e)
								}
							}
						}
					}
				}
			}(c, ch)
		}
	}
	return
}
