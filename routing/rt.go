package routing

import (
	"container/list"
	"errors"
	"log"
	"strings"
	"sync"
)

// Errors
var (
	errInvalidNodeValue      = errors.New("Invalid node value. ")
	errInvalidSubscribeTopic = errors.New("Invalid subscribe topic. ")
	errPatternIsPresenr      = errors.New("Pattern already present in list. ")
)

type (
	// IRoutePattern - запись содержащая маршрут
	IRoutePattern interface {
		Agent() Agent
		Pattern() string
		Match(topic string) bool
		Subscribe(topic string) bool
	}
	// RoutePattern -
	RoutePattern struct {
		broker  Agent
		pattern string
		rt      *list.List
	}
	ISubscriptionTable interface {
		Add(pattern string) error
		Remove(pattern string)
		Match(topic string) bool
	}
	// ISubscriptionPattern -
	ISubscriptionPattern interface {
		Match(topic string) bool
		Pattern() string
	}
	// IRoute -
	IRoute interface {
		Broker() Agent
		Route() string
	}
	// Route -
	Route struct {
		broker Agent
		topic  string
	}
	//
	// IRouteTable - таблица маршрутов
	//
	IRouteTable interface {
		// Remove - удаляет маршрут
		Remove(broker Agent, pattern string)
		// Set - добавляет маршрут в таблицу
		Set(broker Agent, pattern string) error
		// Match - возвращает список маршрутов для
		// топика
		Match(topic string) []IRoute
	}
	Agent string
	// RouteTable -
	RouteTable struct {
		routes map[Agent][]IRoutePattern
	}

	SubscriptionPattern struct {
		p string
		t *list.List
	}

	// SubscriptionTable -
	SubscriptionTable struct {
		sub []ISubscriptionPattern
	}
)

func (r *Route) Broker() Agent {
	return r.broker
}

func (r *Route) Route() string {
	return r.topic
}

func NewRouteTable() IRouteTable {
	t := &RouteTable{routes: make(map[Agent][]IRoutePattern)}
	return t
}

// Set - добавляет в таблицу информацию о брокере и паттерн для роутинга
func (rt *RouteTable) Set(broker Agent, pattern string) (e error) {

	ptrn, e := NewRoutePattern(broker, pattern)

	if e != nil {
		return e
	}
	mu := sync.Mutex{}
	mu.Lock()
	rt.routes[broker] = append(rt.routes[broker], ptrn)
	mu.Unlock()
	return
}

// Remove - Удаляет паттерн роутинга
func (rt *RouteTable) Remove(broker Agent, pattern string) {
	mu := sync.Mutex{}
	mu.Lock()
	if _, ok := rt.routes[broker]; ok {
		var p []IRoutePattern
		for _, v := range rt.routes[broker] {
			if v.Pattern() == pattern {
				continue
			}
			p = append(p, v)
		}
		rt.routes[broker] = p
	}
	mu.Unlock()
}

// Match - находит в таблице маршруты для подключения к брокерам для заданного топика
// возвращает список брокеров для подключения с заданным топиком
func (rt *RouteTable) Match(topic string) (routes []IRoute) {
	mu := sync.Mutex{}
	mu.Lock()
	for _, v := range rt.routes {
		for _, p := range v {
			if p.Match(topic) {
				routes = append(routes, &Route{
					broker: p.Agent(),
					topic:  topic,
				})
			}
		}
	}
	mu.Unlock()
	return
}

//

// NewRoutePattern - создает паттерн для подписки
//
func NewRoutePattern(broker Agent, pattern string) (IRoutePattern, error) {

	r := &RoutePattern{
		broker:  broker,
		pattern: pattern,
		rt:      list.New(),
	}
	if e := r.create(pattern); e != nil {
		return nil, e
	}
	return r, nil
}

// NewSubscribePattern -
func NewSubscribePattern(pattern string) (ISubscriptionPattern, error) {
	s := &SubscriptionPattern{
		p: pattern,
		t: list.New(),
	}
	if e := s.create(pattern); e != nil {
		return nil, e
	}
	return s, nil
}

//
// create - создает структуру для паттерна подписки
func (sp *SubscriptionPattern) create(s string) error {

	sp.t.Init()
	isTail := false

	for _, v := range strings.Split(strings.Trim(s, charDelimiter), charDelimiter) {

		if isTail {
			return errInvalidSubscribeTopic
		}

		if nodeValue, e := createNodeValue(v); e != nil {
			return e
		} else {
			switch nodeValue.Type() {
			case nodeTypeSuffix, nodeTypePrefix:
				return errInvalidSubscribeTopic
			case nodeTypeTail:
				isTail = true
			}
			sp.t.PushBack(nodeValue)
		}
	}
	return nil
}

