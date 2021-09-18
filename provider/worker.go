package provider

// IWorker push data into worker
type IWorker interface {
	Work() chan<- []byte
	Start()
	Stop()
	Next(IWorker)
}

type IWorkHandler interface {
	HandleData(data []byte) error
	HandleError(err error)
	Close() error
}

type IWorkerV2 interface {
	Work() chan<- interface{}
	Start()
	Stop()
	Next(IWorkerV2)
}

type IWorkHandlerV2 interface {
	Name() string
	Size() int
	HandleData(interface{}) (interface{}, error)
	Close() error
}
