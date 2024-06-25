package customvalidator

import (
	"errors"
	"testing"
)

type payload struct {
	Id int 
	Name string `validate:"required"`
	Companies []string `validate:"required"`
	Count int `validate:"required,gt=5"`
	Color string `validate:"required,oneof=red green blue"`
}

var (
	ErrSentinel = errors.New("dummy")
)

func TestValidation(t *testing.T) {
	var tests = []struct {
		name string
		input payload
		want error 
	}{
		{
			name: "Payload should validate successfully", 
			input: payload{Id: 5, Name: "john", Companies: []string{"micron"}, Count: 10, Color: "green"}, 
			want: nil,
		},
		{
			name: "Payload should validate successfully without id", 
			input: payload{Name: "john", Companies: []string{"micron"}, Count: 10, Color: "green"}, 
			want: nil,
		},
		{
			name: "Payload should fail without name and count", 
			input: payload{Companies: []string{"micron"}, Color: "red"}, 
			want: ErrSentinel, 
		},
		{
			name: "Payload should fail with invalid color", 
			input: payload{Name: "john", Companies: []string{"micron"}, Count: 10, Color: "purple"}, 
			want: ErrSentinel, 
		},
		{
			name: "Payload should fail with invalid count", 
			input: payload{Name: "john", Companies: []string{"micron"}, Count: 1, Color: "blue"}, 
			want: ErrSentinel, 
		},
	}

	for _, tt := range tests {
		cv := ProvideValidator()
		t.Run(tt.name, func(t *testing.T) {
			err := cv.Validate(tt.input)
			if err != tt.want {
				if errors.Is(tt.want, ErrSentinel) && err == nil {
					t.Error("expected error, got nil")
				}
			}
		})
	}
}