package master

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/gonum/stat"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/wcharczuk/go-chart"
)

const (
	fmtStr = "----------- Tx Report\n" +
		"\n[Throughput]\t\n" +
		"Throughput(total)\t%.2f/sec\n" +
		"Throughput(success)\t%.2f/sec\n" +
		"\n[Latencies]\t\n" +
		"Average\t%s\n" +
		"Std.Dev\t%.2f\n" +
		"Max\t%s\n" +
		"99 Percentile\t%s\n" +
		"95 Percentile\t%s\n" +
		"90 Percentile\t%s\n" +
		"Median\t%s\n" +
		"Min\t%s\n" +
		"\n[Bytes]\t\n" +
		"Total Bytes In\t%d bytes\n" +
		"Average Bytes In\t%.2f bytes\n" +
		"Total Bytes Out\t%d bytes\n" +
		"Average Bytes Out\t%.2f bytes\n" +
		"\n[Summary]\t\n" +
		"Total Requests\t%d\n" +
		"Success\t%.2f%%\n" +
		"Duration\t[total, attack, wait]\t%s, %s, %s\n" +
		"Status Codes\t[code:count]\t"

	blockFmtStr = "\n----------- Block Report\n" +
		"\n[TPS]\t\n" +
		"Average\t%.2f/sec\n" +
		"Std.Dev\t%.2f\n" +
		"MPS\t%.2f/sec\n" +
		"\n[Block Interval]\t\n" +
		"Average\t%.2f\n" +
		"Std.Dev\t%.2f\n" +
		"\n[Num Txs Per Block]\t\n" +
		"Average\t%.2f\n" +
		"Std.Dev\t%.2f\n"
)

type Reporter struct {
	metrics     *Metrics
	linkService *service.LinkService
	config      types.Config
	startHeight int64
	endHeight   int64
}

func NewReporter(results []vegeta.Results, config types.Config, startHeight, endHeight int64) *Reporter {
	m := NewMetrics()
	for _, res := range results {
		m.AddResults(res)
	}
	m.Close()

	return &Reporter{
		metrics:     m,
		linkService: service.NewLinkService(&http.Client{}, app.MakeCodec(), config.TargetURL),
		config:      config,
		startHeight: startHeight,
		endHeight:   endHeight,
	}
}

func (r *Reporter) Report(outputDir string) error {
	if err := r.getBlockMetrics(); err != nil {
		return err
	}
	if err := r.printReport(); err != nil {
		return err
	}
	if outputDir != "" {
		r.plotLatency(outputDir, r.metrics.timeStamps, r.metrics.latencies)
		r.plotTPS(outputDir, r.metrics.blockMetrics.height, r.metrics.blockMetrics.tps, r.metrics.blockMetrics.blockInterval)
	}
	return nil
}

func (r *Reporter) getBlockMetrics() error {
	block, err := r.linkService.GetBlock(r.startHeight - 1)
	if err != nil {
		return err
	}
	LastBlockTime := block.Block.Time

	for i := r.startHeight; i <= r.endHeight; i++ {
		block, err := r.linkService.GetBlock(i)
		if err != nil {
			return err
		}
		numTxs := float64(len(block.Block.Txs))
		blockTime := block.Block.Time
		blockInterval := blockTime.Sub(LastBlockTime)
		r.metrics.blockMetrics.height = append(r.metrics.blockMetrics.height, float64(i))
		r.metrics.blockMetrics.numTxs = append(r.metrics.blockMetrics.numTxs, numTxs)
		r.metrics.blockMetrics.blockInterval = append(r.metrics.blockMetrics.blockInterval, blockInterval.Seconds())
		r.metrics.blockMetrics.tps = append(r.metrics.blockMetrics.tps, numTxs/blockInterval.Seconds())
		LastBlockTime = blockTime
	}
	return nil
}

func (r *Reporter) printReport() error {
	if err := r.printTxReport(); err != nil {
		return err
	}
	if err := r.printBlockReport(); err != nil {
		return err
	}
	return nil
}

