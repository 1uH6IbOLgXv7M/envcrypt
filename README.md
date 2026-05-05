# envcrypt

A zero-dependency CLI tool to encrypt and version-control `.env` files using age encryption.

---

## Installation

```bash
go install github.com/yourusername/envcrypt@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envcrypt/releases).

---

## Usage

**Encrypt a `.env` file before committing:**

```bash
envcrypt encrypt .env --output .env.age
```

**Decrypt on another machine or in CI:**

```bash
envcrypt decrypt .env.age --output .env
```

**Generate a new key pair:**

```bash
envcrypt keygen
```

Add `.env` to your `.gitignore` and commit `.env.age` safely to version control.

```gitignore
.env
```

> Keys are stored in `~/.config/envcrypt/keys.txt` by default. Set `ENVCRYPT_KEY` as an environment variable to use a custom key path in CI/CD pipelines.

---

## Why envcrypt?

- **Zero dependencies** — single static binary, no runtime required
- **Built on [age](https://age-encryption.org/)** — modern, audited encryption
- **Git-friendly** — encrypted files are safe to commit and diff
- **Simple** — one command to encrypt, one to decrypt

---

## License

MIT © 2024 yourusername