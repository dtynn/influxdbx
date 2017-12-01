package rpc

import "errors"

var (
	// ErrMetaDisabled meta service disabled
	ErrMetaDisabled = errors.New("apis for meta service disabled")
)
