package investments_test

import (
	"fmt"
	"testing"

	i "xiasitao.de/asset-pricing/lib/investments"
)

func TestPerpetuityPresentValue(t *testing.T) {
	cashflow, interest := i.CurrencyUnit(1000.0), 0.2
	got := i.CalculatePerpetuityPresentValue(cashflow, interest)
	expected := i.CurrencyUnit(1000.0 / 0.2)
	assertPresentValue(t, got, expected)
}

func TestLumpSumPresentValue(t *testing.T) {
	lump, interest, periods := i.CurrencyUnit(1000.0), 0.2, 5
	got := i.CalculateLumpSumPresentValue(lump, interest, periods)
	expected := i.CurrencyUnit(1000.0 / 1.2 / 1.2 / 1.2 / 1.2 / 1.2)
	assertPresentValue(t, got, expected)
}

func TestAnnuityPresentValue(t *testing.T) {
	cashflow, interest, periods := i.CurrencyUnit(1000.0), 0.2, 5
	got := i.CalculateAnnuityPresentValue(cashflow, interest, periods)
	expected := cashflow/1.2 + cashflow/1.2/1.2 + cashflow/1.2/1.2/1.2 + cashflow/1.2/1.2/1.2/1.2 + cashflow/1.2/1.2/1.2/1.2/1.2
	assertPresentValue(t, got, expected)
}

func TestGeneralFiniteCashflowPresentValue(t *testing.T) {
	type Case struct {
		name     string
		periods  []i.Period
		expected float64
	}

	cases := []Case{
		{"no period", []i.Period{}, 0.0},
		{
			"single present cashflow with zero interest",
			[]i.Period{{Cashflow: 1000.0}},
			1000.0,
		},
		{
			"single present cashflow with non-zero interest",
			[]i.Period{{Cashflow: 1000.0, Interest: 0.1}},
			1000.0 / 1.1,
		},
		{
			"single future cashflow with zero interest",
			[]i.Period{{}, {}, {Cashflow: 1000.0}},
			1000.0,
		},
		{
			"single future cashflow with non-zero interest",
			[]i.Period{{Interest: 0.1}, {Interest: 0.2}, {Cashflow: 1000.0, Interest: 0.3}},
			1000.0 / (1.1 * 1.2 * 1.3),
		},
		{
			"two cashflows with zero interest",
			[]i.Period{{Cashflow: 1000.0}, {Cashflow: 2000.0}},
			3000.0,
		},
		{
			"two equal cashflows with uniform non-zero interest",
			[]i.Period{{Cashflow: 1000.0, Interest: 0.1}, {Cashflow: 1000.0, Interest: 0.2}},
			1000.0/1.1 + 1000.0/(1.1*1.2),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := i.CalculateGeneralFiniteCashflowPresentValue(test.periods...)
			assertPresentValue(t, got, i.CurrencyUnit(test.expected))
		})
	}
}

func assertPresentValue(t testing.TB, got, expected i.CurrencyUnit) {
	t.Helper()
	toFixedString := func(value i.CurrencyUnit) string {
		return fmt.Sprintf("%.8f", value)
	}
	if toFixedString(got) != toFixedString(expected) {
		t.Errorf("got %f, expected %f", got, expected)
	}
}
