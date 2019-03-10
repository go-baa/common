package nsq

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type wrappedResp struct {
	Status     string      `json:"status_txt"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

// stores the result in the value pointed to by ret(must be a pointer)
func apiRequestNegotiateV1(method string, endpoint string, body io.Reader, ret interface{}) error {
	httpclient := &http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.nsq; version=1.0")

	resp, err := httpclient.Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("got response %s %q", resp.Status, respBody)
	}

	if len(respBody) == 0 {
		respBody = []byte("{}")
	}

	if resp.Header.Get("X-NSQ-Content-Type") == "nsq; version=1.0" {
		return json.Unmarshal(respBody, ret)
	}

	wResp := &wrappedResp{
		Data: ret,
	}

	if err = json.Unmarshal(respBody, wResp); err != nil {
		return err
	}

	// wResp.StatusCode here is equal to resp.StatusCode, so ignore it
	return nil
}
