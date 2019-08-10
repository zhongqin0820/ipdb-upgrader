package ipdb

import (
	"errors"
	"fmt"
	"testing"
)

const STRING_DBNAME = "utils/ipipfree.ipdb"
const STRING_BENCH_IP = "118.28.1.1"

var upgrader *Upgrader
var errUpgrader = errors.New("NIL Upgrader!")

func init() {
	upgrader, _ = NewUpgrader(STRING_DBNAME)
}

// tests driven table
var tests = []struct {
	ip     string
	output error
}{
	// MainLand
	{"124.21.1.6", ERR_NOT_FOUND_MAINLAND_REGION}, // Mainland ip without details
	{"1.2.1.0", ERR_NOT_FOUND_MAINLAND_CITY},      // Mainland ip without details
	{"39.135.129.64", nil},                        // Mainland ip with details
	// HK/TW/MO
	{"61.244.148.166", nil}, // HK ip
	{"23.36.143.47", nil},   // MO ip
	{"114.44.227.87", nil},  // TW ip
	// Foreign
	{"27.116.59.8", nil},   // Foreign ip
	{"57.71.47.25", nil},   // Foreign ip
	{"57.70.191.25", nil},  // Foreign ip
	{"35.248.7.15", nil},   // Foreign ip
	{"38.76.87.25", nil},   // Foreign ip
	{"5.22.191.25", nil},   // Foreign ip
	{"37.0.71.25", nil},    // Foreign ip
	{"57.84.143.25", nil},  // Foreign ip
	{"31.209.135.25", nil}, // Foreign ip
	{"43.245.59.25", nil},  // Foreign ip
	{"14.128.7.25", nil},   // Foreign ip
	{"57.71.47.25", nil},   // Foreign ip
	{"23.208.167.25", nil}, // Foreign ip
	{"31.15.119.2", nil},   // Foreign ip
	// Others
	{"1.2.4.8", ERR_NOT_FOUND_COUNTRY},   // Other ip
	{"127.0.0.1", ERR_NOT_FOUND_COUNTRY}, // Other ip: localhost

}

func TestUpgrader(t *testing.T) {
	if upgrader == nil {
		t.Fatal(errUpgrader)
	}
	// testing
	for i, ts := range tests {
		t.Run(fmt.Sprintf("Example %d", i+1), func(t *testing.T) {
			res, errRes := upgrader.FindCityInfo(ts.ip)
			if errRes != ts.output {
				t.Log(errRes)
				t.Fatal(res)
			}
		})
	}
}

func BenchmarkFindInfo(b *testing.B) {
	if upgrader == nil {
		b.Fatal(errUpgrader)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		upgrader.FindCityInfo(STRING_BENCH_IP)
	}
}
