package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ResponseData struct {
	Data json.RawMessage `json:"data"`
}
type ResponseStatus struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_msg"`
}

// Response API response struct
type Response struct {
	ResponseStatus
	ResponseData
}

type IHttpClient interface {
	Get(url string, data interface{}) error
	Post(url string, body interface{}, data interface{}) error
	HandleAPIError(res Response) error
}

type httpClient struct {
	IHttpClient
}

func NewHttpClient() *httpClient {
	return &httpClient{}
}

//Get send GET HTTP request
func (h *httpClient) Get(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http.Get(url)")
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return errors.Wrapf(err, "json.NewDecoder(resp.Body).Decode(&res)")
	}
	if err = h.HandleAPIError(res); err != nil {
		return err
	}
	err = json.Unmarshal(res.Data, data)
	if err != nil {
		return errors.Wrapf(err, "json.Unmarshal. res.Data: %s", string(res.Data))
	}
	return nil
}

//Post send POST HTTP request
func (h *httpClient) Post(url string, body interface{}, data interface{}) error {
	var res Response
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal(body). body: %+v", body)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.Wrapf(err, "http.Post(url, \"application/json\", bytes.NewBuffer(jsonValue)). jsonValue: %+v", jsonValue)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return errors.Wrapf(err, "json.NewDecoder(resp.Body).Decode(&res)")
	}
	if err := h.HandleAPIError(res); err != nil {
		return err
	}
	err = json.Unmarshal(res.Data, &data)
	if err != nil {
		return errors.Wrapf(err, "json.Unmarshal. res.Data: %s", string(res.Data))
	}

	return nil
}

func (h *httpClient) HandleAPIError(res Response) error {
	if !res.Success {
		return fmt.Errorf("API error: %s", res.ErrorMessage)
	}
	return nil
}