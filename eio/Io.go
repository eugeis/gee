package eio

import "io"

type MapWriter interface {
	WriteMap(data map[string]interface{}) error
}

type WriteCloserMapWriter struct {
	Convert func(map[string]interface{}) (io.Reader, error)
	Out     io.WriteCloser
}

func (o *WriteCloserMapWriter) WriteMap(data map[string]interface{}) (err error) {
	var reader io.Reader
	if reader, err = o.Convert(data); err == nil {
		_, err = io.Copy(o.Out, reader)
	}
	return
}

type CollectMapWriter struct {
	Data [] map[string]interface{}
}

func NewCollectMapWriter() *CollectMapWriter {
	return &CollectMapWriter{Data: make([]map[string]interface{}, 0)}
}

func (o *CollectMapWriter) WriteMap(data map[string]interface{}) (err error) {
	o.Data = append(o.Data, data)
	return
}
