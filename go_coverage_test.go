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
		{"testtesttesttesttesttesttesttesttesttesttestbest", 20, "...testtesttesttestbest"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := trimString(tt.in, tt.trim); got != tt.want {
				t.Errorf("trimString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTrimmedFileName(t *testing.T) {
	type args struct {
		fn   string
		trim bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"trim enabled", args{fn: "test_file_name_test_file_name", trim: true}, "..._name_test_file_name"},
		{"trim enabled", args{fn: "test_file_name_test_file_name", trim: false}, "test_file_name_test_file_name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTrimmedFileName(tt.args.fn, tt.args.trim); got != tt.want {
				t.Errorf("getTrimmedFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
