init:
	sudo apt install ethtool sysbench

call:
	lscpu && cat /proc/cpuinfo | grep 'model name' | uniq
	cat /proc/meminfo | grep MemTotal
	ethtool -i eth0
	tc -s qdisc show dev eth0

net:
	sudo tc qdisc add dev eth0 root netem delay 0.5ms loss 0.01% rate 5gbit 

sysbench:
	sysbench cpu --threads=1 --cpu-max-prime=20000 run | grep 'events per second'

buildapp:
	go build -o ./build/app/main ./cmd/app/main.go

run: buildapp
	./build/app/main

core = 0

test: buildapp
	taskset -c ${core} ./build/app/main

benchtime = 1s

id = 1

bench:
	go test -v -bench=. ./tests/benchmark${id}/... -benchtime=${benchtime}