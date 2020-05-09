package server

import "reflect"

type Server struct {
	addr string
	funcs map[string]reflect.Value
}

// Register a method via name
func (s *Server) Register(name string, f interface{}) {
	if _, ok := s.funcs[name]; ok {
		return
	}
	s.funcs[name] = reflect.ValueOf(f)
}

