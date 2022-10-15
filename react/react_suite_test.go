package react_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
	"github.com/wazofski/store/memory"
	"github.com/wazofski/store/react"
)

var stc store.Store
var ctx context.Context

var _ = BeforeSuite(func() {
	sch := generated.Schema()
	mem := store.New(sch, memory.Factory())
	stc = store.New(sch, react.Factory(mem))
})

func TestReact(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "React Suite")
}
