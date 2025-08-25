package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{2, 4, 6, 8, 10}
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			fmt.Printf("%d^2 = %d\n", x, x*x)
		}(n)
	}

	wg.Wait()
}
