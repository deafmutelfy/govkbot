package main

import (
	"sync"

	"github.com/golang-queue/queue"
)

var queuePool map[string]*queue.Queue = map[string]*queue.Queue{}
var queueIds []int = []int{}

var queuePoolMutex sync.Mutex
var queueIdsMutex sync.Mutex
