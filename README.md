# Shrinkr

Fast image compression & format conversion CLI tool. Batch compress JPG, PNG, TIFF, WebP, HEIC, AVIF, and GIF images with an interactive wizard or direct CLI flags.

## Install

### Homebrew (recommended)

```bash
brew install idrewlong/tap/shrinkr
```

Homebrew handles the `libvips` dependency automatically.

### Go install

Requires Go and `libvips` installed on your system.

```bash
brew install vips       # macOS
go install github.com/idrewlong/shrinkr_cli@latest
```

### From source

```bash
git clone https://github.com/idrewlong/shrinkr_cli
cd shrinkr_cli
make install
```

## Usage

### Interactive wizard

Run `shrinkr` with no arguments to launch the guided wizard:

```bash
shrinkr
```

The wizard will:
1. Auto-detect nearby folders with images — or browse with Finder / enter a path manually
2. Ask for output format (WebP, PNG, JPEG, or AVIF)
3. Ask for compression settings (Recommended, Web Optimized, High Quality, or Custom)
4. Ask where to save the output — defaults to a `compressed/` folder inside your selected input folder
5. Show a summary and confirm before running

> **Tip:** Press `Esc` or `Ctrl+C` at any step to go back one step. Pressing either key repeatedly from the first step exits the program.

### Direct CLI mode

Pass the input folder as an argument for scripting or power use:

```bash
# Compress a folder to WebP at 500 KB target (default)
shrinkr ./photos

# Custom format, size, and output location
shrinkr ./photos -f jpeg -s 300 -o ./compressed-photos

# Recursive scan with quality settings
shrinkr ~/Pictures -r -f webp -s 200 -q 80

# All options
shrinkr <folder> [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `webp` | Output format: `webp`, `png`, `jpeg`, `avif` |
| `--size` | `-s` | `500` | Target file size in KB |
| `--output` | `-o` | `compressed` | Output folder path |
| `--quality` | `-q` | `85` | Initial quality (1–100) |
| `--workers` | `-w` | CPU count | Concurrent workers (auto-detected) |
| `--recursive` | `-r` | `false` | Scan subfolders |
| `--min-quality` | | `60` | Quality floor |
| `--max-quality` | | `90` | Quality ceiling |

## Supported formats

| Input | Output |
|-------|--------|
| JPG / JPEG | WebP |
| PNG | PNG |
| TIFF / TIF | JPEG |
| WebP | AVIF |
| GIF | |
| AVIF | |
| HEIF / HEIC | |

> **AVIF note:** AVIF produces the smallest files but encodes significantly slower than other formats — best for small batches where file size is the top priority.

## How it works

Shrinkr uses a **binary search algorithm** to find the optimal compression quality that hits your target file size:

1. Tries initial quality (default 85)
2. If the result is over target, binary searches between min and max quality
3. Falls back to min quality if the target can't be reached
4. Runs all compressions concurrently using a worker pool sized to your CPU core count

Processing uses [`libvips`](https://libvips.github.io/libvips/) via the `bimg` Go wrapper — one of the fastest image processing libraries available.

## Development

```bash
# Run all tests
make test

# Run format benchmarks (shows per-format encoding speed comparison)
make bench

# Build binary locally
make build

# Build and install globally
make install
```

## Updates

```bash
brew upgrade shrinkr
```

## Requirements

- macOS (Apple Silicon or Intel)
- `libvips` — installed automatically by Homebrew

## License

MIT
