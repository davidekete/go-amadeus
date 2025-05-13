package main

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expiry      int    `json:"expires_in"`
}

//GET request structs
type FlightResponse struct {
	Data []FlightDestination `json:"data"`
}

type FlightDestination struct {
	FlightType    string      `json:"type"`
	Origin        string      `json:"origin"`
	Destination   string      `json:"destination"`
	DepartureDate string      `json:"departureDate"`
	ReturnDate    string      `json:"returnDate"`
	Price         FlightPrice `json:"price"`
}

type FlightPrice struct {
	Total string `json:"total"`
}

// POST request structs
type CancellationResponse struct {
	Data CancellationData `json:"data"`
}

type CancellationData struct {
	ConfirmNbr        string `json:"confirmNbr"`
	ReservationStatus string `json:"reservationStatus"`
}
