package pbdk

import (
	"encoding/base64"
	"testing"
)

func TestNewSalt(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "abbracadabbra!", want: "rUo2KlNbOku4lkxTgFNU6g"},
		{input: "s1mSal5Bim$$", want: "23CXE9KcP46HjoF8TpQuLQ"},
		{input: "maGIKaBul444=", want: "lWELa5xnp+gEhzPkqUgJZA"},
		{input: "ooOOr1uk3nN!!!", want: "A+o1LoKr62W4woL0mz6mtg"},
		{input: "AkKa!n1sc1uNèF355", want: "94o1QxYSM14kqAMXR8IM8w"},
	}

	for _, tc := range tests {
		got, err := NewSalt([]byte(tc.input))
		if err != nil {
			t.Error(err)
		}

		enc := base64.RawStdEncoding.EncodeToString(got)
		if tc.want != enc {
			t.Fatalf("expected: %v, got: %v", tc.want, enc)
		}
	}
}

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "abbracadabbra!", want: "43mi9wbwbyKtvi+0cwz/BEWTMp3wNwOlq6S5EgfTli8="},
		{input: "s1mSal5Bim$$", want: "cr6LMcYnGo2uwa1H6xhRa3vLQSC1pDjo4kko5p0MnOA="},
		{input: "maGIKaBul444=", want: "WMHc2B6Pg6M33ruqWWEjNk7tAIcbcsRP+0ojLnxE5jU="},
		{input: "ooOOr1uk3nN!!!", want: "bVl9bwLpzSmP5LZPWeNgh3K58n09fuu6F6Kv8BGSkj8="},
		{input: "AkKa!n1sc1uNèF355", want: "IGQ1+fni6OCyBmKuMqqLMa6z73WLCc19nJCP3l2W9sY="},
	}

	for _, tc := range tests {
		got, err := DeriveKey([]byte(tc.input))
		if err != nil {
			t.Error(err)
		}

		encoded := base64.StdEncoding.EncodeToString(got[:])
		if tc.want != encoded {
			t.Fatalf("expected: %v, got: %v", tc.want, encoded)
		}
	}
}
