// Package compressutil a plugin to handle compression and decompression tasks
package compressutil

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
)

// ReaderError for notating that an error occurred while reading compressed data.
type ReaderError struct {
	error

	cause error
}

// Decode data in a bytes.Reader based on supplied encoding.
func Decode(byteReader *bytes.Reader, encoding string) (data []byte, err error) {
	reader, err := getRawReader(byteReader, encoding)
	if err != nil {
		return nil, &ReaderError{cause: err}
	}

	return io.ReadAll(reader)
}

func getRawReader(byteReader *bytes.Reader, encoding string) (io.Reader, error) {
	switch encoding {
	case "gzip":
		return gzip.NewReader(byteReader)

	case "deflate":
		return flate.NewReader(byteReader), nil

	default:
		return byteReader, nil
	}
}

// Encode data in a []byte based on supplied encoding.
func Encode(data []byte, encoding string) (readCloser io.ReadCloser, err error) {
	byteBuffer := new(bytes.Buffer)

	compressor, err := getCompressor(byteBuffer, encoding)
	if err != nil {
		// If an error creating the compressor occurs set to nil to allow uncompressed data to be returned
		compressor = nil
	}

	if compressor != nil {
		_, _ = compressor.Write(data)

		if err = compressor.Close(); err == nil {
			return io.NopCloser(byteBuffer), nil
		}
	}

	return nil, err
}

func getCompressor(byteBuffer *bytes.Buffer, encoding string) (compressor io.WriteCloser, err error) {
	switch encoding {
	case "gzip":
		return gzip.NewWriterLevel(byteBuffer, gzip.DefaultCompression)

	case "deflate":
		return flate.NewWriter(byteBuffer, flate.DefaultCompression)

	default:
		return nil, nil
	}
}
