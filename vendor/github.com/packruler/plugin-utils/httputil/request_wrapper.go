package httputil

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/packruler/plugin-utils/compressutil"
	"github.com/packruler/plugin-utils/logger"
)

// RequestWrapper a struct that centralizes request modifications.
type RequestWrapper struct {
	logWriter  logger.LogWriter
	monitoring MonitoringConfig

	http.Request
}

// WrapRequest to get a new instance of RequestWrapper.
func WrapRequest(request *http.Request, monitoringConfig MonitoringConfig, logWriter logger.LogWriter) *RequestWrapper {
	return &RequestWrapper{
		logWriter:  logWriter,
		monitoring: monitoringConfig,
		Request:    *request,
	}
}

// CloneNoEncode create an http.Request that request no encoding.
func (req *RequestWrapper) CloneNoEncode() *http.Request {
	clonedRequest := req.Clone(req.Context())

	clonedRequest.Header.Set("Accept-Encoding", compressutil.Identity)

	return clonedRequest
}

// CloneWithSupportedEncoding create an http.Request that request only supported encoding.
func (req *RequestWrapper) CloneWithSupportedEncoding() *http.Request {
	clonedRequest := req.Clone(req.Context())

	clonedRequest.Header.Set("Accept-Encoding", removeUnsupportedAcceptEncoding(clonedRequest.Header))

	return clonedRequest
}

// GetEncodingTarget get the supported encoding algorithm preferred by request.
func (req *RequestWrapper) GetEncodingTarget() string {
	// Limit Accept-Encoding header to encodings we can handle.
	// acceptEncoding := header.ParseAccept(req.Header, "Accept-Encoding")
	acceptEncoding := parseAcceptEncoding(req.Header)
	filteredEncodings := make([]encodingSpec, 0, len(acceptEncoding))

	for _, a := range acceptEncoding {
		switch a.Value {
		case compressutil.Gzip, compressutil.Deflate:
			filteredEncodings = append(filteredEncodings, a)
		}
	}

	if len(filteredEncodings) == 0 {
		return compressutil.Identity
	}

	sort.Slice(filteredEncodings, func(i, j int) bool {
		return filteredEncodings[i].Quality > filteredEncodings[j].Quality
	})

	return filteredEncodings[0].Value
}

type encodingSpec struct {
	Value   string
	Quality float64
}

func parseAcceptEncoding(header http.Header) []encodingSpec {
	encodingHeader := header.Get("Accept-Encoding")
	if encodingHeader == "*" {
		return []encodingSpec{{Quality: 1.0, Value: compressutil.Gzip}, {Quality: 1.0, Value: compressutil.Deflate}}
	}

	encodingList := strings.Split(encodingHeader, ",")
	result := make([]encodingSpec, 0, len(encodingList))

	for _, encoding := range encodingList {
		result = append(result, parseEncodingItem(encoding))
	}

	return result
}

func parseEncodingItem(encoding string) encodingSpec {
	encoding = strings.TrimSpace(encoding)
	if encoding == "*" {
		return encodingSpec{Value: compressutil.Gzip, Quality: 1.0}
	}

	split := strings.Split(encoding, ";q=")
	quality := 1.0

	if qualitySplitSize := 2; len(split) == qualitySplitSize {
		targetFloat := 64

		parsedQuality, err := strconv.ParseFloat(split[1], targetFloat)
		if err == nil {
			quality = parsedQuality
		}
	}

	return encodingSpec{Value: split[0], Quality: quality}
}

func removeUnsupportedAcceptEncoding(header http.Header) string {
	encodingList := strings.Split(header.Get("Accept-Encoding"), ",")
	result := make([]string, 0, len(encodingList))

	for _, encoding := range encodingList {
		split := strings.Split(strings.TrimSpace(encoding), ";q=")
		switch split[0] {
		case compressutil.Gzip, compressutil.Deflate, compressutil.Identity:
			result = append(result, encoding)
		}
	}

	return strings.Join(result, ",")
}

// SupportsProcessing determine if http.Request is supported by this plugin.
func (req *RequestWrapper) SupportsProcessing() bool {
	acceptHeader := req.Header.Get("Accept")
	isSupported := false

	for _, monitoredType := range req.monitoring.Types {
		if strings.Contains(acceptHeader, monitoredType) {
			isSupported = true
		}
	}

	if !isSupported {
		return false
	}

	isSupported = false

	// Ignore non GET requests
	for _, monitoredMethod := range req.monitoring.Methods {
		if strings.Contains(req.Method, monitoredMethod) {
			isSupported = true
		}
	}

	if !isSupported {
		return false
	}

	if strings.Contains(req.Header.Get("Upgrade"), "websocket") {
		return false
	}

	return true
}
