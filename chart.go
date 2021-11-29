package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func lineBase(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title + "ns/op"}),
	)

	line.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.LineData
		for _, arg := range oBenchArgs {
			items = append(items, opts.LineData{Value: benchResults[name][arg]})
		}
		line.AddSeries(name, items)
	}

	return line
}

func areaBase(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title + "B/op"}),
	)

	line.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.LineData
		for _, arg := range oBenchArgs {
			items = append(items, opts.LineData{Value: benchResults[name][arg]})
		}
		line.AddSeries(name, items).
			SetSeriesOptions(
				charts.WithLabelOpts(opts.Label{
					Show: true,
				}),
				charts.WithAreaStyleOpts(opts.AreaStyle{
					Opacity: 0.2,
				}),
				charts.WithLineChartOpts(opts.LineChart{
					Smooth: true,
				}),
			)
	}

	return line
}

func barBase(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title + "allocs/op"}),
		//charts.WithInitializationOpts(opts.Initialization{
		//	Width:  "1200px",
		//	Height: "600px",
		//}),
	)

	bar.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.BarData
		for _, arg := range oBenchArgs {
			items = append(items, opts.BarData{Value: benchResults[name][arg]})
		}
		bar.AddSeries(name, items)
	}

	return bar
}

func overlap(benchResults map[int]BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := barBase(benchResults[allocsop], oBenchNames, oBenchArgs)
	//bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
	//	opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
	//	opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
	//))
	//bar.Overlap(lineBase(benchResults[nsop], oBenchNames, oBenchArgs))
	//bar.Overlap(areaBase(benchResults[bop], oBenchNames, oBenchArgs))

	return bar
}

/*
func barOverlap() *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "overlap rect-charts"}),
	)

	bar.SetXAxis(weeks).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems()).
		SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
			opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
			opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
		))
	bar.Overlap(lineBase())
	bar.Overlap(scatterBase())
	return bar
}

var (
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
)

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < itemCntLine; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

var (
	itemCntScatter = 6
	sports         = []string{"Swimming", "Surfing", "Shooting ", "Skating", "Wrestling", "Diving"}
)

func generateScatterItems() []opts.ScatterData {
	items := make([]opts.ScatterData, 0)
	for i := 0; i < itemCntScatter; i++ {
		items = append(items, opts.ScatterData{
			Value:        rand.Intn(100),
			Symbol:       "roundRect",
			SymbolSize:   20,
			SymbolRotate: 10,
		})
	}
	return items
}
func scatterBase() *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic scatter example"}),
	)

	scatter.SetXAxis(sports).
		AddSeries("Category A", generateScatterItems()).
		AddSeries("Category B", generateScatterItems())

	return scatter
}

var (
	itemCnt = 7
	weeks   = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
)

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < itemCnt; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}
*/
