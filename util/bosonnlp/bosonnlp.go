package bosonnlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const gateway = "http://api.bosonnlp.com"

// Bosonnlp 博森NLP 接口
type Bosonnlp struct {
	token string
}

// New 创建一个博森NLP实例
func New(token string) *Bosonnlp {
	return &Bosonnlp{
		token: token,
	}
}

// Keywords 关键词
type Keywords struct {
	Score float64
	Word  string
}

// KeywordsAnalysis 关键词分析
func (t *Bosonnlp) KeywordsAnalysis(text string) ([]Keywords, error) {
	buf := bytes.NewBuffer(nil)
	b, err := json.Marshal(text)
	if err != nil {
		return nil, fmt.Errorf("BosonNlp KeywordsAnalysis: %v", err)
	}
	buf.Write(b)
	req, err := http.NewRequest("POST", gateway+"/keywords/analysis", buf)
	if err != nil {
		return nil, fmt.Errorf("BosonNlp KeywordsAnalysis: %v", err)
	}
	header := req.Header
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set("X-Token", t.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BosonNlp KeywordsAnalysis: %v", err)
	}
	buf.Reset()
	n, err := buf.ReadFrom(resp.Body)
	if err != nil || n == 0 {
		return nil, fmt.Errorf("BosonNlp KeywordsAnalysis: %v", err)
	}
	var words []Keywords
	var data [][]interface{}
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		// return nil, fmt.Errorf("BosonNlp KeywordsAnalysis: %v", err)
		words = append(words, Keywords{
			Score: 0.0,
			Word:  text,
		})
		return words, fmt.Errorf("BosonNlp KeywordsAnalysis: %v", err)
	}
	for i := range data {
		words = append(words, Keywords{
			Score: data[i][0].(float64),
			Word:  data[i][1].(string),
		})
	}
	return words, nil
}
