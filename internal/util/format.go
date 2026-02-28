package util

import (
	"fmt"
	"math"
)

// FormatBytes converts bytes to human-readable string (e.g., "4.80 MB").
// Ported from compressor.js formatBytes function.
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 Bytes"
	}

	units := []string{"Bytes", "KB", "MB", "GB", "TB"}
	k := 1024.0
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(k)))
	if i >= len(units) {
		i = len(units) - 1
	}

	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(k, float64(i)), units[i])
}
