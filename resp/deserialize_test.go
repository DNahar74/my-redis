package resp

import (
	"reflect"
	"testing"
)

func TestDeserializeSimpleString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Type
		wantErr bool
	}{
		{
			name:    "Valid SimpleString",
			input:   "+hello\r\n",
			want:    SimpleString{Value: "hello"},
			wantErr: false,
		},
		{
			name:    "Empty Input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "No CRLF",
			input:   "+hello",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deserialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeSimpleString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeSimpleString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializeSimpleError(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Type
		wantErr bool
	}{
		{
			name:    "Valid SimpleError",
			input:   "-Error message\r\n",
			want:    SimpleError{Value: "Error message"},
			wantErr: false,
		},
		{
			name:    "Empty Input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "No CRLF",
			input:   "-Error message",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deserialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeSimpleError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeSimpleError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializeInteger(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Type
		wantErr bool
	}{
		{
			name:    "Valid Integer",
			input:   ":123\r\n",
			want:    Integer{Value: 123},
			wantErr: false,
		},
		{
			name:    "Negative Integer",
			input:   ":-456\r\n",
			want:    Integer{Value: -456},
			wantErr: false,
		},
		{
			name:    "Zero Integer",
			input:   ":0\r\n",
			want:    Integer{Value: 0},
			wantErr: false,
		},
		{
			name:    "Invalid Integer",
			input:   ":abc\r\n",
			wantErr: true,
		},
		{
			name:    "No CRLF",
			input:   ":123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deserialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializeBulkString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Type
		wantErr bool
	}{
		{
			name:    "Valid BulkString",
			input:   "$5\r\nhello\r\n",
			want:    BulkString{Value: "hello", Length: 5},
			wantErr: false,
		},
		{
			name:    "Empty BulkString",
			input:   "$0\r\n\r\n",
			want:    BulkString{Value: "", Length: 0},
			wantErr: false,
		},
		{
			name:    "Invalid Length",
			input:   "$3\r\nhello\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No CRLF after length",
			input:   "$3hello\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No CRLF after data",
			input:   "$5\r\nhello",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deserialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deserialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializeArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Array
		wantErr bool
	}{
		{
			name:  "Valid Array",
			input: "*2\r\n:123\r\n$5\r\nhello\r\n",
			want: Array{
				Length: 2,
				Items:  []Type{Integer{Value: 123}, BulkString{Value: "hello", Length: 5}},
			},
			wantErr: false,
		},
		{
			name:  "Empty Array",
			input: "*0\r\n",
			want: Array{
				Length: 0,
				Items:  []Type{},
			},
			wantErr: false,
		},
		{
			name:  "Nested Array",
			input: "*2\r\n:1\r\n*2\r\n:2\r\n$4\r\ntest\r\n",
			want: Array{
				Length: 2,
				Items: []Type{
					Integer{Value: 1},
					Array{
						Length: 2,
						Items:  []Type{Integer{Value: 2}, BulkString{Value: "test", Length: 4}},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid Array Length - incorrect length",
			input:   "*a\r\n:1\r\n", //invalid length
			wantErr: true,
		},
		{
			name:    "Invalid Array Length - too many elements",
			input:   "*2\r\n:123\r\n$5\r\nhello\r\n:456\r\n", // 3 elements, length is 2
			wantErr: true,
		},
		{
			name:    "Invalid Array Length - too few elements",
			input:   "*3\r\n:123\r\n$5\r\nhello\r\n", // 2 elements, length is 3
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Deserialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && tt.wantErr == false {
				t.Errorf("DeserializeArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
