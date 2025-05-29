package server

import (
	"context"
	"net/url"
)

// Server for abstract interface
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Endpointer is registry endpoint
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
