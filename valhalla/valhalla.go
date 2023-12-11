package valhalla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gps_api/model"
	"io"
	"net/http"
	"os"
)

type Node struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Edge struct {
	PercentAlong  float64 `json:"percent_along"`
	CorrelatedLon float64 `json:"correlated_lon"`
	SideOfStreet  string  `json:"side_of_street"`
	CorrelatedLat float64 `json:"correlated_lat"`
	WayID         int     `json:"way_id"`
}

type Response struct {
	InputLon float64 `json:"input_lon"`
	InputLat float64 `json:"input_lat"`
	Nodes    []Node  `json:"nodes"`
	Edges    []Edge  `json:"edges"`
}

func GetEdge(response *Response) *model.MovementCreate {

	for _, edge := range response.Edges {
		if edge.CorrelatedLon != 0 && edge.CorrelatedLat != 0 {
			return &model.MovementCreate{Latitude: edge.CorrelatedLat, Longitude: edge.CorrelatedLon}
		}
	}
	return nil
}

func SendLocateRequest(lat float64, lon float64) *Response {
	bodyRequest := []byte(fmt.Sprintf(`{
		"locations": [{"lat": %f, "lon": %f, "radius": 1, "search_cutoff": 10}],
		"costing": "auto"
	} `, lat, lon))

	req, err := http.NewRequest("POST", os.Getenv("VALHALLASERVER")+"/locate", bytes.NewBuffer(bodyRequest))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var inputs []Response

	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(responseBody, &inputs)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	fmt.Println(string(responseBody))
	return &inputs[0]
}
