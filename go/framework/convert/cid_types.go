package convert

import (
	"errors"
	"fmt"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/openreserveio/dwn/go/log"
)

func FromUintConverter(node datamodel.Node) (interface{}, error) {

	if node.IsNull() {
		return 0, nil
	}

	uintBytes, err := node.AsBytes()
	if err != nil {
		return nil, err
	}

	if len(uintBytes) != 8 {
		return 0, fmt.Errorf("input byte slice must be exactly 8 bytes long")
	}

	uintValue := uint64(uintBytes[0])<<56 | uint64(uintBytes[1])<<48 | uint64(uintBytes[2])<<40 | uint64(uintBytes[3])<<32 | uint64(uintBytes[4])<<24 | uint64(uintBytes[5])<<16 | uint64(uintBytes[6])<<8 | uint64(uintBytes[7])
	log.Debug("uint value: %d", uintValue)

	return uintValue, nil

}

func ToUintConverter(data interface{}) (datamodel.Node, error) {

	// Perform a type assertion to check if data is of type uint64.
	value, ok := data.(uint64)
	if !ok {
		return nil, errors.New("input data must be of type uint64")
	}

	// Create a builder for a new Int node.
	builder := basicnode.Prototype.Int.NewBuilder()

	// Assign the value to the builder.
	err := builder.AssignInt(int64(value))
	if err != nil {
		return nil, fmt.Errorf("failed to assign int value: %v", err)
	}

	// Build the node.
	node := builder.Build()

	return node, nil
}
