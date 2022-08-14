package testutil_test

import (
	"testing"

	"github.com/ericlagergren/testutil"
)

func TestInlining(t *testing.T) {
	testutil.TestInlining(t, "github.com/ericlagergren/testutil",
		"hasGoBuild",
	)
}
