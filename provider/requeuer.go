package provider

// IRequeuer providers a batch operation
// n is the length of successfully enqueued batch data
// if an error occurs, return n and error immediately
type IRequeuer interface {
	Requeue(batch [][]byte) (n int, err error)
}
