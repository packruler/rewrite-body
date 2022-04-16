package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
)

type ReaderError struct {
	error

	cause error
}

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

func Encode(data []byte, encoding string) (readCloser io.ReadCloser, err error, encoded bool) {
	byteBuffer := new(bytes.Buffer)
	compressor, err := getCompressor(byteBuffer, encoding)
	if err != nil {
		// If an error creating the compressor occurs set to nil to allow uncompressed data to be returned
		compressor = nil
	}

	if compressor != nil {
		compressor.Write(data)

		if err = compressor.Close(); err == nil {
			return io.NopCloser(byteBuffer), nil, true
		}
	}

	return io.NopCloser(bytes.NewReader(data)), err, false
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
