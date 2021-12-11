package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unsafe"

	"github.com/fatih/color"
	"github.com/go-echarts/go-echarts/v2/components"
	"golang.org/x/tools/benchmark/parse"
)

var metrics = []string{
	benchtimes,
	nsop,
	bop,
	allocsop,
}

const (
	benchtimes = "benchtimes"
	nsop       = "ns/op"
	bop        = "B/op"
	allocsop   = "allocs/op"

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

			benchResults[benchtimes][name][arg] = b.N
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

	f, err := os.OpenFile("page.html", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	f.Close()

	log.Println("running server at http://localhost:8090")
	http.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}

func hello(w http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadFile("page.html")
	fmt.Fprint(w, *(*string)(unsafe.Pointer(&b)))
}
