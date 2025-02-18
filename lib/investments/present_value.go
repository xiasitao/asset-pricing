package investments

import "math"

type CurrencyUnit float64

type Period struct {
	Cashflow CurrencyUnit `json:"cashflow"`
	Interest float64      `json:"interest"`
}

func CalculatePerpetuityPresentValue(cashflow CurrencyUnit, interest float64) CurrencyUnit {
	presentValue := CurrencyUnit(float64(cashflow) / interest)
	return presentValue
}

func CalculateLumpSumPresentValue(lump CurrencyUnit, interest float64, periods int) CurrencyUnit {
	presentValue := CurrencyUnit(float64(lump) / math.Pow(1+interest, float64(periods)))
	return presentValue
}

func CalculateAnnuityPresentValue(cashflow CurrencyUnit, interest float64, periods int) CurrencyUnit {
	presentValue := CurrencyUnit(float64(cashflow) * (1 - math.Pow(1+interest, -float64(periods))) / interest)
	return presentValue
}

func CalculateGeneralFiniteCashflowPresentValue(periods ...Period) CurrencyUnit {
	if len(periods) == 0 {
		return 0.0
	}

	currentPeriod := periods[0]
	futurePeriods := periods[1:]

	discountFactor := 1.0 / (1.0 + currentPeriod.Interest)

	discountedCurrentPeriod := float64(currentPeriod.Cashflow) * discountFactor
	discountedFuturePeriods := float64(CalculateGeneralFiniteCashflowPresentValue(futurePeriods...)) * discountFactor

	return CurrencyUnit(discountedCurrentPeriod + discountedFuturePeriods)
}
