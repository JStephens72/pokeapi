package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Map(index *PageTracker) ([]string, error) {
	if index.Next == "" {
		return nil, fmt.Errorf("End of List")
	}
	s := index.Previous.(string)
	rawData, err := getData(s)
	if err != nil {
		return nil, err
	}
	decodedData, err := decodeData(rawData)
	if err != nil {
		return nil, err
	}

	areas := extractLocationNames(decodedData)

	// now update the links in the tracker
	index.Previous = decodedData.Previous
	index.Current = index.Next
	index.Next = decodedData.Next

	return areas, nil
}

func Mapb(index *PageTracker) ([]string, error) {
	s := index.Previous.(string)
	if s == "" {
		return nil, fmt.Errorf("Start of List")
	}

	rawData, err := getData(s)
	if err != nil {
		return nil, err
	}
	decodedData, err := decodeData(rawData)
	if err != nil {
		return nil, err
	}

	areas := extractLocationNames(decodedData)

	// now update the links in the tracker
	index.Previous = decodedData.Previous
	index.Current = index.Next
	index.Next = decodedData.Next

	return areas, nil
}
func getData(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error retrieving locations: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s\n", res.StatusCode, data)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w\n", err)
	}
	return data, nil
}

func decodeData(data []byte) (Location, error) {
	var locationData Location
	if err := json.Unmarshal(data, &locationData); err != nil {
		return locationData, fmt.Errorf("error decoding data: %w", err)
	}
	return locationData, nil
}

func extractLocationNames(data Location) []string {
	areas := []string{}
	for _, result := range data.Results {
		areas = append(areas, result.Name)
	}
	return areas
}
