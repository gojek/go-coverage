package main

import "testing"

func Test_trimString(t *testing.T) {

	tests := []struct {
		in   string
		trim int
		want string
	}{
		{"test", 20, "test"},
		{"test", 3, "...est"},
		{"testtesttesttesttesttest", 20, "...testtesttesttesttest"},
		{"testtesttesttesttesttesttesttesttesttesttestbest", 20,"...testtesttesttestbest"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := trimString(tt.in, tt.trim); got != tt.want {
				t.Errorf("trimString() = %v, want %v", got, tt.want)
			}
		})
	}
}
