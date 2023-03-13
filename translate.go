package deeplx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const endPoint = "https://www2.deepl.com/jsonrpc"

type DeepLData struct {
	Result struct {
		Texts []struct {
			Alternatives []struct {
				Text string `json:"text"`
			} `json:"alternatives"`
			Text string `json:"text"`
		} `json:"texts"`
		Lang              string             `json:"lang"`
		LangIsConfident   bool               `json:"lang_is_confident"`
		DetectedLanguages map[string]float64 `json:"detectedLanguages"`
	} `json:"result"`
}

type Translated struct {
	Detected     Detected `json:"detected"`
	Text         string   `json:"text"`          // translated text
	Alternatives []string `json:"pronunciation"` // pronunciation of translated text
}

// Detected represents language detection result
type Detected struct {
	Lang         string `json:"lang"` // detected language
	IsConfidence bool
	Confidence   float64 `json:"confidence"` // the confidence of detection result (0.00 to 1.00)
}

func Translate(text, source, target string) (*Translated, error) {
	if text == "" {
		return &Translated{}, nil
	}
	postData := initData(source, target, Text{
		Text:                text,
		RequestAlternatives: 3,
	})

	post_byte, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}

	if (postData.ID+5)%29 == 0 || (postData.ID+3)%13 == 0 {
		post_byte = bytes.ReplaceAll(post_byte, []byte("\"method\":\""), []byte("\"method\" : \""))
	} else {
		post_byte = bytes.ReplaceAll(post_byte, []byte("\"method\":\""), []byte("\"method\": \""))
	}

	request, err := http.NewRequest("POST", endPoint, bytes.NewReader(post_byte))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := &DeepLData{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	if len(data.Result.Texts) != 1 {
		return nil, fmt.Errorf("err: %s", body)
	}
	dt := &Translated{}
	dt.Text = data.Result.Texts[0].Text
	dt.Detected.IsConfidence = data.Result.LangIsConfident
	dt.Detected.Lang = data.Result.Lang
	dt.Detected.Confidence = data.Result.DetectedLanguages[data.Result.Lang]
	return dt, nil
}
