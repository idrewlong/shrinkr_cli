# Shrinkr — Landing Page Project Scope

## Project Overview

A marketing website for **Shrinkr**, a fast image compression and format conversion CLI tool for macOS. Two pages: a hero-driven landing page and a comprehensive documentation page. Inspired by the clean, technical aesthetic of Nuxt module sites (nuxt.com, vueuse.org, content.nuxt.com) — dark-first, minimal, strong typographic hierarchy, prominent code snippets. No decorative images — terminal blocks and icons carry the visual weight.

---

## Tech Stack

| Layer | Tool |
|---|---|
| Framework | [Nuxt 4](https://nuxt.com) |
| UI Library | [Nuxt UI v3](https://ui.nuxt.com) — built on **Tailwind CSS v4** + Reka UI. Write Tailwind classes directly alongside the components. |
| Content / Docs | [Nuxt Content v3](https://content.nuxt.com) — Markdown-driven docs page |
| Fonts | [@nuxt/fonts](https://fonts.nuxt.com) — zero-config web font loading |
| Analytics | [nuxt-gtag](https://github.com/johannschopplich/nuxt-gtag) — Google Analytics / GA4 |
| Icons | [@nuxt/icon](https://github.com/nuxt/icon) |
| Package Manager | **Bun** (replaces Node/npm/pnpm) |
| Deployment | Netlify or Vercel (static export) |

> **On Nuxt UI + Tailwind**: Nuxt UI v3 is built on top of Tailwind CSS v4, so you write regular Tailwind utility classes everywhere — no new syntax to learn. Nuxt UI just adds ready-made components (navbar, cards, buttons, hero sections) that match the Nuxt module site aesthetic out of the box.

---

## Getting Started

### Prerequisites

- [Bun](https://bun.sh) — `curl -fsSL https://bun.sh/install | bash`

### Install & Run

```bash
# Clone the repo
git clone https://github.com/idrewlong/shrinkr-site
cd shrinkr-site

# Install dependencies
bun install

# Start dev server (http://localhost:3000)
bun run dev

# Build for production
bun run build

# Preview production build
bun run preview
```

### Project Structure

```
shrinkr-site/
├── app/
│   ├── pages/
│   │   ├── index.vue              # Landing page (all sections)
│   │   └── docs/
│   │       └── [...slug].vue      # Docs catch-all (Nuxt Content)
│   ├── layouts/
│   │   ├── default.vue            # Full-width landing layout
│   │   └── docs.vue               # 3-column docs layout (sidebar + toc)
│   ├── components/
│   │   ├── AppHeader.vue          # Sticky navbar
│   │   ├── AppFooter.vue          # Footer
│   │   └── TerminalSnippet.vue    # Styled dark terminal code block
│   └── app.vue
├── content/
│   └── docs/
│       ├── 1.introduction.md
│       ├── 2.installation.md
│       ├── 3.interactive-wizard.md
│       ├── 4.cli-mode.md
│       ├── 5.flags-reference.md
│       ├── 6.supported-formats.md
│       ├── 7.how-it-works.md
│       └── 8.build-from-source.md
├── public/
│   ├── logo.svg                   # Brand logo SVG
│   └── og-image.png               # Open Graph 1200×630
├── nuxt.config.ts
└── app.config.ts                  # Theme color tokens
```

---

## Branding & Color Palette

### Brand Direction
**Modern Minimalism** — positions Shrinkr as a digital tool first, emphasizing usability and accessibility. The rounded container resembles an application icon, immediately aligning the brand with software: fast, practical, task-oriented. The gradient orange communicates speed, approachability, and modern SaaS perception.

Brand taglines (use across the site):
- *"Smaller files. Faster workflow. Better performance."*
- *"Smaller Files, Better Flow."*
- *"File Compression Made Simple."*

### Logo Variations
Four marks available from the branding kit — use whichever fits the context:

| Variant | Description | Use |
|---|---|---|
| **Wordmark** | "Shrinkr" text in orange rounded-rectangle badge | Navbar, footer, OG image |
| **App icon (S letterform)** | Rounded square, orange gradient, bold S | Favicon, social avatar |
| **S mark with arrows** | S with conversion arrows (stream/switch/shift) | Hero accent, social posts |
| **Standalone S** | Gradient S, no container | Watermarks, pattern use |

Provide all SVG assets in `/public/` before development begins.

### Colors

All values sourced directly from the official Shrinkr Branding Kit (2026):

| Token | Hex | Usage |
|---|---|---|
| **Main gradient start** | `#EC483E` (coral-orange) | Gradient left/top, logo |
| **Main gradient end** | `#F68C31` (amber-orange) | Gradient right/bottom, logo |
| **Primary accent** | `#F06B34` | Buttons, links, highlights (midpoint) |
| **Dark Accent** | `#0E151D` → `#32445B` | Deep navy — page background, dark sections |
| **Dark Shade** | `#313439` | Card backgrounds, sidebar |
| **Light Accent** | `#EEF9FF` | Light mode background |
| **Light Shade** | `#FFFFFF` | Text on dark, light mode surface |

> **Important:** The site background is **navy blue** (`#0E151D` / `#1A2332`), NOT pure black. This matches the brand kit dark palette and distinguishes it from generic dark sites.

```css
/* Main brand gradient — used on logo, buttons, glow effect */
background: linear-gradient(135deg, #EC483E 0%, #F68C31 100%);

/* Hero glow — radial behind terminal block */
background: radial-gradient(ellipse at 60% 0%, rgba(240, 107, 52, 0.18) 0%, transparent 65%);

/* Page dark background */
background-color: #0E151D;

/* Card/surface */
background-color: #1A2332;

/* Border */
border-color: #2A3548;
```

```ts
// app.config.ts — custom Tailwind CSS v4 color tokens
export default defineAppConfig({
  ui: {
    colors: {
      primary: 'orange',
      neutral: 'slate',
    }
  }
})
```

### Typography

The logo wordmark uses a custom **rounded display font**. For the site, use:

```ts
// nuxt.config.ts
fonts: {
  families: [
    { name: 'Inter', provider: 'google' },       // body text
    { name: 'JetBrains Mono', provider: 'google' }, // terminal/code blocks
  ]
}
```

Alternatively, **Geist** (Vercel's font) closely matches the clean technical aesthetic of Nuxt module sites.

---

## Nuxt Config Reference

```ts
// nuxt.config.ts
export default defineNuxtConfig({
  modules: [
    '@nuxt/ui',
    '@nuxt/content',
    '@nuxt/fonts',
    '@nuxt/icon',
    'nuxt-gtag',
  ],
  gtag: {
    id: 'G-XXXXXXXXXX',  // fill in GA4 measurement ID
  },
  app: {
    head: {
      title: 'Shrinkr — Fast Image Compression CLI',
      meta: [
        { name: 'description', content: 'Batch compress JPG, PNG, TIFF, WebP, HEIC, AVIF and GIF images to a target file size. Interactive wizard or direct CLI flags. macOS.' },
        { property: 'og:title', content: 'Shrinkr' },
        { property: 'og:description', content: 'Fast image compression & format conversion CLI tool for macOS.' },
        { property: 'og:image', content: '/og-image.png' },
      ]
    }
  }
})
```

---

## Page 1 — Landing (`/`)

Layout reference: nuxt.com homepage, vueuse.org

---

### Navbar

- Logo (left, SVG) — links to `/`
- Nav links center or right: **Docs** · **GitHub** (with star count badge)
- Dark mode toggle (rightmost)
- Sticky on scroll with `backdrop-blur` + subtle border-bottom

---

### Hero Section

Centered vertical layout. Large typographic headline, subheadline, styled terminal block, two CTAs below. Subtle dot-grid background + orange radial glow behind the terminal block.

**Headline:**
```
Compress smarter.
Ship faster.
```

**Subheadline:**
```
Shrinkr batch-compresses your images to a target file size using a binary search
algorithm — WebP, AVIF, JPEG, PNG. One command. No config required.
```

**Install command (styled terminal block):**
```bash
brew install idrewlong/tap/shrinkr
```

**CTAs:**
- Primary button: `Read the Docs` → `/docs`
- Secondary button (outline): `View on GitHub →` → `https://github.com/idrewlong/shrinkr_cli`

---

### Features Grid

3-column card grid (collapses to 1 on mobile). Each card: icon top-left, title, one-sentence description. Cards have a dark surface background, subtle border, and `hover:border-orange-500/50` transition.

| Icon | Title | Description |
|---|---|---|
| `i-heroicons-bolt` | Blazing fast | Powered by libvips — one of the fastest image processing libraries on the planet |
| `i-heroicons-target` | Hit your target | Binary search compression finds the highest quality that still fits under your KB limit |
| `i-heroicons-arrow-path` | Format conversion | Convert any input to WebP, AVIF, JPEG, or PNG in a single pass |
| `i-heroicons-sparkles` | Interactive wizard | Run `shrinkr` with no args to launch a guided step-by-step terminal wizard |
| `i-heroicons-folder-open` | Batch processing | Process entire folders — or scan subdirectories recursively — with one command |
| `i-heroicons-adjustments-horizontal` | Full control | Fine-tune quality floors, ceilings, target sizes, and worker count via flags |

---

### How It Works

Numbered step list, horizontal on desktop (timeline-style), stacked on mobile. Each step has a number badge, title, and 1-sentence description.

```
Step 1 — Scan
Shrinkr finds every supported image in your folder
(JPG, PNG, TIFF, WebP, GIF, AVIF, HEIF/HEIC).

Step 2 — Target
You choose the output format and maximum file size in KB.
The defaults (WebP, 500 KB) work great for most cases.

Step 3 — Compress
A binary search algorithm tries compression qualities between
your min and max thresholds to find the highest quality that
still meets the target size — up to 10 iterations per image.

Step 4 — Output
Compressed images land in a separate output folder,
preserving the original directory structure.
All images are processed concurrently using a worker pool.
```

---

### Quick Start Section

Two-column layout. Left: headline + copy. Right: `TerminalSnippet` component with three tabbed or sequentially-shown commands.

**Left copy:**
```
Up and running in 30 seconds.

Install via Homebrew — libvips is handled automatically.
Then point Shrinkr at any folder and get back
optimized images. No configuration file required.
```

**Right terminal:**
```bash
# Install (Homebrew handles libvips automatically)
brew install idrewlong/tap/shrinkr

# Compress a folder to WebP at 500 KB target (default)
shrinkr ./photos

# Custom: AVIF, 300 KB target, custom output folder
shrinkr ./photos -f avif -s 300 -o ./compressed
```

---

### Presets Callout (optional section / or fold into features)

In the interactive wizard, Shrinkr ships with three built-in compression presets. Display as a 3-column comparison card or simple table:

| Preset | Target Size | Quality | Use Case |
|---|---|---|---|
| **Recommended** *(default)* | 500 KB | 85 | General use — balanced quality & size |
| **Web Optimized** | 200 KB | 75 | Web assets, thumbnails, fast loading |
| **High Quality** | 2 MB | 95 | Portfolios, print prep, archival |
| **Custom** | You choose | You choose | Full control over all settings |

---

### Supported Formats

Two-column table styled as a card. Add a note that GIF, AVIF, and HEIF/HEIC are supported for input but not output.

| Input Formats | Output Formats |
|---|---|
| JPG / JPEG | WebP |
| PNG | AVIF |
| TIFF / TIF | JPEG |
| WebP | PNG |
| GIF | — |
| AVIF | — |
| HEIF / HEIC | — |

> GIF, AVIF, and HEIF/HEIC can be read and recompressed, but the output format is always one of the four output types above.

---

### Footer

- Logo (small) + tagline: *"Fast image compression for macOS"*
- Links: Docs · GitHub · Issues
- `MIT License · Built with libvips`
- Fill in: city/location line

---

## Page 2 — Documentation (`/docs`)

Layout reference: content.nuxt.com, ui.nuxt.com/getting-started

Three-column layout using Nuxt Content + Nuxt UI:
- **Left sidebar** — auto-generated from content file navigation metadata
- **Center** — rendered Markdown
- **Right sidebar** — `UContentToc` (auto-generated from headings)

Page route: `app/pages/docs/[...slug].vue`

```vue
<script setup lang="ts">
definePageMeta({ layout: 'docs' })
const route = useRoute()
const { data: page } = await useAsyncData(route.path, () =>
  queryCollection('docs').path(route.path).first()
)
const { data: surround } = await useAsyncData(`${route.path}-surround`, () =>
  queryCollectionItemSurroundings('docs', route.path)
)
</script>

<template>
  <UPage>
    <UPageHeader :title="page.title" :description="page.description" />
    <UPageBody>
      <ContentRenderer :value="page" />
      <USeparator />
      <UContentSurround :surround="surround" />
    </UPageBody>
    <template #right>
      <UContentToc :links="page.body.toc.links" />
    </template>
  </UPage>
</template>
```

---

### Docs Sidebar Navigation

```
Getting Started
  ├── Introduction
  └── Installation

Usage
  ├── Interactive Wizard
  ├── Direct CLI Mode
  └── Flags Reference

Reference
  ├── Supported Formats
  ├── How It Works
  └── Updates

Contributing
  └── Build from Source
```

---

### Full Doc Content

All content below should be written verbatim into the markdown files. This gives the AI building the site everything it needs without looking elsewhere.

---

#### `1.introduction.md`

```markdown
---
title: Introduction
description: What Shrinkr is and who it's for.
navigation:
  title: Introduction
---

# Introduction

Shrinkr is a macOS CLI tool for batch-compressing images and converting them to modern web formats. It uses a **binary search algorithm** to hit a target file size while keeping image quality as high as possible — no manual trial-and-error required.

## Who it's for

- **Developers** who need optimized assets before deploying a site
- **Photographers** cleaning up a shoot before delivery
- **Content teams** managing image-heavy pages and blogs
- **Anyone** who wants to shrink a folder of photos without opening Photoshop

## What it does

Given a folder of images and a target file size (in KB), Shrinkr:

1. Scans the folder for all supported image files
2. Compresses each image using a binary search to find the optimal quality setting
3. Converts the output to your chosen format (WebP, AVIF, JPEG, or PNG)
4. Saves the results to a separate output folder — originals are never touched
5. Runs all compressions concurrently using a worker pool sized to your CPU core count

## Requirements

- macOS (Apple Silicon or Intel via Rosetta 2)
- `libvips` — installed automatically when using Homebrew
```

---

#### `2.installation.md`

```markdown
---
title: Installation
description: Three ways to install Shrinkr on macOS.
navigation:
  title: Installation
---

# Installation

## Homebrew (recommended)

The easiest method. Homebrew handles the `libvips` dependency automatically.

```bash
brew install idrewlong/tap/shrinkr
```

After install, verify it works:

```bash
shrinkr --help
```

## Go Install

Requires Go 1.21+ and `libvips` installed separately.

```bash
# Install libvips first
brew install vips

# Then install Shrinkr
go install github.com/idrewlong/shrinkr_cli@latest
```

## Build from Source

See [Build from Source](/docs/build-from-source) for the full walkthrough.

## Updates

```bash
brew upgrade shrinkr
```
```

---

#### `3.interactive-wizard.md`

```markdown
---
title: Interactive Wizard
description: Run shrinkr with no arguments to launch the step-by-step wizard.
navigation:
  title: Interactive Wizard
---

# Interactive Wizard

Running `shrinkr` with no arguments launches an interactive terminal wizard that guides you through all settings step by step.

```bash
shrinkr
```

## Wizard Steps

### Step 1 — Select image folder

Shrinkr scans your current directory and nearby sibling folders for images. If any are found, they appear as quick-select options. You can also browse for a folder via the native macOS Finder dialog, or type a path manually.

If no folders are detected nearby, you'll see two options:
- **Browse for folder…** — opens a Finder dialog
- **Enter folder path manually…** — type or paste a path

> **Tip:** You can drag a folder from Finder into the terminal window and paste its path.

### Step 2 — Output format

Choose the format for your compressed images:

| Option | Notes |
|---|---|
| **WebP** *(recommended)* | Best balance of quality and size for web |
| **AVIF** | Smaller files, slower to encode |
| **JPEG** | Universal compatibility |
| **PNG** | Lossless — file size depends on image content |

### Step 3 — Compression preset

| Preset | Target Size | Initial Quality | Quality Range | Use Case |
|---|---|---|---|---|
| **Recommended** | 500 KB | 85 | 60–90 | General use |
| **Web Optimized** | 200 KB | 75 | 50–80 | Web assets, thumbnails |
| **High Quality** | 2 MB | 95 | 85–98 | Portfolios, archival |
| **Custom** | Your choice | Your choice | Your choice | Full control |

### Step 4 — Custom settings (if Custom preset selected)

If you chose Custom, you'll be asked to enter:

- **Target file size** (KB) — images will be compressed to fit under this size
- **Initial quality** (1–100) — first attempt before binary search kicks in
- **Min quality** (1–100) — compression will never go below this floor
- **Max quality** (1–100) — compression will never exceed this ceiling
- **Worker count** — defaults to your CPU core count (auto-detected)

### Step 5 — Output folder

Choose where to save compressed images:
- **compressed/** *(default)* — created in the current working directory
- **Browse for a custom location** — opens Finder to pick a parent directory, then name the output folder

### Step 6 — Confirm

A summary is shown before running:

```
  Input:    ./photos
  Images:   142 found
  Format:   webp
  Target:   500 KB
  Quality:  85 (range 60–90)
  Workers:  10
  Output:   compressed
```

Confirm with **Let's go!** or go back with **Cancel**.

## Keyboard Controls

| Key | Action |
|---|---|
| `↑` / `↓` | Navigate options |
| `Enter` | Confirm selection |
| `Esc` or `Ctrl+C` | Go back one step |
| `Esc` or `Ctrl+C` at first step | Exit the program |
```

---

#### `4.cli-mode.md`

```markdown
---
title: Direct CLI Mode
description: Pass a folder path as an argument to skip the wizard.
navigation:
  title: Direct CLI Mode
---

# Direct CLI Mode

Pass the input folder as a positional argument to run Shrinkr directly — no wizard.

## Syntax

```bash
shrinkr <folder> [flags]
```

## Examples

```bash
# Compress a folder to WebP at 500 KB (all defaults)
shrinkr ./photos

# Convert to AVIF at 300 KB target
shrinkr ./photos -f avif -s 300

# Custom output location
shrinkr ./photos -f webp -s 500 -o ./compressed-photos

# Recursive scan — include subfolders
shrinkr ~/Pictures -r -f webp -s 200

# All options at once
shrinkr ./photos -f avif -s 300 -q 80 --min-quality 55 --max-quality 85 -w 8 -o ./out

# JPEG output at high quality, large target
shrinkr ./raw-exports -f jpeg -s 2048 -q 95 --min-quality 85 --max-quality 98
```

## Output

Shrinkr prints a live progress bar during compression, then a per-file summary and overall stats when done:

```
  Found 142 image(s)
  Target: 500 KB  |  Format: webp  |  Workers: 10

  ████████████████████░░░░ 83%

  ✓ photo_001.jpg   4.2 MB → 487 KB   (-88%)
  ✓ photo_002.jpg   3.8 MB → 412 KB   (-89%)
  ...

  142 images  |  142 succeeded  |  0 failed
  Total saved: 412 MB → 61 MB  (-85%)
  Time: 8.3s
```

Original files are never modified. Output goes to the folder specified by `--output`.
```

---

#### `5.flags-reference.md`

```markdown
---
title: Flags Reference
description: All CLI flags, short forms, defaults, and descriptions.
navigation:
  title: Flags Reference
---

# Flags Reference

All flags are optional. Shrinkr has sensible defaults for every setting.

## Full Table

| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--format` | `-f` | string | `webp` | Output format: `webp`, `avif`, `jpeg`, `png` |
| `--size` | `-s` | int | `500` | Target file size in KB |
| `--output` | `-o` | string | `compressed` | Output folder path |
| `--quality` | `-q` | int | `85` | Initial quality (1–100) for first compression attempt |
| `--workers` | `-w` | int | CPU count | Number of concurrent workers (auto-detected) |
| `--recursive` | `-r` | bool | `false` | Also scan subdirectories |
| `--min-quality` | — | int | `60` | Quality floor — compression won't go below this |
| `--max-quality` | — | int | `90` | Quality ceiling — compression won't exceed this |

## Notes

- `--workers` defaults to `runtime.NumCPU()` — the number of logical CPU cores on your machine. On Apple Silicon M-series chips this is typically 8–12. You rarely need to change this.
- `--min-quality` and `--max-quality` define the search range for the binary search algorithm. A narrower range means fewer iterations but may not hit the target size.
- `--quality` is the starting point. If the image is already under the target size at this quality, Shrinkr stops immediately (no binary search needed).
- `--recursive` (`-r`) scans all subdirectories. The output folder mirrors the input directory structure.
- `--min-quality` must be less than or equal to `--max-quality`, or Shrinkr will exit with an error.

## Format Notes

| Format | Characteristics |
|---|---|
| `webp` | Best default. Excellent quality-to-size ratio. Supported by all modern browsers. |
| `avif` | Smaller files than WebP but slower to encode. Good for batch jobs where speed isn't critical. |
| `jpeg` | Maximum compatibility. Use for clients or platforms that don't support WebP. |
| `png` | Lossless. File size depends entirely on image complexity — target size may not be reachable for complex images. |
```

---

#### `6.supported-formats.md`

```markdown
---
title: Supported Formats
description: Which image formats Shrinkr can read and write.
navigation:
  title: Supported Formats
---

# Supported Formats

Shrinkr can **read** a wide range of image formats and **write** to four modern output formats.

## Input Formats

| Format | Extensions |
|---|---|
| JPEG | `.jpg`, `.jpeg` |
| PNG | `.png` |
| TIFF | `.tiff`, `.tif` |
| WebP | `.webp` |
| GIF | `.gif` |
| AVIF | `.avif` |
| HEIF / HEIC | `.heif`, `.heic` |

## Output Formats

| Format | Flag value | Notes |
|---|---|---|
| WebP | `webp` | Recommended for web. Lossy by default. |
| AVIF | `avif` | Smallest files. Slower encode. |
| JPEG | `jpeg` | Universal compatibility. |
| PNG | `png` | Lossless. May not reach aggressive KB targets. |

## Notes

- GIF, AVIF, and HEIF/HEIC images can be compressed and converted to any of the four output formats — they cannot be used as output formats themselves.
- Shrinkr uses [`libvips`](https://libvips.github.io/libvips/) via the [`bimg`](https://github.com/h2non/bimg) Go wrapper for all image operations.
- Files with unsupported extensions are silently skipped.
```

---

#### `7.how-it-works.md`

```markdown
---
title: How It Works
description: The binary search compression algorithm and worker pool architecture.
navigation:
  title: How It Works
---

# How It Works

## Binary Search Compression

Shrinkr uses a **binary search algorithm** to find the highest image quality setting that still produces a file under your target size.

### Algorithm walkthrough

1. **Initial attempt** — Compress at your starting quality (default: 85). If the result is already under the target size, done. No binary search needed.

2. **Binary search** — If the initial result is too large, Shrinkr searches between `--min-quality` and `--max-quality`:
   - Try the midpoint quality
   - If the file fits: record it as the best result, try higher quality
   - If the file is too large: try lower quality
   - Repeat up to 10 iterations

3. **Fallback** — If no quality in the range produces a file under the target size (e.g., a PNG with fine detail that can't be compressed enough), Shrinkr falls back to `--min-quality` and saves that result.

This approach guarantees the **best possible quality** for each image at your chosen size constraint, without manual tuning.

### Example

```
Target: 500 KB  |  Initial quality: 85  |  Range: 60–90

Attempt 1  quality=85  →  1.2 MB  (too large)
Binary search:
  quality=75  →  620 KB  (too large)
  quality=67  →  480 KB  ✓ (fits! try higher)
  quality=71  →  550 KB  (too large)
  quality=69  →  495 KB  ✓ (fits! try higher)
  quality=70  →  520 KB  (too large)

Best result: quality=69  →  495 KB
```

## Worker Pool

Shrinkr processes images **concurrently** using a goroutine-based worker pool:

- Pool size defaults to `runtime.NumCPU()` — the number of logical CPU cores on your machine
- Each worker takes a job from the queue, compresses the image, and returns the result
- All jobs are submitted before any results are collected, maximizing throughput
- A live progress bar updates as results come in
- Results are collected and printed after the progress bar completes

This means on a 10-core Apple Silicon machine, Shrinkr compresses 10 images simultaneously. For a folder of 100 images, the total time is roughly `(time per image) × (images / cores)` rather than sequential.

## Technology

Shrinkr is built with:

- **[libvips](https://libvips.github.io/libvips/)** — one of the fastest image processing libraries available. Handles decoding, encoding, and format conversion.
- **[bimg](https://github.com/h2non/bimg)** — a Go wrapper around libvips with idiomatic Go APIs.
- **[huh](https://github.com/charmbracelet/huh)** — a terminal forms library for the interactive wizard UI.
- **[cobra](https://github.com/spf13/cobra)** — CLI flag parsing and command structure.
- **Go** — compiled to a single native binary. No runtime dependencies beyond libvips.
```

---

#### `8.build-from-source.md`

```markdown
---
title: Build from Source
description: Clone the repository and build Shrinkr locally.
navigation:
  title: Build from Source
---

# Build from Source

## Requirements

- macOS (Apple Silicon or Intel — Intel builds run via Rosetta 2)
- Go 1.21+
- `libvips` (installed via Homebrew)

## Steps

```bash
# Clone
git clone https://github.com/idrewlong/shrinkr_cli
cd shrinkr_cli

# Install libvips dependency
brew install vips

# Build binary in current directory
make build

# Or install globally to $GOPATH/bin
make install
```

The `Makefile` targets:

| Target | Command | Description |
|---|---|---|
| `build` | `make build` | Builds `./shrinkr` binary in the project directory |
| `install` | `make install` | Runs `go install .` — installs to `$GOPATH/bin` |
| `clean` | `make clean` | Removes the `./shrinkr` binary and `dist/` folder |
| `run` | `make run` | Runs `go run . ../images` for local testing |
| `release-local` | `make release-local` | Tests GoReleaser config locally without publishing |
| `tag-release` | `make tag-release V=1.2.0` | Tags and pushes a release to trigger GitHub Actions |

## Releasing a New Version

Releases are automated via GitHub Actions + GoReleaser. Pushing a version tag triggers the pipeline:

```bash
make tag-release V=1.2.0
```

This:
1. Creates an annotated git tag `v1.2.0`
2. Pushes it to GitHub
3. GitHub Actions builds the macOS ARM64 binary
4. GoReleaser publishes the release and updates the Homebrew formula automatically
```

---

## Design Notes

Reference sites: nuxt.com, vueuse.org, content.nuxt.com

- **Dark mode first** — navy `#0E151D` background (not black), matching the Shrinkr brand kit dark palette. Light mode uses `#EEF9FF` / `#FFFFFF`.
- **Background** — deep navy `#0E151D` for page, `#1A2332` for cards/surfaces, `#2A3548` for borders
- **Hero glow** — `radial-gradient(ellipse at 60% 0%, rgba(240, 107, 52, 0.18) 0%, transparent 65%)` — orange bloom from top-right behind the terminal block
- **Dot grid** — subtle `radial-gradient(circle, #2A3548 1px, transparent 1px)` tiled over the hero section (navy dots on navy bg — very subtle)
- **Brand gradient** — `linear-gradient(135deg, #EC483E 0%, #F68C31 100%)` — used on primary buttons, badges, accents, and the logo
- **Terminal blocks** — `TerminalSnippet.vue`: background `#0A1018` (darker navy), `font-mono` (JetBrains Mono), border `#2A3548`, top bar with three circles, `$` prompt in orange
- **Feature cards** — `bg-[#1A2332] border border-[#2A3548] rounded-xl p-6 hover:border-[#F06B34]/40 transition-colors`
- **Buttons** — Primary: orange gradient background, white text. Secondary: transparent with `border-[#2A3548]` outline.
- **Spacing** — sections use `py-24 sm:py-32`; generous whitespace
- **Typography** — hero title `text-5xl sm:text-7xl font-bold tracking-tight text-white`; section titles `text-3xl sm:text-4xl`; body `text-[#A8BCCF]` (muted blue-gray, not gray-gray — stays on-brand with the navy)
- **Brand S mark** — can be used as a subtle watermark or section divider (low opacity, large, background layer)
- **No stock imagery** — terminal blocks, format tables, and the S mark carry all visual interest
- **Docs sidebar active state** — left border `border-l-2 border-[#F06B34]`, text brightens to white
- **Mobile nav** — hamburger at `md` breakpoint, slides in from right with navy background

---

## Deliverables Checklist

- [ ] Fill in GA4 measurement ID in `nuxt.config.ts`
- [ ] Add logo SVG to `public/logo.svg`
- [ ] Create OG image `public/og-image.png` (1200×630)
- [ ] `app/pages/index.vue` — all landing sections
- [ ] `app/pages/docs/[...slug].vue` — docs route
- [ ] `app/layouts/docs.vue` — 3-column docs layout
- [ ] All 8 markdown content files in `content/docs/`
- [ ] `AppHeader.vue` with mobile nav
- [ ] `AppFooter.vue`
- [ ] `TerminalSnippet.vue`
- [ ] Dark/light mode toggle working
- [ ] Fully responsive (mobile → desktop)
- [ ] Deployed to Netlify or Vercel
