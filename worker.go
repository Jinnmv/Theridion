package main

import (
	"log"
	"sync"
)

//TODO: rewrite

var IMGDIR = "./img"

type Worker struct {
	feeds   chan *FeedData  // канал для заданий
	pending int             // кол-во оставшихся задач
	index   int             // позиция в куче
	wg      *sync.WaitGroup //указатель на группу ожидания
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
func download(feed *FeedData) {
	log.Println("[DEBUG]: WORKER Processing feed (Regexp here)")
}
