package httphook

import (
	"awesomeProject1/internal/config"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

func Send(conf *config.HttpConfig, client *http.Client) {
	for _, endpoint := range conf.Endpoints {
		go sendToEndpoint(endpoint, client)
	}
}

func sendToEndpoint(endpoint *config.HttpConfigItem, client *http.Client) {
	URL, err := url.Parse(endpoint.URL)
	if err != nil {
		log.Err(err).
			Str("url", endpoint.URL).
			Msg("Error parsing endpoint url")
	}
	req := &http.Request{
		Method: endpoint.Method,
		URL:    URL,
	}
	log.Info().
		Str("url", endpoint.URL).
		Str("method", endpoint.Method).
		Msg("Calling endpoint")
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			log.Err(e).
				Str("url", endpoint.URL).
				Str("method", endpoint.Method).
				Msg("Could not close response body")
		}
	}(resp.Body)
	if err != nil {
		log.Err(err).
			Str("url", endpoint.URL).
			Str("method", endpoint.Method).
			Msg("Could not call endpoint")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		log.Error().
			Str("url", endpoint.URL).
			Str("method", endpoint.Method).
			Int("status", resp.StatusCode).
			Bytes("respBody", body).
			Msg("Endpoint with status code >400")
	}
}
