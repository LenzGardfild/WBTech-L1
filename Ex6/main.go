// main.go
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

func main() {
	_ = flag.CommandLine

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Println("Выберите способ остановки горутины:")
		fmt.Println(" 1) condition     — обычный выход по внутреннему условию (return)")
		fmt.Println(" 2) donechan      — канал уведомления (закрытие done)")
		fmt.Println(" 3) rangeclose    — закрытие рабочего канала (range завершается)")
		fmt.Println(" 4) ctxcancel     — контекст с ручной отменой (WithCancel)")
		fmt.Println(" 5) ctxtimeout    — контекст с таймаутом (WithTimeout)")
		fmt.Println(" 6) goexit        — runtime.Goexit() (жёстко изнутри горутины)")
		fmt.Println(" 7) panic         — panic в горутине + recover в ней же")
		fmt.Println(" 8) after         — локальный таймаут через time.After в select")
		fmt.Println(" 9) atomic        — атомарный флаг остановки (sync/atomic)")
		fmt.Println("10) mainexit      — выход main прерывает прочие (анти-паттерн)")
		fmt.Println(" a) all           — последовательно выполнить все демо")
		fmt.Println(" q) quit          — выход")
		fmt.Print("\nВаш выбор: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		switch choice {
		case "1", "condition":
			runDemo("condition", demoConditionExit, reader)
		case "2", "donechan":
			runDemo("donechan", demoDoneChannel, reader)
		case "3", "rangeclose":
			runDemo("rangeclose", demoRangeClose, reader)
		case "4", "ctxcancel":
			runDemo("ctxcancel", demoContextCancel, reader)
		case "5", "ctxtimeout":
			runDemo("ctxtimeout", demoContextTimeout, reader)
		case "6", "goexit":
			runDemo("goexit", demoGoexit, reader)
		case "7", "panic":
			runDemo("panic", demoPanicRecover, reader)
		case "8", "after":
			runDemo("after", demoTimeAfter, reader)
		case "9", "atomic":
			runDemo("atomic", demoAtomicFlag, reader)
		case "10", "mainexit":
			runDemo("mainexit", demoMainExitStopsOthers, reader)
		case "a", "all":
			runAll(reader)
		case "q", "quit", "exit":
			fmt.Println("Пока!")
			return
		default:
			fmt.Println("Неизвестный выбор:", choice)
			waitEnter(reader)
		}
	}
}

func runDemo(name string, fn func(), reader *bufio.Reader) {
	clearScreen()
	fmt.Println("=== Демонстрация:", name, "===\n")
	fn()
	fmt.Println("\n=== Конец демонстрации:", name, "===\n")
	waitEnter(reader)
}

func runAll(reader *bufio.Reader) {
	demos := []struct {
		name string
		fn   func()
	}{
		{"condition", demoConditionExit},
		{"donechan", demoDoneChannel},
		{"rangeclose", demoRangeClose},
		{"ctxcancel", demoContextCancel},
		{"ctxtimeout", demoContextTimeout},
		{"goexit", demoGoexit},
		{"panic", demoPanicRecover},
		{"after", demoTimeAfter},
		{"atomic", demoAtomicFlag},
		{"mainexit", demoMainExitStopsOthers},
	}
	for _, d := range demos {
		runDemo(d.name, d.fn, reader)
	}
}

func waitEnter(reader *bufio.Reader) {
	fmt.Print("Нажмите Enter, чтобы продолжить…")
	reader.ReadString('\n')
}

func clearScreen() {
	// Простой «очиститель» консоли (визуальный). Работает как отступ.
	fmt.Print("\033[2J\033[H")
}

// 1) Обычный выход по внутреннему условию (без внешних сигналов)
func demoConditionExit() {
	fmt.Println("[condition] старт")
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 1; i <= 5; i++ {
			fmt.Println("[condition] work step", i)
			time.Sleep(120 * time.Millisecond)
		}
		fmt.Println("[condition] достигнут лимит шагов → return")
	}()
	<-done
	fmt.Println("[condition] завершено")
}

// 2) Канал уведомления done (закрытие канала)
func demoDoneChannel() {
	fmt.Println("[donechan] старт")
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("[donechan] получен сигнал завершения → return")
				return
			default:
				fmt.Println("[donechan] работаю…")
				time.Sleep(120 * time.Millisecond)
			}
		}
	}()

	time.Sleep(400 * time.Millisecond)
	close(done)
	time.Sleep(150 * time.Millisecond)
	fmt.Println("[donechan] завершено")
}

