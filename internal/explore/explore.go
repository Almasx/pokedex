package explorepkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/almasx/pokedexcli/internal/api"
	"github.com/almasx/pokedexcli/internal/cli"
)



func fetchLocationAreaPokemons(url string, config *cli.Config) (api.GetLocationAreaPokemons, error) {
	res := api.GetLocationAreaPokemons{}

	if data, ok := config.Cache.Get(url); ok {
		err := json.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	
	data, err := http.Get(url)
	if err != nil {
		return api.GetLocationAreaPokemons{}, err
	}
	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return api.GetLocationAreaPokemons{}, err
	}
	config.Cache.Add(url, body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return api.GetLocationAreaPokemons{}, err
	}
	return res, nil
}

func CommandExplore(config *cli.Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("explore requires a location area")
		return fmt.Errorf("explore requires a location area")
	}
	location_area := args[0]
	if location_area == "" {
		fmt.Println("location area is required")
		return fmt.Errorf("location area is required")	
	}

	fmt.Println("Exploring", location_area, "...")

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location_area)
	location_area_pokemons, err := fetchLocationAreaPokemons(url, config)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range location_area_pokemons.PokemonEncounters {
		fmt.Println(" - ", pokemon.Pokemon.Name)
	}

	return nil
}
