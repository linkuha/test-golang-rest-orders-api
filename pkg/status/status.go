package status

import (
	"fmt"
	"runtime"
	"time"
)

const DateTimeFormat = "2006-01-02 15:04:05 -0700"

var startTime = time.Now()
var startTimeString = time.Now().Format(DateTimeFormat)

type StatusInfo struct {
	MemoryUsage  appMemoryUsage
	NumGoroutine int
	NumCPU       int
	NumCgoCall   int64
	GoVersion    string
	Server       serverInfo
}

type appMemoryUsage struct {
	Alloc        string `json:"Alloc"`
	TotalAlloc   string `json:"TotalAlloc"`
	HeapAlloc    string `json:"HeapAlloc"`
	HeapReleased string `json:"HeapReleased"`
	Sys          string `json:"Sys"`
	NumGC        string `json:"NumGC"`
	LastGC       string `json:"LastGC"`
}

type serverInfo struct {
	Goos        string `json:"goos"`
	ServerTime  string `json:"serverTime"`
	ServerStart string `json:"serverStart"`
	Uptime      string `json:"serverUptime"`
}

func GetStatusInfo() StatusInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	lastGC := time.Unix(0, int64(m.LastGC)).Format(DateTimeFormat)
	return StatusInfo{
		MemoryUsage: appMemoryUsage{
			Alloc:        fmt.Sprintf("%v MiB", bToMb(m.Alloc)),
			TotalAlloc:   fmt.Sprintf("%v MiB", bToMb(m.TotalAlloc)),
			Sys:          fmt.Sprintf("%v MiB", bToMb(m.Sys)),
			HeapAlloc:    fmt.Sprintf("%v MiB", bToMb(m.HeapAlloc)),
			HeapReleased: fmt.Sprintf("%v MiB", bToMb(m.HeapReleased)),
			NumGC:        fmt.Sprintf("%v", m.NumGC),
			LastGC:       lastGC,
		},
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		NumCgoCall:   runtime.NumCgoCall(),
		GoVersion:    runtime.Version(),
		Server: serverInfo{
			Goos:        runtime.GOOS,
			ServerTime:  time.Now().Format(DateTimeFormat),
			ServerStart: startTimeString,
			Uptime:      time.Now().Sub(startTime).String(),
		},
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
