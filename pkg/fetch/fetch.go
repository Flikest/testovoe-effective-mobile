package fetch

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

var client = &http.Client{}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type RequestCountry struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

type RequestAge struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type RequestGender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

func GetCountry(name string) ([]Country, error) {
	res, err := http.Get("https://api.nationalize.io/?name=" + name)
	if err != nil {
		slog.Error("Error sending request: ", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("server side error: ", res.StatusCode)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response: ", err)
		return nil, err
	}

	var requestCountry RequestCountry
	if err = json.Unmarshal(body, &requestCountry); err != nil {
		slog.Error("Error decoding JSON: ", err)
		return nil, err
	}

	return requestCountry.Country, err
}

func GetAge(name string) (int, error) {
	res, err := http.Get("https://api.agify.io/?name=" + name)
	if err != nil {
		slog.Error("Error sending request: ", err)
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("server side error: ", res.StatusCode)
		return 0, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response: ", err)
		return 0, err
	}

	var requestAge RequestAge
	if err = json.Unmarshal(body, &requestAge); err != nil {
		slog.Error("Error decoding JSON: ", err)
		return 0, err
	}

	return requestAge.Age, err
}

func GetGender(name string) (string, error) {
	res, err := http.Get("https://api.genderize.io/?name=" + name)
	if err != nil {
		slog.Error("Error sending request:", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("server side error: ", err)
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response:", err)
		return "", nil
	}

	var requestGender RequestGender
	if err = json.Unmarshal(body, &requestGender); err != nil {
		slog.Error("Error decoding JSON: ", err)
		return "", err
	}

	return requestGender.Gender, err
}
