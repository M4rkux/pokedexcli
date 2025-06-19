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

type LocationAreaEncounter struct {
	LocationArea   []LocationArea `json:"location_area"`
	VersionDetails interface{}    `json:"version_details"`
}

type ListLocationAreas struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Pokemon struct {
	ID                     int                `json:"id"`
	Name                   string             `json:"name"`
	BaseExperience         int                `json:"base_experience"`
	Height                 int                `json:"height"`
	IsDefault              bool               `json:"is_default"`
	Order                  int                `json:"order"`
	Weight                 int                `json:"weight"`
	Abilities              []Ability          `json:"abilities"`
	Forms                  []NamedAPIResource `json:"forms"`
	GameIndices            []GameIndex        `json:"game_indices"`
	HeldItems              []HeldItem         `json:"held_items"`
	LocationAreaEncounters string             `json:"location_area_encounters"`
	Moves                  []Move             `json:"moves"`
	Species                NamedAPIResource   `json:"species"`
	Sprites                Sprites            `json:"sprites"`
	Cries                  Cries              `json:"cries"`
	Stats                  []Stat             `json:"stats"`
	Types                  []TypeSlot         `json:"types"`
	PastTypes              []PastType         `json:"past_types"`
	PastAbilities          []PastAbility      `json:"past_abilities"`
}

// Abilities
type Ability struct {
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
	Ability  NamedAPIResource `json:"ability"`
}

// Game indices
type GameIndex struct {
	GameIndex int              `json:"game_index"`
	Version   NamedAPIResource `json:"version"`
}

// Held items
type HeldItem struct {
	Item           NamedAPIResource `json:"item"`
	VersionDetails []VersionDetail  `json:"version_details"`
}

type VersionDetail struct {
	Rarity  int              `json:"rarity"`
	Version NamedAPIResource `json:"version"`
}

// Moves
type Move struct {
	Move                NamedAPIResource     `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}

type VersionGroupDetail struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	VersionGroup    NamedAPIResource `json:"version_group"`
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
	Order           int              `json:"order"`
}

// Sprites
type Sprites struct {
	BackDefault      *string                      `json:"back_default"`
	BackFemale       *string                      `json:"back_female"`
	BackShiny        *string                      `json:"back_shiny"`
	BackShinyFemale  *string                      `json:"back_shiny_female"`
	FrontDefault     *string                      `json:"front_default"`
	FrontFemale      *string                      `json:"front_female"`
	FrontShiny       *string                      `json:"front_shiny"`
	FrontShinyFemale *string                      `json:"front_shiny_female"`
	Other            OtherSprites                 `json:"other"`
	Versions         map[string]map[string]Sprite `json:"versions"`
}

type OtherSprites struct {
	DreamWorld      Sprite `json:"dream_world"`
	Home            Sprite `json:"home"`
	OfficialArtwork Sprite `json:"official-artwork"`
	Showdown        Sprite `json:"showdown"`
}

type Sprite struct {
	FrontDefault     *string `json:"front_default"`
	FrontFemale      *string `json:"front_female,omitempty"`
	FrontShiny       *string `json:"front_shiny,omitempty"`
	FrontShinyFemale *string `json:"front_shiny_female,omitempty"`
	BackDefault      *string `json:"back_default,omitempty"`
	BackFemale       *string `json:"back_female,omitempty"`
	BackShiny        *string `json:"back_shiny,omitempty"`
	BackShinyFemale  *string `json:"back_shiny_female,omitempty"`
}

// Cries
type Cries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

// Stats
type Stat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

// Types
type TypeSlot struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

// Past types
type PastType struct {
	Generation NamedAPIResource `json:"generation"`
	Types      []TypeSlot       `json:"types"`
}

// Past abilities
type PastAbility struct {
	Generation NamedAPIResource `json:"generation"`
	Abilities  []Ability        `json:"abilities"`
}

var cache = internal.NewCache(5 * time.Second)

func GetListLocationAreas(paramUrl string) (ListLocationAreas, error) {
	var url string
	if paramUrl != "" {
		url = paramUrl
	} else {
		url = base_url + "location-area/"
	}

	locationAreas, err := callEndpoint[ListLocationAreas](url)
	if err != nil {
		return ListLocationAreas{}, err
	}
	return locationAreas, nil
}

func GetListPokemonsInArea(area string) ([]PokemonEncounter, error) {
	url := base_url + "location-area/" + area

	locationArea, err := callEndpoint[LocationArea](url)

	if err != nil {
		return []PokemonEncounter{}, err
	}

	return locationArea.PokemonEncounters, nil
}

func GetPokemonLocationArea(pokemonName string) (LocationAreaEncounter, error) {
	url := base_url + "pokemon/" + pokemonName + "/encounters"

	locationAreaEncounter, err := callEndpoint[LocationAreaEncounter](url)
	if err != nil {
		return LocationAreaEncounter{}, nil
	}

	return locationAreaEncounter, nil
}

func GetPokemon(pokemonName string) (Pokemon, error) {
	url := base_url + "pokemon/" + pokemonName

	pokemon, err := callEndpoint[Pokemon](url)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}

func callEndpoint[T any](url string) (T, error) {
	var response T

	if data, ok := cache.Get(url); ok {
		err := json.Unmarshal(data, &response)
		if err != nil {
			return response, err
		}
		return response, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Request failed with status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
