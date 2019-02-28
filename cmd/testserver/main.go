package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// make N buckets (N%2==0) in [mean-domain/2.0, ..., mean, ..., mean+domain/2.0],
	normMean   = float64(10.0)
	normDomain = float64(20.0)
	normCount  = 20

	rpcDurationsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "rpc_durations_histogram_seconds",
		Help:    "RPC latency distributions.",
		Buckets: prometheus.LinearBuckets(normMean-(normDomain/2.0), normDomain/float64(normCount), normCount),
	})
)

func main() {
	prometheus.MustRegister(rpcDurationsHistogram)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			<-c
			fmt.Println("adding a new sample. what should its value be?")

			str, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("failed to read input; err= %q\n", err)
				continue
			}

			val, err := strconv.ParseFloat(str[:len(str)-1], 64)
			if err != nil {
				fmt.Printf("cannot parse float; err= %q\n", err)
				continue
			}

			rpcDurationsHistogram.Observe(val)
			fmt.Printf("added sample with value: %v\n", val)
		}
	}()

	fmt.Println("running server")
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("cannot run server; err= %q\n")
	}
}
