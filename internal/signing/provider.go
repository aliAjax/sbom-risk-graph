package signing

import (
	"context"
	"fmt"
)

type Signer interface {
	Sign(context.Context, []byte) (string, error)
}

func chooseSigner(primary, fallback Signer) (Signer, error) {
	if primary != nil {
		return primary, nil
	}
	if fallback != nil {
		return fallback, nil
	}
	return nil, fmt.Errorf("no signing provider available")
}
