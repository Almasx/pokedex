package cli

import (
	"github.com/almasx/pokedexcli/internal/api"
	"github.com/almasx/pokedexcli/internal/pokecache"
)

type Config struct {
	Next string
	Prev string
	Cache *pokecache.Cache
	Pokedex map[string]api.GetPokemon
}
