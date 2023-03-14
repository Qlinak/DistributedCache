package geeCache

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewMap(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	// 6 -> 6, 16, 26
	// 4 -> 4, 14, 24
	// 2 -> 2, 12, 22
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if i := hash.Get(k); i != v {
			t.Errorf("Asking for %s, should have yielded %s, but got %s", k, v, i)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if i := hash.Get(k); i != v {
			t.Errorf("Asking for %s, should have yielded %s, but got %s", k, v, i)
		}
	}

}
