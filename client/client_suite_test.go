package client_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store"
	"github.com/wazofski/store/client"
	"github.com/wazofski/store/generated"
)

var stc store.Store
var ctx context.Context

var _ = BeforeSuite(func() {
	stc = store.New(
		generated.Schema(),
		client.Factory(
			"http://localhost:8000/",
			client.Header("test", "header")))
})

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}
