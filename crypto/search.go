package crypto

import "strings"

// SearchResult holds a matched key and its value from a vault entry.
type SearchResult struct {
	Version   int
	Key       string
	Value     string
	Snippet   string
}

// SearchVault searches all decrypted vault entries for keys or values matching
// the given query string (case-insensitive). privateKeyPath is used to decrypt
// each entry before searching.
func SearchVault(dir, privateKeyPath, query string) ([]SearchResult, error) {
	vault, err := LoadVault(dir)
	if err != nil {
		return nil, err
	}

	privKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var results []SearchResult

	for _, entry := range vault.Entries {
		plaintext, err := Decrypt(privKey, entry.Ciphertext)
		if err != nil {
			continue
		}

		pairs, err := ParseEnvBytes(plaintext)
		if err != nil {
			continue
		}

		for _, kv := range pairs {
			kLow := strings.ToLower(kv.Key)
			vLow := strings.ToLower(kv.Value)
			if strings.Contains(kLow, q) || strings.Contains(vLow, q) {
				snippet := kv.Key + "=" + kv.Value
				if len(snippet) > 60 {
					snippet = snippet[:57] + "..."
				}
				results = append(results, SearchResult{
					Version: entry.Version,
					Key:     kv.Key,
					Value:   kv.Value,
					Snippet: snippet,
				})
			}
		}
	}

	return results, nil
}
