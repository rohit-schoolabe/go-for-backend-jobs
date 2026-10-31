package main

import "fmt"

func main() {
	maturity := 70127.586535
	years := 5

	fmt.Println("Maturity:", maturity, "after", years, "years")
	// fmt.Println("commented out, so it does not run")
	/* a block comment
	   can span lines */

	fmt.Printf("%v years gives %v\n", years, maturity)
	fmt.Printf("%f\n", maturity)      // 70127.586535
	fmt.Printf("%.2f\n", maturity)    // 70127.59
	fmt.Printf("%.0f\n", maturity)    // 70128
	fmt.Printf("%8.2f|\n", 42.5)      //    42.50|
	fmt.Printf("%T\n", maturity)      // float64
	fmt.Printf("%d\n", years)         // 5
	fmt.Printf("%s %q\n", "hi", "hi") // hi "hi"
	fmt.Printf("%t\n", true)          // true
	fmt.Printf("7%%\n")               // 7%

	summary := fmt.Sprintf("Maturity: ₹%.2f", maturity)
	fmt.Println(summary)

	report := `FD summary
----------
Deposit  : %.2f
Maturity : %.2f
`
	fmt.Printf(report, 50000.0, maturity)
	fmt.Println(`C:\Users\priya\n stays as typed`)
}
