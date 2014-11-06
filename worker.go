package main

import (
	"sync"
)

type Worker struct {
	feeds   chan *FeedConfig // канал для заданий
	pending int              // кол-во оставшихся задач
	index   int              // позиция в куче
	id      int
	wg      *sync.WaitGroup //указатель на группу ожидания
}

type WorkerFunc func(feed) 

func (w *Worker) work(done chan *Worker, fn WorkerFunc) {
	for {
		feed := <-w.feeds //читаем следующее задание
		w.wg.Add(1)       //инкриминируем счетчик группы ожидания

		fn(feed)
		
		//price := PriceList{}
		//price.Parse(feed)

		w.wg.Done() //сигнализируем группе ожидания что закончили
		done <- w   //показываем что завершили работу
	}
}

