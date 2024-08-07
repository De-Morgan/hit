package main

import (
	"io"
	"testing"
)

type parseArgsTest struct {
	name    string
	args    []string
	want    config
	wantErr bool
}

func Test_parseArgs(t *testing.T) {
	t.Parallel()
	tests := []parseArgsTest{
		{
			name: "all_flags",
			args: []string{"-n=10", "-c=5", "-m=GET", "-rps=5",
				"http://test"},
			want: config{n: 10, c: 5, rps: 5,
				method: GET,
				url:    "http://test"},
		},
		{
			name: "only_n",
			args: []string{
				"-n=4", "https://go.dev",
			},
			want: config{
				url: "https://go.dev",
				n:   4,
			},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got config
			err := parseArgs(&got, tt.args, io.Discard)
			if (err == nil) && tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
			} else if got != tt.want {
				t.Errorf("flags = %+v, want %+v,", got, tt.want)
			}

		})

	}

}
