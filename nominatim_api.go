package mapquest

import (
	"fmt"
	"log"
	"net/url"

	"github.com/google/go-querystring/query"
)

var _ = log.Print

const (
	// NominatimPathPrefix is the default path prefix for the Nominatim API.
	NominatimPathPrefix = "/nominatim/v1"
)

// NominatimAPI is a geographic search service that relies solely on the
// data contributed to OpenStreetMap.
// See http://open.mapquestapi.com/nominatim/ for details.
type NominatimAPI struct {
	c *Client
}

// Search searches for details given an address.
func (api *NominatimAPI) Search(req *NominatimSearchRequest) (*NominatimSearchResponse, error) {
	u, err := api.buildSearchURL(req)
	if err != nil {
		return nil, err
	}

	res := new(NominatimSearchResponse)
	res.Results = make([]*NominatimSearchResult, 0)

	if err := api.c.getJSON(u, &res.Results); err != nil {
		return nil, err
	}

	return res, nil
}

// buildSearchURL returns the complete URL for the request,
// including the key to query the MapQuest API.
func (api *NominatimAPI) buildSearchURL(req *NominatimSearchRequest) (string, error) {
	urls := fmt.Sprintf("%s%s/search.php", api.c.BaseURL(), NominatimPathPrefix)
	u, err := url.Parse(urls)
	if err != nil {
		return "", err
	}

	// Add key and other parameters to the query string
	if req.Query != "" {
		req.Street = ""
		req.City = ""
		req.County = ""
		req.State = ""
		req.Country = ""
		req.PostalCode = ""
	}
	v, _ := query.Values(req)

	v.Set("format", "json")
	v.Set("addressdetails", "1")
	v.Set("key", api.c.key)
	u.RawQuery = v.Encode()
	return u.String(), nil
}

type NominatimViewBox struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func (vb *NominatimViewBox) EncodeValues(key string, v *url.Values) error {
	v.Set(key, fmt.Sprintf("%f,%f,%f,%f", vb.Left, vb.Top, vb.Right, vb.Bottom))
	return nil
}

type NominatimSearchRequest struct {
	Query           string            `url:"q,omitempty"`
	Street          string            `url:"street,omitempty"`
	City            string            `url:"city,omitempty"`
	County          string            `url:"county,omitempty"`
	State           string            `url:"state,omitempty"`
	Country         string            `url:"country,omitempty"`
	PostalCode      string            `url:"postalcode,omitempty"`
	Limit           int               `url:"limit,omitempty"`
	CountryCodes    []string          `url:"coutrycodes,comma,omitempty"`
	ViewBox         *NominatimViewBox `url:"viewbox,omitempty"`
	ExcludePlaceIds []string          `url:"exclude_place_ids,comma,omitempty"`
	Bounded         *bool             `url:"bounded,omitempty"`
	RouteWidth      *float64          `url:"routewidth,omitempty"`
	OSMType         string            `url:"osm_type,omitempty"`
	OSMId           string            `url:"osm_id,omitempty"`
}

type NominatimSearchResponse struct {
	Results []*NominatimSearchResult
}

type NominatimSearchResult struct {
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
	//BoundingBox []float64 `json:"boundingbox,omitempty"`
	Class       string  `json:"class,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	Importance  float64 `json:"importance,omitempty"`
	Latitude    float64 `json:"lat,string,omitempty"`
	Longitude   float64 `json:"lon,string,omitempty"`
	OSMId       string  `json:"osm_id,omitempty"`
	OSMType     string  `json:"osm_type,omitempty"`
	PlaceId     string  `json:"place_id,omitempty"`
	Type        string  `json:"type,omitempty"`
	License     string  `json:"licence,omitempty"` // typo in API?
}
