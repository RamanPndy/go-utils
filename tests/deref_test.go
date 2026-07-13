package goutils_test

import (
	"testing"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestDeref(t *testing.T) {
	value := 42
	if got := goutils.Deref(&value); got != value {
		t.Errorf("Deref(&value) = %v, want %v", got, value)
	}

	if got := goutils.Deref[int](nil); got != 0 {
		t.Errorf("Deref[int](nil) = %v, want 0", got)
	}

	if got := goutils.Deref[string](nil); got != "" {
		t.Errorf("Deref[string](nil) = %q, want empty string", got)
	}
}

func TestDerefOr(t *testing.T) {
	if got := goutils.DerefOr[int](nil, 7); got != 7 {
		t.Errorf("DerefOr[int](nil, 7) = %v, want 7", got)
	}

	value := true
	if got := goutils.DerefOr(&value, false); got != true {
		t.Errorf("DerefOr(&value, false) = %v, want true", got)
	}
}

func TestDerefSliceWrappers(t *testing.T) {
	stringSlice := goutils.DerefStringSlice(nil)
	if stringSlice == nil || len(stringSlice) != 0 {
		t.Errorf("DerefStringSlice(nil) = %#v, want empty non-nil slice", stringSlice)
	}

	intSlice := goutils.DerefIntSlice(nil)
	if intSlice == nil || len(intSlice) != 0 {
		t.Errorf("DerefIntSlice(nil) = %#v, want empty non-nil slice", intSlice)
	}
}
