package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/m4rkux/pokedexcli/internal"
)

const base_url = "https://pokeapi.co/api/v2/"

type Language struct {
	id       int
	name     string
	official bool
	iso639   string
	names    []Name
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource         `json:"encounter_method"`
	VersionDetails  []EncounterVersionDetail `json:"version_details"`
}

type EncounterVersionDetail struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type Name struct {
	Language NamedAPIResource `json:"language"`
	Name     string           `json:"name"`
}

type PokemonEncounter struct {
	Pokemon        NamedAPIResource          `json:"pokemon"`
	VersionDetails []PokemonEncounterVersion `json:"version_details"`
}

type PokemonEncounterVersion struct {
	MaxChance        int               `json:"max_chance"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
	Version          NamedAPIResource  `json:"version"`
}

type EncounterDetail struct {
	Chance          int                `json:"chance"`
	ConditionValues []NamedAPIResource `json:"condition_values"`
	MaxLevel        int                `json:"max_level"`
	Method          NamedAPIResource   `json:"method"`
	MinLevel        int                `json:"min_level"`
}

type LocationArea struct {
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	GameIndex            int                   `json:"game_index"`
	ID                   int                   `json:"id"`
	Location             NamedAPIResource      `json:"location"`
	Name                 string                `json:"name"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type ListLocationAreas struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

var cache = internal.NewCache(5 * time.Second)

func GetListLocationAreas(paramUrl string) (ListLocationAreas, error) {
	var url string
	if paramUrl != "" {
		url = paramUrl
	} else {
		url = base_url + "location-area/"
	}

	var locationAreas ListLocationAreas

	if data, ok := cache.Get(url); ok {
		json.Unmarshal(data, &locationAreas)
		return locationAreas, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return ListLocationAreas{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ListLocationAreas{}, fmt.Errorf("Request failed with status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ListLocationAreas{}, err
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return ListLocationAreas{}, err
	}

	return locationAreas, nil
}

func GetListPokemonsInArea(area string) ([]PokemonEncounter, error) {

	url := base_url + "location-area/" + area

	var locationArea LocationArea

	if data, ok := cache.Get(url); ok {
		json.Unmarshal(data, &locationArea)
		return locationArea.PokemonEncounters, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return []PokemonEncounter{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []PokemonEncounter{}, fmt.Errorf("Request failed with status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []PokemonEncounter{}, err
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return []PokemonEncounter{}, err
	}

	return locationArea.PokemonEncounters, nil
}
