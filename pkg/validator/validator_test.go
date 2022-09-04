package validator_test

import (
	"authstore/pkg/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	testCases := []struct {
		name    string
		field   func() *string
		wantErr bool
	}{
		{
			name: "nil field",
			field: func() *string {
				return nil
			},
			wantErr: true,
		},
		{
			name: "string field",
			field: func() *string {
				test := "test"
				return &test
			},
			wantErr: false,
		},
		{
			name: "string zero value field",
			field: func() *string {
				test := ""
				return &test
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := validator.Required(tc.field())
			err := fn("test")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	testCases := []struct {
		name      string
		field     func() *string
		minLength uint
		wantErr   bool
	}{
		{
			name: "nil field",
			field: func() *string {
				return nil
			},
			minLength: 6,
			wantErr:   true,
		},
		{
			name: "4 symbols string",
			field: func() *string {
				test := "test"
				return &test
			},
			minLength: 6,
			wantErr:   true,
		},
		{
			name: "string zero value field",
			field: func() *string {
				test := ""
				return &test
			},
			minLength: 6,
			wantErr:   true,
		},
		{
			name: "7 symbols string",
			field: func() *string {
				test := "ttuiop["
				return &test
			},
			minLength: 6,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := validator.MinLength(tc.field(), tc.minLength)
			err := fn("test")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	testCases := []struct {
		name      string
		field     func() *string
		maxLength uint
		wantErr   bool
	}{
		{
			name: "nil field",
			field: func() *string {
				return nil
			},
			maxLength: 6,
			wantErr:   false,
		},
		{
			name: "4 symbols string",
			field: func() *string {
				test := "test"
				return &test
			},
			maxLength: 6,
			wantErr:   false,
		},
		{
			name: "string zero value field",
			field: func() *string {
				test := ""
				return &test
			},
			maxLength: 6,
			wantErr:   false,
		},
		{
			name: "7 symbols string",
			field: func() *string {
				test := "ttuiop["
				return &test
			},
			maxLength: 6,
			wantErr:   true,
		},
		{
			name: "6 symbols string",
			field: func() *string {
				test := "ttuiop"
				return &test
			},
			maxLength: 6,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := validator.MaxLength(tc.field(), tc.maxLength)
			err := fn("test")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithoutSymbols(t *testing.T) {
	testCases := []struct {
		name    string
		field   func() *string
		symbols []rune
		wantErr bool
	}{
		{
			name: "nil field",
			field: func() *string {
				return nil
			},
			// symbols: []rune{'+', '-', ')', '(', '=', '_', '\'', '"', '/', '?', '.', ',', '|', '\\', ':', ';', '>', '<'},
			symbols: []rune("+-()_='\";:[]{}\\|/.,?><`*&^%$#@!~"),
			wantErr: false,
		},
		{
			name: "with error +",
			field: func() *string {
				str := "fuck the police +"
				return &str
			},
			symbols: []rune("+-()_='\";:[]{}\\|/.,?><`*&^%$#@!~"),
			wantErr: true,
		},
		{
			name: "no error",
			field: func() *string {
				str := "fuck the police"
				return &str
			},
			symbols: []rune("+-()_='\";:[]{}\\|/.,?><`*&^%$#@!~"),
			wantErr: false,
		},
		{
			name: "with error \\",
			field: func() *string {
				str := "fuck the police \\"
				return &str
			},
			symbols: []rune("+-()_='\";:[]{}\\|/.,?><`*&^%$#@!~"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := validator.WithoutSymbols(tc.field(), tc.symbols...)
			err := fn("test")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExistSymbolInString(t *testing.T) {
	testCases := []struct {
		name   string
		symbol rune
		str    string
		want   bool
	}{
		{
			name:   "symbol exist",
			symbol: 'u',
			str:    "fuck the police",
			want:   true,
		},
		{
			name:   "symbol not exist",
			symbol: 'm',
			str:    "fuck the police",
			want:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isExist := validator.ExistSymbolInString(tc.symbol, tc.str)
			assert.Equal(t, tc.want, isExist)
		})
	}
}

func TestExistSymbolsInString(t *testing.T) {
	testCases := []struct {
		name    string
		symbols []rune
		str     string
		want    bool
	}{
		{
			name:    "symbols exist",
			symbols: []rune{'s', 'l'},
			str:     "fuck the police",
			want:    true,
		},
		{
			name:    "symbols not exist",
			symbols: []rune{'m', 'b', '+', '2'},
			str:     "fuck the police",
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isExist := validator.ExistSymbolsInString(tc.symbols, tc.str)
			assert.Equal(t, tc.want, isExist)
		})
	}
}

func TestJoinRunes(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []rune
		separator string
		want      string
	}{
		{
			name:      "with comma separator",
			slice:     []rune{'s', 'l'},
			separator: ", ",
			want:      "s, l",
		},
		{
			name:      "with complex separator",
			slice:     []rune{'s', 'l', 'p'},
			separator: " f+-)(",
			want:      "s f+-)(l f+-)(p",
		},
		{
			name:      "with one elem",
			slice:     []rune{'s'},
			separator: " f+-)(",
			want:      "s",
		},
		{
			name:      "with zero elem",
			slice:     []rune{},
			separator: " f+-)(",
			want:      "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.JoinRunes(tc.slice, tc.separator)
			assert.Equal(t, tc.want, result)
		})
	}
}
