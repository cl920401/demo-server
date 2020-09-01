package worker

type Job struct {
	Data interface{}
	Proc func(interface{})
}

//var JobQueue chan Job = make(chan Job, 10)

//Worker,用来从Job队列中取出Job执行
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	Quit       chan bool
}

func (worker Worker) Stop() {
	go func() {
		worker.Quit <- true
	}()
}

type Dispatcher struct {
	MaxWorker  int
	WorkerPool chan chan Job
}

func (worker Worker) Start() {
	go func() {
		for {
			worker.WorkerPool <- worker.JobChannel
			select {
			case job := <-worker.JobChannel:
				job.Proc(job.Data)
			case quit := <-worker.Quit:
				if quit {
					return
				}
			}
		}
	}()
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
	}
}

func (dispatcher *Dispatcher) Run(JobQueue chan Job) {
	for i := 0; i < dispatcher.MaxWorker; i++ {
		worker := NewWorker(dispatcher.WorkerPool)
		worker.Start()
	}
	go dispatcher.dispatch(JobQueue)
}

func (dispatcher *Dispatcher) dispatch(JobQueue chan Job) {
	for job := range JobQueue {
		jobChannel := <-dispatcher.WorkerPool
		jobChannel <- job
	}
}

func NewDispatcher(maxWorker int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorker)
	return &Dispatcher{
		WorkerPool: workerPool,
		MaxWorker:  maxWorker,
	}
}
