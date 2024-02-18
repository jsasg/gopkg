package errors

import (
	"fmt"
	"testing"
)
import stderrors "errors"

func TestUnwrap(t *testing.T) {
	err := Wrap(stderrors.New("1"), "2")
	t.Log(Unwrap(err))

	err1 := fmt.Errorf("%w; %s", stderrors.New("3"), stderrors.New("4"))
	t.Log(Unwrap(err1))

	err2 := stderrors.New("5")
	t.Log(Unwrap(err2))
}