// 3) Выход по закрытию рабочего канала (range заканчивается)
func demoRangeClose() {
	fmt.Println("[rangeclose] старт")
	jobs := make(chan int)

	go func() {
		for j := range jobs {
			fmt.Println("[rangeclose] got job", j)
			time.Sleep(80 * time.Millisecond)
		}
		fmt.Println("[rangeclose] канал jobs закрыт и опустошён → return")
	}()

	for i := 1; i <= 4; i++ {
		jobs <- i
	}
	close(jobs)
	time.Sleep(300 * time.Millisecond)
	fmt.Println("[rangeclose] завершено")
}

// 4) Контекст: ручная отмена
func demoContextCancel() {
	fmt.Println("[ctxcancel] старт")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("[ctxcancel] ctx.Done():", ctx.Err(), "→ return")
				return
			default:
				fmt.Println("[ctxcancel] работаю…")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	time.Sleep(350 * time.Millisecond)
	cancel()
	time.Sleep(150 * time.Millisecond)
	fmt.Println("[ctxcancel] завершено")
}

// 5) Контекст: таймаут
func demoContextTimeout() {
	fmt.Println("[ctxtimeout] старт")
	ctx, cancel := context.WithTimeout(context.Background(), 380*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("[ctxtimeout] ctx.Done():", ctx.Err(), "→ return")
				return
			default:
				fmt.Println("[ctxtimeout] работаю…")
				time.Sleep(110 * time.Millisecond)
			}
		}
	}()

	<-done
	fmt.Println("[ctxtimeout] завершено")
}

// 6) Жёсткое завершение текущей горутины: runtime.Goexit()
func demoGoexit() {
	fmt.Println("[goexit] старт")
	done := make(chan struct{})
	go func() {
		defer func() {
			fmt.Println("[goexit] defer выполнен перед завершением")
			close(done)
		}()
		fmt.Println("[goexit] вызываю runtime.Goexit()")
		runtime.Goexit() // завершит ТЕКУЩУЮ горутину, выполнив defer
		// сюда выполнение не дойдёт
	}()
	<-done
	fmt.Println("[goexit] завершено")
}

// 7) Останов через panic в горутине (с recover в этой же горутине)
func demoPanicRecover() {
	fmt.Println("[panic] старт")
	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("[panic] поймано:", r, "→ завершаемся корректно")
			}
			close(done)
		}()
		fmt.Println("[panic] имитируем аварию через 200ms…")
		time.Sleep(200 * time.Millisecond)
		panic("что-то пошло не так")
	}()
	<-done
	fmt.Println("[panic] завершено")
}

// 8) Таймаут через time.After в select
func demoTimeAfter() {
	fmt.Println("[after] старт")
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-time.After(350 * time.Millisecond):
				fmt.Println("[after] истёк локальный таймер → return")
				return
			default:
				fmt.Println("[after] работаю…")
				time.Sleep(90 * time.Millisecond)
			}
		}
	}()
	<-done
	fmt.Println("[after] завершено")
}

// 9) Атомарный флаг остановки
func demoAtomicFlag() {
	fmt.Println("[atomic] старт")
	var stop int32 // 0 — работать, 1 — остановиться

	go func() {
		for {
			if atomic.LoadInt32(&stop) == 1 {
				fmt.Println("[atomic] получен флаг stop=1 → return")
				return
			}
			fmt.Println("[atomic] работаю…")
			time.Sleep(120 * time.Millisecond)
		}
	}()

	time.Sleep(400 * time.Millisecond)
	atomic.StoreInt32(&stop, 1)
	time.Sleep(150 * time.Millisecond)
	fmt.Println("[atomic] завершено")
}

// 10) Выход main прерывает прочие горутины (анти-паттерн)
func demoMainExitStopsOthers() {
	fmt.Println("[mainexit] старт (анти-паттерн)")
	go func() {
		for i := 0; ; i++ {
			fmt.Println("[mainexit] рабочая горутина:", i)
			time.Sleep(120 * time.Millisecond)
		}
	}()
	time.Sleep(350 * time.Millisecond)
	fmt.Println("[mainexit] main завершается прямо сейчас → процесс остановит все горутины")
	// main возвращается — программа завершится, а горутина будет прервана рантайм
}
