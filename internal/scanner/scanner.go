package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// SupportedExtensions maps all recognized input image extensions.
// .jpg, .jpeg, .png, .tiff, .tif, .webp, .gif, .avif, .heif, .heic
var SupportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".tiff": true,
	".tif":  true,
	".webp": true,
	".gif":  true,
	".avif": true,
	".heif": true,
	".heic": true,
}

// IsSupportedImage checks if a filename has a supported image extension.
func IsSupportedImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return SupportedExtensions[ext]
}

// FindImages walks the inputFolder and returns paths to all supported images.
// If recursive is false, only the top-level directory is scanned.
func FindImages(inputFolder string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.WalkDir(inputFolder, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && IsSupportedImage(path) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(inputFolder)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				fullPath := filepath.Join(inputFolder, entry.Name())
				if IsSupportedImage(fullPath) {
					files = append(files, fullPath)
				}
			}
		}
	}

	return files, nil
}
