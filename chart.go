package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func lineChart(benchResults map[string]BenchNameSet, metric string, oBenchNames, oBenchArgs stringList) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeRoma}),
		charts.WithTitleOpts(opts.Title{Title: title + metric}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "900px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	line.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.LineData
		for _, arg := range oBenchArgs {
			items = append(items, opts.LineData{Value: benchResults[metric][name][arg]})
		}
		line.AddSeries(name, items)
	}

	return line
}

func areaChart(benchResults map[string]BenchNameSet, metric string, oBenchNames, oBenchArgs stringList) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{Title: title + metric}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "900px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	line.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.LineData
		for _, arg := range oBenchArgs {
			items = append(items, opts.LineData{Value: benchResults[metric][name][arg]})
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

func barChart(benchResults map[string]BenchNameSet, metric string, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeShine}),
		charts.WithTitleOpts(opts.Title{Title: title + metric}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "900px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	bar.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.BarData
		for _, arg := range oBenchArgs {
			items = append(items, opts.BarData{Value: benchResults[metric][name][arg]})
		}
		bar.AddSeries(name, items)
	}

	if showBarMaxLine {
		bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
			opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
		))
	}

	if showBarAvgLine {
		bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
			opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
		))
	}
	//bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
	//	opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
	//	opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
	//))

	return bar
}

func scatterChart(benchResults map[string]BenchNameSet, metric string, oBenchNames, oBenchArgs stringList) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeShine}),
		charts.WithTitleOpts(opts.Title{Title: title + metric}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "900px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true, Right: "top", Feature: &opts.ToolBoxFeature{SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: true}}}),
	)

	scatter.SetXAxis(oBenchArgs)

	for _, name := range oBenchNames {
		var items []opts.ScatterData
		for _, arg := range oBenchArgs {
			items = append(items, opts.ScatterData{
				Value:        benchResults[metric][name][arg],
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

func overlapChart(benchResults map[string]BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := barChart(benchResults, bop, oBenchNames, oBenchArgs)
	//bar.SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
	//	opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max"},
	//	opts.MarkLineNameTypeItem{Name: "Avg", Type: "average"},
	//))
	bar.Overlap(scatterChart(benchResults, allocsop, oBenchNames, oBenchArgs))
	bar.Overlap(lineChart(benchResults, nsop, oBenchNames, oBenchArgs))

	return bar
}

func barStackChart(benchResults map[string]BenchNameSet, oBenchNames, oBenchArgs stringList) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeShine}),
		charts.WithTitleOpts(opts.Title{Title: title + allocsop}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "900px",
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
