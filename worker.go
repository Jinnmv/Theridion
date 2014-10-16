package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

//TODO: rewrite

var IMGDIR = "./img"

type Worker struct {
	urls    chan string     // канал для заданий
	pending int             // кол-во оставшихся задач
	index   int             // позиция в куче
	wg      *sync.WaitGroup //указатель на группу ожидания
}

func (w *Worker) work(done chan *Worker) {
	for {
		url := <-w.urls //читаем следующее задание
		w.wg.Add(1)     //инкриминируем счетчик группы ожидания
		download(url)   //загружаем файл
		w.wg.Done()     //сигнализируем группе ожидания что закончили
		done <- w       //показываем что завершили работу
	}
}

//Загрузка изображения
func download(url string) {
	fileName := IMGDIR + "/" + url[strings.LastIndex(url, "/")+1:]
	output, err := os.Create(fileName)
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	io.Copy(output, response.Body)
}
