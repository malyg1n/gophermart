package validation

// IsLunh check number by lunh algo.
func IsLunh(number string) bool {
	var checksum int

	numberLen := len(number)
	for i := numberLen - 1; i >= 0; i -= 2 {
		n := number[i] - '0'
		checksum += int(n)
	}
	for i := numberLen - 2; i >= 0; i -= 2 {
		n := number[i] - '0'
		n *= 2
		if n > 9 {
			n -= 9
		}
		checksum += int(n)
	}

	return checksum%10 == 0
}
