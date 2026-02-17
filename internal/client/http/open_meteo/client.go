package open_meteo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	CurrentWeathe struct {
		Temperature float64 `json:"temperature"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *client {
	return &client{
		httpClient: httpClient,
	}
}

func (c *client) GetTemperature(latitude, longitude float64) (Response, error) {
	res, err := c.httpClient.Get(
		fmt.Sprintf("http://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
			latitude,
			longitude,
		),
	)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("status code %d", res.StatusCode)
	}
	var tempResp Response
	if err = json.NewDecoder(res.Body).Decode(&tempResp); err != nil {
		return Response{}, err
	}
	fmt.Println(tempResp)
	return tempResp, nil
}
