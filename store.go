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
 */

// Package fact_totem_client abstracts the web interface.
package fact_totem_client

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type Fact struct {
	Id        ulid.ULID
	Timestamp time.Time
	Content   interface{}
}

type TailResponse struct {
	Aggregate string `json:"aggregate"`
	Entity    string `json:"entity,omitempty"`
	Fact      Fact   `json:"fact"`
	Total     uint   `json:"total"`
}

type ReadResponse struct {
	Aggregate string `json:"aggregate"`
	Entity    string `json:"entity"`
	Facts     []Fact `json:"facts"`
	Total     uint   `json:"total"`
	PageSize  int    `json:"page-size"`
}

type ScanResponse struct {
	Aggregate string   `json:"aggregate"`
	Entities  []string `json:"entities"`
	Total     uint     `json:"total"`
}

// EventStore provides an interface to store events for a topic, and retrieve them later.
type EventStore interface {
	// Append an event to the event store for the fact
	Append(aggregate string, entity string, content interface{}) (*TailResponse, error)
	// Tail gets the last event id
	Tail(aggregate string, entity string) (*TailResponse, error)
	// Read the events for an aggregate from the identified event id
	Read(aggregate string, entity string, originEventId string, maxCount int) (*ReadResponse, error)
	// Scan will list all keys in the aggregate (excluding individual events)
	Scan(aggregate string) (*ScanResponse, error)
}
