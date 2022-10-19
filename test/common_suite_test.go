package common_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/client"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/react"
	"github.com/wazofski/storz/sql"
	"github.com/wazofski/storz/store"
)

var clt store.Store
var ctx context.Context

// os.Remove("test.sqlite")

var stores []store.Store = []store.Store{
	store.New(
		generated.Schema(),
		memory.Factory()),

	store.New(
		generated.Schema(),
		react.ReactFactory(
			store.New(
				generated.Schema(),
				memory.Factory()))),

	store.New(
		generated.Schema(),
		client.Factory("http://localhost:8000/")),

	store.New(
		generated.Schema(),
		logger.StoreFactory("SQLite",
			store.New(
				generated.Schema(),
				sql.Factory(sql.SqliteConnection("test.sqlite"))))),

	store.New(
		generated.Schema(),
		logger.StoreFactory("mySQL",
			store.New(
				generated.Schema(),
				sql.Factory(sql.MySqlConnection(
					"root:qwerasdf@tcp(127.0.0.1:3306)/test"))))),
}

func TestNegative(t *testing.T) {
	RegisterFailHandler(Fail)

	argKey := "store="
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, argKey) {
			sarg, err := strconv.Atoi(
				strings.Split(arg, "=")[1])
			if err != nil {
				log.Println(err)
				return
			}

			clt = stores[sarg]
			RunSpecs(t, fmt.Sprintf("Common Suite %d", sarg))

			break
		}
	}
}
