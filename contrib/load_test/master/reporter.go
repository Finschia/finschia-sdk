package master

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/gonum/stat"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/types"
	"github.com/olekukonko/tablewriter"
	tmtypes "github.com/tendermint/tendermint/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/wcharczuk/go-chart"
	"golang.org/x/sync/errgroup"
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
		"Std.Dev\t%.2f\n" +
		"\n[Failed Txs]\t\n" +
		"Total Failed Tx\t%.d\n"
)

var mutex = &sync.Mutex{}

type Reporter struct {
	metrics     *Metrics
	linkService *service.LinkService
	slaves      []types.Slave
	config      types.Config
	startHeight int64
	endHeight   int64
	thresholds  types.Thresholds
}

func NewReporter(results []vegeta.Results, slaves []types.Slave, config types.Config, startHeight,
	endHeight int64, thresholds types.Thresholds) *Reporter {
	m := NewMetrics()
	for _, res := range results {
		m.AddResults(res)
	}
	m.Close()

	return &Reporter{
		metrics:     m,
		linkService: service.NewLinkService(&http.Client{Timeout: 30 * time.Second}, app.MakeCodec(), config.TargetURL),
		slaves:      slaves,
		config:      config,
		startHeight: startHeight,
		endHeight:   endHeight,
		thresholds:  thresholds,
	}
}

func (r *Reporter) Report(outputDir string) error {
	log.Println("Start to analyze results")
	if r.isTxScenario() {
		if err := r.getBlockMetrics(); err != nil {
			return err
		}
		if err := r.checkMissingTx(); err != nil {
			return err
		}
	}
	if err := r.printReport(); err != nil {
		return err
	}
	if outputDir != "" {
		if err := r.plotLatency(outputDir, r.metrics.timeStamps, r.metrics.latencies); err != nil {
			return err
		}
		if r.isTxScenario() {
			if err := r.plotTPS(outputDir, r.metrics.blockMetrics.height, r.metrics.blockMetrics.tps,
				r.metrics.blockMetrics.blockInterval); err != nil {
				return err
			}
		}
	}

	return r.checkThresholds()
}

func (r *Reporter) isTxScenario() bool {
	for _, slave := range r.slaves {
		if strings.HasPrefix(slave.Scenario, "tx") {
			return true
		}
	}
	return false
}

func (r *Reporter) getBlockMetrics() error {
	log.Println("Get block metrics")
	block, err := r.linkService.GetBlock(r.startHeight - 1)
	for err != nil {
		log.Println("Retrying to get block")
		block, err = r.linkService.GetBlock(r.startHeight - 1)
	}
	LastBlockTime := block.Block.Time

	blocks := make([]*tmtypes.Block, r.endHeight-r.startHeight+2)
	blockGas := make([]int64, r.endHeight-r.startHeight+2)
	for i := r.startHeight; i <= r.endHeight+1; i++ {
		blockWithTxResults, err := r.linkService.GetBlocksWithTxResults(i, 1)
		if err != nil {
			return err
		}

		blocks[i-r.startHeight] = blockWithTxResults[0].ResultBlock.Block

		for j, txResponse := range blockWithTxResults[0].TxResponses {
			r.metrics.processedTxs[txResponse.TxHash] = true

			if j == 0 {
				r.metrics.maxUsedGas = txResponse.GasUsed
				r.metrics.minUsedGas = txResponse.GasUsed
			} else {
				if r.metrics.maxUsedGas < txResponse.GasUsed {
					r.metrics.maxUsedGas = txResponse.GasUsed
				}
				if r.metrics.minUsedGas > txResponse.GasUsed {
					r.metrics.minUsedGas = txResponse.GasUsed
				}
			}

			blockGas[i-r.startHeight] += txResponse.GasUsed
			if txResponse.Code != 0 {
				r.metrics.blockMetrics.numFailedTxs++
				r.metrics.blockMetrics.failedTxLogs = append(r.metrics.blockMetrics.failedTxLogs, txResponse.RawLog)
			}
		}
	}
	for i, block := range blocks {
		if r.metrics.txSize == 0 && len(block.Txs) > 0 {
			r.metrics.txSize = len(block.Txs[0])
		}
		numTxs := float64(len(block.Txs))
		blockTime := block.Time

		if i < len(blocks)-1 {
			r.metrics.blockMetrics.height = append(r.metrics.blockMetrics.height, float64(block.Height))
			r.metrics.blockMetrics.numTxs = append(r.metrics.blockMetrics.numTxs, numTxs)
			r.metrics.blockMetrics.blockGas = append(r.metrics.blockMetrics.blockGas, float64(blockGas[i]))
		}
		if i > 0 {
			blockInterval := blockTime.Sub(LastBlockTime)
			r.metrics.blockMetrics.blockInterval = append(r.metrics.blockMetrics.blockInterval, blockInterval.Seconds())
			r.metrics.blockMetrics.tps = append(r.metrics.blockMetrics.tps, numTxs/blockInterval.Seconds())
			r.metrics.blockMetrics.round = append(r.metrics.blockMetrics.round, float64(block.LastCommit.Round))
		}
		LastBlockTime = blockTime
	}

	return nil
}

