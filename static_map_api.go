package mapquest

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"strings"

	// static map returns a raw gif, jpeg or png object
	// to bring support out of the box, we require them here
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/go-querystring/query"
)

const (
	StaticMapPrefix  = "staticmap"
	StaticMapVersion = "v5"
)

// StaticMapAPI enables users to request static map images via the
// MapQuest API. See http://open.mapquestapi.com/staticmap/ for details.
type StaticMapAPI struct {
	c *Client
}

func (api *StaticMapAPI) Map(req *StaticMapRequest) (image.Image, error) {
	reader, err := api.MapReader(req)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	return img, err
}

func (api *StaticMapAPI) MapReader(req *StaticMapRequest) (io.ReadCloser, error) {
	q, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	q.Set("key", api.c.key)
	u := apiURL(StaticMapPrefix, StaticMapVersion, "map")
	u.RawQuery = q.Encode()

	httpRequest, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", UserAgent)
	httpResponse, err := api.c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return httpResponse.Body, nil
}

type StaticMapSize struct {
	Width  int
	Height int
	Retina bool
}

func (s *StaticMapSize) EncodeValues(key string, v *url.Values) error {
	if s.Width > 1920 {
		return ErrDimensionToLarge
	}
	if s.Height > 1920 {
		return ErrDimensionToLarge
	}

	size := ""

	if s.Width > 0 && s.Height > 0 {
		size = fmt.Sprintf("%d,%d", s.Width, s.Height)
	}

	if s.Retina {
		size += "@2"
	}

	v.Set(key, size)
	return nil
}

type StaticMapGeoPoint struct {
	Latitude  float64
	Longitude float64
}

func (s *StaticMapGeoPoint) String() string {
	return fmt.Sprintf("%f,%f", s.Latitude, s.Longitude)
}

func (s *StaticMapGeoPoint) EncodeValues(key string, v *url.Values) error {
	v.Set(key, s.String())
	return nil
}

type StaticMapBoundingBox struct {
	TopLeft     StaticMapGeoPoint
	BottomRight StaticMapGeoPoint
}

func (s *StaticMapBoundingBox) EncodeValues(key string, v *url.Values) error {
	v.Set(key, fmt.Sprintf("%v,%v", s.TopLeft.String(), s.BottomRight.String()))
	return nil
}

type StaticMapFormat string

const (
	StaticMapFormatPNG   StaticMapFormat = "png"
	StaticMapFormatGIF                   = "gif"
	StaticMapFormatJPEG                  = "jpeg"
	StaticMapFormatJPG                   = "jpg"
	StaticMapFormatJPG70                 = "jpg70"
	StaticMapFormatJPG80                 = "jpg80"
	StaticMapFormatJPG90                 = "jpg90"
)

type StaticMapType string

const (
	StaticMapTypeDark      StaticMapType = "dark"
	StaticMapTypeLight                   = "light"
	StaticMapTypeMap                     = "map"
	StaticMapTypeHybrid                  = "hyb"
	StaticMapTypeSatellite               = "sat"
)

type StaticMapScalebar struct {
	Enable   bool
	Position string
}

func (s *StaticMapScalebar) EncodeValues(key string, v *url.Values) error {
	if s.Enable {
		val := "true"
		if s.Position != "" {
			val += "|" + s.Position
		}
		v.Set(key, val)
	} else {
		v.Set(key, "false")
	}

	return nil
}

type StaticMapLocation struct {
	Location string
	Marker   string
}

func (s *StaticMapLocation) String() string {
	str := s.Location
	if s.Marker != "" {
		str += "|" + s.Marker
	}
	return str
}
func (s *StaticMapLocation) EncodeValues(key string, v *url.Values) error {
	v.Set(key, s.String())
	return nil
}

type StaticMapLocations []StaticMapLocation

func (s StaticMapLocations) EncodeValues(key string, v *url.Values) error {
	if len(s) > 0 {
		locs := make([]string, len(s))
		for i, l := range s {
			locs[i] = l.String()
		}

		v.Set(key, strings.Join(locs, "||"))
	} else {
		v.Set(key, "")
	}

	return nil
}

