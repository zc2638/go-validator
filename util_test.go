/**
 * Created by zc on 2020/7/7.
 */
package validator

import (
	"testing"
)

func Test_buildPath(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one",
			args: args{
				paths: []string{"user"},
			},
			want: "user",
		},
		{
			name: "two",
			args: args{
				paths: []string{"user", "name"},
			},
			want: "user.name",
		},
		{
			name: "three",
			args: args{
				paths: []string{"user", "addr", "name"},
			},
			want: "user.addr.name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildPath(tt.args.paths...); got != tt.want {
				t.Errorf("buildPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildSlicePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one",
			args: args{
				path: "user.cates." + SignSlice + "0",
			},
			want: "user.cates." + SignSlice,
		},
		{
			name: "two",
			args: args{
				path: "user.cates." + SignSlice + "0.name",
			},
			want: "user.cates." + SignSlice + ".name",
		},
		{
			name: "three",
			args: args{
				path: "user.cates." + SignSlice + "1.set." + SignSlice + "0.name",
			},
			want: "user.cates." + SignSlice + ".set." + SignSlice + ".name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSlicePath(tt.args.path); got != tt.want {
				t.Errorf("buildSlicePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_camelToUnderline(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one",
			args: args{
				s: "User",
			},
			want: "user",
		},
		{
			name: "two",
			args: args{
				s: "UserCenter",
			},
			want: "user_center",
		},
		{
			name: "three",
			args: args{
				s: "UserCenterSuper",
			},
			want: "user_center_super",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelToUnderline(tt.args.s); got != tt.want {
				t.Errorf("camelToUnderline() = %v, want %v", got, tt.want)
			}
		})
	}
}
