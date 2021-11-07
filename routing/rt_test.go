package routing

import (
	"log"
	"testing"
)

//
//
//
func Test_NewRouterPattern(t *testing.T) {

	ptrn := "/topic1/topic2/+/+suffix1/prefix1+/#"
	topics := []string{
		"/topic1/topic2/sample/samplesuffix1/prefix1sample/test1/test2",
		"/topic1/topic2/sample/samplesuffix1/prefix1sample",
		"/topic1/topic2/sample/suffix1/prefix1sample/test1/test2",
		"/topic1/topic2/sample/samplesuffix1/prefix1/test1/test2",
		"/topic1/topic2/sample/samplesuffix1",
		"/topic1/topic2/suffix1/prefix1sample/test1/test2",
		"/topic1/topic2/sample/samplesuffix1/test1/test2",
		"/topic1/topic2/sample",
	}

	subs := []string{
		"topic1/+/+/1suffix1/",
		"+/+/#",
		"topic1/topic2/topic3/+/#",
		"topic1/topic3/+/#",
	}

	if p, e := NewRoutePattern("1", ptrn); e == nil {

		for _, v := range topics {
			log.Printf("%v %v %v", ptrn, v, p.Match(v))
		}
		for _, v := range subs {
			log.Printf("%v %v %v", ptrn, v, p.Subscribe(v))
		}
	} else {
		log.Printf("Error: %v", e)
	}
}

//
//
//
func Test_RouteTable(t *testing.T) {

}

//
//
//
func Test_SubscribeTable(t *testing.T) {
	table := NewSubscriptionTable()
	topicsLst := []string{
		"/topic1/topic2/topic3/topic4/topic5/topic6",
		"/topic1/topic2/topic3/topic4/topic5/topic7",
		"/topic1/topic2/topic3/topic4/topic5/topic9",
		"/topic1/topic3/topic2/topic4/topic5/topic6",
		"/topic1/topic2/topic6",
	}

	subPtrns := []string{
		"topic1/topic2/+/topic4",
		"topic1/topic2/topic3/+/topic5",
		"topic1/+/+/+/#",
		"topic1/++/#",
		"topic1/#/+/#",
	}
	for _, v := range subPtrns {
		if e := table.Add(v); e != nil {
			log.Printf("%v Error : %v", v, e)
		}
	}

	for _, v := range topicsLst {
		log.Printf("topic:%v  find: %v", v, table.Match(v))
	}

}
