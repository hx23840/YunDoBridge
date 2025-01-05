package audio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hx23840/YunDoBridge/internal/types"
)

// CallsSession represents a Calls session
type CallsSession struct {
	SessionID string
	Headers   map[string]string
	Endpoint  string
}

// NewTracks creates new tracks for the session
func (s *CallsSession) NewTracks(body interface{}) (*types.NewTracksResponse, error) {
	url := fmt.Sprintf("%s/sessions/%s/tracks/new?streamDebug", s.Endpoint, s.SessionID)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var newTracksResponse types.NewTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&newTracksResponse); err != nil {
		return nil, err
	}

	return &newTracksResponse, nil
}

// Renegotiate performs session renegotiation
func (s *CallsSession) Renegotiate(sdp types.SessionDescription) error {
	url := fmt.Sprintf("%s/sessions/%s/renegotiate?streamDebug", s.Endpoint, s.SessionID)

	body := map[string]types.SessionDescription{
		"sessionDescription": sdp,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return nil
}

// NewSession creates a new Calls session
func NewSession(baseURL, appID, appToken string, thirdparty bool) (*CallsSession, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", appToken),
		"Content-Type":  "application/json",
	}

	endpoint := fmt.Sprintf("%s/%s", baseURL, appID)
	newSessionURL := fmt.Sprintf("%s/sessions/new?streamDebug", endpoint)

	if thirdparty {
		newSessionURL += "&thirdparty=true"
	}

	req, err := http.NewRequest("POST", newSessionURL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var sessionResponse types.NewSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		return nil, err
	}

	return &CallsSession{
		SessionID: sessionResponse.SessionID,
		Headers:   headers,
		Endpoint:  endpoint,
	}, nil
}
