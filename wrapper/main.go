package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Есть функция, работающая неопределённо долго и возвращающая число.
// Её тело нельзя изменять (представим, что внутри сетевой запрос).
func unpredictableFunc() int64 {
	rnd := rand.Int63n(5000)
	time.Sleep(time.Duration(rnd) * time.Millisecond)
	return rnd
}

// Нужно изменить функцию обёртку, которая будет работать с заданным таймаутом (например, 1 секунду).
// Если "длинная" функция отработала за это время - отлично, возвращаем результат.
// Если нет - возвращаем ошибку. Результат работы в этом случае нам не важен.
//
// Дополнительно нужно измерить, сколько выполнялась эта функция (просто вывести в лог).
// Сигнатуру функцию обёртки менять можно.
func predictableFunc(timeout time.Duration, fn func() int64) (int64, error) {
	start := time.Now()
	defer func() { fmt.Println("func took:", time.Since(start)) }()

	resCh := make(chan int64, 1)
	go func() {
		resCh <- fn()
	}()

	cancelCh := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		cancelCh <- struct{}{}
	}()

	select {
	case res := <-resCh:
		return res, nil
	case <-cancelCh:
		return 0, fmt.Errorf("timeout exceeded")
	}
}

func main() {
	fmt.Println("started")
	res, err := predictableFunc(time.Millisecond*2500, unpredictableFunc)
	fmt.Println(res, err)
}
