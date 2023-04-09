package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
)

type MockHTTPClient struct{}

func (c MockHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(LETTERS)),
	}, nil
}

func TestCount(t *testing.T) {
	testCases := []struct {
		name     string
		file     []byte
		expected [26]uint32
		mode     int
	}{
		{
			name: "test with lowercase letters only",
			file: []byte(LETTERS),
			expected: [26]uint32{
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			},
			mode: ATOMIC,
		},
		{
			name: "test with uppercase letters only",
			file: []byte(LETTERS),
			expected: [26]uint32{
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			},
			mode: MUTEX,
		},
		{
			name: "test with mixed case letters",
			file: []byte(LETTERS),
			expected: [26]uint32{
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			},
			mode: ATOMIC,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := new(Score)
			var wg sync.WaitGroup
			wg.Add(1)
			count(bytes.ToLower(tc.file), score, tc.mode, &wg)
			wg.Wait()
			if score.counts != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, score.counts)
			}
		})
	}
}

func BenchmarkCountAtomic(b *testing.B) {
	score := new(Score)
	var wg sync.WaitGroup
	files := getFiles(new(MockHTTPClient))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, f := range files {
			wg.Add(1)
			go count(bytes.ToLower(f), score, ATOMIC, &wg)
		}
		wg.Wait()
	}
}

func BenchmarkCountMutex(b *testing.B) {
	score := new(Score)
	var wg sync.WaitGroup
	files := getFiles(new(MockHTTPClient))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, f := range files {
			wg.Add(1)
			go count(bytes.ToLower(f), score, MUTEX, &wg)
		}
		wg.Wait()
	}
}
