package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/hx23840/YunDoBridge/internal/audio"
	"github.com/hx23840/YunDoBridge/internal/types"
)

// checkNewTracksResponse checks if there are any errors in the new tracks response
func checkNewTracksResponse(resp *types.NewTracksResponse, sdpExpected bool) error {
	if resp.ErrorCode != "" {
		return errors.New(resp.ErrorDescription)
	}
	if len(resp.Tracks) > 0 && resp.Tracks[0].ErrorDescription != "" {
		return errors.New(resp.Tracks[0].ErrorDescription)
	}
	if sdpExpected && resp.SessionDescription == nil {
		return errors.New("empty sdp from Calls for session")
	}
	return nil
}

func handleExchange(sessionA *audio.CallsSession, sessionB *audio.CallsSession, newTracksResponseA *types.NewTracksResponse, r *http.Request, env *types.Env) {
	// Create new tracks request for session B
	newTracksResponseB, err := sessionB.NewTracks(map[string]interface{}{
		"tracks": []map[string]interface{}{
			{
				"location":                 "local",
				"trackName":                "ai-generated-voice",
				"bidirectionalMediaStream": true,
				"kind":                     "audio",
			},
		},
	})
	if err != nil {
		log.Printf("Error in exchange step B: %v", err)
		return
	}

	if err := checkNewTracksResponse(newTracksResponseB, true); err != nil {
		log.Printf("Error checking tracks response B: %v", err)
		return
	}

	// Request OpenAI service
	openaiAnswer, err := requestOpenAIService(r, newTracksResponseB.SessionDescription, env)
	if err != nil {
		log.Printf("Error requesting OpenAI service: %v", err)
		return
	}

	// Complete negotiation with OpenAI
	if err := sessionB.Renegotiate(openaiAnswer); err != nil {
		log.Printf("Error in renegotiation: %v", err)
		return
	}

	// Exchange step one: Get AI-generated voice from session B and send to session A
	exchangeStepOne, err := sessionA.NewTracks(map[string]interface{}{
		"tracks": []map[string]interface{}{
			{
				"location":  "remote",
				"sessionId": sessionB.SessionID,
				"trackName": "ai-generated-voice",
				"mid":       "#user-mic",
			},
		},
	})
	if err != nil {
		log.Printf("Error in exchange step one: %v", err)
		return
	}

	if err := checkNewTracksResponse(exchangeStepOne, false); err != nil {
		log.Printf("Error checking exchange step one: %v", err)
		return
	}

	// Exchange step two: Get user voice from session A and send to session B
	exchangeStepTwo, err := sessionB.NewTracks(map[string]interface{}{
		"tracks": []map[string]interface{}{
			{
				"location":  "remote",
				"sessionId": sessionA.SessionID,
				"trackName": "user-mic",
				"mid":       "#ai-generated-voice",
			},
		},
	})
	if err != nil {
		log.Printf("Error in exchange step two: %v", err)
		return
	}

	if err := checkNewTracksResponse(exchangeStepTwo, false); err != nil {
		log.Printf("Error checking exchange step two: %v", err)
		return
	}

	log.Println("Exchange process completed successfully")
}

// requestOpenAIService sends a request to OpenAI service and returns the session description
func requestOpenAIService(r *http.Request, offer *types.SessionDescription, env *types.Env) (types.SessionDescription, error) {
	endpointURL, err := url.Parse(env.OpenAIModelEndpoint)
	if err != nil {
		return types.SessionDescription{}, err
	}

	originalURL := url.URL{RawQuery: r.URL.RawQuery}
	originalParams := originalURL.Query()
	endpointParams := endpointURL.Query()

	// Merge parameters, giving priority to original request parameters
	for key, value := range endpointParams {
		if !originalParams.Has(key) {
			originalParams[key] = value
		}
	}
	endpointURL.RawQuery = originalParams.Encode()

	resp, err := http.NewRequest("POST", endpointURL.String(), bytes.NewBufferString(offer.SDP))
	if err != nil {
		return types.SessionDescription{}, err
	}

	resp.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.OpenAIAPIKey))
	resp.Header.Set("Content-Type", "application/sdp")

	client := &http.Client{}
	response, err := client.Do(resp)
	if err != nil {
		return types.SessionDescription{}, err
	}
	defer response.Body.Close()

	answerSDP, err := io.ReadAll(response.Body)
	if err != nil {
		return types.SessionDescription{}, err
	}

	return types.SessionDescription{
		Type: "answer",
		SDP:  string(answerSDP),
	}, nil
}

// HandleRequest handles the main WebRTC endpoint request
func HandleRequest(w http.ResponseWriter, r *http.Request, env *types.Env) {
	if r.Method == "OPTIONS" {
		handleOptions(w)
		return
	}

	// Create session A (connection between end-user and Calls)
	sessionA, err := audio.NewSession(env.CallsBaseURL, env.CallsAppID, env.CallsAppToken, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read user's SDP
	userSDP, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create new tracks request A
	newTracksResponseA, err := sessionA.NewTracks(map[string]interface{}{
		"sessionDescription": types.SessionDescription{
			SDP:  string(userSDP),
			Type: "offer",
		},
		"tracks": []map[string]interface{}{
			{
				"location":                 "local",
				"trackName":                "user-mic",
				"bidirectionalMediaStream": true,
				"kind":                     "audio",
				"mid":                      "0",
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := checkNewTracksResponse(newTracksResponseA, true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session B (connection between Calls and OpenAI)
	sessionB, err := audio.NewSession(env.CallsBaseURL, env.CallsAppID, env.CallsAppToken, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle exchange process asynchronously
	go handleExchange(sessionA, sessionB, newTracksResponseA, r, env)

	// Return session A's SDP
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newTracksResponseA.SessionDescription.SDP))
}

// handleOptions handles the OPTIONS request
func handleOptions(w http.ResponseWriter) {
	headers := w.Header()
	headers.Set("Accept-Post", "application/sdp")
	headers.Set("Access-Control-Allow-Credentials", "true")
	headers.Set("Access-Control-Allow-Headers", "content-type,authorization,if-match")
	headers.Set("Access-Control-Allow-Methods", "PATCH,POST,PUT,DELETE,OPTIONS")
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Access-Control-Expose-Headers", "x-thunderclap,location,link,accept-post,accept-patch,etag")
	headers.Set("Link", "<stun:stun.cloudflare.com:3478>; rel=\"ice-server\"")
	w.WriteHeader(http.StatusNoContent)
}
