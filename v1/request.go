package v1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var prefix = "/api/v1"

// GetRequest implements GET Request
func (c *MgClient) GetRequest(urlWithParameters string) ([]byte, int, error) {
	var res []byte

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), nil)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("X-Transport-Token", c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return res, 0, err
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		err = fmt.Errorf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, err
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		return res, resp.StatusCode, err
	}

	return res, resp.StatusCode, err
}

// PostRequest implements POST Request
func (c *MgClient) PostRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"POST",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

// PutRequest implements PUT Request
func (c *MgClient) PutRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"PUT",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

// DeleteRequest implements DELETE Request
func (c *MgClient) DeleteRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"DELETE",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

func makeRequest(reqType, url string, buf *bytes.Buffer, c *MgClient) ([]byte, int, error) {
	var res []byte
	req, err := http.NewRequest(reqType, url, buf)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transport-Token", c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return res, 0, err
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		err = fmt.Errorf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, err
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		return res, 0, err
	}

	return res, resp.StatusCode, err
}

func buildRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}
