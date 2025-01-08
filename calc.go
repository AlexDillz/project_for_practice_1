package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Проверка операторов
func Operators(char byte) bool {
	return char == '+' || char == '-' || char == '*' || char == '/' || char == '%'
}

// Приоритет операторов
func Priority(char byte) int {
	if char == '*' || char == '/' || char == '%' {
		return 2
	}
	if char == '+' || char == '-' {
		return 1
	}
	return 0
}

// Выполнение арифметических операций
func OperatorsUsing(a, b float64, operator byte) (float64, error) {
	switch operator {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, errors.New("dividing by zero")
		}
		return a / b, nil
	case '%':
		if b == 0 {
			return 0, errors.New("dividing by zero in modulus")
		}
		return math.Mod(a, b), nil
	default:
		return 0, errors.New("invalid operator")
	}
}

// Выполнение тригонометрических и логарифмических функций
func MathFunctions(funcName string, param float64) (float64, error) {
	switch funcName {
	case "sin":
		return math.Sin(param), nil
	case "cos":
		return math.Cos(param), nil
	case "tan":
		return math.Tan(param), nil
	case "log":
		if param <= 0 {
			return 0, errors.New("logarithm undefined for non-positive values")
		}
		return math.Log(param), nil
	case "log10":
		if param <= 0 {
			return 0, errors.New("logarithm undefined for non-positive values")
		}
		return math.Log10(param), nil
	default:
		return 0, errors.New("unknown function: " + funcName)
	}
}

// Вычисление выражения
func Calculation(input string) (float64, error) {
	var nums []float64
	var ops []byte

	i := 0
	for i < len(input) {
		char := input[i]

		// Исключение пробелов
		if unicode.IsSpace(rune(char)) {
			i++
			continue
		}

		// Обработка числа
		if unicode.IsDigit(rune(char)) || char == '.' {
			numStart := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(input[numStart:i], 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %v", err)
			}
			nums = append(nums, num)
			continue
		}

		// Обработка функции
		if unicode.IsLetter(rune(char)) {
			funcStart := i
			for i < len(input) && unicode.IsLetter(rune(input[i])) {
				i++
			}
			funcName := input[funcStart:i]

			if i < len(input) && input[i] == '(' {
				openParen := i
				i++
				paramStart := i
				fmt.Printf("Обнаружена открывающая скобка на позиции %d\n", openParen)
				for i < len(input) && input[i] != ')' {
					i++
				}
				if i == len(input) {
					return 0, errors.New("mismatched parentheses in function")
				}
				param, err := Calculation(input[paramStart:i])
				if err != nil {
					return 0, err
				}
				result, err := MathFunctions(funcName, param)
				if err != nil {
					return 0, err
				}
				nums = append(nums, result)
				i++ // Пропускаем закрывающую скобку
				continue
			} else {
				return 0, errors.New("invalid function syntax")
			}
		}

		// Обработка операторов
		if Operators(char) {
			for len(ops) > 0 && Priority(ops[len(ops)-1]) >= Priority(char) {
				b := nums[len(nums)-1]
				a := nums[len(nums)-2]
				nums = nums[:len(nums)-2]
				result, err := OperatorsUsing(a, b, ops[len(ops)-1])
				if err != nil {
					return 0, err
				}
				ops = ops[:len(ops)-1]
				nums = append(nums, result)
			}
			ops = append(ops, char)
		} else if char == '(' {
			ops = append(ops, char)
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				b := nums[len(nums)-1]
				a := nums[len(nums)-2]
				nums = nums[:len(nums)-2]
				result, err := OperatorsUsing(a, b, ops[len(ops)-1])
				if err != nil {
					return 0, err
				}
				ops = ops[:len(ops)-1]
				nums = append(nums, result)
			}
			if len(ops) == 0 || ops[len(ops)-1] != '(' {
				return 0, errors.New("mismatched parentheses")
			}
			ops = ops[:len(ops)-1]
		} else {
			return 0, errors.New("invalid character in input")
		}

		i++
	}

	// Выполнение оставшихся операций
	for len(ops) > 0 {
		b := nums[len(nums)-1]
		a := nums[len(nums)-2]
		nums = nums[:len(nums)-2]
		result, err := OperatorsUsing(a, b, ops[len(ops)-1])
		if err != nil {
			return 0, err
		}
		ops = ops[:len(ops)-1]
		nums = append(nums, result)
	}

	if len(nums) != 1 {
		return 0, errors.New("invalid expression")
	}

	return nums[0], nil
}

// Основная функция
func main() {
	fmt.Println("Введите математическое выражение (или 'exit' для выхода):")
	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)

		// Проверка команды выхода
		if strings.ToLower(input) == "exit" {
			fmt.Println("Завершение программы.")
			break
		}

		// Вычисление результата
		result, err := Calculation(input)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	}
}
