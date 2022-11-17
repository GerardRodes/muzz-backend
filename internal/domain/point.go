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

type mariadbGeo uint32

const (
	mariadbGeoPoint mariadbGeo = iota + 1
	mariadbGeoLineString
	mariadbGeoPolygon
	mariadbGeoMultiPoint
	mariadbGeoMultiLineString
	mariadbGeoMultiPolygon
	mariadbGeoGeometryCollection
)

// Scan parses a geographical MariaDB Point:
// https://mariadb.com/kb/en/well-known-binary-wkb-format/
func (g *Point) Scan(src any) error {
	data, ok := src.([]uint8)
	if !ok {
		return fmt.Errorf("expected []byte for Point, got  %T", src)
	}

	// It is not explained in the docs but it is preceded by 4bytes
	// containing the SRID, which makes a total of 25 bytes
	if len(data) != 25 {
		return fmt.Errorf("expected []bytes with length 25, got %d", len(data))
	}

	// we don't care about SRID
	// would need it if we wanted to load the
	// value again to MariaDB
	data = data[4:]

	// The first byte indicates the byte order. 00 for big endian, or 01 for little endian.
	var order binary.ByteOrder
	if data[0] == 1 {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}
	data = data[1:]

	var geoType uint32
	if err := binary.Read(bytes.NewReader(data[:4]), order, &geoType); err != nil {
		return err
	}
	data = data[4:]

	if mariadbGeo(geoType) != mariadbGeoPoint {
		return fmt.Errorf("expected geo type Point %q but found %q", mariadbGeoPoint, geoType)
	}

	// first goes lng, then lat
	// Point(x, y) = Point(lng, lat)
	if err := binary.Read(bytes.NewReader(data[:8]), order, &g.Lng); err != nil {
		return err
	}
	data = data[8:]

	if err := binary.Read(bytes.NewReader(data[:8]), order, &g.Lat); err != nil {
		return err
	}

	return nil
}
