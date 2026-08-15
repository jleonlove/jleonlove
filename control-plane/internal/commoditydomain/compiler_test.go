package commoditydomain

import "testing"

func TestSeedAndResolve(t *testing.T) {
	c := New()
	if err := SeedCore(c); err != nil {
		t.Fatal(err)
	}
	if len(c.IDs()) != 5 {
		t.Fatalf("packs=%d", len(c.IDs()))
	}
	p, err := c.Resolve("sulphur")
	if err != nil || p.Pack.Commodity != "Sulfur" {
		t.Fatalf("resolve=%+v err=%v", p, err)
	}
	w, _ := c.Resolve("milling wheat")
	if len(w.Pack.Documents) == 0 {
		t.Fatal("missing document rules")
	}
}
func TestRejectInvalidAndDuplicate(t *testing.T) {
	c := New()
	if _, err := c.Compile(Pack{}); err != ErrInvalidPack {
		t.Fatalf("err=%v", err)
	}
	p := Pack{ID: "x", Family: "f", Commodity: "c", Version: "1"}
	if _, err := c.Compile(p); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Compile(p); err != ErrDuplicatePack {
		t.Fatalf("err=%v", err)
	}
}
