package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numReaders    = 5
	numWriters    = 5
	readDuration  = 100 * time.Millisecond
	writeDuration = 150 * time.Millisecond
)

type Resource struct {
	data  int
	mutex sync.RWMutex
}

func main() {
	resource := &Resource{}
	var wg sync.WaitGroup

	// Запускаем читателей
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go reader(i, resource, &wg)
	}

	// Запускаем писателей
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go writer(i, resource, &wg)
	}

	wg.Wait()
}

func reader(id int, resource *Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(readDuration)
		resource.mutex.RLock()
		fmt.Printf("Reader %d is reading: %d\n", id, resource.data)
		resource.mutex.RUnlock()
	}
}

func writer(id int, resource *Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(writeDuration)
		resource.mutex.Lock()
		resource.data++
		fmt.Printf("Writer %d is writing: %d\n", id, resource.data)
		resource.mutex.Unlock()
	}
}
