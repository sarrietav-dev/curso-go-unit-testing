package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
	c := require.New(t)

	pokemon, err := GetPokemonFromPokeApi("bulbasaur")
	c.NoError(err)

	body, err := ioutil.ReadFile("samples/poke_api_readed.json")
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(pokemon, expected)
}

func TestGetPokemonFromPokeApiSuccessWithMocks(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "bulbasaur"

	body, err := ioutil.ReadFile("samples/poke_api_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder(http.MethodGet, request, httpmock.NewBytesResponder(200, body))

	pokemon, err := GetPokemonFromPokeApi(id)
	c.NoError(err)

	var expectedResponse models.PokeApiPokemonResponse
	json.Unmarshal([]byte(body), &expectedResponse)

	c.Equal(expectedResponse, pokemon)
}

func TestGetPokemonFromPokeApiServerInternalError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "bulbasaur"

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder(http.MethodGet, request, httpmock.NewStringResponder(500, ""))

	_, err := GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(err, ErrPokeApiFailure.Error())
}

func TestGetPokemonFromPokeApiNotFound(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "najsndkjnaksjdnkjan"

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder(http.MethodGet, request, httpmock.NewStringResponder(404, ""))

	_, err := GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(err, ErrPokemonNotFound.Error())
}
