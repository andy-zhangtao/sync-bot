package utils

import "testing"

func TestGrabCommand(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "keep content",
			args: args{msg: "/docker-name vikings"},
			want: "vikings",
		},
		{
			name: "only has command",
			args: args{msg: "/docker-name"},
			want: "/docker-name",
		},
		{
			name: "has some multiple spaces",
			args: args{msg: "/docker-name  "},
			want: "/docker-name  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GrabCommand(tt.args.msg); got != tt.want {
				t.Errorf("GrabCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
