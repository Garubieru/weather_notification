package value_objects

import "github.com/nrednav/cuid2"

type ID struct {
	Value string
}

func NewID() ID {
	return ID{Value: cuid2.Generate()}
}

func RecoverID(value string) ID {
	return ID{Value: value}
}
