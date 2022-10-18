package sql_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/sql"
	"github.com/wazofski/storz/store"
)

func TestSql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sql Suite")
}

var stc store.Store
var ctx context.Context

var _ = BeforeSuite(func() {
	sch := generated.Schema()
	const path = "test.sqlite"
	stc = store.New(sch, sql.Factory(sql.SqliteConnection(path)))
})

var _ = AfterSuite(func() {
	const path = "test.sqlite"
	os.Remove(path)
})
