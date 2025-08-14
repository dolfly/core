package handler

import (
	"context"
	"net"

	"github.com/dolfly/core/hop"
	"github.com/dolfly/core/metadata"
)

type Handler interface {
	Init(metadata.Metadata) error
	Handle(context.Context, net.Conn, ...HandleOption) error
}

type Forwarder interface {
	Forward(hop.Hop)
}
