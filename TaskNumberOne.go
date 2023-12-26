package main

import (
	"LubNumberFour/SpinLock"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numThreads = 5
	raceLength = 50
)

var (
	raceTrack = make([]int, raceLength)
	mu        sync.Mutex
	sem       = make(chan struct{}, numThreads)
	semSlim   = make(chan struct{}, numThreads)
	barrier   = sync.WaitGroup{}
	spinLock  = SpinLock.Locker{}
	monitor   = sync.Mutex{}
)

func SpinWait(duration time.Duration) {
	startTime := time.Now()
	for time.Since(startTime) < duration {
		// Ничего не делаем - активное ожидание
	}
}

func main() {
	runRaceWithPrimitive("Mutex", runRaceWithMutex)
	runRaceWithPrimitive("Semaphore", runRaceWithSemaphore)
	runRaceWithPrimitive("SemaphoreSlim", runRaceWithSemaphoreSlim)
	runRaceWithPrimitive("Barrier", runRaceWithBarrier)
	runRaceWithPrimitive("SpinLock", runRaceWithSpinLock)
	runRaceWithPrimitive("SpinWait", runRaceWithSpinWait)
	runRaceWithPrimitive("Monitor", runRaceWithMonitor)
}

func runRaceWithPrimitive(primitiveName string, raceFunc func(wg *sync.WaitGroup, threadID int)) {
	fmt.Printf("Race start with %s!\n", primitiveName)
	var wg sync.WaitGroup
	wg.Add(numThreads)

	startTime := time.Now()
	for i := 0; i < numThreads; i++ {
		go raceFunc(&wg, i)
	}
	wg.Wait()
	elapsed := time.Since(startTime)

	// Вывод результатов гонки
	fmt.Printf("Results with %s: ", primitiveName)
	for i := 0; i < raceLength; i++ {
		fmt.Print(raceTrack[i], " ")
	}
	fmt.Printf("\nRace finished in %s\n\n", elapsed)
}

func runRaceWithMutex(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		mu.Lock()
		raceTrack[i] = threadID
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10))) // Simulating random speed
		mu.Unlock()
	}
}

func runRaceWithSemaphore(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		sem <- struct{}{}
		raceTrack[i] = threadID
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10))) // Simulating random speed
		<-sem
	}
}

func runRaceWithSemaphoreSlim(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		semSlim <- struct{}{}
		raceTrack[i] = threadID
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10))) // Simulating random speed
		<-semSlim
	}
}

func runRaceWithBarrier(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		barrier.Wait()
		raceTrack[i] = threadID
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10))) // Simulating random speed
	}
}

func runRaceWithSpinLock(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		spinLock.Lock()
		raceTrack[i] = threadID
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10))) // Simulating random speed
		spinLock.Unlock()
	}
}

func runRaceWithSpinWait(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		spinLock.Lock()
		raceTrack[i] = threadID
		spinLock.Unlock()

		// Симулируем случайную скорость
		SpinWait(time.Millisecond * time.Duration(rand.Intn(10)))
	}
}

func runRaceWithMonitor(wg *sync.WaitGroup, threadID int) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < raceLength; i++ {
		monitor.Lock()
		raceTrack[i] = threadID
		monitor.Unlock()

		// Симулируем случайную скорость
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
	}
}
