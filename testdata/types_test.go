package testdata

import (
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestName(t *testing.T) {
	var age uint8 = 18
	p := PersonUpdate{
		Name:  proto.String("John Doe"),
		Age:   &age,
		Email: proto.String("e@email.com"),
	}
	cols, args := p.Update()
	t.Log(cols)
	t.Log(args)
}
