package work

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

const (
	barChar = "∎"
)

type Report struct {
	avgTotal float64
	fastest  float64
	slowest  float64
	average  float64
	rps      float64

	avgConn   float64
	avgDNS    float64
	avgReq    float64
	avgRes    float64
	avgDelay  float64
	connLats  []float64
	dnsLats   []float64
	reqLats   []float64
	resLats   []float64
	delayLats []float64

	results chan *Result
	total   time.Duration

	errorDist      map[string]int
	statusCodeDist map[int]int
	lats           []float64
	sizeTotal      int64

	output string

	w io.Writer
}

func NewReport(w io.Writer, size int, results chan *Result, output string, total time.Duration) *Report {
	return &Report{
		output:         output,
		results:        results,
		total:          total,
		statusCodeDist: make(map[int]int),
		errorDist:      make(map[string]int),
		w:              w,
	}
}

func (r *Report) Finalize() {
	for res := range r.results {
		if res.Err != nil {
			r.errorDist[res.Err.Error()]++
		} else {
			r.lats = append(r.lats, res.Duration.Seconds())
			r.avgTotal += res.Duration.Seconds()
			r.avgConn += res.ConnDuration.Seconds()
			r.avgDelay += res.DelayDuration.Seconds()
			r.avgDNS += res.DnsDuration.Seconds()
			r.avgReq += res.ReqDuration.Seconds()
			r.avgRes += res.ResDuration.Seconds()
			r.connLats = append(r.connLats, res.ConnDuration.Seconds())
			r.dnsLats = append(r.dnsLats, res.DnsDuration.Seconds())
			r.reqLats = append(r.reqLats, res.ReqDuration.Seconds())
			r.delayLats = append(r.delayLats, res.DelayDuration.Seconds())
			r.resLats = append(r.resLats, res.ResDuration.Seconds())
			r.statusCodeDist[res.StatusCode]++
			if res.ContentLength > 0 {
				r.sizeTotal += res.ContentLength
			}
		}
	}
	r.rps = float64(len(r.lats)) / r.total.Seconds()
	r.average = r.avgTotal / float64(len(r.lats))
	r.avgConn = r.avgConn / float64(len(r.lats))
	r.avgDelay = r.avgDelay / float64(len(r.lats))
	r.avgDNS = r.avgDNS / float64(len(r.lats))
	r.avgReq = r.avgReq / float64(len(r.lats))
	r.avgRes = r.avgRes / float64(len(r.lats))
	r.print()
}

func (r *Report) printCSV() {
	r.printf("response-time,DNS+dialup,DNS,Request-write,Response-delay,Response-read\n")
	for i, val := range r.lats {
		r.printf("%4.4f,%4.4f,%4.4f,%4.4f,%4.4f,%4.4f\n",
			val, r.connLats[i], r.dnsLats[i], r.reqLats[i], r.delayLats[i], r.resLats[i])
	}
}

func (r *Report) print() {
	if r.output == "csv" {
		r.printCSV()
		return
	}

	if len(r.lats) > 0 {
		sort.Float64s(r.lats)
		r.fastest = r.lats[0]
		r.slowest = r.lats[len(r.lats)-1]
		r.printf("Summary:\n")
		r.printf("  Total:\t%4.4f secs\n", r.total.Seconds())
		r.printf("  Slowest:\t%4.4f secs\n", r.slowest)
		r.printf("  Fastest:\t%4.4f secs\n", r.fastest)
		r.printf("  Average:\t%4.4f secs\n", r.average)
		r.printf("  Requests/sec:\t%4.4f\n", r.rps)
		if r.sizeTotal > 0 {
			r.printf("  Total data:\t%d bytes\n", r.sizeTotal)
			r.printf("  Size/request:\t%d bytes\n", r.sizeTotal/int64(len(r.lats)))
		}
		r.printHistogram()
		r.printLatencies()
		r.printf("\nDetails (average, fastest, slowest):")
		r.printSection("DNS+dialup", r.avgConn, r.connLats)
		r.printSection("DNS-lookup", r.avgDNS, r.dnsLats)
		r.printSection("req write", r.avgReq, r.reqLats)
		r.printSection("resp wait", r.avgDelay, r.delayLats)
		r.printSection("resp read", r.avgRes, r.resLats)
		r.printStatusCodes()
	}
	if len(r.errorDist) > 0 {
		r.printErrors()
	}
	r.printf("\n")
}

// printSection prints details for http-trace fields
func (r *Report) printSection(tag string, avg float64, lats []float64) {
	sort.Float64s(lats)
	fastest, slowest := lats[0], lats[len(lats)-1]
	r.printf("\n  %s:\t", tag)
	r.printf(" %4.4f secs, %4.4f secs, %4.4f secs", avg, fastest, slowest)
}

// printLatencies prints percentile latencies.
func (r *Report) printLatencies() {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(r.lats) && j < len(pctls); i++ {
		current := i * 100 / len(r.lats)
		if current >= pctls[j] {
			data[j] = r.lats[i]
			j++
		}
	}
	r.printf("\nLatency distribution:\n")
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			r.printf("  %v%% in %4.4f secs\n", pctls[i], data[i])
		}
	}
}

func (r *Report) printHistogram() {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (r.slowest - r.fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = r.fastest + bs*float64(i)
	}
	buckets[bc] = r.slowest
	var bi int
	var max int
	for i := 0; i < len(r.lats); {
		if r.lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	r.printf("\nResponse time histogram:\n")
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if max > 0 {
			barLen = (counts[i]*40 + max/2) / max
		}
		r.printf("  %4.3f [%v]\t|%v\n", buckets[i], counts[i], strings.Repeat(barChar, barLen))
	}
}

// printStatusCodes prints status code distribution.
func (r *Report) printStatusCodes() {
	r.printf("\n\nStatus code distribution:\n")
	for code, num := range r.statusCodeDist {
		r.printf("  [%d]\t%d responses\n", code, num)
	}
}

func (r *Report) printErrors() {
	r.printf("\nError distribution:\n")
	for err, num := range r.errorDist {
		r.printf("  [%d]\t%s\n", num, err)
	}
}

func (r *Report) printf(s string, v ...interface{}) {
	fmt.Fprintf(r.w, s, v...)
}
