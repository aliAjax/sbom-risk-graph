package signing

import (
	"context"
	"fmt"
	"reflect"
)

type Signer interface {
	Sign(context.Context, []byte) (string, error)
}

func signerAvailable(signer Signer) bool {
	if signer == nil {
		return false
	}
	value := reflect.ValueOf(signer)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return !value.IsNil()
	default:
		return true
	}
}

func chooseSigner(primary, fallback Signer) (Signer, error) {
	if signerAvailable(primary) {
		return primary, nil
	}
	if signerAvailable(fallback) {
		return fallback, nil
	}
	return nil, fmt.Errorf("no signing provider available")
}
