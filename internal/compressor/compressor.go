package compressor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/h2non/bimg"
)

// Job represents a single compression task.
type Job struct {
	InputPath      string
	OutputPath     string
	TargetSizeKB   int
	Format         OutputFormat
	MinQuality     int
	MaxQuality     int
	InitialQuality int
	MaxIterations  int
}

// Result represents the outcome of a compression job.
type Result struct {
	InputPath      string
	OutputPath     string
	Filename       string // relative filename for display
	Success        bool
	OriginalSize   int64
	CompressedSize int64
	Quality        int
	CompressionPct float64 // percentage saved (e.g., 91.28)
	Error          string
}

// Compress performs the binary-search compression for a single image.
// This is a faithful port of the algorithm from compressor.js.
//
// Algorithm:
//  1. Read input, try initial quality
//  2. If under target size, done immediately
//  3. Binary search between minQuality and maxQuality to find the
//     highest quality that still meets the target file size
//  4. Fallback to minQuality if target can never be met
func Compress(job Job) Result {
	if job.MaxIterations == 0 {
		job.MaxIterations = 10 // match Node.js default
	}

	// Guard against loading enormous files into memory
	const maxFileSizeBytes = 500 * 1024 * 1024 // 500 MB
	info, err := os.Stat(job.InputPath)
	if err != nil {
		return errorResult(job, fmt.Sprintf("failed to stat file: %v", err))
	}
	if info.Size() > maxFileSizeBytes {
		return errorResult(job, fmt.Sprintf("file too large (%d MB), skipping — 500 MB limit", info.Size()/1024/1024))
	}

	// Read input file
	inputBytes, err := os.ReadFile(job.InputPath)
	if err != nil {
		return errorResult(job, fmt.Sprintf("failed to read file: %v", err))
	}

	originalSize := int64(len(inputBytes))
	targetSizeBytes := int64(job.TargetSizeKB) * 1024

	// Try initial compression with starting quality
	quality := job.InitialQuality
	buffer, err := encodeImage(inputBytes, quality, job.Format)
	if err != nil {
		return errorResult(job, fmt.Sprintf("compression failed: %v", err))
	}

	// If already under target size, we're done
	if int64(len(buffer)) <= targetSizeBytes {
		if err := writeOutput(job.OutputPath, buffer); err != nil {
			return errorResult(job, fmt.Sprintf("failed to write output: %v", err))
		}
		return successResult(job, originalSize, int64(len(buffer)), quality)
	}

	// Binary search for optimal quality
	low := job.MinQuality
	high := job.MaxQuality
	var bestBuffer []byte
	bestQuality := quality
	iterations := 0

	for low <= high && iterations < job.MaxIterations {
		quality = (low + high) / 2
		buffer, err = encodeImage(inputBytes, quality, job.Format)
		if err != nil {
			return errorResult(job, fmt.Sprintf("compression failed at quality %d: %v", quality, err))
		}

		if int64(len(buffer)) <= targetSizeBytes {
			// Size is good, try higher quality
			bestBuffer = buffer
			bestQuality = quality
			low = quality + 1
		} else {
			// Size too large, reduce quality
			high = quality - 1
		}

		iterations++
	}

	// If we still can't get under target size with minimum quality, use minimum
	if bestBuffer == nil || int64(len(bestBuffer)) > targetSizeBytes {
		bestBuffer, err = encodeImage(inputBytes, job.MinQuality, job.Format)
		if err != nil {
			return errorResult(job, fmt.Sprintf("compression failed at min quality: %v", err))
		}
		bestQuality = job.MinQuality
	}

	// Write the best result
	if err := writeOutput(job.OutputPath, bestBuffer); err != nil {
		return errorResult(job, fmt.Sprintf("failed to write output: %v", err))
	}

	return successResult(job, originalSize, int64(len(bestBuffer)), bestQuality)
}

// encodeImage compresses the input bytes to the specified format and quality.
func encodeImage(inputBytes []byte, quality int, format OutputFormat) ([]byte, error) {
	img := bimg.NewImage(inputBytes)
	opts := format.EncodeOptions(quality)
	return img.Process(opts)
}

// writeOutput creates parent directories and writes the buffer to disk.
func writeOutput(outputPath string, buffer []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, buffer, 0644)
}

func successResult(job Job, originalSize, compressedSize int64, quality int) Result {
	pct := 0.0
	if originalSize > 0 {
		pct = (1.0 - float64(compressedSize)/float64(originalSize)) * 100
	}
	return Result{
		InputPath:      job.InputPath,
		OutputPath:     job.OutputPath,
		Filename:       filepath.Base(job.InputPath),
		Success:        true,
		OriginalSize:   originalSize,
		CompressedSize: compressedSize,
		Quality:        quality,
		CompressionPct: pct,
	}
}

func errorResult(job Job, errMsg string) Result {
	return Result{
		InputPath: job.InputPath,
		Filename:  filepath.Base(job.InputPath),
		Success:   false,
		Error:     errMsg,
	}
}
