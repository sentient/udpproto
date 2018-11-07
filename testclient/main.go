package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/golang/protobuf/proto"
	protobuf "github.com/sentient/udpproto/testclient/generated"
)

//nolint:errcheck
func main() {
	fmt.Println("start test client")
	to, err := net.ResolveUDPAddr("udp", ":8127")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, to)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	col := &collector{}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-c:
			return
		case <-ticker.C:
			col.sendMetrics(conn)
		}
	}

}

type collector struct {
	lastNumGC uint32
}

func (c *collector) sendMetrics(conn *net.UDPConn) {

	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)

	g := &protobuf.GolangHeap{
		TimestampNano:        time.Now().UnixNano(),
		AllocationsMallocs:   ms.Mallocs,
		AllocationsFrees:     ms.Frees,
		AllocationsObjects:   ms.HeapObjects,
		AllocationsTotal:     ms.TotalAlloc,
		AllocationsAllocated: ms.HeapAlloc,
		AllocationsIdle:      ms.HeapIdle,
		AllocationsActive:    ms.HeapInuse,

		SystemTotal:    ms.Sys,
		SystemObtained: ms.HeapSys,
		SystemStack:    ms.StackSys,
		SystemReleased: ms.HeapReleased,
	}

	// garbage collector summary
	var duration, maxDuration, avgDuration, count uint64
	// collect last gc run stats
	if c.lastNumGC < ms.NumGC {
		delta := ms.NumGC - c.lastNumGC
		start := c.lastNumGC
		if delta > 256 {
			logp.Debug("golang", "Missing %v gc cycles", delta-256)
			start = ms.NumGC - 256
			delta = 256
		}

		end := start + delta
		for i := start; i < end; i++ {
			idx := i % 256
			d := ms.PauseNs[idx]
			count++
			duration += d
			if d > maxDuration {
				maxDuration = d
			}
		}

		avgDuration = duration / count
		c.lastNumGC = ms.NumGC
	}

	g.GcNextGcLimit = ms.NextGC
	g.GcTotalCount = ms.NumGC
	g.GcCpuFraction = ms.GCCPUFraction
	g.GcTotalPauseNs = ms.PauseTotalNs
	g.GcPauseCount = count
	g.GcPauseSumNs = duration
	g.GcPauseAvgNs = avgDuration
	g.GcPauseMaxNs = maxDuration

	b, err := proto.Marshal(g)
	if err != nil {
		panic(err)
	}

	n, err := conn.Write(b)
	// msg := "test message"
	// n, err := fmt.Fprintf(conn, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("I wrote %d \n", n)

	g = nil
	ms = nil

}
