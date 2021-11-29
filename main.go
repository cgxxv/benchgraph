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

const (
	nsop = iota
	bop
	allocsop

	title = "Graph: Benchmark results in "
)

func main() {

	var oBenchNames, oBenchArgs stringList

	// graph elements will be ordered as in benchmark output by default - unless the order was specified here
	flag.Var(&oBenchNames, "obn", "comma-separated list of benchmark names")
	flag.Var(&oBenchArgs, "oba", "comma-separated list of benchmark arguments")
	flag.Parse()

	var skipBenchNamesParsing, skipBenchArgsParsing bool

	if oBenchNames.Len() > 0 {
		skipBenchNamesParsing = true
	}
	if oBenchArgs.Len() > 0 {
		skipBenchArgsParsing = true
	}

	benchResults := make(map[int]BenchNameSet)

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

			for i := nsop; i <= allocsop; i++ {
				if _, ok := benchResults[i]; !ok {
					benchResults[i] = make(BenchNameSet)
				}

				if _, ok := benchResults[i][name]; !ok {
					benchResults[i][name] = make(BenchArgSet)
				}
			}

			benchResults[nsop][name][arg] = b.NsPerOp
			benchResults[bop][name][arg] = b.AllocedBytesPerOp
			benchResults[allocsop][name][arg] = b.AllocsPerOp
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
	page.AddCharts(
		//overlap(benchResults, oBenchNames, oBenchArgs),
		lineBase(benchResults[nsop], oBenchNames, oBenchArgs),
		areaBase(benchResults[bop], oBenchNames, oBenchArgs),
		barBase(benchResults[allocsop], oBenchNames, oBenchArgs),
	)
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
