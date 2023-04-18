package util

import (
	"catching-pokemons/models"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	c := require.New(t)

	var pokeApiResponse models.PokeApiPokemonResponse

	unmarshalJson("samples/pokeapi_response.json", &pokeApiResponse, c)

	parsedPokemon, err := ParsePokemon(pokeApiResponse)
	c.NoError(err)

	var pokemon models.Pokemon

	unmarshalJson("samples/api_response.json", pokemon, c)

	c.Equal(parsedPokemon, pokemon)
}

func unmarshalJson(path string, v any, c *require.Assertions) {
	body, err := ioutil.ReadFile(path)
	c.NoError(err)

	err = json.Unmarshal([]byte(body), v)
	c.NoError(err)
}
