package conflictident_test

import (
	"testing"

	"github.com/mackee/conflictident"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "varspec", "funcarg")
}

func TestVarassign(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "varassign")
}

func TestPackageIdent(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "pkgident")
}

func TestMultiFile(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "multifile")
}

func TestReceiver(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "receiver")
}

func TestUse(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "use")
}

func TestBuiltin(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "builtintest")
}

func TestConst(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "consttest")
}