func (r *Reporter) checkMissingTx() error {
	log.Println("Check Missing Tx")

	sem := make(chan int, r.config.MaxWorkers)
	var eg errgroup.Group
	for _, tx := range r.metrics.transferredTxs {
		tx := tx
		sem <- 1
		eg.Go(func() error {
			defer func(sem *chan int) {
				<-*sem
				if err := recover(); err != nil {
					log.Println("Failed to get tx:", err)
					log.Println(string(debug.Stack()))
				}
			}(&sem)

			if _, ok := r.metrics.processedTxs[tx]; ok {
				return nil
			}

			_, err := r.linkService.GetTx(tx)
			if err != nil {
				reqFailedError, ok := err.(types.RequestFailed)
				if ok && reqFailedError.Status == "404 Not Found" {
					mutex.Lock()
					r.metrics.missingTxs = append(r.metrics.missingTxs, tx)
					mutex.Unlock()
				} else {
					return err
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (r *Reporter) printReport() error {
	if err := r.printTxReport(os.Stdout); err != nil {
		return err
	}
	if r.isTxScenario() {
		if err := r.printBlockReport(os.Stdout); err != nil {
			return err
		}
		if err := r.printMsgReport(os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

// Modified the code in the following link.
// https://github.com/tsenart/vegeta/blob/master/lib/reporters.go
func (r *Reporter) printTxReport(output io.Writer) error {
	m := r.metrics
	tw := tabwriter.NewWriter(output, 0, 8, 2, ' ', tabwriter.StripEscape)
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

	if err := r.printErrorLogs(tw, "Error Set:", m.Errors); err != nil {
		return err
	}

	if err := r.printErrorLogs(tw, "Failed Tx Logs:", m.failedTxLogs); err != nil {
		return err
	}

	if err := r.printErrorLogs(tw, "Missing Txs:", m.missingTxs); err != nil {
		return err
	}

	return tw.Flush()
}

func (r *Reporter) printErrorLogs(tw *tabwriter.Writer, name string, logs []string) error {
	if _, err := fmt.Fprintf(tw, "\nNum %s %d\n", name, len(logs)); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(tw, "%s\n", name); err != nil {
		return err
	}

	counter := make(map[string]int)
	for _, e := range logs {
		counter[e]++
	}

	errors := make([]string, 0, len(counter))
	for e := range counter {
		errors = append(errors, e)
	}
	sort.Strings(errors)

	for _, e := range errors {
		if counter[e] > 1 {
			if _, err := fmt.Fprintf(tw, "%s (%d logs)\n", e, counter[e]); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(tw, e); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Reporter) printBlockReport(output io.Writer) error {
	m := r.metrics.blockMetrics
	tw := tabwriter.NewWriter(output, 0, 8, 2, ' ', tabwriter.StripEscape)
	meanTPS := stat.Mean(m.tps, nil)
	if _, err := fmt.Fprintf(tw, blockFmtStr,
		meanTPS,
		stat.StdDev(m.tps, nil),
		r.getNumMsgs()*meanTPS,
		stat.Mean(m.blockInterval, nil),
		stat.StdDev(m.blockInterval, nil),
		stat.Mean(m.numTxs, nil),
		stat.StdDev(m.numTxs, nil),
		m.numFailedTxs,
	); err != nil {
		return err
	}

	if err := r.printErrorLogs(tw, "Failed Tx Logs:", m.failedTxLogs); err != nil {
		return err
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(output, "\n[Block Details]"); err != nil {
		return err
	}

	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"Block Height", "Num Txs", "Block Used Gas", "Block Interval", "Round"})

	for i := 0; i < len(r.metrics.blockMetrics.height); i++ {
		v := []string{
			strconv.FormatInt(int64(r.metrics.blockMetrics.height[i]), 10),
			strconv.FormatInt(int64(r.metrics.blockMetrics.numTxs[i]), 10),
			strconv.FormatInt(int64(r.metrics.blockMetrics.blockGas[i]), 10),
			strconv.FormatInt(int64(r.metrics.blockMetrics.blockInterval[i]), 10),
			strconv.FormatInt(int64(r.metrics.blockMetrics.round[i]), 10),
		}
		table.Append(v)
	}
	table.Render()

	return nil
}

func (r *Reporter) getNumMsgs() float64 {
	totalMsgs := 0
	for _, slave := range r.slaves {
		switch slave.Scenario {
		case types.TxToken:
			totalMsgs += r.roundUpNumMsgs(6)

		case types.TxCollection:
			totalMsgs += r.roundUpNumMsgs(8)

		case types.TxAll:
			totalMsgs += r.roundUpNumMsgs(29)

		default:
			totalMsgs += r.config.MsgsPerTxLoadTest
		}
	}
	totalMsgs /= len(r.slaves)
	return float64(totalMsgs)
}

func (r *Reporter) roundUpNumMsgs(numMsgIncrement int) int {
	repeatCount := (r.config.MsgsPerTxLoadTest + numMsgIncrement - 1) / numMsgIncrement
	return numMsgIncrement * repeatCount
}

func (r *Reporter) printMsgReport(output io.Writer) error {
	scenario := r.slaves[0].Scenario
	for i := 1; i < len(r.slaves); i++ {
		if r.slaves[i].Scenario != scenario {
			return nil
		}
	}

	if _, err := fmt.Fprint(output, "\n----------- Msg Report\n\n"); err != nil {
		return err
	}
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"Msg", "Num Msgs Per Tx", "Tx Size (Bytes)", "Max Used Gas", "Min Used Gas",
		"Average Used Gas"})
	v := []string{
		r.slaves[0].Scenario,
		strconv.Itoa(r.config.MsgsPerTxLoadTest),
		strconv.Itoa(r.metrics.txSize),
		strconv.FormatInt(r.metrics.maxUsedGas, 10),
		strconv.FormatInt(r.metrics.minUsedGas, 10),
		fmt.Sprintf("%f", r.getAverageUsedGas()),
	}
	table.Append(v)
	table.Render()
	return nil
}

func (r *Reporter) getAverageUsedGas() float64 {
	totalGas := 0.0
	for _, gas := range r.metrics.blockMetrics.blockGas {
		totalGas += gas
	}

	totalTxs := 0.0
	for _, numTx := range r.metrics.blockMetrics.numTxs {
		totalTxs += numTx
	}
	return totalGas / totalTxs
}

func (r *Reporter) plotLatency(outputDir string, xValues []time.Time, yValues []float64) error {
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
	if err := r.drawGraph(graph, outputDir, name); err != nil {
		return err
	}
	return nil
}

func (r *Reporter) plotTPS(outputDir string, xValues, yValuesTPS, yValuesInterval []float64) error {
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
	var yAxis chart.YAxis
	if yValuesTPS[len(yValuesTPS)/2] == 0 {
		yAxis = chart.YAxis{
			Name:      "tps",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Range:     &chart.ContinuousRange{Min: 0, Max: 1},
		}
	} else {
		yAxis = chart.YAxis{
			Name:      "tps",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		}
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
		YAxis: yAxis,
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

	if err := r.drawGraph(graph, outputDir, name); err != nil {
		return err
	}
	return nil
}

func (r *Reporter) drawGraph(graph chart.Chart, outputDir, name string) error {
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}

	f, err := os.Create(fmt.Sprintf("%s/%s.png", outputDir, name))
	if err != nil {
		return types.FailedToCreateFile{Err: err}
	}
	defer f.Close()
	if err := graph.Render(chart.PNG, f); err != nil {
		return types.FailedToRenderGraph{Err: err}
	}
	return nil
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

func (r *Reporter) checkThresholds() error {
	throughputThreshold := float64(r.thresholds.Throughput)
	if throughputThreshold > 0 && r.metrics.Throughput < throughputThreshold {
		return types.LowThroughputError{
			Throughput: r.metrics.Throughput,
			Threshold:  throughputThreshold,
		}
	}

	if r.thresholds.Latency > 0 && r.metrics.Latencies.Mean > r.thresholds.Latency {
		return types.HighLatencyError{
			Latency:   r.metrics.Latencies.Mean,
			Threshold: r.thresholds.Latency,
		}
	}

	tps := stat.Mean(r.metrics.blockMetrics.tps, nil)
	tpsThreshold := float64(r.thresholds.TPS)
	if r.thresholds.TPS > 0 && tps < tpsThreshold {
		return types.LowTPSError{
			TPS:       tps,
			Threshold: tpsThreshold,
		}
	}
	return nil
}
