package crypto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// CompressEnvMap compresses a serialized env map using gzip.
// This is useful before encryption to reduce ciphertext size.
func CompressEnvMap(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("compress: write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress: close failed: %w", err)
	}
	return buf.Bytes(), nil
}

// DecompressEnvMap decompresses gzip-compressed env data.
func DecompressEnvMap(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decompress: reader failed: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decompress: read failed: %w", err)
	}
	return out, nil
}

// CompressedSize returns the gzip-compressed size of the given byte slice.
func CompressedSize(data []byte) (int, error) {
	compressed, err := CompressEnvMap(data)
	if err != nil {
		return 0, err
	}
	return len(compressed), nil
}

// CompressionRatio returns the ratio of compressed to original size.
// A value < 1.0 means compression saved space.
func CompressionRatio(data []byte) (float64, error) {
	if len(data) == 0 {
		return 1.0, nil
	}
	csz, err := CompressedSize(data)
	if err != nil {
		return 0, err
	}
	return float64(csz) / float64(len(data)), nil
}
