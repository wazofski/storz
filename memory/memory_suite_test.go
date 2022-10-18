package memory_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/store"
)

func TestMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Memory Suite")
}

var client store.Store
var ctx context.Context

var _ = BeforeSuite(func() {
	client = store.New(
		generated.Schema(),
		memory.Factory())
})