type StaticMapBannerSize string

const (
	StaticMapBannerSizeSmall  StaticMapBannerSize = "sm"
	StaticMapBannerSizeMedium                     = "md"
	StaticMapBannerSizeLarge                      = "lg"
)

type StaticMapBanner struct {
	Text            string
	Size            StaticMapBannerSize
	OnTop           bool
	TextColor       int // find a way to set 0x000000 without a helper function maybe?
	BackgroundColor int
}

func (s *StaticMapBanner) appendOption(line string, option string, hasOption bool) (string, bool) {
	if hasOption {
		line = line + "-" + option
	} else {
		line = line + "|" + option
	}

	return line, true
}

func (s *StaticMapBanner) EncodeValues(key string, v *url.Values) error {
	val := s.Text
	hasOptions := false
	hasTextColor := false

	if s.Size != "" {
		val, hasOptions = s.appendOption(val, string(s.Size), hasOptions)
	}

	// bottom is default
	if s.OnTop {
		val, hasOptions = s.appendOption(val, "top", hasOptions)
	}

	if s.TextColor > 0 {
		val, hasOptions = s.appendOption(val, fmt.Sprintf("%06x", s.TextColor), hasOptions)
		hasTextColor = true
	}

	if s.BackgroundColor > 0 {
		if !hasTextColor {
			val, hasOptions = s.appendOption(val, "ffffff", hasOptions)
		}

		val, hasOptions = s.appendOption(val, fmt.Sprintf("%06x", s.TextColor), hasOptions)
	}

	v.Set(key, val)
	return nil
}

type StaticMapColor struct {
	R int
	G int
	B int
	A int
}

func (s *StaticMapColor) EncodeValues(key string, v *url.Values) error {
	if s.A == 0 {
		v.Set(key, fmt.Sprint("%d,%d,%d", s.R, s.G, s.B))
	} else {
		v.Set(key, fmt.Sprint("%d,%d,%d,%d", s.R, s.G, s.B, s.A))
	}

	return nil
}

func StaticMapColorHex(h int) *StaticMapColor {
	c := new(StaticMapColor)

	c.R = (h >> 16) & 0xff
	c.G = (h >> 8) & 0xff
	c.B = (h) & 0xff

	return c
}

func StaticMapColorHexAlpha(h int) *StaticMapColor {
	c := StaticMapColorHex((h >> 8) & 0xffffff)
	c.A = h & 0xff
	return c
}

type StaticMapShape struct {
}

type StaticMapRequest struct {
	Size        *StaticMapSize        `url:"size,omitempty"`
	Center      string                `url:"center,omitempty"`
	BoundingBox *StaticMapBoundingBox `url:"boundingBox,omitempty"`
	Margin      int                   `url:"margin,omitempty"`
	Zoom        int                   `url:"zoom,omitempty"`
	Format      StaticMapFormat       `url:"format,omitempty"`
	Type        StaticMapType         `url:"type,omitempty"`
	Scalebar    *StaticMapScalebar    `url:"scalebar,omitempty"`

	// additional location options
	Locations     StaticMapLocations `url:"locations,omitempty"`
	Declutter     bool               `url:"declutter,omitempty"`
	DefaultMarker string             `url:"defaultMarker,omitempty"`

	// banner
	Banner *StaticMapBanner `url:"banner,omitempty"`

	// routes
	Start *StaticMapLocation `url:"start,omitempty"`
	End   *StaticMapLocation `url:"end,omitempty"`
	// TODO: https://developer.mapquest.com/documentation/open/static-map-api/v5/map/#request_parameters-session
	RotueArc   bool            `url:"routeArc,omitempty"`
	RouteWidth int             `url:"routeWidth,omitempty"`
	RouteColor *StaticMapColor `url:"routeColor,omitempty"`

	// shapemagics
	// TODO: https://developer.mapquest.com/documentation/open/static-map-api/v5/map/#request_parameters-shape
}
