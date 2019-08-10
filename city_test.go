package ipdb

import (
	"errors"
	"fmt"
	"testing"
)

var db *City
var errDB = errors.New("NIL DB")

func init() {
	db, _ = NewCity(FILENAME_DATABASE)
}

func TestNewCity(t *testing.T) {
	if db == nil {
		t.Fatal(errDB)
	}
	// testing
	for i, ts := range tests {
		t.Run(fmt.Sprintf("Example %d", i+1), func(t *testing.T) {
			res, errRes := db.FindInfo(ts.ip, db.Languages()[0])
			if errRes != nil {
				t.Fatal(res)
			}
		})
	}
}

func BenchmarkCity_Find(b *testing.B) {
	if db == nil {
		b.Fatal(errDB)
	}
	for i := 0; i < b.N; i++ {
		db.Find(STRING_BENCH_IP, "CN")
	}
}

func BenchmarkCity_FindMap(b *testing.B) {
	if db == nil {
		b.Fatal(errDB)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.FindMap(STRING_BENCH_IP, "CN")
	}
}

func BenchmarkCity_FindInfo(b *testing.B) {
	if db == nil {
		b.Fatal(errDB)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.FindInfo(STRING_BENCH_IP, "CN")
	}
}
