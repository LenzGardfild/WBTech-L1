package main

import (
	"fmt"
	"time"
)

func main() {
	// время работы программы
	N := 5 * time.Second

	ch := make(chan int)

	// горутина-ввод
	go func() {
		i := 1
		for {
			ch <- i
			i++
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// горутина-вывод
	go func() {
		for val := range ch {
			fmt.Println("Получено:", val)
		}
	}()

	// ждем N секунд
	<-time.After(N)
	fmt.Println("Время вышло, завершаем")

	// закрываем канал
	close(ch)
}
