package rewrite_body

import (
	"log"

	"github.com/packruler/rewrite-body/compressutil"
)

// func (wrappedWriter *ResponseWrapper) getHeaderContent() (encoding string, contentType string, isSupported bool) {
// 	encoding = wrappedWriter.GetContentEncoding()
// 	contentType = wrappedWriter.GetContentType()

// 	// If content type does not match return values with false
// 	if contentType != "" && !strings.Contains(contentType, "text") {
// 		return encoding, contentType, false
// 	}

// 	// If content type is supported validate encoding as well
// 	switch encoding {
// 	case "gzip":
// 		fallthrough
// 	case "deflate":
// 		fallthrough
// 	case "identity":
// 		fallthrough
// 	case "":
// 		return encoding, contentType, true
// 	default:
// 		return encoding, contentType, false
// 	}
// }

func (wrappedWriter *ResponseWrapper) decompressBody(encoding string) ([]byte, bool) {
	compressed, err := compressutil.Decode(wrappedWriter.GetBuffer(), encoding)
	if err != nil {
		return nil, false
	}

	// switch encoding {
	// case "gzip":
	// 	return getBytesFromGzip(wrappedWriter.GetBuffer())

	// case "deflate":
	// 	return getBytesFromZlib(wrappedWriter.GetBuffer())

	// default:
	// 	return wrappedWriter.GetBuffer().Bytes(), true
	// }
	return compressed, true
}

// func getBytesFromZlib(buffer *bytes.Buffer) ([]byte, bool) {
// 	zlibReader, err := zlib.NewReader(buffer)
// 	if err != nil {
// 		log.Printf("Failed to load body reader: %v", err)

// 		return buffer.Bytes(), false
// 	}

// 	bodyBytes, err := io.ReadAll(zlibReader)
// 	if err != nil {
// 		log.Printf("Failed to read body: %s", err)

// 		return buffer.Bytes(), false
// 	}

// 	err = zlibReader.Close()

// 	if err != nil {
// 		log.Printf("Failed to close reader: %v", err)

// 		return buffer.Bytes(), false
// 	}

// 	return bodyBytes, true
// }

// func getBytesFromGzip(buffer *bytes.Buffer) ([]byte, bool) {
// 	gzipReader, err := gzip.NewReader(buffer)
// 	if err != nil {
// 		log.Printf("Failed to load body reader: %v", err)

// 		return buffer.Bytes(), false
// 	}

// 	bodyBytes, err := io.ReadAll(gzipReader)
// 	if err != nil {
// 		log.Printf("Failed to read body: %s", err)

// 		return buffer.Bytes(), false
// 	}

// 	err = gzipReader.Close()

// 	if err != nil {
// 		log.Printf("Failed to close reader: %v", err)

// 		return buffer.Bytes(), false
// 	}

// 	return bodyBytes, true
// }

func prepareBodyBytes(bodyBytes []byte, encoding string) ([]byte, error) {
	data, err := compressutil.Encode(bodyBytes, encoding)
	if err != nil {
		log.Printf("Unable to encode data: %v", err)

		return nil, err
	}
	// switch encoding {
	// case "gzip":
	// 	return compressWithGzip(bodyBytes)

	// case "deflate":
	// 	return compressWithZlib(bodyBytes)

	// default:
	// 	return bodyBytes
	// }
	return data, nil
}

// func compressWithGzip(bodyBytes []byte) []byte {
// 	var buf bytes.Buffer
// 	gzipWriter := gzip.NewWriter(&buf)

// 	if _, err := gzipWriter.Write(bodyBytes); err != nil {
// 		log.Printf("unable to recompress rewrited body: %v", err)

// 		return bodyBytes
// 	}

// 	if err := gzipWriter.Close(); err != nil {
// 		log.Printf("unable to close gzip writer: %v", err)

// 		return bodyBytes
// 	}

// 	return buf.Bytes()
// }

// func compressWithZlib(bodyBytes []byte) []byte {
// 	var buf bytes.Buffer
// 	zlibWriter := zlib.NewWriter(&buf)

// 	if _, err := zlibWriter.Write(bodyBytes); err != nil {
// 		log.Printf("unable to recompress rewrited body: %v", err)

// 		return bodyBytes
// 	}

// 	if err := zlibWriter.Close(); err != nil {
// 		log.Printf("unable to close zlib writer: %v", err)

// 		return bodyBytes
// 	}

// 	return buf.Bytes()
// }
