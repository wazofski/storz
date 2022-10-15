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
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
	"github.com/wazofski/store/memory"
	"github.com/wazofski/store/react"
)

var clt store.Store
var ctx context.Context

var stores []store.Store = []store.Store{
	store.New(
		generated.Schema(),
		memory.Factory()),

	store.New(
		generated.Schema(),
		react.Factory(
			store.New(
				generated.Schema(),
				memory.Factory()))),
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
