package utils

import (
	"bytes"
	"fmt"
)

type Gen struct {
	buf bytes.Buffer
}

func (g *Gen) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

func (g *Gen) Println(args ...interface{}) {
	_, _ = fmt.Fprintln(&g.buf, args...)
}

func (g Gen) String() string { return g.buf.String() }

func (g Gen) Bytes() []byte { return g.buf.Bytes() }
