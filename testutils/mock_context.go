package testutils

import (
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
)

// NewMockContext returns a context mock
func NewMockContext() *lepton.Context {
	return lepton.NewContext(types.NewConfig())
}
