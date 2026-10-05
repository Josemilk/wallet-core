package ratelimit

import("sync";"time")
type Limiter struct { mu sync.Mutex; hits map[string][]time.Time; Limit int; Window time.Duration }
func New(limit int,window time.Duration)*Limiter{return &Limiter{hits:map[string][]time.Time{},Limit:limit,Window:window}}
func(l *Limiter) Allow(key string)bool{l.mu.Lock();defer l.mu.Unlock();now:=time.Now();cut:=now.Add(-l.Window); xs:=l.hits[key]; n:=0;for _,t:=range xs{if t.After(cut){xs[n]=t;n++}};xs=xs[:n];if len(xs)>=l.Limit{return false};l.hits[key]=append(xs,now);return true}
