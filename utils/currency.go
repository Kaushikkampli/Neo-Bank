package utils

const (
	USD = "USD"
	INR = "INR"
	EUR = "EUR"
)

func IsSupportedCurrency(curr string) bool {
	switch curr {
	case USD, INR, EUR:
		return true
	}
	return false
}
