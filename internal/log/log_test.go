package log

import (
	"log"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	got := New()
	want := &log.Logger{}

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("Expected %+v got %+v", want, got)
	}
}
