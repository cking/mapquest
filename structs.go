package mapquest

import (
	"fmt"
	"net/url"
)

type GeoPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func (s *GeoPoint) String() string {
	return fmt.Sprintf("%f,%f", s.Latitude, s.Longitude)
}

func (s *GeoPoint) EncodeValues(key string, v *url.Values) error {
	v.Set(key, s.String())
	return nil
}

type BoundingBox struct {
	TopLeft     GeoPoint
	BottomRight GeoPoint
}

func (s *BoundingBox) EncodeValues(key string, v *url.Values) error {
	v.Set(key, fmt.Sprintf("%v,%v", s.TopLeft.String(), s.BottomRight.String()))
	return nil
}
