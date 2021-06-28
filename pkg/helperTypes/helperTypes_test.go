package helperTypes

import "testing"

func TestStringSlice_Contains(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		slice StringSlice
		args  args
		want  bool
	}{
		{
			name:  "nil slice",
			slice: nil,
			args: args{
				str: "foo",
			},
			want: false,
		},
		{
			name:  "slice contains str",
			slice: []string{"foo", "bar", "baz"},
			args: args{
				str: "foo",
			},
			want: true,
		},
		{
			name:  "slice does not contain str",
			slice: []string{"foo", "bar", "baz"},
			args: args{
				str: "qux",
			},
			want: false,
		},
		{
			name:  "empty str",
			slice: []string{"foo", "bar", "baz"},
			args: args{
				str: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.slice.Contains(tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSlice_IndexOf(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		slice StringSlice
		args  args
		want  int
	}{
		{
			name:  "nil slice",
			slice: nil,
			args: args{
				str: "foo",
			},
			want: -1,
		},
		{
			name:  "slice contains str",
			slice: []string{"foo", "bar", "baz"},
			args: args{
				str: "bar",
			},
			want: 1,
		},
		{
			name:  "slice does not contain str",
			slice: []string{"foo", "bar", "baz"},
			args: args{
				str: "qux",
			},
			want: -1,
		},
		{
			name:  "empty str",
			slice: []string{"foo", "bar", "baz"},
			args: args{
				str: "",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.slice.IndexOf(tt.args.str); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
