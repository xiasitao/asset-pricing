package investments

type CurrencyUnit float64

type Period struct {
	Cashflow CurrencyUnit `json:"cashflow"`
	Interest float64      `json:"interest"`
}

func CalculatePresentValue(periods ...Period) CurrencyUnit {
	if len(periods) == 0 {
		return 0.0
	}

	currentPeriod := periods[0]
	futurePeriods := periods[1:]

	discountFactor := 1.0 / (1.0 + currentPeriod.Interest)

	discountedCurrentPeriod := float64(currentPeriod.Cashflow) * discountFactor
	discountedFuturePeriods := float64(CalculatePresentValue(futurePeriods...)) * discountFactor

	return CurrencyUnit(discountedCurrentPeriod + discountedFuturePeriods)
}
