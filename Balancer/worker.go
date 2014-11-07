package Balancer

import (
	"sync"
)

type Worker struct {
	jobs    chan interface{} // channel for jobs
	pending int              // rest jobs count
	index   int              // position in heap
	id      int              // worker ID
	wg      *sync.WaitGroup
}

type JobFunc func(interface{})

func (w *Worker) work(done chan *Worker, fn JobFunc) {
	for {
		job := <-w.jobs //читаем следующее задание
		w.wg.Add(1)     //инкриминируем счетчик группы ожидания

		fn(job)

		//price := PriceList{}
		//price.Parse(feed)

		w.wg.Done() //сигнализируем группе ожидания что закончили
		done <- w   //показываем что завершили работу
	}
}
