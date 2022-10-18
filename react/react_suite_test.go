package react_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/react"
	"github.com/wazofski/storz/store"
)

var stc store.Store
var ctx context.Context

var _ = BeforeSuite(func() {
	sch := generated.Schema()
	mem := store.New(sch, memory.Factory())
	mhr := store.New(sch, react.MetaHHandlerFactory(mem))
	rct := store.New(sch, react.ReactFactory(mhr))
	stc = store.New(sch, react.StatusStripperFactory(rct))
})

func TestReact(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "React Suite")
}
