package view

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"rcgreed.bid/ics/mgr"
)

func JsonView() mgr.View {
	return func(w io.Writer, d any) error {
		en := json.NewEncoder(w)
		if err := en.Encode(d); err != nil {
			return err
		}
		return nil
	}
}
func StaticPage(path string) mgr.View {
	return func(w io.Writer, d any) error {
		var buf []byte
		var err error
		if buf, err = ioutil.ReadFile(path); err != nil {
			return err
		}
		if _, err = w.Write(buf); err != nil {
			return nil
		}
		return nil

	}
}
