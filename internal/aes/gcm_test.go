package aes

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/lucasepe/kvs/internal/pbdk"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		input  string
		secret string
	}{
		{
			input:  "The Force will be with you",
			secret: "abbracadabbra!",
		},
		{
			input:  "A long time ago in a galaxy far, far away",
			secret: "s1mSal5Bim$$",
		},
		{
			input:  "Do. Or do not. There is no try.",
			secret: "maGIKaBul444=",
		},
		{
			input:  "Never tell me the odds!",
			secret: "ooOOr1uk3nN!!!",
		},
		{
			input:  "Chewie, we’re home.",
			secret: "AkKa!n1sc1uNèF355",
		},
	}

	for _, tc := range tests {
		key, err := pbdk.DeriveKey([]byte(tc.secret))
		if err != nil {
			t.Error(err)
		}

		enc, err := GcmEncrypt([]byte(tc.input), key)
		if err != nil {
			t.Error(err)
		}

		dec, _ := GcmDecrypt(enc, key)
		if !reflect.DeepEqual(tc.input, string(dec)) {
			t.Fatalf("expected: %v, got: %v", tc.input, string(dec))
		}
	}
}

func TestGcmDecrypt(t *testing.T) {
	t.Parallel() // marks TestGcmDecrypt as capable of running in parallel with other tests
	tests := []struct {
		input  string
		secret string
		want   string
	}{
		{
			secret: "abbracadabbra!",
			input:  "+D7jTP9GxUGOcMnK2J4I7kHxH1dM+1PVPh/RB6FVQPFZ9rvtf3wTEX23f5By8KtArOqq/cNZ",
			want:   "The Force will be with you",
		},
		{
			secret: "s1mSal5Bim$$",
			input:  "ra66MZ/I0mEyczfdPsMThHZTIjbpRCOCEFahbIiIgcFz3teDCN/PPqNhBysP/c3avy9ofOgt+GpdfxgccyVlV1Yu0rnZ",
			want:   "A long time ago in a galaxy far, far away",
		},
		{
			secret: "maGIKaBul444=",
			input:  "84UXhttoDv/0PkQmMBYg+otDlKP5TX5FKWK5SlI8YHadbRM2+dZkUvi69vbkfJwmyw21WsIvZog3GCI",
			want:   "Do. Or do not. There is no try.",
		},
		{
			secret: "ooOOr1uk3nN!!!",
			input:  "8mSEC7Yx2CFRb/iA6G2o28zSONpNigXIWevj9tT/tWFHXmvcv4Ag4w3cERb88sFty7xc",
			want:   "Never tell me the odds!",
		},
		{
			secret: "AkKa!n1sc1uNèF355",
			input:  "TeL09re/GDeG8AmkGd3BlldOXv7XB1CCNLHEmDM7P37KIA0+fkYK0LAW7S2mnfHI7Q",
			want:   "Chewie, we’re home.",
		},
	}

	for _, tc := range tests {
		key, err := pbdk.DeriveKey([]byte(tc.secret))
		if err != nil {
			t.Error(err)
		}

		in, err := base64.RawStdEncoding.DecodeString(tc.input)
		if err != nil {
			t.Error(err)
		}

		dec, err := GcmDecrypt(in, key)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual([]byte(tc.want), dec) {
			t.Fatalf("want: %q, got: %q", tc.want, dec)
		}
	}
}
