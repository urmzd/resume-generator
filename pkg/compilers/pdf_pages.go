package compilers

import "bytes"

// CountPDFPages returns the number of pages in a PDF byte slice.
// It counts occurrences of "/Type /Page" that are NOT "/Type /Pages"
// in the PDF cross-reference objects (which are always uncompressed).
func CountPDFPages(data []byte) int {
	marker := []byte("/Type /Page")
	count := 0
	offset := 0
	for {
		idx := bytes.Index(data[offset:], marker)
		if idx == -1 {
			break
		}
		pos := offset + idx + len(marker)
		// Ensure this is "/Type /Page" and not "/Type /Pages"
		if pos >= len(data) || data[pos] != 's' {
			count++
		}
		offset = pos
	}
	return count
}
