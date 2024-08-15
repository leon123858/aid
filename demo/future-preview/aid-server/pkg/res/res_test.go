package res

import (
	"reflect"
	"testing"
)

func TestGenerateResponse(t *testing.T) {
	type args struct {
		result  bool
		content string
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		{
			name: "TestGenerateResponse",
			args: args{
				result:  true,
				content: "content",
			},
			want: Response{
				Result:  true,
				Content: "content",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateResponse(tt.args.result, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
