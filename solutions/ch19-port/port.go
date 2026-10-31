// Package port shows the listen address rule for
// platforms that set PORT.
package port

// Addr prefers ADDR, then PORT, then :8080.
func Addr(getenv func(string) string) string {
	if a := getenv("ADDR"); a != "" {
		return a
	}
	if p := getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}
