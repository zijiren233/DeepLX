package deeplx

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	id = getRandomNumber()
	l  = &sync.Mutex{}
)

func getID() int64 {
	l.Lock()
	defer l.Unlock()
	id += 1
	return id
}

func getRandomNumber() int64 {
	rand.Seed(time.Now().Unix())
	num := rand.Int63n(99999) + 8300000
	return num * 1000
}

func init() {
	rand.Seed(time.Now().Unix())
}

type Lang struct {
	SourceLangUserSelected string `json:"source_lang_user_selected"`
	TargetLang             string `json:"target_lang"`
}

type CommonJobParams struct {
	WasSpoken    bool   `json:"wasSpoken"`
	TranscribeAS string `json:"transcribe_as"`
	// RegionalVariant string `json:"regionalVariant"`
}

type Params struct {
	Texts           []Text          `json:"texts"`
	Splitting       string          `json:"splitting"`
	Lang            Lang            `json:"lang"`
	Timestamp       int64           `json:"timestamp"`
	CommonJobParams CommonJobParams `json:"commonJobParams"`
}

type Text struct {
	Text                string `json:"text"`
	RequestAlternatives int    `json:"requestAlternatives"`
}

type PostData struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int64  `json:"id"`
	Params  Params `json:"params"`
}

func initData(sourceLang string, targetLang string, t Text) *PostData {
	data := &PostData{
		Jsonrpc: "2.0",
		ID:      getID(),
		Method:  "LMT_handle_texts",
		Params: Params{
			Texts:     []Text{t},
			Splitting: "sentences", // newlines,sentences,paragraphs
			Lang: Lang{
				SourceLangUserSelected: parseToDeeplSupportedLanguage(sourceLang),
				TargetLang:             parseToDeeplSupportedLanguage(targetLang),
			},
			CommonJobParams: CommonJobParams{
				WasSpoken:    false,
				TranscribeAS: "",
				// RegionalVariant: "en-US",
			},
		},
	}
	ts := time.Now().UnixMilli()
	iCount := int64(strings.Count(t.Text, "i"))
	if iCount != 0 {
		iCount += 1
		ts = ts - ts%iCount + iCount
	}
	data.Params.Timestamp = ts
	return data
}

func parseToDeeplSupportedLanguage(lang string) string {
	lang = strings.ToUpper(lang)
	switch lang {
	case "ZH", "CHS", "ZH-CN", "ZH-HANS", "CHT", "ZH-TW", "ZH-HK", "ZH-HANT":
		return "ZH"
	case "AUTO":
		return ""
	default:
		return lang
	}
}
