package crypto

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompressDecompress_RoundTrip(t *testing.T) {
	original := []byte("DB_HOST=localhost\nDB_PORT=5432\nDB_NAME=mydb\nDB_USER=admin\nDB_PASS=secret\n")
	compressed, err := CompressEnvMap(original)
	if err != nil {
		t.Fatalf("compress failed: %v", err)
	}
	decompressed, err := DecompressEnvMap(compressed)
	if err != nil {
		t.Fatalf("decompress failed: %v", err)
	}
	if !bytes.Equal(original, decompressed) {
		t.Errorf("round-trip mismatch: got %q, want %q", decompressed, original)
	}
}

func TestCompressEnvMap_ReducesSize(t *testing.T) {
	// Repetitive content compresses well
	original := []byte(strings.Repeat("KEY_ALPHA=value\nKEY_BETA=value\n", 20))
	compressed, err := CompressEnvMap(original)
	if err != nil {
		t.Fatalf("compress failed: %v", err)
	}
	if len(compressed) >= len(original) {
		t.Errorf("expected compression to reduce size: %d -> %d", len(original), len(compressed))
	}
}

func TestDecompressEnvMap_InvalidData(t *testing.T) {
	_, err := DecompressEnvMap([]byte("not gzip data"))
	if err == nil {
		t.Error("expected error for invalid gzip data")
	}
}

func TestCompressedSize_Empty(t *testing.T) {
	sz, err := CompressedSize([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sz == 0 {
		t.Error("expected non-zero compressed size even for empty input (gzip header overhead)")
	}
}

func TestCompressionRatio_EmptyInput(t *testing.T) {
	ratio, err := CompressionRatio([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ratio != 1.0 {
		t.Errorf("expected ratio 1.0 for empty input, got %f", ratio)
	}
}

func TestCompressionRatio_RepetitiveData(t *testing.T) {
	data := []byte(strings.Repeat("REPEATED_KEY=repeated_value\n", 30))
	ratio, err := CompressionRatio(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ratio >= 1.0 {
		t.Errorf("expected ratio < 1.0 for repetitive data, got %f", ratio)
	}
}
