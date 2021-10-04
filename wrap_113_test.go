// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build go1.13

package ferrors_test

import (
	"errors"
	"testing"

	"github.com/shasderias/ferrors"
)

func TestErrorsIs(t *testing.T) {
	var errSentinel = errors.New("sentinel")

	got := errors.Is(ferrors.Errorf("%w", errSentinel), errSentinel)
	if !got {
		t.Error("got false, want true")
	}

	got = errors.Is(ferrors.Errorf("%w: %s", errSentinel, "foo"), errSentinel)
	if !got {
		t.Error("got false, want true")
	}
}
