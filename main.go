package main

import (
	"celeryCliStats/minheap"
	"celeryCliStats/model"
	"context"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ctx = context.Background()

func main() {
	redisAddr := flag.String("redis-addr", "localhost:6379", "Address of the Redis server")
	redisPassword := flag.String("redis-password", "", "Password for the Redis server (default is no password)")
	redisDB := flag.Int("redis-db", 0, "Redis database number to use")
	queueName := flag.String("queue-name", "production_default", "Name of the Redis queue")
	topK := flag.Int("top-k", 100, "Number of top frequent elements to return")
	delay := flag.Int("delay", -1, "The delay to apply while polling results, if -1 won't loop")

	flag.Parse()

	// Setup a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: *redisPassword,
		DB:       *redisDB,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	for true {
		go doPollQueue(ctx, redisClient, queueName, topK)
		if *delay > 0 {
			time.Sleep(time.Duration(*delay) * time.Millisecond)
		} else {
			wg.Done()
			break
		}
		select {
		case _ = <-signalChan:
			cancel()
			wg.Done()
		default:
		}
	}
	wg.Wait()
}

func doPollQueue(ctx context.Context, redisClient *redis.Client, queueName *string, topK *int) {
	result, err := redisClient.LRange(ctx, *queueName, 0, -1).Result()
	if err != nil {
		log.Fatalf("Error fetching from Redis: %v", err)
	}

	// Count the occurrences of each task
	taskCounts := make(map[string]int)
	for _, l := range result {
		var message model.Message
		err := jsoniter.Unmarshal([]byte(l), &message)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			continue
		}
		taskCounts[message.Headers.Task]++
	}

	// Build a min-heap
	h := minheap.NewMinHeap()

	for value, freq := range taskCounts {
		h.Push(&minheap.Element{Value: value, Frequency: freq})
	}

	// Get the top K frequent elements
	topKFrequent := h.PopTopFrequent(*topK)
	fmt.Printf("Top %d tasks:\n", *topK)
	for _, elem := range topKFrequent {
		fmt.Printf("Task: %s, Count: %d\n", elem.Value, elem.Frequency)
	}
}