func (sp *SubscriptionPattern) Pattern() string {
	return sp.p
}

// Match - сравнивает топик с паттерном подписки
func (sp *SubscriptionPattern) Match(topic string) (result bool) {

	tn := strings.Split(strings.Trim(topic, charDelimiter), charDelimiter)
	if len(tn) < 1 {
		return false
	}
	ptr := sp.t.Front()

	for _, v := range tn {
		if ptr == nil {
			return false
		}
		switch ptr.Value.(INode).Type() {
		case nodeTypeString:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		case nodeTypePlus:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		case nodeTypeTail:
			return ptr.Value.(INode).Validate(v)

		default:
			return false
		}
		ptr = ptr.Next()
	}
	return ptr.Value == nil
	return
}

func NewSubscriptionTable() ISubscriptionTable {
	return &SubscriptionTable{
		sub: nil,
	}
}

//
func (st *SubscriptionTable) Match(topic string) bool {
	for _, v := range st.sub {
		if v.Match(topic) {
			return true
		}
	}
	return false
}

//
func (st *SubscriptionTable) Add(pattern string) (e error) {
	if p, e := NewSubscribePattern(pattern); e == nil {
		for _, v := range st.sub {
			if v.Pattern() == pattern {
				return errPatternIsPresenr
			}
		}
		st.sub = append(st.sub, p)
	}
	return
}

//
func (st *SubscriptionTable) Remove(pattern string) {
	if len(st.sub) > 0 {
		var tmp []ISubscriptionPattern
		mu := sync.Mutex{}
		mu.Lock()
		for _, v := range st.sub {
			if v.Pattern() != pattern {
				tmp = append(tmp, v)
			}
		}
		st.sub = tmp
		mu.Unlock()
	}
}

// Agent -
func (rp *RoutePattern) Agent() Agent {
	return Agent(rp.broker)
}

// Pattern -
func (rp *RoutePattern) Pattern() string {
	return rp.pattern
}

// Match - проверяет возможность подписки для топика
func (rp *RoutePattern) Match(topic string) bool {
	if len(topic) > 1 {
		if rp.rt != nil {
			return rp.find(topic)
		}
	}
	return false
}

// Subscribe -
func (rp *RoutePattern) Subscribe(topic string) bool {
	if len(topic) > 1 {
		if rp.rt != nil {
			return rp.subscribe(topic)
		}
	}
	return false
}

// create -
func (rp *RoutePattern) create(route string) error {
	//ptr := list.New()
	rp.rt.Init()

	for _, v := range strings.Split(strings.Trim(route, charDelimiter), charDelimiter) {
		if nodeValue, e := createNodeValue(v); e != nil {
			return e
		} else {
			rp.rt.PushBack(nodeValue)
		}
	}
	return nil
}

// find - найти соответствие
func (rp *RoutePattern) find(topic string) bool {
	tn := strings.Split(strings.Trim(topic, charDelimiter), charDelimiter)
	if len(tn) < 1 {
		return false
	}
	ptr := rp.rt.Front()

	for _, v := range tn {
		switch ptr.Value.(INode).Type() {
		case nodeTypeString:
			//log.Printf("nodeTypeString %v", v )
			if !ptr.Value.(INode).Validate(v) {
				return false
			}

		case nodeTypePlus:
			//log.Printf("nodeTypePlus %v", v)
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		case nodeTypePrefix:
			//log.Printf("nodeTypePrefix %v", v)
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		case nodeTypeSuffix:
			//log.Printf("nodeTypeSuffix %v", v)
			if !ptr.Value.(INode).Validate(v) {
				return false
			}

		case nodeTypeTail:
			//log.Printf("nodeTypeTail %v", v)
			//if !ptr.Value.(INode).Validate(v){
			//	return false
			//}
			return ptr.Value.(INode).Validate(v)
		default:
			log.Printf("Error. Unknown node type %v", v)
			return false
		}
		ptr = ptr.Next()
	}
	return ptr.Value == nil
}

// subscribe - проверить возможность подписки
func (rp *RoutePattern) subscribe(topic string) bool {
	ptr := rp.rt.Front()

	for _, v := range strings.Split(strings.Trim(topic, "/ "), "/") {
		switch v {
		case "+":
			ptr = ptr.Next()
		case "#":
			return true
		default:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		}
		ptr = ptr.Next()

		if ptr == nil {
			return false
		}

	}
	return true
}

//func NewRouteTree(pattern string) IRouteTree {
//
//}
