package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("return nil from New")
	} else {
		tracer.Trace("hello, trace package")
		if buf.String() != "hello, trace package\n" {
			t.Errorf("'%s' :error printed ", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("data")
}
