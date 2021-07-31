package main

import (
	"os"
	"reflect"
	"testing"
)

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

func Test_fmtFuncInfo(t *testing.T) {
	type args struct {
		x       *funcInfo
		covered int64
		total int64
		trim    bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"returns function details without trim when trim is false",
			args{
				x: &funcInfo{
					fileName: "test",
					functionName: "test_func_name_test_func_name_test_func_name",
					functionStartLine: 10,
					functionEndLine:   20,
					uncoveredLines:    0},
				covered: 50,
				total: 50,
				trim:    false},
			[]string{"test", "test_func_name_test_func_name_test_func_name", "10", "20", "0", "0.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmtFuncInfo(tt.args.x, tt.args.covered, tt.args.total, tt.args.trim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fmtFuncInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateCoverage(t *testing.T) {
	type args struct {
		covered int64
		total   int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"fully covered", args{100, 100}, 100.0},
		{"partially covered", args{50, 100}, 50.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateCoverage(tt.args.covered, tt.args.total); got != tt.want {
				t.Errorf("calculateCoverage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findFile(t *testing.T) {
	path, _ := os.Getwd()

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"file from this package", args{file: "github.com/gojek/go-coverage/go_coverage.go"}, path + "/go_coverage.go", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("findFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}