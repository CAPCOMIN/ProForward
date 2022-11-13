package Controllers

import (
	"bytes"
	"fmt"
	"forward-server/Controllers/BaseCtrl"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/utils"
	"github.com/dlclark/regexp2"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"
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
type grInfo struct {
	Id       int
	Status   string
	Duration string
	Detail   string
}

//func (a *PerformanceController) registerHttpServer(svr *HttpServer) {
//	a.servers = append(a.servers, svr)
//}

func (a *PerformanceController) TrimEmptyLinesAndDelFirstLine(str string) string {
	strs := strings.Split(str, "\n")
	str = ""
	strs = strs[1:]
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + "\n"
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}

func (a *PerformanceController) PerfGoroutine() {
	var (
		//format = r.Form.Get("format")
		//data   = make(map[interface{}]interface{})
		result bytes.Buffer
	)
	var grOne grInfo
	var grList []grInfo
	ProcessInput("lookup goroutine", &result)
	goroutineInfo := result.String()
	goroutineInfoList := strings.Split(goroutineInfo, "\n\n")

	grStatusReg := regexp2.MustCompile(`(?<=\[).+(?=,)|(?<=\[).+(?=])`, 0)
	grIdReg := regexp2.MustCompile(`(?<=goroutine )\d+`, 0)
	grDurationReg := regexp2.MustCompile(`\d+ minutes`, 0)

	for i, num := 0, len(goroutineInfoList); i < num; i++ {
		idGroup, _ := grIdReg.FindStringMatch(goroutineInfoList[i])
		grOne.Id, _ = strconv.Atoi(idGroup.String())

		statusGrp, _ := grStatusReg.FindStringMatch(goroutineInfoList[i])
		grOne.Status = statusGrp.String()

		durationGrp, _ := grDurationReg.FindStringMatch(goroutineInfoList[i])
		if durationGrp == nil {
			grOne.Duration = "暂无"
		} else {
			grOne.Duration = durationGrp.String()
		}

		grOne.Detail = a.TrimEmptyLinesAndDelFirstLine(goroutineInfoList[i])

		grList = append(grList, grOne)
		logs.Warn("PerfGoroutine", grOne.Id, grOne.Status, grOne.Duration)
	}

	//logs.Warn("PerfGoroutine ", len(goroutineInfoList))

	a.Data["Content"] = template.HTMLEscapeString(goroutineInfo)

	a.Data["Title"] = template.HTMLEscapeString("lookup goroutine")

	a.Data["goroutineList"] = grList

	//writeTemplate(rw, data, "performance/goroutine.html")
	a.TplName = "performance/goroutine.html"
}

func (a *PerformanceController) PerfCPU() {

	a.Data["Title"] = "CPU调用图分析"
	a.TplName = "performance/cpu.html"
}

func (a *PerformanceController) PerfHeap() {
	var (
		result bytes.Buffer
	)
	ProcessInput("lookup heap", &result)
	heapInfo := template.HTMLEscapeString(result.String())
	heapInfoList := strings.Split(heapInfo, "\n\n")

	heapDetailList := strings.Split(heapInfoList[len(heapInfoList)-1], "\n")
	a.Data["Summary"] = heapDetailList[1:]
	a.Data["Detail"] = heapInfoList[:len(heapInfoList)-1]

	a.TplName = "performance/heap.html"
}

func (a *PerformanceController) MonitorCPU() {
	a.TplName = "performance/monitorCPU.html"
}

type MonitorResult struct {
	Code     int
	Msg      string
	Filename string
	Content  string
}

func (a *PerformanceController) DoMonitorCPU() {
	var result bytes.Buffer
	secFloat, err0 := strconv.ParseFloat(a.GetString("sec"), 10)
	//logs.Warn("DoMonitorCPU", secFloat)
	sec := int(math.Floor(secFloat + 0.5))
	//logs.Warn("DoMonitorCPU", sec)
	err, filename := GetCPUProfile(sec, &result)
	content := template.HTMLEscapeString(result.String())

	if err != nil || err0 != nil {
		logs.Error("MonitorCPU", err)
		a.Data["json"] = MonitorResult{Code: 1, Msg: "获取CPU调用失败，" + err.Error() + err0.Error(),
			Filename: filename, Content: content}
	} else {
		a.Data["json"] = MonitorResult{Code: 0, Msg: "已成功获取CPU调用信息",
			Filename: filename, Content: content}
	}
	logs.Warn("DoMonitorCPU", filename, content)
	a.ServeJSON()
}

func (a *PerformanceController) PprofWebAnalysis() {
	var cmd *exec.Cmd
	operate := a.GetString("operate")

	cmd = exec.Command("dir")

	if operate == "start" {
		err := cmd.Start()
		if err != nil {
			logs.Error("PprofWebAnalysis 1", err)
		}
		logs.Warn("state =>", cmd.Dir)
		logs.Warn("result =>", cmd.ProcessState.String())
		//a.Data["json"]
	} else if operate == "stop" {
		err := cmd.Process.Kill()
		if err != nil {
			logs.Error("PprofWebAnalysis 2", err)
		}
	} else if operate == "status" {
		//todo
	}

}

func (a *PerformanceController) PerfThreadBlockGC() {
	var (
		result1 bytes.Buffer
		result2 bytes.Buffer
		result3 bytes.Buffer
	)
	ProcessInput("lookup threadcreate", &result1)
	threadInfo := template.HTMLEscapeString(result1.String())

	ProcessInput("lookup block", &result2)
	blockInfo := template.HTMLEscapeString(result2.String())

	ProcessInput("gc summary", &result3)
	gcInfo := template.HTMLEscapeString(result3.String())

	a.Data["thread"] = threadInfo
	a.Data["block"] = blockInfo
	a.Data["gc"] = gcInfo

	a.TplName = "performance/tbgc.html"
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
	//case "get cpuprof":
	//GetCPUProfile(w)
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
func GetCPUProfile(sec int, w io.Writer) (error, string) {
	//sec := 10
	//filename := "cpu-" + strconv.Itoa(pid) + ".pprof"
	filename := "cpu-0.pprof"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(w, "Could not enable CPU profiling: %s\n", err)
		log.Fatal("record cpu profile failed: ", err)
	}
	err2 := pprof.StartCPUProfile(f)
	time.Sleep(time.Duration(sec) * time.Second)
	pprof.StopCPUProfile()

	fmt.Fprintf(w, "create cpu profile %s \n", filename)
	_, fl := path.Split(os.Args[0])
	fmt.Fprintf(w, "Now you can use this to check it: go tool pprof %s %s\n", fl, filename)

	return err2, filename
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
