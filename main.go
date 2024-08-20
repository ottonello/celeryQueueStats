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
)

var ctx = context.Background()

func main() {
	redisAddr := flag.String("redis-addr", "localhost:6379", "Address of the Redis server")
	redisPassword := flag.String("redis-password", "", "Password for the Redis server (default is no password)")
	redisDB := flag.Int("redis-db", 0, "Redis database number to use")
	queueName := flag.String("queue-name", "production_default", "Name of the Redis queue")
	topK := flag.Int("top-k", 100, "Number of top frequent elements to return")

	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: *redisPassword,
		DB:       *redisDB,
	})

	// Fetch the list from Redis
	result, err := rdb.LRange(ctx, *queueName, 0, -1).Result()
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
	fmt.Printf("Top %d frequent elements with their counts:\n", *topK)
	for _, elem := range topKFrequent {
		fmt.Printf("Element: %s, Count: %d\n", elem.Value, elem.Frequency)
	}
}
