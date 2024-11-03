package luhn

import "strconv"

// VerifyLuhn is an implementation of Luhn algorithm.
func VerifyLuhn(number string) bool {
	sum := 0
	parity := len(number) % 2

	for i := 0; i < len(number); i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
