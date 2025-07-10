package gasper

import (
	"github.com/mysteryforge/gasper/k6/integrity"
	"github.com/mysteryforge/gasper/k6/loadtest"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/gasper/loadtest", loadtest.New())
	modules.Register("k6/x/gasper/integrity", integrity.New())
}
