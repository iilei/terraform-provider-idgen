package provider

import "testing"

func TestParseBigSeed(t *testing.T) {
	tests := []struct {
		name   string
		seed   string
		format string
		want   string // decimal representation; empty means an error is expected
	}{
		{"hex lowercase", "ff", seedFormatHex, "255"},
		{"hex uppercase", "FF", seedFormatHex, "255"},
		{"hex 0x prefix", "0xff", seedFormatHex, "255"},
		{"hex 0X prefix", "0XFF", seedFormatHex, "255"},
		{"hex 0x without letters", "0x10", seedFormatHex, "16"},
		{"hex digits only are still hex", "10", seedFormatHex, "16"},
		{"hex leading zeros", "0000", seedFormatHex, "0"},
		{"decimal", "10", seedFormatDecimal, "10"},
		{"decimal 2^256-1", "115792089237316195423570985008687907853269984665640564039457584007913129639935", seedFormatDecimal,
			"115792089237316195423570985008687907853269984665640564039457584007913129639935"},
		{"hex negative", "-ff", seedFormatHex, ""},
		{"decimal negative", "-5", seedFormatDecimal, ""},
		{"plus sign", "+5", seedFormatDecimal, ""},
		{"empty", "", seedFormatHex, ""},
		{"prefix only", "0x", seedFormatHex, ""},
		{"not hex", "xyz", seedFormatHex, ""},
		{"hex letters as decimal", "ff", seedFormatDecimal, ""},
		{"0x prefix as decimal", "0x10", seedFormatDecimal, ""},
		{"whitespace", " ff", seedFormatHex, ""},
		{"underscore", "f_f", seedFormatHex, ""},
		{"unknown format", "ff", "base64", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, errMsg := parseBigSeed(tc.seed, tc.format)
			if tc.want == "" {
				if errMsg == "" {
					t.Fatalf("parseBigSeed(%q, %q) = %s, want error", tc.seed, tc.format, got)
				}
				return
			}
			if errMsg != "" {
				t.Fatalf("parseBigSeed(%q, %q) error: %s", tc.seed, tc.format, errMsg)
			}
			if got.String() != tc.want {
				t.Fatalf("parseBigSeed(%q, %q) = %s, want %s", tc.seed, tc.format, got, tc.want)
			}
		})
	}
}
