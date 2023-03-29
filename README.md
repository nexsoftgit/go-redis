# Redis client for Golang




Supports:

- Redis 3 commands except QUIT, MONITOR, SLOWLOG and SYNC.
- Automatic connection pooling with [circuit breaker](https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern) support.
- [Pub/Sub](https://godoc.org/github.com/nexsoftgit/go-redis#PubSub).
- [Transactions](https://godoc.org/github.com/nexsoftgit/go-redis#Multi).
- [Pipeline](https://godoc.org/github.com/nexsoftgit/go-redis#example-Client-Pipeline) and [TxPipeline](https://godoc.org/github.com/nexsoftgit/go-redis#example-Client-TxPipeline).
- [Scripting](https://godoc.org/github.com/nexsoftgit/go-redis#Script).
- [Timeouts](https://godoc.org/github.com/nexsoftgit/go-redis#Options).
- [Redis Sentinel](https://godoc.org/github.com/nexsoftgit/go-redis#NewFailoverClient).
- [Redis Cluster](https://godoc.org/github.com/nexsoftgit/go-redis#NewClusterClient).
- [Cluster of Redis Servers](https://godoc.org/github.com/nexsoftgit/go-redis#example-NewClusterClient--ManualSetup) without using cluster mode and Redis Sentinel.
- [Ring](https://godoc.org/github.com/nexsoftgit/go-redis#NewRing).
- [Instrumentation](https://godoc.org/github.com/nexsoftgit/go-redis#ex-package--Instrumentation).
- [Cache friendly](https://github.com/go-redis/cache).
- [Rate limiting](https://github.com/go-redis/redis_rate).
- [Distributed Locks](https://github.com/bsm/redis-lock).


## Installation

Install:

```shell
go get -u github.com/nexsoftgit/go-redis
```

Import:

```go
import "github.com/nexsoftgit/go-redis"
```

## Quickstart

```go
func ExampleNewClient() {
	opts := &redis.Options{
		CircuitBreaker: optCB,
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client := redis.NewClient(opts)

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}
```
### With Circuit Breaker
```go
func WithCircuitBreaker() {
	
	optCB := &hystrix.CommandConfig{
		Timeout:                10000,
		RequestVolumeThreshold: 2,
		SleepWindow:            500,
		ErrorPercentThreshold:  5,
	}
	// Please, read more about command config in hystrix-go doc.
	//https://godoc.org/github.com/afex/hystrix-go/hystrix#pkg-variables)

	opts := &redis.Options{
		CircuitBreaker: optCB,
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client := redis.NewClient(opts)

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

```
Any changes in breaker state will generate metrics for monitoring. Below is a list of the metrics.
``` yaml
Name: "circuit_breaker"
Help: "A total number of redis client make a request to redis server with circuit breaker state."
Labels: 
	- command : "redis command" 
	- service : "service name"
	- status : "fail/ok"
	- state: 
		- "circuit breaker open": A total number of circuit breaker state open. This happens due to the circuit being measured as unhealthy.
		- "max_concurency": A total number of client executed at the same time and exceeded max concurrency.
		- "timeout": A total number of request exceeded timeout duration

```
## Howto

Please go through [examples](example_test.go) to get an idea how to use this package.

## Look and feel

Some corner cases:

    SET key value EX 10 NX
    set, err := client.SetNX("key", "value", 10*time.Second).Result()

    SORT list LIMIT 0 2 ASC
    vals, err := client.Sort("list", redis.Sort{Offset: 0, Count: 2, Order: "ASC"}).Result()

    ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
    vals, err := client.ZRangeByScoreWithScores("zset", redis.ZRangeBy{
        Min: "-inf",
        Max: "+inf",
        Offset: 0,
        Count: 2,
    }).Result()

    ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
    vals, err := client.ZInterStore("out", redis.ZStore{Weights: []int64{2, 3}}, "zset1", "zset2").Result()

    EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
    vals, err := client.Eval("return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()

## Benchmark

go-redis vs redigo:

```
BenchmarkSetGoRedis10Conns64Bytes-4 	  200000	      7621 ns/op	     210 B/op	       6 allocs/op
BenchmarkSetGoRedis100Conns64Bytes-4	  200000	      7554 ns/op	     210 B/op	       6 allocs/op
BenchmarkSetGoRedis10Conns1KB-4     	  200000	      7697 ns/op	     210 B/op	       6 allocs/op
BenchmarkSetGoRedis100Conns1KB-4    	  200000	      7688 ns/op	     210 B/op	       6 allocs/op
BenchmarkSetGoRedis10Conns10KB-4    	  200000	      9214 ns/op	     210 B/op	       6 allocs/op
BenchmarkSetGoRedis100Conns10KB-4   	  200000	      9181 ns/op	     210 B/op	       6 allocs/op
BenchmarkSetGoRedis10Conns1MB-4     	    2000	    583242 ns/op	    2337 B/op	       6 allocs/op
BenchmarkSetGoRedis100Conns1MB-4    	    2000	    583089 ns/op	    2338 B/op	       6 allocs/op
BenchmarkSetRedigo10Conns64Bytes-4  	  200000	      7576 ns/op	     208 B/op	       7 allocs/op
BenchmarkSetRedigo100Conns64Bytes-4 	  200000	      7782 ns/op	     208 B/op	       7 allocs/op
BenchmarkSetRedigo10Conns1KB-4      	  200000	      7958 ns/op	     208 B/op	       7 allocs/op
BenchmarkSetRedigo100Conns1KB-4     	  200000	      7725 ns/op	     208 B/op	       7 allocs/op
BenchmarkSetRedigo10Conns10KB-4     	  100000	     18442 ns/op	     208 B/op	       7 allocs/op
BenchmarkSetRedigo100Conns10KB-4    	  100000	     18818 ns/op	     208 B/op	       7 allocs/op
BenchmarkSetRedigo10Conns1MB-4      	    2000	    668829 ns/op	     226 B/op	       7 allocs/op
BenchmarkSetRedigo100Conns1MB-4     	    2000	    679542 ns/op	     226 B/op	       7 allocs/op
```

Redis Cluster:

```
BenchmarkRedisPing-4                	  200000	      6983 ns/op	     116 B/op	       4 allocs/op
BenchmarkRedisClusterPing-4         	  100000	     11535 ns/op	     117 B/op	       4 allocs/op
```

## See also

- [Golang PostgreSQL ORM](https://github.com/go-pg/pg)
- [Golang msgpack](https://github.com/vmihailenco/msgpack)
- [Golang message task queue](https://github.com/go-msgqueue/msgqueue)
