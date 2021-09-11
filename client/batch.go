package client

import (
	"fmt"
)

// BatchOperation is a kind of operation to perform
type BatchOperation int

// types of operations
const (
	OpRead  BatchOperation = 0
	OpWrite BatchOperation = 1
)

// how many worker threads to use for batch operations
const (
	VaultConcurency = 5
)

// BatchOperation can perform reads or writes with concurrency
func (client *Client) BatchOperation(absolutePaths []string, op BatchOperation, secretsIn []*Secret) (secrets []*Secret, err error) {
	readQueue := make(chan string, len(absolutePaths))
	writeQueue := make(chan *Secret, len(absolutePaths))
	results := make(chan *secretOperation, len(absolutePaths))

	// load up queue for operation
	switch op {
	case OpRead:
		for _, path := range absolutePaths {
			readQueue <- path
		}
	case OpWrite:
		for _, secret := range secretsIn {
			writeQueue <- secret
		}
	default:
		return nil, fmt.Errorf("invalid batch operation")
	}

	// fire off goroutines for operation
	for i := 0; i < VaultConcurency; i++ {
		client.waitGroup.Add(1)
		switch op {
		case OpRead:
			go client.readWorker(readQueue, results)
		case OpWrite:
			go client.writeWorker(writeQueue, results)
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
