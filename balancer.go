package main

import (
	"container/heap"
	"log"
	"sync"
)

//Балансировщик
type Balancer struct {
	pool     Pool            //Наша "куча" рабочих
	done     chan *Worker    //Канал уведомления о завершении для рабочих
	requests chan *FeedData  //Канал для получения новых заданий
	flowctrl chan bool       //Канал для PMFC
	queue    int             //Количество незавершенных заданий переданных рабочим
	wg       *sync.WaitGroup //Группа ожидания для рабочих
}

// TODO: rewrite
var (
	WORKERS    = 5 //количество рабочих
	WORKERSCAP = 3 //размер очереди каждого рабочего
	//ENDMESSAGE = "basta"
)

//Инициализируем балансировщик. Аргументом получаем канал по которому приходят задания
func (b *Balancer) init(in chan *FeedData) {
	b.requests = make(chan *FeedData)
	b.flowctrl = make(chan bool)
	b.done = make(chan *Worker)
	b.wg = new(sync.WaitGroup)

	//Запускаем наш Flow Control:
	go func() {
		for {
			b.requests <- <-in //получаем новое задание и пересылаем его на внутренний канал
			<-b.flowctrl       //а потом ждем получения подтверждения
		}
	}()

	//Инициализируем кучу и создаем рабочих:
	heap.Init(&b.pool)
	for i := 0; i < WORKERS; i++ {
		w := &Worker{
			feeds:   make(chan *FeedData, WORKERSCAP),
			index:   0,
			pending: 0,
			wg:      b.wg,
		}
		go w.work(b.done)     //запускаем рабочего
		heap.Push(&b.pool, w) //и заталкиваем его в кучу
	}
}

//Рабочая функция балансировщика получает аргументом канал уведомлений от главного цикла
func (b *Balancer) balance(quit chan bool) {
	lastjobs := false //Флаг завершения, поднимаем когда кончились задания
	for {
		select { //В цикле ожидаем коммуникации по каналам:

		case <-quit: //пришло указание на остановку работы
			b.wg.Wait()  //ждем завершения текущих загрузок рабочими..
			quit <- true //..и отправляем сигнал что закончили

		case feed, ok := <-b.requests: //Получено новое задание (от flow controller)
			if !ok || feed == nil { //Проверяем - а не кодовая ли это фраза?
				log.Printf("[DEBUG]: BALANCER End of inputs pool size: %d, pool queue: %d\n", len(b.pool), b.queue)
				if b.queue == 0 { // TODO: Refactor
					b.finalize(quit)
				}
				lastjobs = true // если да, поднимаем флаг завершения
			} else {
				log.Println("[DEBUG]: BALANCER New job received to Balancer", feed.FeedConfig.Url)
				b.dispatch(feed) //иначе то отправляем рабочим
			}

		case w := <-b.done: //пришло уведомление, что рабочий закончил загрузку
			log.Printf("[DEBUG]: BALANCER Worker #%d has completed a task, pool size: %d, pool queue: %d\n", w.index, len(b.pool), b.queue)
			b.completed(w) //обновляем его данные
			if lastjobs {
				log.Println("[DEBUG]: BALANCER Finalization started")
				if w.pending == 0 { //если у рабочего кончились задания..
					heap.Remove(&b.pool, w.index) //то удаляем его из кучи
					log.Println("[DEBUG]: BALANCER Worker #%d has completed ALL tasks - removing, pool size:", w.index, len(b.pool))
				}
				if len(b.pool) == 0 { //а если куча стала пуста
					//значит все рабочие закончили свои очереди
					quit <- true //и можно отправлять сигнал подтверждения готовности к останову
					log.Println("[DEBUG]: BALANCER pool is empty - sending quit message")
				}
			}
		}
	}
}

func (b *Balancer) finalize(quit chan bool) {
	for _, worker := range b.pool {
		if worker.pending == 0 { //если у рабочего кончились задания..
			heap.Remove(&b.pool, worker.index) //то удаляем его из кучи
			log.Printf("[DEBUG]: BALANCER Worker #%d has completed ALL tasks - removing, pool size: %d", worker.index, len(b.pool))
		}
		if len(b.pool) == 0 { //а если куча стала пуста
			//значит все рабочие закончили свои очереди
			quit <- true //и можно отправлять сигнал подтверждения готовности к останову
			log.Println("[DEBUG]: BALANCER pool is empty - sending quit message")
		}
	}
}

// Функция отправки задания
func (b *Balancer) dispatch(feed *FeedData) {
	w := heap.Pop(&b.pool).(*Worker) //Берем из кучи самого незагруженного рабочего..
	w.feeds <- feed                  //..и отправляем ему задание.
	w.pending++                      //Добавляем ему "весу"..
	heap.Push(&b.pool, w)            //..и отправляем назад в кучу
	if b.queue++; b.queue < WORKERS*WORKERSCAP {
		b.flowctrl <- true
	}
}

//Обработка завершения задания
func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
	if b.queue--; b.queue == WORKERS*WORKERSCAP-1 {
		b.flowctrl <- true
	}
}
