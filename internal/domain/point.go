package domain

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (g *Point) Scan(src interface{}) error {
	switch src := src.(type) {
	case []uint8:
		if len(src) != 25 {
			return fmt.Errorf("expected []bytes with length 25, got %d", len(src))
		}
		var longitude float64
		var latitude float64
		buf := bytes.NewReader(src[9:17])
		err := binary.Read(buf, binary.LittleEndian, &longitude)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(src[17:25])
		err = binary.Read(buf, binary.LittleEndian, &latitude)
		if err != nil {
			return err
		}
		*g = Point{Lng: longitude, Lat: latitude}
	default:
		return fmt.Errorf("expected []byte for Point, got  %T", src)
	}
	return nil
}
