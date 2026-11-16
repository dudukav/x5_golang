package config

import "testing"

func TestLoadHTTPAddr(t *testing.T) {
	t.Setenv("HTTP_ADDR", "127.0.0.1:9000")

	got := Load()
	if got.HTTPAddr != "127.0.0.1:9000" {
		t.Fatalf("HTTPAddr = %q, want %q", got.HTTPAddr, "127.0.0.1:9000")
	}
}

func TestLoadUsesDefaultHTTPAddr(t *testing.T) {
	for _, value := range []string{"", "   "} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("HTTP_ADDR", value)

			got := Load()
			if got.HTTPAddr != DefaultHTTPAddr {
				t.Fatalf("HTTPAddr = %q, want default %q", got.HTTPAddr, DefaultHTTPAddr)
			}
		})
	}
}
