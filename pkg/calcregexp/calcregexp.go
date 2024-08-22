package calcregexp

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	rePriorFirst       = regexp.MustCompile(`([-+]?(\d+(\.\d*)?|\.\d+))([\*\/])([-+]?(\d+(\.\d*)?|\.\d+))\s*`)
	rePriorSecond      = regexp.MustCompile(`([-+]?(\d+(\.\d*)?|\.\d+))([\+\-])([-+]?(\d+(\.\d*)?|\.\d+))([\+\-\)])\s*`)
	reParentheses      = regexp.MustCompile(`([\+\-\*])([\(])([-+]?(\d+(\.\d*)?|\.\d+)([\)]))\s*`)
	reParenthesesStart = regexp.MustCompile(`([\(])([-+]?(\d+(\.\d*)?|\.\d+)([\)]))([\+\-\*\/])\s*`)
	reFloatParentheses = regexp.MustCompile(`^\(([-+]?\d+(\.\d+)?)\)$`)
	reFloat            = regexp.MustCompile(`^([-+]?\d+(\.\d+)?)$`)
	reCheckFormat      = regexp.MustCompile(`^[\d.+\-*/()]+\=\?$`)
	reCheckParentheses = regexp.MustCompile(`^\(.*\)$`)
	symbolEnd          = "=?"
)

type Calc struct{}

func NewCalc() Calc {
	return Calc{}
}

// функция содержит функции проверки на соответствие строки, вычисление, вывод результатов
func (c *Calc) Сalculate(inputSlice []string, fullResponse bool) []string {

	response := make([]string, 0)

	for _, input := range inputSlice {
		//проверяем
		val, err := c.calcCheckFormat(input)
		if err == nil {
			//вычисляем
			val, err = c.calcStepByStep(val)
			if err != nil {
				if fullResponse {
					response = append(response, fmt.Sprintf("ошибка в строке [%v]-> %v", input, err))
				}
			} else {
				f, err := strconv.ParseFloat(val, 64)
				if err != nil && fullResponse {
					response = append(response, fmt.Sprintf("ошибка при преобразовании в float64 [%v]-> %v", input, err))
				} else {
					//выводим результат
					if math.Mod(f, 1.0) == 0 {
						response = append(response, fmt.Sprintf("%v=%.0f", strings.Trim(input, symbolEnd), f))
					} else {
						response = append(response, fmt.Sprintf("%v=%v", strings.Trim(input, symbolEnd), strconv.FormatFloat(f, 'f', -1, 64)))
					}
				}
			}
		} else {
			if fullResponse {
				response = append(response, fmt.Sprintf("ошибка в строке [%v]-> %v", input, err))
			}
		}
	}
	//fmt.Printf("%#v->%#v->\n", inputSlice, response)
	return response
}

// проверяем строук на соответствие исходным условиям
// строка должна содержать: числа(целые или вещественные), знаки +-*/, скобки() и заканчиваться =?
func (c *Calc) calcCheckFormat(input string) (string, error) {
	//удаляем пробелы
	valOut := strings.ReplaceAll(input, " ", "")
	//проверяем на соответствие основным условиям
	b := reCheckFormat.MatchString(valOut)
	if !b {
		return "", fmt.Errorf("выражение [%v] не соответствует формату: num1[*/+-]num2=?", input)
	}
	//удаляем =?
	valOut = strings.Trim(valOut, symbolEnd)
	//проверяем на скобки (.....), если нет добавляем
	b = reCheckParentheses.MatchString(valOut)
	if !b {
		valOut = "(" + valOut + ")"
	}

	return valOut, nil
}

// проходим по строке находим совпадения в определенном порядке, вычисляем
func (c *Calc) calcStepByStep(input string) (string, error) {

	valOut := input

	//вычисления закончены удачно
	//если получем (+ХХ.ХХ)
	matches := reFloatParentheses.FindStringSubmatch(valOut)
	if len(matches) == 3 {
		return matches[1], nil
	}
	//если получем +ХХ.ХХ
	matches = reFloat.FindStringSubmatch(valOut)
	if len(matches) == 3 {
		return matches[1], nil
	}
	//1. ищем ХХ.ХХ*/ХХ.ХХ, вычисляем
	val, err := c.findAllPrior(valOut, rePriorFirst)
	if err != nil {
		return valOut, err
	}
	//2. ищем ХХ.ХХ+-ХХ.ХХ, вычисляем
	val, err = c.findAllPrior(val, rePriorSecond)
	if err != nil {
		return valOut, err
	}
	//3. ищем +(+ХХ.ХХ), раскрываем скобки
	val, err = c.findAllParentheses(val, reParentheses)
	if err != nil {
		return valOut, err
	}
	//4. ищем (+ХХ.ХХ)*/+-, раскрываем скобки
	val, err = c.findAllParentheses(val, reParenthesesStart)
	if err != nil {
		return valOut, err
	}

	valOut = val
	//все деиствия выполнили данные не изменились, то ошибка
	if valOut == input {
		return valOut, fmt.Errorf("неизвестный формат выражения [%v]", valOut)
	}

	return c.calcStepByStep(valOut)
}

