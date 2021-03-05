package main

import (
	"fmt"
	"go-pool/pool"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxGoroutines   = 10 //Total number of go routines to run
	pooledResources = 2  // total number of resources in the pool
)

type dbConnection struct {
	ID int32
}

// Implementation of the Closer interface so that we can "clean up"
// this resource

func (db *dbConnection) Close() error {
	fmt.Println("Close: Connection", db.ID)
	return nil
}

var idCounter int32

// Factory function called by Pool to create a new connection
func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	fmt.Println("--> Create: New Connection", id)
	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	p, err := pool.New(createConnection, pooledResources)

	if err != nil {
		fmt.Println(err)
	}

	for query := 0; query < maxGoroutines; query++ {
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}

	wg.Wait()

	//Close the pool down
	fmt.Println("Shutdown Program.")
	p.Close()

}

func performQueries(query int, p *pool.Pool) {
	//Acquire connection from pool

	conn, err := p.Acquire()
	if err != nil {
		fmt.Println(err)
		return
	}

	// When the function returns, release the connection back to the pool
	defer p.Release(conn)

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Printf("--> QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)
}
