/*
 * Copyright (c) 2021.  D-Haven.org
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package fact_totem_client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Action    Action      `json:"action"`
	Aggregate string      `json:"aggregate"`
	Entity    string      `json:"entity,omitempty"`
	Content   interface{} `json:"content,omitempty"`
	Origin    string      `json:"origin,omitempty"`
	PageSize  int         `json:"page-size,omitempty"`
}

type StandardClient struct {
	Token        string
	FactTotemUrl string
	http         *http.Client
}

func (c *StandardClient) Append(aggregate string, entity string, content interface{}) (*TailResponse, error) {
	req := Request{
		Action:    Append,
		Aggregate: aggregate,
		Entity:    entity,
		Content:   content,
	}

	trJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.post(bytes.NewBuffer(trJson))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed: %d %s\n%s", req.Action.String(), resp.StatusCode, resp.Status, string(body))
	}

	response := TailResponse{}
	err = json.Unmarshal(body, &response)

	return &response, err
}

func (c *StandardClient) Tail(aggregate string, entity string) (*TailResponse, error) {
	req := Request{
		Action:    Tail,
		Aggregate: aggregate,
		Entity:    entity,
	}

	trJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.post(bytes.NewBuffer(trJson))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed: %d %s\n%s", req.Action.String(), resp.StatusCode, resp.Status, string(body))
	}

	response := TailResponse{}
	err = json.Unmarshal(body, &response)

	return &response, err
}

func (c *StandardClient) Read(aggregate string, entity string, originEventId string, maxCount int) (*ReadResponse, error) {
	req := Request{
		Action:    Read,
		Aggregate: aggregate,
		Entity:    entity,
		Origin:    originEventId,
		PageSize:  maxCount,
	}

	trJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.post(bytes.NewBuffer(trJson))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed: %d %s\n%s", req.Action.String(), resp.StatusCode, resp.Status, string(body))
	}

	response := ReadResponse{}
	err = json.Unmarshal(body, &response)

	return &response, err
}

func (c *StandardClient) Scan(aggregate string) (*ScanResponse, error) {
	req := Request{
		Action:    Scan,
		Aggregate: aggregate,
	}

	trJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.post(bytes.NewBuffer(trJson))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed: %d %s\n%s", req.Action.String(), resp.StatusCode, resp.Status, string(body))
	}

	response := ScanResponse{}
	err = json.Unmarshal(body, &response)

	return &response, err
}

func (c *StandardClient) post(body io.Reader) (*http.Response, error) {
	if c.http == nil {
		err := c.Refresh()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(http.MethodPost, c.FactTotemUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	return c.http.Do(req)
}

func (c *StandardClient) Refresh() error {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return err
	}

	c.http = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	return nil
}
