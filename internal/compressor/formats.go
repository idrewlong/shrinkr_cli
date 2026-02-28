package compressor

import (
	"fmt"
	"strings"

	"github.com/h2non/bimg"
)

// OutputFormat represents a supported output image format.
type OutputFormat string

const (
	FormatWebP OutputFormat = "webp"
	FormatAVIF OutputFormat = "avif"
	FormatJPEG OutputFormat = "jpeg"
	FormatPNG  OutputFormat = "png"
)

// ParseFormat validates and returns an OutputFormat from a string.
func ParseFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "webp":
		return FormatWebP, nil
	case "avif":
		return FormatAVIF, nil
	case "jpeg", "jpg":
		return FormatJPEG, nil
	case "png":
		return FormatPNG, nil
	default:
		return "", fmt.Errorf("unsupported output format: %s (supported: webp, avif, jpeg, png)", s)
	}
}

// FileExtension returns the file extension for the format (e.g., ".webp").
func (f OutputFormat) FileExtension() string {
	switch f {
	case FormatWebP:
		return ".webp"
	case FormatAVIF:
		return ".avif"
	case FormatJPEG:
		return ".jpg"
	case FormatPNG:
		return ".png"
	default:
		return ".webp"
	}
}

// BimgType returns the bimg.ImageType for encoding.
func (f OutputFormat) BimgType() bimg.ImageType {
	switch f {
	case FormatWebP:
		return bimg.WEBP
	case FormatAVIF:
		return bimg.AVIF
	case FormatJPEG:
		return bimg.JPEG
	case FormatPNG:
		return bimg.PNG
	default:
		return bimg.WEBP
	}
}

// EncodeOptions returns bimg.Options configured for the given quality and format.
// Mirrors the Node.js Sharp settings: effort 6, smartSubsample for WebP.
func (f OutputFormat) EncodeOptions(quality int) bimg.Options {
	opts := bimg.Options{
		Quality: quality,
		Type:    f.BimgType(),
	}

	// Per-format tuning to match/exceed the Node.js Sharp settings
	switch f {
	case FormatWebP:
		// Sharp equivalent: { quality, effort: 6, smartSubsample: true }
		// bimg doesn't expose effort/smartSubsample directly, but libvips
		// uses good defaults. Quality is the primary control.
	case FormatAVIF:
		// AVIF: quality controls compression, speed is handled by libvips
	case FormatJPEG:
		// JPEG: strip metadata for smaller files
		opts.StripMetadata = true
	case FormatPNG:
		// PNG is lossless — quality maps to compression level
		// Higher compression = smaller but slower
		opts.Compression = pngCompressionFromQuality(quality)
	}

	return opts
}

// pngCompressionFromQuality maps a quality value (1-100) to PNG compression level (0-9).
// Lower quality = higher compression for PNG.
func pngCompressionFromQuality(quality int) int {
	// Invert: quality 100 = compression 0 (fastest, largest)
	//         quality 1   = compression 9 (slowest, smallest)
	level := 9 - (quality * 9 / 100)
	if level < 0 {
		level = 0
	}
	if level > 9 {
		level = 9
	}
	return level
}
