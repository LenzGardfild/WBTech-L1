package main

import "fmt"

// Родительская структура
type Human struct {
    Name string
    Age  int
}

// Метод для Human
func (h Human) SayHello() {
    fmt.Printf("Привет, меня зовут %s, мне %d лет.\n", h.Name, h.Age)
}

// Метод для Human
func (h Human) IsAdult() bool {
    return h.Age >= 18
}

// Создаем структуру Action, которая встраивает структуру Human
type Action struct {
    Human // встроенная структура
    Role  string
}

func main() {
    // Создаём Action с вложенным Human
    a := Action{
        Human: Human{Name: "Атлухан", Age: 19},
        Role:  "Программист",
    }

    // Используем методы Human напрямую через Action
    a.SayHello()

    if a.IsAdult() {
        fmt.Printf("%s достаточно взрослый, чтобы работать как %s.\n", a.Name, a.Role)
    } else {
        fmt.Printf("%s ещё слишком молод.\n", a.Name)
    }
}
