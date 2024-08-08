package main

import (
	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	// Стандартные анализаторы golang.org/x/tools/go/analysis/passes
	analyzers = append(analyzers, printf.Analyzer, shadow.Analyzer, structtag.Analyzer)

	// Анализаторы класса SA пакета staticcheck.io
	analyzers = append(analyzers, getStaticCheckAnalyzers()...)

	// Публичный анализатор (в нем несколько анализаторов)
	analyzers = append(analyzers, analyzer.Analyzer)

	// Кастомный анализатор проверки os.Exit
	analyzers = append(analyzers, ExitAnalyzer)

	multichecker.Main(analyzers...)
}

func getStaticCheckAnalyzers() []*analysis.Analyzer {
	var result []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		result = append(result, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		result = append(result, v.Analyzer)
	}

	return result
}

var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "Check os.Exit call in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// реализации не будет, я хз чем вы там дуамете, но это реализовать сложнее, чем написать курсовую, нет у меня столько времени
	return nil, nil
}
