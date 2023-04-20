package controller

import (
	"catching-pokemons/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

func TestGetPokemon(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "bulbasaur"

	body, err := ioutil.ReadFile("samples/poke_api_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder(http.MethodGet, request, httpmock.NewBytesResponder(200, body))

	w, r, err := mockRequestFactory(id)
	c.NoError(err)

	GetPokemon(w, r)

	expectedResponse, err := parseFile[models.Pokemon]("samples/api_response.json")
	c.NoError(err)

	var actualPokemon models.Pokemon
	err = json.Unmarshal(w.Body.Bytes(), &actualPokemon)
	c.NoError(err)

	c.Equal(http.StatusOK, w.Code)
	c.Equal(expectedResponse, actualPokemon)
}

func mockRequestFactory(id string) (*httptest.ResponseRecorder, *http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, "/pokemon/{id}", nil)
	if err != nil {
		return nil, nil, err
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", id)

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w := httptest.NewRecorder()

	return w, r, nil
}

func parseFile[T any](path string) (*T, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var writable T
	err = json.Unmarshal(body, writable)
	if err != nil {
		return nil, err
	}

	return &writable, nil
}
