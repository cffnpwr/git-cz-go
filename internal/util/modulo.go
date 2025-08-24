package util

func GenMod(divisor int) func(dividend int) int {
	return func(dividend int) int {
		r := dividend % divisor
		if r < 0 {
			r += divisor
		}
		return r
	}
}
