package minrevpro

import (
	"net/http"
	"strconv"
)

// InformativeResponseWriter is a simple wrapper of http.ReponseWriter,
// just a litle bit more informative about the status code
type InformativeResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten uint64
}

func NewInformativeResponseWriter(w http.ResponseWriter) *InformativeResponseWriter {
	return &InformativeResponseWriter{ResponseWriter: w}
}

func (rw *InformativeResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	// better save than sorry, check for nil
	if rw.ResponseWriter != nil {
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *InformativeResponseWriter) Write(data []byte) (cnt int, err error) {
	cnt, err = rw.ResponseWriter.Write(data)
	rw.bytesWritten += uint64(cnt)
	return
}

func (rw *InformativeResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *InformativeResponseWriter) StatusCodeString() string {
	return strconv.Itoa(rw.statusCode)
}

func (rw *InformativeResponseWriter) StatusText() string {
	return http.StatusText(rw.statusCode)
}

func (rw *InformativeResponseWriter) BytesWritten() uint64 {
	return rw.bytesWritten
}
