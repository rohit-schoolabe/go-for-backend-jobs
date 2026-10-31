package main

import "fmt"

type adjustFn func(float64) float64

func adjustAll(fines *[]float64, adjust adjustFn) []float64 {
	out := make([]float64, 0, len(*fines))
	for _, f := range *fines {
		out = append(out, adjust(f))
	}
	return out
}

func half(f float64) float64   { return f / 2 }
func double(f float64) float64 { return f * 2 }

func ruleFor(memberType string) adjustFn {
	if memberType == "student" {
		return half
	}
	return func(f float64) float64 { return f }
}

func capAt(limit float64) adjustFn {
	return func(f float64) float64 {
		return min(f, limit)
	}
}

func total(label string, fines ...float64) string {
	sum := 0.0
	for _, f := range fines {
		sum += f
	}
	return fmt.Sprintf("%s: %.2f", label, sum)
}

func main() {
	fines := []float64{18, 4, 50, 12}

	fmt.Println(adjustAll(&fines, half))   // [9 2 25 6]
	fmt.Println(adjustAll(&fines, double)) // [36 8 100 24]

	rule := ruleFor("student")
	fmt.Println(adjustAll(&fines, rule)) // [9 2 25 6]

	plusTwo := adjustAll(&fines, func(f float64) float64 {
		return f + 2
	})
	fmt.Println(plusTwo) // [20 6 52 14]

	capAt30 := capAt(30)
	capAt10 := capAt(10)
	fmt.Println(adjustAll(&fines, capAt30)) // [18 4 30 12]
	fmt.Println(adjustAll(&fines, capAt10)) // [10 4 10 10]

	fmt.Println(total("Today", 18, 4))  // Today: 22.00
	fmt.Println(total("None"))          // None: 0.00
	fmt.Println(total("All", fines...)) // All: 84.00
}
