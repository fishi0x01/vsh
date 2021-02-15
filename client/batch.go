package client

import (
	"fmt"
)

// BatchOperation is a kind of operation to perform
type BatchOperation int

// types of operations
const (
	OP_READ  BatchOperation = 0
	OP_WRITE BatchOperation = 1
)

// how many worker threads to use for batch operations
const (
	VAULT_CONCURENCY = 5
)

// BatchOperation can perform reads or writes with concurrency
func (client *Client) BatchOperation(absolutePaths []string, op BatchOperation, secretsIn []*Secret) (secrets []*Secret, err error) {
	read_queue := make(chan string, len(absolutePaths))
	write_queue := make(chan *Secret, len(absolutePaths))
	results := make(chan *secretOperation, len(absolutePaths))

	// load up queue for operation
	switch op {
	case OP_READ:
		for _, path := range absolutePaths {
			read_queue <- path
		}
	case OP_WRITE:
		for _, secret := range secretsIn {
			write_queue <- secret
		}
	default:
		return nil, fmt.Errorf("invalid batch operation")
	}

	// fire off goroutines for operation
	for i := 0; i < VAULT_CONCURENCY; i++ {
		client.waitGroup.Add(1)
		switch op {
		case OP_READ:
			go client.readWorker(read_queue, results)
		case OP_WRITE:
			go client.writeWorker(write_queue, results)
		}
	}
	client.waitGroup.Wait()
	close(results)

	// read results from the queue and return as array
	for result := range results {
		err = result.Error
		if err != nil {
			return secrets, err
		}
		if result.Result != nil {
			secrets = append(secrets, result.Result)
		}
	}
	return secrets, nil
}

// readWorker fetches paths to be read from the queue until empty
func (client *Client) readWorker(queue chan string, out chan *secretOperation) {
	defer client.waitGroup.Done()
readFromQueue:
	for {
		select {
		case path := <-queue:
			s, err := client.Read(path)
			out <- &secretOperation{Result: s, Path: path, Error: err}
		default:
			break readFromQueue
		}
	}
}

// writeWorker writes secrets to Vault in parallel
func (client *Client) writeWorker(queue chan *Secret, out chan *secretOperation) {
	defer client.waitGroup.Done()
readFromQueue:
	for {
		select {
		case secret := <-queue:
			err := client.Write(secret.Path, secret)
			out <- &secretOperation{Result: nil, Path: secret.Path, Error: err}
		default:
			break readFromQueue
		}
	}
}
