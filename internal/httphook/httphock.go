package httphook

import (
	"awesomeProject1/internal/config"
	"io"
	"net/http"
	neturl "net/url"

	"github.com/rs/zerolog/log"
)

func Send(conf *config.HTTPConfig, client *http.Client) {
	for _, endpoint := range conf.Endpoints {
		go sendToEndpoint(endpoint, client)
	}
}

func sendToEndpoint(endpoint *config.HTTPConfigItem, client *http.Client) {
	url, err := neturl.Parse(endpoint.URL)
	if err != nil {
		log.Err(err).
			Str("url", endpoint.URL).
			Msg("Error parsing endpoint url")
	}
	req := &http.Request{
		Method: endpoint.Method,
		URL:    url,
	}
	log.Info().
		Str("url", endpoint.URL).
		Str("method", endpoint.Method).
		Msg("Calling endpoint")
	resp, err := client.Do(req)
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			log.Err(e).
				Str("url", endpoint.URL).
				Str("method", endpoint.Method).
				Msg("Could not close response body")
		}
	}()
	if err != nil {
		log.Err(err).
			Str("url", endpoint.URL).
			Str("method", endpoint.Method).
			Msg("Could not call endpoint")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).
			Str("url", endpoint.URL).
			Str("method", endpoint.Method).
			Int("status", resp.StatusCode).
			Msg("Could not read response body")
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error().
			Str("url", endpoint.URL).
			Str("method", endpoint.Method).
			Int("status", resp.StatusCode).
			Bytes("respBody", body).
			Msg("Endpoint with status code >400")
	}
}
