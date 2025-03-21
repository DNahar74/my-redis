package resp

import (
	"testing"
)

func TestSimpleStringSerialize(t *testing.T) {
	tests := []struct {
		name    string
		input   SimpleString
		want    string
		wantErr bool
	}{
		{
			name:    "Valid SimpleString",
			input:   SimpleString{Value: "hello"},
			want:    "+hello\r\n",
			wantErr: false,
		},
		{
			name:    "Empty SimpleString",
			input:   SimpleString{Value: ""},
			wantErr: true,
		},
		{
			name:    "SimpleString with \\r",
			input:   SimpleString{Value: "hello\rworld"},
			wantErr: true,
		},
		{
			name:    "SimpleString with \\n",
			input:   SimpleString{Value: "hello\nworld"},
			wantErr: true,
		},
		{
			name:    "SimpleString with spaces",
			input:   SimpleString{Value: "  hello   "},
			want:    "+hello\r\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleString.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SimpleString.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleErrorSerialize(t *testing.T) {
	tests := []struct {
		name    string
		input   SimpleError
		want    string
		wantErr bool
	}{
		{
			name:    "Valid SimpleError",
			input:   SimpleError{Value: "Error message"},
			want:    "-Error message\r\n",
			wantErr: false,
		},
		{
			name:    "Empty SimpleError",
			input:   SimpleError{Value: ""},
			wantErr: true,
		},
		{
			name:    "SimpleError with \\r",
			input:   SimpleError{Value: "Error\rMessage"},
			wantErr: true,
		},
		{
			name:    "SimpleError with \\n",
			input:   SimpleError{Value: "Error\nMessage"},
			wantErr: true,
		},
		{
			name:    "SimpleError with spaces",
			input:   SimpleError{Value: "   Error    "},
			want:    "-Error\r\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleError.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SimpleError.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegerSerialize(t *testing.T) {
	tests := []struct {
		name  string
		input Integer
		want  string
	}{
		{
			name:  "Positive Integer",
			input: Integer{Value: 123},
			want:  ":123\r\n",
		},
		{
			name:  "Negative Integer",
			input: Integer{Value: -456},
			want:  ":-456\r\n",
		},
		{
			name:  "Zero Integer",
			input: Integer{Value: 0},
			want:  ":0\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.input.Serialize()
			if got != tt.want {
				t.Errorf("Integer.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBulkStringSerialize(t *testing.T) {
	tests := []struct {
		name  string
		input BulkString
		want  string
	}{
		{
			name:  "Valid BulkString",
			input: BulkString{Value: "hello"},
			want:  "$5\r\nhello\r\n",
		},
		{
			name:  "Empty BulkString",
			input: BulkString{Value: ""},
			want:  "$0\r\n\r\n",
		},
		{
			name:  "BulkString with spaces",
			input: BulkString{Value: "  hello world  "},
			want:  "$15\r\n  hello world  \r\n",
		},
		{
			name:  "BulkString with special characters",
			input: BulkString{Value: "test\r\nstring"},
			want:  "$12\r\ntest\r\nstring\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.input.Serialize()
			if got != tt.want {
				t.Errorf("BulkString.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArraySerialize(t *testing.T) {
	tests := []struct {
		name    string
		input   Array
		want    string
		wantErr bool
	}{
		{
			name: "Valid Array",
			input: Array{
				Items: []RSPEType{
					Integer{Value: 123},
					BulkString{Value: "hello"},
				},
			},
			want: "*2\r\n:123\r\n$5\r\nhello\r\n",
		},
		{
			name: "Empty Array",
			input: Array{
				Items: []RSPEType{},
			},
			want: "*0\r\n",
		},
		{
			name: "Nested Array",
			input: Array{
				Items: []RSPEType{
					Integer{Value: 1},
					Array{Items: []RSPEType{Integer{Value: 2}, BulkString{Value: "test"}}},
				},
			},
			want: "*2\r\n:1\r\n*2\r\n:2\r\n$4\r\ntest\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Array.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Array.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
