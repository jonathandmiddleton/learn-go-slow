package app

import (
	"net/http"
	"regexp"
)

type CORSMiddleware struct {
	OriginRule string
}

func NewOrigin(origin string) *CORSMiddleware {
	return &CORSMiddleware{
		OriginRule: origin,
	}
}

func (m *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "OPTIONS" {
		if m.allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	next(w, r)
}

func (m *CORSMiddleware) allowedOrigin(origin string) bool {
	if m.OriginRule == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(m.OriginRule, origin); matched {
		return true
	}
	return false
}