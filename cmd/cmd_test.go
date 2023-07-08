package cmd

import (
	"reflect"
	"testing"
)

func Test_getEnv(t *testing.T) {
	type args struct {
		in     Input
		curEnv []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "no env specified",
			want: nil,
		},
		{
			name: "keep current env intact",
			args: args{
				curEnv: []string{"bar=baz"},
			},
			want: []string{"bar=baz"},
		},
		{
			name: "new env for empty current",
			args: args{
				in:     Input{Env: map[string]string{"foo": "bar"}},
				curEnv: nil,
			},
			want: []string{"foo=bar"},
		},
		{
			name: "new env for non-empty current",
			args: args{
				in:     Input{Env: map[string]string{"foo": "bar"}},
				curEnv: []string{"baz=qux"},
			},
			want: []string{"baz=qux", "foo=bar"},
		},
		{
			name: "Override env",
			args: args{
				in:     Input{Env: map[string]string{"foo": "bar"}},
				curEnv: []string{"foo=baz"},
			},
			want: []string{"foo=bar"},
		},
		{
			name: "Handle env with equal signs",
			args: args{
				curEnv: []string{"foo=baz=foo"},
				in:     Input{Env: map[string]string{"bar": "qux"}},
			},
			want: []string{"foo=baz=foo", "bar=qux"},
		},
		{
			name: "Set path",
			args: args{
				in:     Input{Env: map[string]string{"PATH": "foo:bar:baz"}},
				curEnv: []string{},
			},
			want: []string{"PATH=foo:bar:baz"},
		},
		{
			name: "Append to path",
			args: args{
				in:     Input{Env: map[string]string{"PATH": "foo:bar:baz"}},
				curEnv: []string{"PATH=qux"},
			},
			want: []string{"PATH=qux:foo:bar:baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEnv(tt.args.in, tt.args.curEnv)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
