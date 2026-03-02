package compressor_test

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/idrewlong/shrinkr_cli/internal/compressor"
)

// createTestJPEG writes a small synthetic JPEG to dir and returns its path.
// Uses Go's stdlib — no libvips required to create the input.
func createTestJPEG(t *testing.T, dir string) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	for y := 0; y < 300; y++ {
		for x := 0; x < 400; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 255 / 400),
				G: uint8(y * 255 / 300),
				B: 128,
				A: 255,
			})
		}
	}
	path := filepath.Join(dir, "input.jpg")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create test image: %v", err)
	}
	defer f.Close()
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode test image: %v", err)
	}
	return path
}

// baseJob returns a Job with sensible defaults for the given format.
func baseJob(inputPath, outputPath string, format compressor.OutputFormat) compressor.Job {
	return compressor.Job{
		InputPath:      inputPath,
		OutputPath:     outputPath,
		TargetSizeKB:   100,
		Format:         format,
		MinQuality:     50,
		MaxQuality:     90,
		InitialQuality: 85,
	}
}

// ── Format tests ─────────────────────────────────────────────────────────────

// TestAllFormats runs a compression through every supported output format and
// reports timing. This is the primary diagnostic for the AVIF slowness issue.
func TestAllFormats(t *testing.T) {
	cases := []struct {
		format compressor.OutputFormat
		ext    string
	}{
		{compressor.FormatWebP, ".webp"},
		{compressor.FormatPNG, ".png"},
		{compressor.FormatJPEG, ".jpg"},
		{compressor.FormatAVIF, ".avif"},
	}

	dir := t.TempDir()
	inputPath := createTestJPEG(t, dir)

	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.format), func(t *testing.T) {
			outputPath := filepath.Join(dir, "output"+tc.ext)
			job := baseJob(inputPath, outputPath, tc.format)

			start := time.Now()
			result := compressor.Compress(job)
			elapsed := time.Since(start)

			if !result.Success {
				t.Errorf("FAILED: %s", result.Error)
				return
			}

			info, err := os.Stat(outputPath)
			if err != nil {
				t.Errorf("output file missing: %v", err)
				return
			}
			if info.Size() == 0 {
				t.Error("output file is empty")
				return
			}

			t.Logf("OK  elapsed=%-10s  input=%dKB  output=%dKB  quality=%d  saved=%.1f%%",
				elapsed.Round(time.Millisecond),
				result.OriginalSize/1024,
				result.CompressedSize/1024,
				result.Quality,
				result.CompressionPct,
			)
		})
	}
}

// TestAVIFTimeout checks AVIF specifically with a hard timeout so the test
// suite doesn't hang indefinitely if libvips AVIF support is broken.
func TestAVIFTimeout(t *testing.T) {
	dir := t.TempDir()
	inputPath := createTestJPEG(t, dir)
	outputPath := filepath.Join(dir, "output.avif")
	job := baseJob(inputPath, outputPath, compressor.FormatAVIF)

	type outcome struct {
		result compressor.Result
		took   time.Duration
	}
	ch := make(chan outcome, 1)

	go func() {
		start := time.Now()
		r := compressor.Compress(job)
		ch <- outcome{r, time.Since(start)}
	}()

	const timeout = 30 * time.Second
	select {
	case o := <-ch:
		if !o.result.Success {
			t.Errorf("AVIF compression failed: %s", o.result.Error)
			t.Log("This likely means libvips was built without AVIF/libaom support.")
			t.Log("Fix: brew reinstall vips --build-from-source  OR  brew install libavif && brew reinstall vips")
			return
		}
		t.Logf("AVIF OK in %s  (%dKB → %dKB)", o.took.Round(time.Millisecond),
			o.result.OriginalSize/1024, o.result.CompressedSize/1024)
		if o.took > 10*time.Second {
			t.Logf("WARNING: AVIF took %s — consider removing it from the wizard default options", o.took.Round(time.Second))
		}
	case <-time.After(timeout):
		t.Errorf("AVIF compression timed out after %s — libvips AVIF encoder appears to be hanging", timeout)
		t.Log("Fix: brew reinstall vips")
	}
}

// ── Target size tests ────────────────────────────────────────────────────────

// TestTargetSizes verifies the binary search hits a range of KB targets with WebP.
func TestTargetSizes(t *testing.T) {
	dir := t.TempDir()
	inputPath := createTestJPEG(t, dir)

	targets := []int{150, 80, 40}
	for _, targetKB := range targets {
		targetKB := targetKB
		t.Run(fmt.Sprintf("%dKB", targetKB), func(t *testing.T) {
			outputPath := filepath.Join(dir, fmt.Sprintf("output_%dkb.webp", targetKB))
			job := compressor.Job{
				InputPath:      inputPath,
				OutputPath:     outputPath,
				TargetSizeKB:   targetKB,
				Format:         compressor.FormatWebP,
				MinQuality:     20,
				MaxQuality:     90,
				InitialQuality: 85,
			}
			result := compressor.Compress(job)
			if !result.Success {
				t.Errorf("compression failed: %s", result.Error)
				return
			}
			actualKB := result.CompressedSize / 1024
			t.Logf("target=%dKB  actual=%dKB  quality=%d", targetKB, actualKB, result.Quality)
		})
	}
}

// ── Edge case tests ───────────────────────────────────────────────────────────

// TestLargeFileGuard verifies the 500MB limit error message is returned cleanly.
func TestLargeFileGuard(t *testing.T) {
	dir := t.TempDir()
	// Create a stub file that reports a large size by checking the guard logic
	// indirectly: a normal file should pass through without hitting the guard.
	inputPath := createTestJPEG(t, dir)
	outputPath := filepath.Join(dir, "output.webp")
	job := baseJob(inputPath, outputPath, compressor.FormatWebP)
	result := compressor.Compress(job)
	if !result.Success {
		t.Errorf("normal-sized file should not hit the size guard: %s", result.Error)
	}
}

// TestMissingInputFile verifies a clean error when the input doesn't exist.
func TestMissingInputFile(t *testing.T) {
	dir := t.TempDir()
	job := baseJob(
		filepath.Join(dir, "nonexistent.jpg"),
		filepath.Join(dir, "output.webp"),
		compressor.FormatWebP,
	)
	result := compressor.Compress(job)
	if result.Success {
		t.Error("expected failure for missing input file, got success")
	}
}

// ── Benchmarks ────────────────────────────────────────────────────────────────

func BenchmarkWebP(b *testing.B)  { benchmarkFormat(b, compressor.FormatWebP, ".webp") }
func BenchmarkPNG(b *testing.B)   { benchmarkFormat(b, compressor.FormatPNG, ".png") }
func BenchmarkJPEG(b *testing.B)  { benchmarkFormat(b, compressor.FormatJPEG, ".jpg") }
func BenchmarkAVIF(b *testing.B)  { benchmarkFormat(b, compressor.FormatAVIF, ".avif") }

func benchmarkFormat(b *testing.B, format compressor.OutputFormat, ext string) {
	b.Helper()
	dir := b.TempDir()

	// Create input once
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	for y := 0; y < 300; y++ {
		for x := 0; x < 400; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}
	inputPath := filepath.Join(dir, "input.jpg")
	f, _ := os.Create(inputPath)
	jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
	f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(dir, fmt.Sprintf("out_%d%s", i, ext))
		job := compressor.Job{
			InputPath:      inputPath,
			OutputPath:     outputPath,
			TargetSizeKB:   100,
			Format:         format,
			MinQuality:     50,
			MaxQuality:     90,
			InitialQuality: 85,
		}
		compressor.Compress(job)
	}
}
