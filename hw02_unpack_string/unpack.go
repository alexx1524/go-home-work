package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack функция используется для примитивной распаковки строки.
func Unpack(input string) (string, error) {
	var sb strings.Builder

	runes := []rune(input)
	length := len(runes)

	for i := 0; i < length; {
		switch {
		// если цифра в любой позиции стоит первее символа, то возвращаем ошибку
		case unicode.IsDigit(runes[i]):
			return "", ErrInvalidString

		// обработка экранированных символов
		case string(runes[i]) == `\`:
			symbolPosition := i + 1
			quantifierPosition := i + 2
			if symbolPosition < length && (unicode.IsDigit(runes[symbolPosition]) || string(runes[symbolPosition]) == `\`) {
				if quantifierPosition < length && unicode.IsDigit(runes[quantifierPosition]) {
					quantifier, err := strconv.Atoi(string(runes[quantifierPosition]))
					if err != nil {
						return "", err
					}

					// обработали руну экранирования, символа и количества, увеличиваем переменную цикла на 3
					sb.WriteString(strings.Repeat(string(runes[symbolPosition]), quantifier))
					i += 3
					continue
				}

				// обработали руну символа и руну количества, увеличиваем переменную цикла на 2
				sb.WriteRune(runes[symbolPosition])
				i += 2
				continue
			}
			return "", ErrInvalidString

		// обработка отдельного неэкранированного символа с количеством или без
		default:
			quantifierPosition := i + 1
			if quantifierPosition < length && unicode.IsDigit(runes[quantifierPosition]) {
				quantifier, err := strconv.Atoi(string(runes[quantifierPosition]))
				if err != nil {
					return "", err
				}

				// обработали руну символа и руну количества, увеличиваем переменную цикла на 2
				sb.WriteString(strings.Repeat(string(runes[i]), quantifier))
				i += 2
				continue
			}
			// обработали только одну руну символа без количества, увеличиваем переменную цикла на 1
			sb.WriteRune(runes[i])
			i++
		}
	}

	return sb.String(), nil
}
