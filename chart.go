package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func line(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeRoma}),
		charts.WithTitleOpts(opts.Title{Title: title + "ns/op"}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
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

func area(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{Title: title + "B/op"}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	line.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.LineData
		for _, arg := range oBenchArgs {
			items = append(items, opts.LineData{Value: benchResults[name][arg]})
		}
		line.AddSeries(name, items).
			SetSeriesOptions(
				charts.WithLabelOpts(
					opts.Label{
						Show: true,
					}),
				charts.WithAreaStyleOpts(
					opts.AreaStyle{
						Opacity: 0.2,
					}),
			)
	}

	return line
}

func bar(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeShine}),
		charts.WithTitleOpts(opts.Title{Title: title + "allocs/op"}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	bar.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.BarData
		for _, arg := range oBenchArgs {
			items = append(items, opts.BarData{Value: benchResults[name][arg]})
		}
		bar.AddSeries(name, items)
	}
	bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
		opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
		opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
	))

	return bar
}

func scatter(benchResults BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeShine}),
		charts.WithTitleOpts(opts.Title{Title: title + "allocs/op"}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	scatter.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.ScatterData
		for _, arg := range oBenchArgs {
			items = append(items, opts.ScatterData{
				Value:        benchResults[name][arg],
				Symbol:       "roundRect",
				SymbolSize:   20,
				SymbolRotate: 10,
			})
		}
		scatter.AddSeries(name, items).
			SetSeriesOptions(charts.WithLabelOpts(
				opts.Label{
					Show:     true,
					Position: "right",
				}),
			)
	}

	return scatter
}

func overlap(benchResults map[int]BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := bar(benchResults[bop], oBenchNames, oBenchArgs)
	//bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
	//	opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
	//	opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
	//))
	bar.Overlap(scatter(benchResults[allocsop], oBenchNames, oBenchArgs))
	bar.Overlap(line(benchResults[nsop], oBenchNames, oBenchArgs))

	return bar
}

func barStack(benchResults map[int]BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeShine}),
		charts.WithTitleOpts(opts.Title{Title: title + "allocs/op"}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	bar.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items1 []opts.BarData
		for _, arg := range oBenchArgs {
			items1 = append(items1, opts.BarData{Value: benchResults[bop][name][arg]})
		}
		bar.AddSeries(name, items1)
		var items2 []opts.BarData
		for _, arg := range oBenchArgs {
			items2 = append(items2, opts.BarData{Value: benchResults[allocsop][name][arg]})
		}
		bar.AddSeries(name+"-stack", items2)
		bar.SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
			Type:  "bar",
			Stack: name + "-stack",
		}))
	}
	//bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
	//	opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
	//	opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
	//))
	//bar.Overlap(lineBase(benchResults[nsop], oBenchNames, oBenchArgs))

	return bar
}
