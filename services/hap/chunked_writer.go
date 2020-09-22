package hap

import "io"

type chunkedWriter struct {
	wr    io.Writer
	chunk int
}

func NewChunkedWriter(writer io.Writer, chunk int) *chunkedWriter {
	return &chunkedWriter{wr: writer, chunk: chunk}
}

func (w *chunkedWriter) Write(p []byte) (int, error) {
	var max = len(p)
	var nn int
	var end int
	for nn < max {
		end = nn + w.chunk
		if end > max {
			end = max
		}
		n, err := w.wr.Write(p[nn:end])
		if err != nil {
			return nn, err
		}
		nn += n
	}

	return nn, nil
}
