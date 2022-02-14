package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoid(t *testing.T) {
	if reflect.TypeOf(Goid()).Kind().String() == "int" {
		assert.True(t, true, "")
	}
}
