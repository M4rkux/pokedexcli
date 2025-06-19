package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const base_url = "https://pokeapi.co/api/v2/"

type Language struct {
	id       int
	name     string
	official bool
	iso639   string
	names    []Name
}

type Name struct {
	name     string
	language Language
}

type EncounterMethod struct {
	id    int
	name  string
	order int
	names []Name
}

type Generation struct{}

type MoveLearnMethod struct{}

type Pokedex struct{}

type VersionGroup struct {
	id                 int
	name               string
	order              int
	generation         Generation
	move_learn_methods []MoveLearnMethod
	pokedexes          []Pokedex
	regions            []struct{}
	versions           []Version
}

type Version struct {
	id            int
	name          string
	names         []Name
	version_group VersionGroup
}

type EncounterVersionDetails struct {
	rate    int
	version Version
}

type EncounterMethodRate struct {
	encounter_method []EncounterMethod
	version_details  []EncounterVersionDetails
}

type LocationArea struct {
	Id                   int      `json:"id"`
	Name                 string   `json:"name"`
	GameIndex            int      `json:"game_index"`
	EncounterMethodRates struct{} `json:"encounter_method_rates"`
}

type ListLocationAreas struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func GetListLocationAreas(paramUrl string) (ListLocationAreas, error) {
	var url string
	if paramUrl != "" {
		url = paramUrl
	} else {
		url = base_url + "location-area/"
	}

	resp, err := http.Get(url)
	if err != nil {
		return ListLocationAreas{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ListLocationAreas{}, errors.New(fmt.Sprintf("Request failed with status: %v", resp.Status))
	}

	var locationAreas ListLocationAreas
	err = json.NewDecoder(resp.Body).Decode(&locationAreas)

	if err != nil {
		return ListLocationAreas{}, err
	}

	return locationAreas, nil
}
