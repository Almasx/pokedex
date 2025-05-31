package pokecache

import (
	"testing"
	"time"
)

type testCase struct {
	cacheInterval time.Duration
	sleepDuration time.Duration
}

func TestCache(t *testing.T) {
	cases := []struct {
        input    testCase
        expected bool
    }{
        {
            input:    testCase{
				cacheInterval: time.Second * 2,
				sleepDuration: time.Second * 1,
				
			},
			expected: true,
        },
		{
			input:    testCase{
				cacheInterval: time.Second * 1,
				sleepDuration: time.Second * 2,
			},
			expected: false,
		},
    }

	for _, c := range cases {
		cache := NewCache(c.input.cacheInterval)
		cache.Add("test", []byte("test"))
		time.Sleep(c.input.sleepDuration)
		data, ok := cache.Get("test")
		if ok != c.expected {
			t.Errorf("expected data to be found: %v, got %v", c.expected, ok)
		}	
		if ok && string(data) != "test" {
			t.Errorf("expected data to be 'test', got %v", data)
		}
	}

}