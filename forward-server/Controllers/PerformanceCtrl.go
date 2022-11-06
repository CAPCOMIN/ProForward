package Controllers

import (
	"bytes"
	"fmt"
	"forward-server/Controllers/BaseCtrl"
	"github.com/beego/beego/v2/core/utils"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"text/template"
	"time"
)

var (
	startTime = time.Now()
	pid       int
)

type PerformanceController struct {
	BaseCtrl.ConsoleCtrl
	//servers []*HttpServer
}

//func (a *PerformanceController) registerHttpServer(svr *HttpServer) {
//	a.servers = append(a.servers, svr)
//}

func (a *PerformanceController) PerfGoroutine() {
	//rw, r := a.Ctx.ResponseWriter, a.Ctx.Request
	//err := r.ParseForm()
	//if err != nil {
	//	return
	//}
	//command := a.GetString("command")
	//if command == "" {
	//	return
	//}

	var (
		//format = r.Form.Get("format")
		//data   = make(map[interface{}]interface{})
		result bytes.Buffer
	)
	ProcessInput("lookup goroutine", &result)
	a.Data["Content"] = template.HTMLEscapeString(result.String())

	//if format == "json" && command == "gc summary" {
	//	dataJSON, err := json.Marshal(data)
	//	if err != nil {
	//		http.Error(rw, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	writeJSON(rw, dataJSON)
	//	return
	//}

	a.Data["Title"] = template.HTMLEscapeString("lookup goroutine")

	//writeTemplate(rw, data, "performance/goroutine.html")
	a.TplName = "performance/goroutine.html"
}

func ProcessInput(input string, w io.Writer) {
	switch input {
	case "lookup goroutine":
		p := pprof.Lookup("goroutine")
		p.WriteTo(w, 2)
	case "lookup heap":
		p := pprof.Lookup("heap")
		p.WriteTo(w, 2)
	case "lookup threadcreate":
		p := pprof.Lookup("threadcreate")
		p.WriteTo(w, 2)
	case "lookup block":
		p := pprof.Lookup("block")
		p.WriteTo(w, 2)
	case "get cpuprof":
		GetCPUProfile(w)
	case "get memprof":
		MemProf(w)
	case "gc summary":
		PrintGCSummary(w)
	}
}

// MemProf record memory profile in pprof
func MemProf(w io.Writer) {
	filename := "mem-" + strconv.Itoa(pid) + ".memprof"
	if f, err := os.Create(filename); err != nil {
		fmt.Fprintf(w, "create file %s error %s\n", filename, err.Error())
		log.Fatal("record heap profile failed: ", err)
	} else {
		runtime.GC()
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Fprintf(w, "create heap profile %s \n", filename)
		_, fl := path.Split(os.Args[0])
		fmt.Fprintf(w, "Now you can use this to check it: go tool pprof %s %s\n", fl, filename)
	}
}

// GetCPUProfile start cpu profile monitor
func GetCPUProfile(w io.Writer) {
	sec := 30
	filename := "cpu-" + strconv.Itoa(pid) + ".pprof"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(w, "Could not enable CPU profiling: %s\n", err)
		log.Fatal("record cpu profile failed: ", err)
	}
	pprof.StartCPUProfile(f)
	time.Sleep(time.Duration(sec) * time.Second)
	pprof.StopCPUProfile()

	fmt.Fprintf(w, "create cpu profile %s \n", filename)
	_, fl := path.Split(os.Args[0])
	fmt.Fprintf(w, "Now you can use this to check it: go tool pprof %s %s\n", fl, filename)
}

// PrintGCSummary print gc information to io.Writer
func PrintGCSummary(w io.Writer) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	gcstats := &debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(gcstats)

	printGC(memStats, gcstats, w)
}

func printGC(memStats *runtime.MemStats, gcstats *debug.GCStats, w io.Writer) {
	if gcstats.NumGC > 0 {
		lastPause := gcstats.Pause[0]
		elapsed := time.Since(startTime)
		overhead := float64(gcstats.PauseTotal) / float64(elapsed) * 100
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "NumGC:%d Pause:%s Pause(Avg):%s Overhead:%3.2f%% Alloc:%s Sys:%s Alloc(Rate):%s/s Histogram:%s %s %s \n",
			gcstats.NumGC,
			utils.ToShortTimeFormat(lastPause),
			utils.ToShortTimeFormat(avg(gcstats.Pause)),
			overhead,
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)),
			utils.ToShortTimeFormat(gcstats.PauseQuantiles[94]),
			utils.ToShortTimeFormat(gcstats.PauseQuantiles[98]),
			utils.ToShortTimeFormat(gcstats.PauseQuantiles[99]))
	} else {
		// while GC has disabled
		elapsed := time.Since(startTime)
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "Alloc:%s Sys:%s Alloc(Rate):%s/s\n",
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)))
	}
}

func avg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

// toH format bytes number friendly
func toH(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}