// 3. ищем +(+ХХ.ХХ), раскрываем скобки; 4. ищем (+ХХ.ХХ)*/+-, раскрываем скобки
func (c *Calc) findAllParentheses(input string, re *regexp.Regexp) (string, error) {
	expressions := re.FindAllStringSubmatch(input, -1)
	expressionsPlace := re.Split(input, -1)
	expressionsSlice := make([]float64, 0)
	expressionsSliceZnAfter := make([]string, 0)
	expressionsSliceZnBefore := make([]string, 0)
	expressionsRet := ""

	if len(expressions) == 0 {
		return input, nil
	}
	//разбираем данные согласно группе захвата
	for _, expression := range expressions {
		val, err := strconv.ParseFloat(expression[4], 64)
		if err != nil {
			return "", err
		}
		if expression[0][:3] == "-(+" || expression[0][:3] == "+(-" {
			val *= -1
		}
		if expression[0][:1] == "*" || expression[0][:1] == "/" {
			expressionsSliceZnBefore = append(expressionsSliceZnBefore, expression[0][:1])
			expressionsSliceZnAfter = append(expressionsSliceZnAfter, "")
		} else if expression[0][:1] == "(" {
			val, err = strconv.ParseFloat(expression[3], 64)
			expressionsSliceZnBefore = append(expressionsSliceZnBefore, "")
			expressionsSliceZnAfter = append(expressionsSliceZnAfter, expression[0][len(expression[0])-1:])
			if err != nil {
				return "", err
			}
		} else {
			expressionsSliceZnBefore = append(expressionsSliceZnBefore, "")
			expressionsSliceZnAfter = append(expressionsSliceZnAfter, "")
		}
		expressionsSlice = append(expressionsSlice, val)
	}
	//собираем обратно данные согласно группе захвата
	for i, expressionPlace := range expressionsPlace {
		express := ""
		if i < len(expressionsPlace)-1 {
			express = expressionsSliceZnBefore[i] + fmt.Sprintf("%+f", expressionsSlice[i]) + expressionsSliceZnAfter[i]
		}
		expressionsRet += expressionPlace + express
	}
	//fmt.Println("findAllParentheses", input, "->", expressionsRet)
	return c.findAllParentheses(expressionsRet, re)
}

// 1. ищем ХХ.ХХ*/ХХ.ХХ, вычисляем; 2. ищем ХХ.ХХ+-ХХ.ХХ, вычисляем
func (c *Calc) findAllPrior(input string, re *regexp.Regexp) (string, error) {

	expressions := re.FindAllStringSubmatch(input, -1)
	expressionsPlace := re.Split(input, -1)
	expressionsSlice := make([]float64, 0)
	expressionsSliceZn := make([]string, 0)
	expressionsRet := ""

	if len(expressions) == 0 {
		return input, nil
	}
	//разбираем данные согласно группе захвата
	for _, expression := range expressions {
		//вычисляем
		val, err := c.calculateLineTwoArg(expression)
		if err != nil {
			return "", err
		}
		expressionsSlice = append(expressionsSlice, val)
		if len(expression) == 9 {
			expressionsSliceZn = append(expressionsSliceZn, expression[len(expression)-1])
		} else {
			expressionsSliceZn = append(expressionsSliceZn, "")
		}
	}
	//собираем обратно данные согласно группе захвата
	for i, expressionPlace := range expressionsPlace {
		express := ""
		if i < len(expressionsPlace)-1 {
			express = fmt.Sprintf("%+f", expressionsSlice[i]) + expressionsSliceZn[i]
		}
		expressionsRet += expressionPlace + express
	}

	//fmt.Println("findAllPrior", input, "->", expressionsRet)

	//fmt.Printf("%#v->%#v->\n", input, expressionsRet)
	return c.findAllPrior(expressionsRet, re)
}

// вычисляем математическое выражение
func (c *Calc) calculateLineTwoArg(expression []string) (float64, error) {

	if len(expression) < 8 {
		return 0, fmt.Errorf("выражение [%v] не соответствует формату num1[*/+-]num2", expression)
	}

	a, err := strconv.ParseFloat(expression[1], 64)
	if err != nil {
		return 0, err
	}

	b, err := strconv.ParseFloat(expression[5], 64)
	if err != nil {
		return 0, err
	}

	operator := string(expression[4][0])

	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("некорректное выражение - деление на ноль")
		} else {
			return a / b, nil
		}
	default:
		return 0, fmt.Errorf("неподдерживаемый оператор [%v]", operator)
	}
}

// читаем данные из файла
func (c *Calc) ReadLinesFromFile(filename string) ([]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]string, 0)
	f := bufio.NewReader(file)
	for {
		line, _, err := f.ReadLine()
		if err != nil {
			break
		}
		data = append(data, string(line))
	}

	return data, nil
}

// записываем данные в файл
func (c *Calc) WriteLinesToFile(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
