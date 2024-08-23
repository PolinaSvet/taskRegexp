package calcregexp

import (
	"reflect"
	"regexp"
	"testing"
)

func TestCalc_calculateLineTwoArg(t *testing.T) {
	type args struct {
		expression []string
	}
	tests := []struct {
		name    string
		c       *Calc
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "[calculateLineTwoArg]_Good Addition",
			c:    &Calc{},
			args: args{
				expression: []string{"", "-1", "", "", "+", "9", "", ""},
			},
			want:    8,
			wantErr: false,
		},
		{
			name: "[calculateLineTwoArg]_Good Subtraction",
			c:    &Calc{},
			args: args{
				expression: []string{"", "-1", "", "", "-", "9", "", ""},
			},
			want:    -10,
			wantErr: false,
		},
		{
			name: "[calculateLineTwoArg]_Good Multiplication",
			c:    &Calc{},
			args: args{
				expression: []string{"", "-1", "", "", "*", "9", "", ""},
			},
			want:    -9,
			wantErr: false,
		},
		{
			name: "[calculateLineTwoArg]_Good Division",
			c:    &Calc{},
			args: args{
				expression: []string{"", "-10", "", "", "/", "2", "", ""},
			},
			want:    -5,
			wantErr: false,
		},
		{
			name: "[calculateLineTwoArg]_Error Division by Zero",
			c:    &Calc{},
			args: args{
				expression: []string{"", "-1", "", "", "/", "0", "", ""},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "[calculateLineTwoArg]_Error Invalid Expression",
			c:    &Calc{},
			args: args{
				expression: []string{"", "5", "", "+", "3"},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "[calculateLineTwoArg]_Error Invalid Operator",
			c:    &Calc{},
			args: args{
				expression: []string{"", "-1", "", "", "%", "9", "", ""},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "[calculateLineTwoArg]_Error Empty list",
			c:    &Calc{},
			args: args{
				expression: []string{},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.calculateLineTwoArg(tt.args.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("calc.calculateLineTwoArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calc.calculateLineTwoArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalc_findAllPrior(t *testing.T) {
	type args struct {
		input string
		re    *regexp.Regexp
	}
	tests := []struct {
		name    string
		c       *Calc
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "[findAllPrior]_Good Multiplication get result with ()",
			c:    &Calc{},
			args: args{
				input: "(+2*+5)",
				re:    rePriorFirst,
			},
			want:    "(+10.000000)",
			wantErr: false,
		},
		{
			name: "[findAllPrior]_Good Multiplication get result",
			c:    &Calc{},
			args: args{
				input: "+2*-5",
				re:    rePriorFirst,
			},
			want:    "-10.000000",
			wantErr: false,
		},
		{
			name: "[findAllPrior]_Good Add not result return input",
			c:    &Calc{},
			args: args{
				input: "+2*-5",
				re:    rePriorSecond,
			},
			want:    "+2*-5",
			wantErr: false,
		},
		{
			name: "[findAllPrior]_Good Multiplication",
			c:    &Calc{},
			args: args{
				input: "(3+6533.000000*+106.000000)*+10.111109*+1.000000",
				re:    rePriorFirst,
			},
			want:    "(3+692498.000000)*+10.111109",
			wantErr: false,
		},
		{
			name: "[findAllPrior]_Error empty",
			c:    &Calc{},
			args: args{
				input: "",
				re:    rePriorSecond,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.findAllPrior(tt.args.input, tt.args.re)
			if (err != nil) != tt.wantErr {
				t.Errorf("calc.findAllPrior() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calc.findAllPrior() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalc_findAllParentheses(t *testing.T) {
	type args struct {
		input string
		re    *regexp.Regexp
	}
	tests := []struct {
		name    string
		c       *Calc
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "[findAllParentheses]_Good with result with (()*)",
			c:    &Calc{},
			args: args{
				input: "((+692501.000000)*+10.111109)",
				re:    reParenthesesStart,
			},
			want:    "(+692501.000000*+10.111109)",
			wantErr: false,
		},
		{
			name: "[findAllParentheses]_Good with result with ()*",
			c:    &Calc{},
			args: args{
				input: "(+692501.000000)*+10.111109",
				re:    reParenthesesStart,
			},
			want:    "+692501.000000*+10.111109",
			wantErr: false,
		},
		{
			name: "[findAllParentheses]_Good without result ()*",
			c:    &Calc{},
			args: args{
				input: "(+692501.000000)*+10.111109",
				re:    reParentheses,
			},
			want:    "(+692501.000000)*+10.111109",
			wantErr: false,
		},
		{
			name: "[findAllParentheses]_Good with result (+(+))*",
			c:    &Calc{},
			args: args{
				input: "(3+(+6533.000000)*(+106.000000))*(+10.111109)",
				re:    reParentheses,
			},
			want:    "(3+6533.000000*+106.000000)*+10.111109",
			wantErr: false,
		},
		{
			name: "[findAllParentheses]_Good with result (-(-))*",
			c:    &Calc{},
			args: args{
				input: "(3-(-6533.000000)*(+106.000000))*(+10.111109)",
				re:    reParentheses,
			},
			want:    "(3+6533.000000*+106.000000)*+10.111109",
			wantErr: false,
		},
		{
			name: "[findAllParentheses]_Good with result (-(+))*",
			c:    &Calc{},
			args: args{
				input: "(3-(+6533.000000)*(+106.000000))*(+10.111109)",
				re:    reParentheses,
			},
			want:    "(3-6533.000000*+106.000000)*+10.111109",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.findAllParentheses(tt.args.input, tt.args.re)
			if (err != nil) != tt.wantErr {
				t.Errorf("calc.findAllParentheses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calc.findAllParentheses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalc_calcStepByStep(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		c       *Calc
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "[calcStepByStep]_Good variant 1",
			c:    &Calc{},
			args: args{
				input: "(5+4)",
			},
			want:    "+9.000000",
			wantErr: false,
		},
		{
			name: "[calcStepByStep]_Good variant 2",
			c:    &Calc{},
			args: args{
				input: "(3+5+6+8+9)",
			},
			want:    "+31.000000",
			wantErr: false,
		},
		{
			name: "[calcStepByStep]_Good variant 3",
			c:    &Calc{},
			args: args{
				input: "(3+(5-(7+1)*-(2+2-8*9)*3*-4)*(4-(-3.0*8.0+7)*6))*(12-8/9*(12+8)+(1*8/9-5+4*5))",
			},
			want:    "+7001953.093609",
			wantErr: false,
		},
		{
			name: "[calcStepByStep]_Error variant 3",
			c:    &Calc{},
			args: args{
				input: "(((3+5+6+8+9)",
			},
			want:    "(((+31.000000)",
			wantErr: true,
		},
		{
			name: "[calcStepByStep]_Error variant 4",
			c:    &Calc{},
			args: args{
				input: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.calcStepByStep(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("calc.calcStepByStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calc.calcStepByStep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalc_calcCheckFormat(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		c       *Calc
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "[calcCheckFormat]_Good variant 1",
			c:    &Calc{},
			args: args{
				input: "2-4=?",
			},
			want:    "(2-4)",
			wantErr: false,
		},
		{
			name: "[calcCheckFormat]_Good variant 2",
			c:    &Calc{},
			args: args{
				input: "(3+5+6)=?",
			},
			want:    "(3+5+6)",
			wantErr: false,
		},
		{
			name: "[calcCheckFormat]_Error variant 1",
			c:    &Calc{},
			args: args{
				input: "(3+5+6)=",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "[calcCheckFormat]_Error variant 2",
			c:    &Calc{},
			args: args{
				input: "(3+5a+6)=",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "[calcCheckFormat]_Error variant 3",
			c:    &Calc{},
			args: args{
				input: "(3+5%6)=?",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "[calcCheckFormat]_Error variant 4",
			c:    &Calc{},
			args: args{
				input: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "[calcCheckFormat]_Error variant 5",
			c:    &Calc{},
			args: args{
				input: "рпрмтпмр",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "[calcCheckFormat]_Error variant 6",
			c:    &Calc{},
			args: args{
				input: "=?",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.calcCheckFormat(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("calc.calcCheckFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calc.calcCheckFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalc_Сalculate(t *testing.T) {
	type args struct {
		inputSlice   []string
		fullResponse bool
	}
	tests := []struct {
		name string
		c    *Calc
		args args
		want []string
	}{
		{
			name: "[Сalculate]_Error variant 1",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{},
				fullResponse: false,
			},
			want: []string{},
		},
		{
			name: "[Сalculate]_Good variant 2",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"3+5+6+8+9=?"},
				fullResponse: false,
			},
			want: []string{"3+5+6+8+9=31"},
		},
		{
			name: "[Сalculate]_Error variant 3",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"9/0"},
				fullResponse: false,
			},
			want: []string{},
		},
		{
			name: "[Сalculate]_Error variant 4",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"9/0=?"},
				fullResponse: false,
			},
			want: []string{},
		},
		{
			name: "[Сalculate]_Error variant 5",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"sVSFv"},
				fullResponse: true,
			},
			want: []string{"ошибка в строке [sVSFv]-> выражение [sVSFv] не соответствует формату: num1[*/+-]num2=?"},
		},
		{
			name: "[Сalculate]_Good variant 6",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"(3+(5-(7+1)*-(2+2-8*9)*3*-4)*(4-(-3.0*8.0+7)*6))*(12-8/9*(12+8)+(1*8/9-5+4*5))=?"},
				fullResponse: true,
			},
			want: []string{"(3+(5-(7+1)*-(2+2-8*9)*3*-4)*(4-(-3.0*8.0+7)*6))*(12-8/9*(12+8)+(1*8/9-5+4*5))=7001953.093609"},
		},
		{
			name: "[Сalculate]_Good variant 7",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"(3+(5-(7+1)*-(2+2-8*9)*3*-4)*(4-(-3.0*8.0+7)*6))*(12-8/9*(12+8)+(1*8/9-5+4*5))*(-1)=?"},
				fullResponse: true,
			},
			want: []string{"(3+(5-(7+1)*-(2+2-8*9)*3*-4)*(4-(-3.0*8.0+7)*6))*(12-8/9*(12+8)+(1*8/9-5+4*5))*(-1)=-7001953.093609"},
		},
		{
			name: "[Сalculate]_Good variant 8",
			c:    &Calc{},
			args: args{
				inputSlice:   []string{"(-1)*+8.000000=?"},
				fullResponse: true,
			},
			want: []string{"(-1)*+8.000000=-8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Сalculate(tt.args.inputSlice, tt.args.fullResponse); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calc.Сalculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
