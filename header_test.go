package conventionalcommitparser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHeader(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name string
		args args
		want Header
	}{
		{
			name: "commom header",
			args: args{txt: "commom header"},
			want: Header{Subject: "commom header"},
		},
		{
			name: "feat: valid header",
			args: args{txt: "feat: valid header"},
			want: Header{
				Type:    "feat",
				Subject: "valid header",
			},
		},
		{
			name: "feat: valid header with multiple spaces",
			args: args{txt: "feat:   valid header"},
			want: Header{
				Type:    "feat",
				Subject: "valid header",
			},
		},
		{
			name: "feat-with-dash: valid header with multiple spaces",
			args: args{txt: "feat-with-dash: valid header"},
			want: Header{
				Type:    "feat-with-dash",
				Subject: "valid header",
			},
		},
		{
			name: "feat with space: valid header with multiple spaces",
			args: args{txt: "feat with space: valid header"},
			want: Header{
				Type:    "feat with space",
				Subject: "valid header",
			},
		},
		{
			name: " feat with space prefix and suffix: valid header with multiple spaces",
			args: args{txt: " feat with space prefix and suffix : valid header"},
			want: Header{
				Type:    "feat with space prefix and suffix",
				Subject: "valid header",
			},
		},
		{
			name: "feat(scope): valid header",
			args: args{txt: "feat(scope): valid header"},
			want: Header{
				Type:    "feat",
				Scope:   "scope",
				Subject: "valid header",
			},
		},
		{
			name: "feat(scope-1): valid header",
			args: args{txt: "feat(scope-1): valid header"},
			want: Header{
				Type:    "feat",
				Scope:   "scope-1",
				Subject: "valid header",
			},
		},
		{
			name: "feat(scope with space): valid header",
			args: args{txt: "feat(scope with space): valid header"},
			want: Header{
				Type:    "feat",
				Scope:   "scope with space",
				Subject: "valid header",
			},
		},
		{
			name: "feat( scope with space prefix and suffix ): valid header",
			args: args{txt: "feat( scope with space prefix and suffix ): valid header"},
			want: Header{
				Type:    "feat",
				Scope:   "scope with space prefix and suffix",
				Subject: "valid header",
			},
		},
		{
			name: "feat(scope-with-dash): valid header",
			args: args{txt: "feat(scope-with-dash): valid header"},
			want: Header{
				Type:    "feat",
				Scope:   "scope-with-dash",
				Subject: "valid header",
			},
		},
		{
			name: "feat(scope)!: valid header",
			args: args{txt: "feat(scope)!: valid header"},
			want: Header{
				Type:      "feat",
				Scope:     "scope",
				Subject:   "valid header",
				Important: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseHeader(tt.args.txt); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
