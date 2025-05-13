package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var BaseURL string = "https://test.api.amadeus.com/v1/"

func getAccessToken() (AccessTokenResponse, error) {
	var respData AccessTokenResponse
	endpoint := fmt.Sprintf("%s/security/oauth2/token", BaseURL)

	// Read and validate credentials
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		return respData, fmt.Errorf("API_KEY not set")
	}
	apiSecret, ok := os.LookupEnv("API_SECRET")
	if !ok {
		return respData, fmt.Errorf("API_SECRET not set")
	}

	// Build form body
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {apiKey},
		"client_secret": {apiSecret},
	}

	// Create request
	req, err := http.NewRequest("POST",
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return respData, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Do the request with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return respData, err
	}
	defer resp.Body.Close()

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return respData, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	// Decode JSON into your struct
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return respData, err
	}

	return respData, nil
}

// GET REQUEST
func GetRequest(origin, maxPrice string) (*FlightResponse, error) {
	// Build the full URL for the endpoint.
	endpoint := fmt.Sprintf("%s/shopping/flight-destinations", BaseURL)

	// Create a new HTTP GET request so you can add headers and params.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Attach query parameters: “origin” and “maxPrice”.
	q := req.URL.Query()
	q.Add("origin", origin)
	q.Add("maxPrice", maxPrice)
	req.URL.RawQuery = q.Encode()

	//Get accessToken
	accessToken, err := getAccessToken()

	if err != nil {
		return nil, err
	}

	//  Set the Authorization header with your bearer token.
	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

	// Create a new HTTP client to execute the request.
	client := &http.Client{}

	// Perform the HTTP request and capture the response.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Ensure the response body is closed after processing to avoid resource leaks.
	defer resp.Body.Close()

	// Check for non-2xx status codes and return an error if found.
	//    This prevents you from trying to parse error HTML or JSON.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %s", resp.Status)
	}

	// Decode the JSON response directly into your struct.
	//    This maps the API’s JSON into Go types for you to work with.
	var result FlightResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Return the parsed response to the caller.
	return &result, nil
}

func PostRequest(orderID string, confirmationNumber string) (*CancellationResponse, error) {
	// Build the endpoint URL using the base path and the provided order ID.
	endpoint := fmt.Sprintf("%s/ordering/transfer-orders/%s/transfers/cancellation", BaseURL, orderID)

	// Prepare query parameters: add the confirmation number as 'confirmNbr'.
	queryParams := url.Values{}
	queryParams.Add("confirmNbr", confirmationNumber)

	// Parse the endpoint URL string into a url.URL object for safe manipulation.
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	// Attach the encoded query parameters to the URL (e.g., ?confirmNbr=12345).
	u.RawQuery = queryParams.Encode()

	// Create a new HTTP POST request with no body (nil) since parameters are in the URL.
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Get accessToken
	accessToken, err := getAccessToken()

	if err != nil {
		return nil, err
	}

	//  Set the Authorization header with your bearer token.
	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

	// Initialize an HTTP client.
	client := &http.Client{}

	// Execute the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	// Ensure the response body is closed to free resources.
	defer resp.Body.Close()

	// Check for HTTP status codes outside the 2xx success range.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %s", resp.Status)
	}

	// Decode the JSON response body into the CancellationResponse struct.
	var result CancellationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Return the parsed result and nil error on success.
	return &result, nil
}

func DeleteRequest(flightOrderID string) {
	// Construct the endpoint URL using the base path and the provided flight order ID.
	endpoint := fmt.Sprintf("%s/booking/flight-orders/%s", BaseURL, flightOrderID)

	// Create a new HTTP DELETE request with no body (nil).
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Get the access token required for authorization.
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	// Set the Authorization header with the bearer token.
	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

	// Initialize an HTTP client to execute the request.
	client := &http.Client{}

	// Perform the HTTP DELETE request and capture the response.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	// Ensure the response body is closed after processing to avoid resource leaks.
	defer resp.Body.Close()

	// Read the response body for logging or debugging purposes.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Reading response failed: %v", err)
	}

	// Print the HTTP status and response body to the console.
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", body)
}
