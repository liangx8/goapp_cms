package view

import (
	"encoding/json"
	"io"

	"rcgreed.bid/ics/mgr"
)

func JsonView(d any) mgr.View {
	return func(w io.Writer) error {
		en := json.NewEncoder(w)
		if err := en.Encode(d); err != nil {
			return err
		}
		return nil
	}
}
