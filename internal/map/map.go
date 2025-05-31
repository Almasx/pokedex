package mappkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/almasx/pokedexcli/internal/api"
	"github.com/almasx/pokedexcli/internal/cli"
)


const BASE_URL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

func fetchMapData(url string, config *cli.Config) (api.GetLocationAreas, error) {
	res := api.GetLocationAreas{}

	if data, ok := config.Cache.Get(url); ok {
		err := json.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	
	data, err := http.Get(url)
	if err != nil {
		return api.GetLocationAreas{}, err
	}
	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return api.GetLocationAreas{}, err
	}
	config.Cache.Add(url, body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return api.GetLocationAreas{}, err
	}
	return res, nil
}

func CommandMap(config *cli.Config, args []string) error {
	mapData := api.GetLocationAreas{}
	url := config.Next
	if url == "" {
		url = BASE_URL
	}
	
	mapData, err := fetchMapData(url, config)
	if err != nil {
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	
	config.Next = mapData.Next
	config.Prev = mapData.Previous

	return nil
}

func CommandMapb(config *cli.Config, args []string) error {
	mapData := api.GetLocationAreas{}
	url := config.Prev
	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	
	mapData, err := fetchMapData(url, config)
	if err != nil {
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	
	config.Prev = mapData.Previous
	config.Next = mapData.Next

	return nil
}
