package mapquest

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

const (
	GeocodingPrefix  = "geocoding"
	GeocodingVersion = "v1"
)

// GeocodingAPI enables users to request geocoding searches via the
// MapQuest API. See https://developer.mapquest.com/documentation/open/geocoding-api/ for details.
// Batch API is not implemented, 5 point queries are not supported
type GeocodingAPI struct {
	c *Client
}

func (api *GeocodingAPI) SimpleAddress(location string, limit int) (*GeocodeAddressResponse, error) {
	if limit < 0 {
		limit = 0
	}
	return api.Address(&GeocodeAddressRequest{Location: location, Limit: limit})
}

func (api *GeocodingAPI) Address(req *GeocodeAddressRequest) (*GeocodeAddressResponse, error) {
	q, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	q.Set("key", api.c.key)
	q.Set("outFormat", "json")
	u := apiURL(GeocodingPrefix, GeocodingVersion, "address")
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
	defer httpResponse.Body.Close()

	var res *GeocodeAddressResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (api *GeocodingAPI) SimpleReverse(lat, long float64) (*GeocodeAddressResponse, error) {
	return api.Reverse(&GeocodeReverseRequest{Location: &GeoPoint{Latitude: lat, Longitude: long}})
}

func (api *GeocodingAPI) Reverse(req *GeocodeReverseRequest) (*GeocodeAddressResponse, error) {
	q, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	q.Set("key", api.c.key)
	q.Set("outFormat", "json")
	u := apiURL(GeocodingPrefix, GeocodingVersion, "reverse")
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
	defer httpResponse.Body.Close()

	var res *GeocodeAddressResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

type GeocodeAddressRequest struct {
	Location           string       `url:"location"`
	BoundingBox        *BoundingBox `url:"boundingBox,omitempty"`
	IgnoreLatLongInput bool         `url:"ignoreLatLngInput,omitempty"`
	ThumbMaps          bool         `url:"thumbMaps"` // dont omit, omitempty works on false, default is true though
	Limit              int          `url:"maxResults,omitempty"`

	// fields ignored:
	// - delimiter => csv output only
	// - intlMode => feel free to implement
}

type GeocodeAddressResponse struct {
	Info *struct {
		StatusCode int `json:"statuscode,omitempty"` // https://developer.mapquest.com/documentation/geocoding-api/status-codes
		Copyright  *struct {
			Text         string `json:"text,omitempty"`
			ImageURL     string `json:"imageUrl,omitempty"`
			ImageAltText string `json:"imageAltText,omitempty"`
		} `json:"copyright,omitempty"`
		Messages []string `json:"messages,omitempty"`
	} `json:"info,omitempty"`
	Options *struct {
		MaxResults         int  `json:"maxResults,omitempty"`
		ThumbMaps          bool `json:"thumbMaps"` // dont omit, omitempty works on false, default is true though
		IgnoreLatLongInput bool `json:"ignoreLatLngInput,omitempty"`
	} `json:"options,omitempty"`

	Results []*GeocodeAddressResponseEntry `json:"results,omitempty"`
}

type GeocodeAddressResponseEntry struct {
	ProvidedLocation *struct {
		Location string    `json:"location,omitempty"`
		LatLong  *GeoPoint `json:"latLng,omitempty"`
	} `json:"providedLocation,omitempty"` // this needs probably 5point support
	Locations []*GeocodeAddressResponseLocationEntry `json:"location,omitempty"`
}

type GeocodeType string

const (
	GeocodeTypeStop GeocodeType = "s"
	GeocodeTypeVia              = "v"
)

type GeocodeAddressResponseLocationEntry struct {
	LatLong        *GeoPoint `json:"latLng,omitempty"`
	DisplayLatLong *GeoPoint `json:"displayLatLng,omitempty"`
	MapURL         string    `json:"mapUrl,omitempty"`

	Street         string      `json:"street,omitempty"`
	PostalCode     string      `json:"postalCode,omitempty"`
	Type           GeocodeType `json:"type,omitempty"`
	AdminArea6     string      `json:"adminArea6,omitempty"`
	AdminArea6Type string      `json:"adminArea6Type,omitempty"`
	AdminArea5     string      `json:"adminArea5,omitempty"`
	AdminArea5Type string      `json:"adminArea5Type,omitempty"`
	AdminArea4     string      `json:"adminArea4,omitempty"`
	AdminArea4Type string      `json:"adminArea4Type,omitempty"`
	AdminArea3     string      `json:"adminArea3,omitempty"`
	AdminArea3Type string      `json:"adminArea3Type,omitempty"`
	AdminArea2     string      `json:"adminArea2,omitempty"`
	AdminArea2Type string      `json:"adminArea2Type,omitempty"`
	AdminArea1     string      `json:"adminArea1,omitempty"`
	AdminArea1Type string      `json:"adminArea1Type,omitempty"`

	GeocodeQuality     string `json:"geocodeQuality,omitempty"` // https://developer.mapquest.com/documentation/geocoding-api/quality-codes
	GeocodeQualityCode string `json:"geocodeQualityCode,omitempty"`

	UnkownInput string `json:"unkownInput,omitempty"`

	RoadMetadata *struct {
		SpeedLimitUnits string            `json:"speedLimitUnits,omitempty"`
		TollRoad        []json.RawMessage `json:"TollRoad,omitempty"` // unkown data type, can be nullable
		SpeedLimit      int               `json:"speedLimit,omitempty"`
	} `json:"roadMetadata,omitempty"`

	NearestIntersection *struct {
		StreetDisplayName string    `json:"streetDisplayName,omitempty"`
		DistanceMeters    string    `json:"distanceMeters,omitempty"`
		LatLng            *GeoPoint `json:"latLng,omitempty"`
		Label             string    `json:"label,omitempty"`
	} `json:"nearestIntersection,omitempty"`

	// ignored fields
	// - dragPoint => matters only for dragroute calls
	// - linkId => String that identifies the closest road to the address for routing purposes.
	// - sideOfStreet => Specifies the side of street. (Left, Right, Mixed, None)
}

type GeocodeReverseRequest struct {
	Location                   *GeoPoint `url:"location"`
	ThumbMaps                  bool      `url:"thumbMaps"` // dont omit, omitempty works on false, default is true though
	IncludeNearestIntersection bool      `url:"includeNearestIntersection,omitempty"`
	IncludeRoadMetadata        bool      `url:"includeRoadMetadata,omitempty"`
}
