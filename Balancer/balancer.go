package Balancer

import (
	"container/heap"
	"log"
	"sync"
)

// Balancer
type Balancer struct {
	pool         Pool             // Workers heap
	done         chan *Worker     //Канал уведомления о завершении для рабочих
	requests     chan interface{} //Канал для получения новых заданий
	flowctrl     chan bool        //Канал для PMFC
	queue        int              //Количество незавершенных заданий переданных рабочим
	wg           *sync.WaitGroup  //Группа ожидания для рабочих
	workersCap   byte
	workersCount byte
	fn           JobFunc
}

func NewBalancer() *Balancer {
	return &Balancer{}
}

//Инициализируем балансировщик. Аргументом получаем канал по которому приходят задания
func (b *Balancer) Init(in chan interface{}, fn JobFunc, workersCount, workersCap byte) {
	b.requests = make(chan interface{})
	b.flowctrl = make(chan bool)
	b.done = make(chan *Worker)
	b.wg = new(sync.WaitGroup)
	b.workersCount = workersCount
	b.workersCap = workersCap
	b.fn = fn

	//Запускаем наш Flow Control:
	go func() {
		for {
			b.requests <- <-in //получаем новое задание и пересылаем его на внутренний канал
			<-b.flowctrl       //а потом ждем получения подтверждения
		}
	}()

	//Инициализируем кучу и создаем рабочих:
	heap.Init(&b.pool)
	for i := 0; i < int(b.workersCount); i++ {
		w := &Worker{
			jobs:    make(chan interface{}, b.workersCap),
			index:   0,
			pending: 0,
			id:      i,
			wg:      b.wg,
		}
		//go w.work(b.done)     //запускаем рабочего
		go w.work(b.done, b.fn) //запускаем рабочего
		heap.Push(&b.pool, w)   //и заталкиваем его в кучу
	}
}

//Рабочая функция балансировщика получает аргументом канал уведомлений от главного цикла
func (b *Balancer) Balance(quit chan bool) {
	lastjob := false //Флаг завершения, поднимаем когда кончились задания
	for {
		select { //В цикле ожидаем коммуникации по каналам:

		case <-quit: //пришло указание на остановку работы
			b.wg.Wait()  //ждем завершения текущих загрузок рабочими..
			quit <- true //..и отправляем сигнал что закончили

		case job, ok := <-b.requests: //Получено новое задание (от flow controller)
			if !ok || job == nil { //Проверяем - а не кодовая ли это фраза?
				log.Printf("[DEBUG]: BALANCER End of inputs. queue: %d\n", b.queue)
				lastjob = true // если да, поднимаем флаг завершения
			} else {
				b.dispatch(job) //иначе то отправляем рабочим
			}

		case w := <-b.done: // пришло уведомление, что рабочий закончил загрузку
			log.Printf("[DEBUG]: BALANCER Worker #%d has completed a task, pool size: %d, pool queue: %d\n", w.id, len(b.pool), b.queue)
			b.completed(w) //обновляем его данные
			if lastjob {
				log.Println("[DEBUG]: BALANCER Finalization started")
				if w.pending == 0 { //если у рабочего кончились задания..
					log.Printf("[DEBUG]: BALANCER Worker #%d has completed ALL tasks - removing, pool size: %d", w.id, len(b.pool))
					heap.Remove(&b.pool, w.index) //то удаляем его из кучи
				}

				b.flush()

				if len(b.pool) == 0 { // All workers has finished their tasks
					quit <- true // sending ready to quit
					log.Println("[DEBUG]: BALANCER pool is empty - sending quit message")
				}
			}
		}
	}
}

// Функция отправки задания
func (b *Balancer) dispatch(job interface{}) {
	w := heap.Pop(&b.pool).(*Worker) //Берем из кучи самого незагруженного рабочего..
	w.jobs <- job                    //..и отправляем ему задание.
	w.pending++                      //Добавляем ему "весу"..
	heap.Push(&b.pool, w)            //..и отправляем назад в кучу

	if b.queue++; b.queue < int(b.workersCount*b.workersCap) {
		b.flowctrl <- true
	}
}

// Task completed handler
func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)

	if b.queue--; b.queue == int(b.workersCount*b.workersCap-1) {
		b.flowctrl <- true
	}
}

// Remove Empty workers and return count of removed workers
func (b *Balancer) flush() (removed uint) {
	for _, w := range b.pool {
		if w.pending == 0 {
			heap.Remove(&b.pool, w.index)
			removed++
		}
	}
	return removed
}
