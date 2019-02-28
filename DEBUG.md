# Debug

This is a guide to run the "debug-enabled" Prometheus and a testserver that
exports metrics, to better understand how things work.

## Go version

Prometheus recently switched to use Go modules, so you would need it to build.
Download Go 1.11 and above.

## Building Prometheus and testserver

```
$> GO111MODULE=on go build -mod=vendor -o prometheus ./cmd/prometheus
$> GO111MODULE=on go build -mod=vendor -o testserver ./cmd/testserver
```

## Running test setup:

You will need to 3-4 windows/tabs/tmux panes:

#### 1.testserver

run testserver, listen to "0.0.0.0:8080/metrics" and export metric:
`rpc_durations_histogram_seconds`.

```
$> ./testserver
```

#### 2.Prometheus

run a modified Prometheus (only scrape on when receiving SIGUSR1 signal,
instead of periodically scrape the data).

```
$> ./prometheus --config.file=./config.d/prometheus.yml
```

#### 3.Signal testserver to add a sample to the histogram

```
$> kill -s SIGUSR1 `ps -ef | grep '\.\/testserver' | awk '{print $2}'`
```

#### 4.Signal Prometheus to scrape metrics

```
$> kill -s SIGUSR1 `ps -ef | grep '\.\/prometheus' | awk '{print $2}'`
```

#### 5. (optional), look at output of testserver

```
$> curl -sX GET '127.0.0.1:8080/metrics' -i | grep rpc_durations
```

## Observe

You can go to the Prometheus UI at <http://localhost:9090> to run queries.
