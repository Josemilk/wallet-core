package marketdata

import (
    "encoding/json"
    "sync"
)

type Tick struct { Market string `json:"market"`; Price string `json:"price"`; Quantity string `json:"quantity"`; Sequence uint64 `json:"sequence"` }
type Subscriber chan []byte

type Stream struct { mu sync.RWMutex; subs map[string]map[Subscriber]struct{} }
func NewStream() *Stream { return &Stream{subs:make(map[string]map[Subscriber]struct{})} }
func (s *Stream) Subscribe(market string, ch Subscriber) { s.mu.Lock(); defer s.mu.Unlock(); if s.subs[market]==nil { s.subs[market]=map[Subscriber]struct{}{} }; s.subs[market][ch]=struct{}{} }
func (s *Stream) Unsubscribe(market string, ch Subscriber) { s.mu.Lock(); defer s.mu.Unlock(); delete(s.subs[market],ch) }
func (s *Stream) Publish(t Tick) { b,err:=json.Marshal(t); if err!=nil{return}; s.mu.RLock(); defer s.mu.RUnlock(); for ch:=range s.subs[t.Market] { select { case ch<-b: default: } } }
