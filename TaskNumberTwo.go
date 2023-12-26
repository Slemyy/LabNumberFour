package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Apartment Структура для представления квартир
type Apartment struct {
	Address  string
	Rooms    int
	Cost     float64
	Distance float64
}

// Для вывода с квартирами
func main() {
	// Задаем данные
	arraySize := 50000
	numThreads := 4
	maxDistance := 1.0 // Максимальное расстояние до метро в километрах

	// Генерируем массив квартир
	apartments := generateApartments(arraySize)

	// Измеряем время выполнения без многопоточности
	startTime := time.Now()
	resultsSequential := findApartmentsSequential(apartments, maxDistance)
	sequentialTime := time.Since(startTime)

	// Измеряем время выполнения с использованием многопоточности
	startTime = time.Now()
	resultsParallel := findApartmentsParallel(apartments, maxDistance, numThreads)
	parallelTime := time.Since(startTime)

	// Вывод результатов и времени выполнения
	fmt.Println("Results without concurrency:")
	printApartments(resultsSequential)
	fmt.Printf("\nElapsed time (sequential): %s\n\n", sequentialTime)

	fmt.Println("Results with concurrency (using", numThreads, "threads):")
	printApartments(resultsParallel)
	fmt.Printf("\nElapsed time (parallel): %s\n", parallelTime)
}

//// Для вывода только времени
//func main() {
//	// Задаем данные
//	arraySize := 50000
//	numThreads := 4
//	maxDistance := 1.0 // Максимальное расстояние до метро в километрах
//
//	// Генерируем массив квартир
//	apartments := generateApartments(arraySize)
//
//	// Измеряем время выполнения без многопоточности
//	startTime := time.Now()
//	_ = findApartmentsSequential(apartments, maxDistance)
//	sequentialTime := time.Since(startTime)
//
//	// Измеряем время выполнения с использованием многопоточности
//	startTime = time.Now()
//	_ = findApartmentsParallel(apartments, maxDistance, numThreads)
//	parallelTime := time.Since(startTime)
//
//	// Вывод времени выполнения
//	fmt.Println("Results without concurrency:")
//	fmt.Printf("\nElapsed time (sequential): %s\n\n", sequentialTime)
//
//	fmt.Println("Results with concurrency (using", numThreads, "threads):")
//	fmt.Printf("\nElapsed time (parallel): %s\n", parallelTime)
//}

// Генерация случайных квартир
func generateApartments(size int) []Apartment {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	apartments := make([]Apartment, size)
	for i := 0; i < size; i++ {
		address := fmt.Sprintf("Address%d", i)
		rooms := i%5 + 1
		cost := float64(rand.Intn(10000) + 500) // Генерация случайной стоимости от 500 до 10500
		distance := rand.Float64() * 5          // Генерация случайного расстояния от 0 до 5
		apartments[i] = Apartment{Address: address, Rooms: rooms, Cost: cost, Distance: distance}
	}
	return apartments
}

// Поиск квартир с расстоянием до метро меньше заданного без многопоточности
func findApartmentsSequential(apartments []Apartment, maxDistance float64) []Apartment {
	var results []Apartment
	var totalCost float64

	for _, apartment := range apartments {
		if apartment.Distance < maxDistance {
			results = append(results, apartment)
			totalCost += apartment.Cost
		}
	}

	// Фильтруем квартиры со стоимостью ниже средней
	avgCost := totalCost / float64(len(results))
	var filteredResults []Apartment
	for _, apartment := range results {
		if apartment.Cost < avgCost {
			filteredResults = append(filteredResults, apartment)
		}
	}

	return filteredResults
}

// Поиск квартир с расстоянием до метро меньше заданного с использованием многопоточности
func findApartmentsParallel(apartments []Apartment, maxDistance float64, numThreads int) []Apartment {
	var results []Apartment
	var totalCost float64
	var mu sync.Mutex
	wg := sync.WaitGroup{}

	// Функция поиска в отдельном потоке
	findInThread := func(threadID int, startIdx int, endIdx int) {
		defer wg.Done()
		localResults := make([]Apartment, 0)
		localTotalCost := 0.0

		for i := startIdx; i < endIdx; i++ {
			if apartments[i].Distance < maxDistance {
				localResults = append(localResults, apartments[i])
				localTotalCost += apartments[i].Cost
			}
		}

		// Записываем результаты из локального массива в общий с защитой мьютексом
		mu.Lock()
		results = append(results, localResults...)
		totalCost += localTotalCost
		mu.Unlock()
	}

	// Разбиваем массив на части для обработки в потоках
	partSize := len(apartments) / numThreads
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		startIdx := i * partSize
		endIdx := (i + 1) * partSize
		if i == numThreads-1 {
			endIdx = len(apartments)
		}
		go findInThread(i, startIdx, endIdx)
	}

	wg.Wait()

	// Фильтруем квартиры со стоимостью ниже средней
	avgCost := totalCost / float64(len(results))
	var filteredResults []Apartment
	for _, apartment := range results {
		if apartment.Cost < avgCost {
			filteredResults = append(filteredResults, apartment)
		}
	}

	return filteredResults
}

// Вывод списка квартир
func printApartments(apartments []Apartment) {
	for _, apartment := range apartments {
		fmt.Printf("Address: %s, Rooms: %d, Cost: %.2f, Distance: %.2f km\n", apartment.Address, apartment.Rooms, apartment.Cost, apartment.Distance)
	}
}
