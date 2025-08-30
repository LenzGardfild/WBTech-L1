package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Проверка на количество воркеров
	if len(os.Args) < 2 {
		fmt.Println("Укажите количество воркеров.")
		return
	}

	// Парсим количество воркеров из аргумента
	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || numWorkers <= 0 {
		fmt.Println("Неверное количество воркеров.")
		return
	}

	dataChannel := make(chan string)

	var wg sync.WaitGroup

	// Запускаем воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, dataChannel, &wg)
	}

	go producer(dataChannel)

	wg.Wait()
}

// producer записывает данные в канал
func producer(dataChannel chan string) {
	for {
		// Записываем данные в канал каждую секунду
		dataChannel <- fmt.Sprintf("Данные от продюсера: %s", time.Now().Format(time.RFC3339))
		time.Sleep(1 * time.Second)
	}
}

// worker читает данные из канала и выводит в stdout
func worker(id int, dataChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range dataChannel {
		fmt.Printf("Воркер %d получил: %s\n", id, data)
	}
}
