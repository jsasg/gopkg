package snowflake

import "testing"

func TestSnowFlake(t *testing.T) {
	sf, _ := New(0)
	for i := 0; i < 100; i++ {
		t.Log(sf.Next())
	}
}

func BenchmarkSnowflake(b *testing.B) {
	sf, _ := New(0)
	for i := 0; i < b.N; i++ {
		sf.Next()
	}
}
