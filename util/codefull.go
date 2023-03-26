package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

const openaiAPIKey = "<your_openai_api_key>"
const openaiAPIURL = "https://api.openai.com/v1/completions"

type openaiCompletionRequest struct {
    Prompt     string `json:"prompt"`
    MaxTokens  int    `json:"max_tokens"`
    Temperature float64 `json:"temperature"`
    Model      string `json:"model"`
}

type openaiCompletionResponse struct {
    Choices []struct {
        Text string `json:"text"`
    } `json:"choices"`
}

func main() {
    http.HandleFunc("/", handleCompletion)
    http.ListenAndServe(":8080", nil)
}

func handleCompletion(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    prompt := r.FormValue("prompt")

    if prompt == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    requestBody := openaiCompletionRequest{
        Prompt:     prompt,
        MaxTokens:  10,
        Temperature: 0.5,
        Model:      "text-davinci-002",
    }

    requestBodyJSON, err := json.Marshal(requestBody)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    client := &http.Client{}
    req, err := http.NewRequest(http.MethodPost, openaiAPIURL, bytes.NewReader(requestBodyJSON))

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiAPIKey))

    res, err := client.Do(req)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    responseBody, err := ioutil.ReadAll(res.Body)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    var response openaiCompletionResponse

    if err := json.Unmarshal(responseBody, &response); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    var completionText strings.Builder

    for _, choice := range response.Choices {
        completionText.WriteString(choice.Text)
    }

    w.Write([]byte(completionText.String()));
}