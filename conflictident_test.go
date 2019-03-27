package conflictident_test

import (
	"testing"

	"github.com/mackee/conflictident"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, conflictident.Analyzer, "varspec", "varassign", "funcarg")
}
