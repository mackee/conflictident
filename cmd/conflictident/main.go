package main

import (
	"github.com/mackee/conflictident"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(conflictident.Analyzer) }
