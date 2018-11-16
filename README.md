# udpproto

udpproto is a beat based on metricbeat which was generated with metricbeat/metricset generator.
The beat provides an alternative to get golang heap metrics from a service to elasticsearch. 
It does not fetch the data, but listens on a port to receive binary packages (in proto3 format) and pushes the data into ElasticSearch.
It uses the same document format as the `golang` module so the dashbards can be re-used.

    Note: metricbeats has a module `golang` that _fetches_ the heap metrics from a golang application by querying an http port.



```
    Your application 
    -> write udp message 
    UdpProto server
    -> udpproto receives message 
    -> posts the golang heap data to ElasticSearch
    ElasticSearch
    -> use existing dashboard [Meticbeat Golang]
    
```

See the testclient as an example how your application could send the data.

# The data

| proto field name 	                | elastic field 	                |
|--------------	                    |-------------------	            |
| timestamp_nano                 	| @timestamp 						|
| process_name 	                    | system.process.name               |
| process_id 	                    | system.process.pid                |
| process_cmdline                   | system.process.cmdline            |
| allocations_active 	            | golang.heap.allocations.active 	|
| allocations_allocated         	| golang.heap.allocations.allocated |
| allocations_frees             	| golang.heap.allocations.frees     |
| allocations_idle              	| golang.heap.allocations.idle  	|
| allocations_mallocs           	| golang.heap.allocations.mallocs  	|
| allocations_objects           	| golang.heap.allocations.obects    |
| allocations_total             	| golang.heap.allocations.total     |
| gc_cpu_fraction               	| golang.heap.gc.cpu.fraction       |
| gc_next_gc_limit              	| golang.heap.gc.next.gc_limit      |
| gc_pause_avg_ns               	| golang.heap.gc.pause.avg_ns       |
| gc_pause_count                	| golang.heap.gc.pause.count        |
| gc_pause_max_ns               	| golang.heap.gc.pause.max_ns       |
| gc_pause_sum_ns               	| golang.heap.gc.pause.sum_ns       |
| gc_total_count                	| golang.heap.gc.total.count        |
| gc_total_pause_ns             	| golang.heap.gc.total.pause_ns     |
| system_obtained               	| golang.heap.system.obtaines       |
| system_released               	| golang.heap.system.released       |
| system_stack                  	| golang.heap.system.stack          |
| system_total                  	| golang.heap.system.total          |


## Motivation for this beat

We did not want to open a port on our running services. This is the Hollywood approach: don't call us, we call you.


# Installation

Something I have to figure out yet.... Ideally I like to have it part of just the metricbeat, but the udpproto service running in a seperate process. /etc/udpproto

