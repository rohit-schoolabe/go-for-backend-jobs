package bufpool

import "testing"

func TestSameText(t *testing.T) {
	a := Render("Asha", "Optics", 3)
	b := RenderPooled("Asha", "Optics", 3)
	c := RenderPooled("Asha", "Optics", 3)
	if a != b || b != c {
		t.Fatalf("texts differ:\n%s\n%s", a, c)
	}
}

func BenchmarkRender(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Render("Asha", "Optics", i)
	}
}

func BenchmarkRenderPooled(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RenderPooled("Asha", "Optics", i)
	}
}
