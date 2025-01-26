package investments_test

import (
	"fmt"
	"testing"

	i "xiasitao.de/asset-pricing/lib/investments"
)

func TestPresentValue(t *testing.T) {
	assertPresentValue := func(t testing.TB, got, expected i.CurrencyUnit) {
		t.Helper()
		toFixedString := func(value i.CurrencyUnit) string {
			return fmt.Sprintf("%.8f", value)
		}
		if toFixedString(got) != toFixedString(expected) {
			t.Errorf("got %f, expected %f", got, expected)
		}
	}

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
			got := i.CalculatePresentValue(test.periods...)
			assertPresentValue(t, got, i.CurrencyUnit(test.expected))
		})
	}
}
