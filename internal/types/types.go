package types

// Env holds configuration settings
type Env struct {
	OpenAIAPIKey        string
	OpenAIModelEndpoint string
	CallsBaseURL        string
	CallsAppID          string
	CallsAppToken       string
}

// SessionDescription represents a WebRTC session description
type SessionDescription struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

// NewSessionResponse represents a new Calls session response
type NewSessionResponse struct {
	SessionID string `json:"sessionId"`
}

// NewTrackResponse represents a new track response
type NewTrackResponse struct {
	TrackName        string `json:"trackName"`
	Mid              string `json:"mid"`
	ErrorCode        string `json:"errorCode,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
}

// NewTracksResponse represents a collection of new track responses
type NewTracksResponse struct {
	Tracks             []NewTrackResponse  `json:"tracks"`
	SessionDescription *SessionDescription `json:"sessionDescription,omitempty"`
	ErrorCode          string              `json:"errorCode,omitempty"`
	ErrorDescription   string              `json:"errorDescription,omitempty"`
}

// TrackLocator represents a track location identifier
type TrackLocator struct {
	Location  string `json:"location"`
	SessionID string `json:"sessionId"`
	TrackName string `json:"trackName"`
}
