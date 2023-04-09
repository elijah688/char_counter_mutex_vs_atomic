package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	FILE_COUNT = 2000
	LETTERS    = "abcdefghijklmnopqrstuvwxyz"
	URL_FORMAT = "https://www.rfc-editor.org/rfc/rfc%d.txt"
	ATOMIC     = iota
	MUTEX
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func getFile(fileNumber int, httpClient HTTPClient) []byte {
	res, err := httpClient.Get(fmt.Sprintf(URL_FORMAT, fileNumber))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	return body
}
func getFiles(httpClient HTTPClient) [][]byte {
	// we use a channel here to not need to
	// employ both a mutex and wait groups

	ch := make(chan []byte)
	for i := 1; i <= FILE_COUNT; i++ {
		go func(i int) {
			ch <- getFile(i, httpClient)
		}(i)
	}
	res := make([][]byte, 0, FILE_COUNT)
	for i := 0; i < FILE_COUNT; i++ {
		res = append(res, <-ch)
	}
	return res
}

func count(file []byte, score *Score, mode int, wg *sync.WaitGroup) {

	defer wg.Done()
	for _, b := range file {
		if idx := strings.IndexByte(LETTERS, b); idx >= 0 {
			switch mode {
			case MUTEX:
				score.incrementWithMutex(idx)
			case ATOMIC:
				score.incrementWithAtomic(idx)
			}
		}
	}
}

type Score struct {
	counts [26]uint32
	mutex  sync.Mutex
}

func (s *Score) incrementWithMutex(idx int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counts[idx]++
}

func (s *Score) incrementWithAtomic(idx int) {
	atomic.AddUint32(&s.counts[idx], 1)
}

func main() {
	score := new(Score)
	wg := new(sync.WaitGroup)
	files := getFiles(new(http.Client))

	// Using a mutex variables
	start := time.Now()
	for _, f := range files {
		wg.Add(1)
		go count(bytes.ToLower(f), score, MUTEX, wg)
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("MIUTEX took %.3v seconds\n", time.Since(start).Seconds()))

	// Using atomic variables
	start = time.Now()
	for _, f := range files {
		wg.Add(1)
		go count(bytes.ToLower(f), score, ATOMIC, wg)
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("ATOMIC took %.3v seconds\n", time.Since(start).Seconds()))
}
