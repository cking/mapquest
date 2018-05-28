package mapquest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	NominatimPrefix  = "nominatim"
	NominatimVersion = "v1"
)

// NominatimAPI enables users to request nominatim searches via the
// MapQuest API. See https://developer.mapquest.com/documentation/open/nominatim-search/ for details.
type NominatimAPI struct {
	c *Client
}

func (api *NominatimAPI) SimpleSearch(query string, limit int) (*NominatimSearchResponse, error) {
	if limit < 0 {
		limit = 0
	}
	return api.Search(&NominatimSearchRequest{Query: query, Limit: limit})
}

func (api *NominatimAPI) Search(req *NominatimSearchRequest) (*NominatimSearchResponse, error) {
	q, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	q.Set("key", api.c.key)
	q.Set("format", "json")
	u := apiURL(NominatimPrefix, NominatimVersion, "search.api")
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

	var res *NominatimSearchResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (api *NominatimAPI) SimpleReverse(lat, long float64) (*NominatimSearchResponseEntry, error) {
	return api.Reverse(&NominatimReverseRequest{Latitude: lat, Longitude: long})
}

func (api *NominatimAPI) Reverse(req *NominatimReverseRequest) (*NominatimSearchResponseEntry, error) {
	q, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	q.Set("key", api.c.key)
	q.Set("format", "json")
	u := apiURL(NominatimPrefix, NominatimVersion, "search.api")
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

	var res *NominatimSearchResponseEntry
	if err := json.NewDecoder(httpResponse.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

type NominatimViewBox struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (s *NominatimViewBox) EncodeValues(key string, v *url.Values) error {
	v.Set(key, fmt.Sprintf("%f,%f,%f,%f", s.Left, s.Top, s.Right, s.Bottom))
	return nil
}

type NominatimOSMType string

const (
	OSMTypeNode     NominatimOSMType = "N"
	OSMTypeWay                       = "W"
	OSMTypeRelation                  = "R"
)

type NominatimSearchRequest struct {
	Query           string            `url:"q"`
	AddressDetails  bool              `url:"addressdetails,omitempty"`
	Limit           int               `url:"limit,omitempty"`
	CountryCodes    []string          `url:"countrycodes,comma,omitempty"`
	ViewBox         *NominatimViewBox `url:"viewbox,omitempty"` // let,top,right,bottom => todo
	ExcludePlaceIDs []string          `url:"exclude_place_ids,comma,omitempty"`
	RouteWidth      float64           `url:"routewidth,omitempty"`
	OSMType         NominatimOSMType  `url:"osm_type,omitempty"`
	OSMID           string            `url:"osm_id,omitempty"`
}

type NominatimSearchResponse struct {
	Results []*NominatimSearchResponseEntry
}

type NominatimSearchResponseEntry struct {
	Address *struct {
		City          string `json:"city,omitempty"`
		CityDistrict  string `json:"city_district,omitempty"`
		Continent     string `json:"continent,omitempty"`
		Country       string `json:"country,omitempty"`
		CountryCode   string `json:"country_code,omitempty"`
		County        string `json:"county,omitempty"`
		Hamlet        string `json:"hamlet,omitempty"`
		HouseNumber   string `json:"house_number,omitempty"`
		Pedestrian    string `json:"pedestrian,omitempty"`
		Neighbourhood string `json:"neighbourhood,omitempty"`
		PostCode      string `json:"postcode,omitempty"`
		Road          string `json:"road,omitempty"`
		State         string `json:"state,omitempty"`
		StateDistrict string `json:"state_district,omitempty"`
		Suburb        string `json:"suburb,omitempty"`
	} `json:"address,omitempty"`
	BoundingBox []float64 `json:"boundingbox,omitempty"`
	Class       string    `json:"class,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	Importance  float64   `json:"importance,omitempty"`
	Latitude    float64   `json:"lat,string,omitempty"`
	Longitude   float64   `json:"lon,string,omitempty"`
	OSMId       string    `json:"osm_id,omitempty"`
	OSMType     string    `json:"osm_type,omitempty"`
	PlaceID     string    `json:"place_id,omitempty"`
	Type        string    `json:"type,omitempty"`
	License     string    `json:"licence,omitempty"` // typo in API
	Icon        string    `json:"icon,omitempty"`
}

type NominatimReverseRequest struct {
	Latitude  float64          `url:"lat"`
	Longitude float64          `url:"long"`
	OSMType   NominatimOSMType `url:"osm_type,omitempty"`
	OSMID     string           `url:"osm_id,omitempty"`
}
