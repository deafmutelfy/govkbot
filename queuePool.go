package main

import "github.com/golang-queue/queue"

var queuePool map[string]*queue.Queue = map[string]*queue.Queue{}
var queueIds []int = []int{}
