package serverinterface

import "golang.org/x/net/context"

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}
