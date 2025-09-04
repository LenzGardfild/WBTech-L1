package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: завершение...\n", id)
			return
		default:
			fmt.Printf("Worker %d: работаю\n", id)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// Контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())

	// Ловим SIGINT (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	workers := 3

	// Запускаем воркеров
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	// Ждём сигнал
	<-sigChan
	fmt.Println("\nПолучен сигнал, завершаем...")
	cancel()

	// Ждём завершения всех воркеров
	wg.Wait()
	fmt.Println("Все воркеры завершены. Выход.")
}

/* Для корректного завершения всех воркеров при нажатии Ctrl+C удобно использовать context.WithCancel
Когда программа получает сигнал, мы вызываем cancel()
Все воркеры слушают ctx.Done() и выходят из цикла при отмене
Контекст — стандартный инструмент для управления временем жизни горутин
*/