// Modified the code in the following link.
// https://github.com/tsenart/vegeta/blob/master/lib/reporters.go
func (r *Reporter) printTxReport() error {
	m := r.metrics
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', tabwriter.StripEscape)
	if _, err := fmt.Fprintf(tw, fmtStr,
		m.Rate, m.Throughput,
		round(m.Latencies.Mean),
		stat.StdDev(m.latencies, nil),
		round(m.Latencies.Max),
		round(m.Latencies.P99),
		round(m.Latencies.P95),
		round(m.Latencies.P90),
		round(m.Latencies.P50),
		round(m.Latencies.Min),
		m.BytesIn.Total, m.BytesIn.Mean,
		m.BytesOut.Total, m.BytesOut.Mean,
		m.Requests,
		m.Success*100,
		round(m.Duration+m.Wait),
		round(m.Duration),
		round(m.Wait),
	); err != nil {
		return err
	}

	codes := make([]string, 0, len(m.StatusCodes))
	for code := range m.StatusCodes {
		codes = append(codes, code)
	}

	sort.Strings(codes)

	for _, code := range codes {
		count := m.StatusCodes[code]
		if _, err := fmt.Fprintf(tw, "%s:%d  ", code, count); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(tw, "\nError Set:"); err != nil {
		return err
	}

	for _, e := range m.Errors {
		if _, err := fmt.Fprintln(tw, e); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func (r *Reporter) printBlockReport() error {
	m := r.metrics.blockMetrics
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', tabwriter.StripEscape)
	meanTPS := stat.Mean(m.tps, nil)
	if _, err := fmt.Fprintf(tw, blockFmtStr,
		meanTPS,
		stat.StdDev(m.tps, nil),
		float64(r.config.MsgsPerTxLoadTest)*meanTPS,
		stat.Mean(m.blockInterval, nil),
		stat.StdDev(m.blockInterval, nil),
		stat.Mean(m.numTxs, nil),
		stat.StdDev(m.numTxs, nil),
	); err != nil {
		return err
	}

	return tw.Flush()
}

func (r *Reporter) plotLatency(outputDir string, xValues []time.Time, yValues []float64) {
	name := "Latency"
	mainSeries := chart.TimeSeries{
		Name: name,
		Style: chart.Style{
			StrokeColor: chart.ColorBlue,
			Show:        true,
		},
		XValues: xValues,
		YValues: yValues,
	}
	graph := chart.Chart{
		Width:  1280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
		},
		XAxis: chart.XAxis{
			Style:          chart.StyleShow(),
			ValueFormatter: chart.TimeValueFormatterWithFormat("15:04:05"),
		},
		YAxis: chart.YAxis{
			Name:      "Elapsed Millis",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d ms", int(v.(float64)*1000.0))
			},
		},
		Series: []chart.Series{
			mainSeries,
		},
	}
	r.drawGraph(graph, outputDir, name)
}

func (r *Reporter) plotTPS(outputDir string, xValues, yValuesTPS, yValuesInterval []float64) {
	name := "TPS"
	mainSeries := chart.ContinuousSeries{
		Name: name,
		Style: chart.Style{
			StrokeColor: chart.ColorBlue,
			Show:        true,
		},
		XValues: xValues,
		YValues: yValuesTPS,
	}
	intervalSeries := chart.ContinuousSeries{
		Name: "Block Interval",
		Style: chart.Style{
			StrokeColor:     chart.ColorAlternateBlue,
			StrokeDashArray: []float64{5.0, 5.0},
			Show:            true,
		},
		XValues: xValues,
		YValues: yValuesInterval,
		YAxis:   chart.YAxisSecondary,
	}
	graph := chart.Chart{
		Width:  1280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top:  50,
				Left: 50,
			},
		},
		XAxis: chart.XAxis{
			Name:      "height",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		YAxis: chart.YAxis{
			Name:      "tps",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxisSecondary: chart.YAxis{
			Name:      "block interval",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			mainSeries,
			intervalSeries,
		},
	}

	r.drawGraph(graph, outputDir, name)
}

func (r *Reporter) drawGraph(graph chart.Chart, outputDir, name string) {
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}

	f, err := os.Create(fmt.Sprintf("%s/%s.png", outputDir, name))
	if err != nil {
		log.Fatalf("Failed to create file:%s", err.Error())
	}
	defer f.Close()
	if err := graph.Render(chart.PNG, f); err != nil {
		log.Fatalf("Failed to render graph:%s", err.Error())
	}
}

var durations = [...]time.Duration{
	time.Hour,
	time.Minute,
	time.Second,
	time.Millisecond,
	time.Microsecond,
	time.Nanosecond,
}

func round(d time.Duration) time.Duration {
	for i, unit := range durations {
		if d >= unit && i < len(durations)-1 {
			return d.Round(durations[i+1])
		}
	}
	return d
}
