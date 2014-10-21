package main

import (
	"log"
	"sync"
)

type Worker struct {
	feeds   chan *FeedConfig // канал для заданий
	pending int              // кол-во оставшихся задач
	index   int              // позиция в куче
	wg      *sync.WaitGroup  //указатель на группу ожидания
}

func (w *Worker) work(done chan *Worker) {
	for {
		feed := <-w.feeds //читаем следующее задание
		w.wg.Add(1)       //инкриминируем счетчик группы ожидания
		download(feed)    //загружаем файл
		w.wg.Done()       //сигнализируем группе ожидания что закончили
		done <- w         //показываем что завершили работу
	}
}

//Загрузка изображения
func download(feed *FeedConfig) {
	log.Printf("[DEBUG]: WORKER Processing feed (Regexp here) [%dbytes]", len(feed.Html))
}
