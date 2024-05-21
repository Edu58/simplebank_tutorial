package util

const (
	USD = "USD"
	KES = "KES"
	EUR = "EUR"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, KES, EUR:
		return true
	}

	return false
}
