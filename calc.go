// Простой калькулятор на командной строке
// Это классика.
// Ты можешь начать с калькулятора для простых арифметических операций,
// а потом добавить дополнительные функции, например, работу с тригонометрическими или логарифмическими функциям
package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// функция для проверки операторов
func Operators(char byte) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

// функция для приоритета операторов
func Priority(char byte) int {
	if char == '*' || char == '/' {
		return 2
	}
	if char == '+' || char == '-' {
		return 1
	}
	return 0
}

// функция для выполнения операций
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
			return 0, errors.New("Dividing by zero")
		}
		return a / b, nil
	default:
		return 0, errors.New("Invalid operator")
	}
}

// функция для вычисления выражений
func Calculation(input string) (float64, error) {
	var nums []float64
	var ops []byte

	i := 0
	for i < len(input) {
		char := input[i]

		// исключение пробелов
		if unicode.IsSpace(rune(char)) {
			i++
			continue
		}

		// чтение числа
		if unicode.IsDigit(rune(char)) || char == '.' {
			numStart := i
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(input[numStart:i], 64)
			if err != nil {
				return 0, fmt.Errorf("Invalid number: %v", err)
			}
			nums = append(nums, num)
			continue
		}

		// операторы
		if Operators(char) {
			for len(ops) > 0 && Priority(ops[len(ops)-1]) >= Priority(char) {
				if len(nums) < 2 {
					return 0, errors.New("Invalid expression")
				}
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
		}

		// обработка скобок
		if char == '(' {
			ops = append(ops, char)
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(nums) < 2 {
					return 0, errors.New("Invalid expression")
				}
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
				return 0, errors.New("Mismatched parentheses")
			}
			ops = ops[:len(ops)-1]
		}

		i++
	}

	// выполнение оставшихся операций
	for len(ops) > 0 {
		if len(nums) < 2 {
			return 0, errors.New("Invalid expression")
		}
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
		return 0, errors.New("Invalid expression")
	}

	return nums[0], nil
}

func main() {
	for {
		var input string
		fmt.Println("Введите математическое выражение (или 'exit' для выхода):")
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		result, err := Calculation(input)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	}
}
