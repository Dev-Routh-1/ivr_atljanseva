package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type lookupResponse struct {
	Found       bool              `json:"found"`
	ChannelVars map[string]string `json:"_channel_vars"`
}

func main() {
	env := readAGIEnv()

	callerID := env["agi_callerid"]
	if callerID == "" {
		log.Fatal("agi_callerid not found in AGI environment")
	}

	apiURL := os.Getenv("AGI_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	url := fmt.Sprintf("%s/ivr/citizen/%s", apiURL, callerID)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	var data lookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("JSON decode failed: %v", err)
	}

	for k, v := range data.ChannelVars {
		fmt.Printf("SET VARIABLE %s \"%s\"\n", k, v)
	}
}

func readAGIEnv() map[string]string {
	env := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	return env
}
