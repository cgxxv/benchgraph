package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/go-echarts/go-echarts/v2/components"
	"golang.org/x/tools/benchmark/parse"
)

var metrics = []string{
	nsop,
	bop,
	allocsop,
}

const (
	nsop     = "n/op"
	bop      = "B/op"
	allocsop = "allocs/op"

	bar = iota
	line
	area

	title = "Graph: Benchmark results in "
)

var (
	showBarMaxLine bool
	showBarAvgLine bool
)

func main() {

	var (
		oBenchNames, oBenchArgs, extra stringList
		shape                          int
	)

	// graph elements will be ordered as in benchmark output by default - unless the order was specified here
	flag.Var(&oBenchNames, "obn", "comma-separated list of benchmark names")
	flag.Var(&oBenchArgs, "oba", "comma-separated list of benchmark arguments")
	flag.IntVar(&shape, "shape", 0, "result of charts: 0 => bar, 1 => line, 2 => area, default: 0")
	flag.BoolVar(&showBarMaxLine, "max", false, "show max line for bar chart, default: false")
	flag.BoolVar(&showBarAvgLine, "avg", false, "show avg line for bar chart, default: false")
	flag.Parse()

	var skipBenchNamesParsing, skipBenchArgsParsing bool

	if oBenchNames.Len() > 0 {
		skipBenchNamesParsing = true
	}
	if oBenchArgs.Len() > 0 {
		skipBenchArgsParsing = true
	}

	benchResults := make(map[string]BenchNameSet)

	// parse Golang benchmark results, line by line
	scan := bufio.NewScanner(os.Stdin)
	green := color.New(color.FgGreen).SprintfFunc()
	red := color.New(color.FgRed).SprintFunc()
	for scan.Scan() {
		line := scan.Text()

		mark := green("√")

		b, err := parse.ParseLine(line)
		if err != nil {
			mark = red("?")
		}

		// read bench name and arguments
		if b != nil {
			name, arg, _, err := parseNameArgThread(b.Name)
			if err != nil {
				mark = red("!")
				fmt.Printf("%s %s\n", mark, line)
				continue
			}

			if !skipBenchNamesParsing && !oBenchNames.stringInList(name) {
				oBenchNames.Add(name)
			}

			if !skipBenchArgsParsing && !oBenchArgs.stringInList(arg) {
				oBenchArgs.Add(arg)
			}

			for _, v := range metrics {
				if _, ok := benchResults[v]; !ok {
					benchResults[v] = make(BenchNameSet)
				}

				if _, ok := benchResults[v][name]; !ok {
					benchResults[v][name] = make(BenchArgSet)
				}
			}

			for v := range b.Extra {
				if _, ok := benchResults[v]; !ok {
					benchResults[v] = make(BenchNameSet)
				}

				if _, ok := benchResults[v][name]; !ok {
					benchResults[v][name] = make(BenchArgSet)
				}
			}

			benchResults[nsop][name][arg] = b.NsPerOp
			benchResults[bop][name][arg] = b.AllocedBytesPerOp
			benchResults[allocsop][name][arg] = b.AllocsPerOp

			for k, v := range b.Extra {
				benchResults[k][name][arg] = v
				if !extra.stringInList(k) {
					extra.Add(k)
				}
			}
		}

		fmt.Printf("%s %s\n", mark, line)
	}

	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v", err)
		os.Exit(1)
	}

	if len(benchResults) == 0 {
		fmt.Fprintf(os.Stderr, "no data to show.\n\n")
		os.Exit(1)
	}

	page := components.NewPage()

	switch shape {
	case line:
		for _, v := range metrics {
			page.AddCharts(
				lineChart(benchResults, v, oBenchNames, oBenchArgs),
			)
		}
	case area:
		for _, v := range metrics {
			page.AddCharts(
				areaChart(benchResults, v, oBenchNames, oBenchArgs),
			)
		}
	default:
		for _, v := range metrics {
			page.AddCharts(
				barChart(benchResults, v, oBenchNames, oBenchArgs),
			)
		}
	}

	fmt.Printf("%#v\n", extra)

	for _, v := range extra {
		switch shape {
		case line:
			page.AddCharts(
				lineChart(benchResults, v, oBenchNames, oBenchArgs),
			)
		case area:
			page.AddCharts(
				areaChart(benchResults, v, oBenchNames, oBenchArgs),
			)
		default:
			page.AddCharts(
				barChart(benchResults, v, oBenchNames, oBenchArgs),
			)
		}
	}

	f, err := os.Create("asset/page.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

	fs := http.FileServer(http.Dir("asset"))
	log.Println("running server at http://localhost:8090")
	log.Fatal(http.ListenAndServe("localhost:8090", logRequest(fs)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
