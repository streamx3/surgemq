// Copyright (c) 2014 The SurgeMQ Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package topics deals with MQTT topic names, topic filters and subscriptions.
// - "Topic name" is a / separated string that could contain #, * and $
// - / in topic name separates the string into "topic levels"
// - # is a multi-level wildcard, and it must be the last character in the
//   topic name. It represents the parent and all children levels.
// - + is a single level wildwcard. It must be the only character in the
//   topic level. It represents all names in the current level.
// - $ is a special character that says the topic is a system level topic
package topics

import (
	"fmt"

	"github.com/troian/surgemq/message"
	persistenceTypes "github.com/troian/surgemq/persistence/types"
	"github.com/troian/surgemq/systree"
	"github.com/troian/surgemq/types"
)

const (
	// MWC is the multi-level wildcard
	MWC = "#"

	// SWC is the single level wildcard
	SWC = "+"

	// SEP is the topic level separator
	//SEP = "/"

	// SYS is the starting character of the system level topics
	//SYS = "$"

	// Both wildcards
	//BWC = "#+"
)

var (
	// ErrAuthFailure is returned when the user/pass supplied are invalid
	//ErrAuthFailure = errors.New("auth: Authentication failure")

	// ErrAuthProviderNotFound is returned when the requested provider does not exist.
	// It probably hasn't been registered yet.
	//ErrAuthProviderNotFound = errors.New("auth: Authentication provider not found")

	providers = make(map[string]Provider)
)

// Provider interface
type Provider interface {
	Configure(stat systree.TopicsStat, persist persistenceTypes.Retained) error
	Subscribe(topic string, qos message.QosType, subscriber *types.Subscriber) (message.QosType, error)
	UnSubscribe(topic string, subscriber *types.Subscriber) error
	Publish(msg *message.PublishMessage) error
	Retain(msg *message.PublishMessage) error
	Retained(topic string, msgs *[]*message.PublishMessage) error
	Close() error
}

// Register topic provider
func Register(name string, provider Provider) {
	if provider == nil {
		panic("topics: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		panic("topics: Register called twice for provider " + name)
	}

	providers[name] = provider
}

// UnRegister topic provider
func UnRegister(name string) {
	delete(providers, name)
}

// Config of topics manager
type Config struct {
	Name    string
	Stat    systree.TopicsStat
	Persist persistenceTypes.Retained
}

// Manager of topics
type Manager struct {
	//config *Config
	p Provider
}

// NewManager add new manager
func NewManager(config Config) (*Manager, error) {
	p, ok := providers[config.Name]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", config.Name)
	}

	if err := p.Configure(config.Stat, config.Persist); err != nil {
		return nil, err
	}

	return &Manager{p: p}, nil
}

// Subscribe to topic
func (m *Manager) Subscribe(topic string, qos message.QosType, subscriber *types.Subscriber) (message.QosType, error) {
	return m.p.Subscribe(topic, qos, subscriber)
}

// UnSubscribe from topic
func (m *Manager) UnSubscribe(topic string, subscriber *types.Subscriber) error {
	return m.p.UnSubscribe(topic, subscriber)
}

// Publish message
func (m *Manager) Publish(msg *message.PublishMessage) error {
	return m.p.Publish(msg)
}

// Retain messages
func (m *Manager) Retain(msg *message.PublishMessage) error {
	return m.p.Retain(msg)
}

// Retained messages
func (m *Manager) Retained(topic string, msgs *[]*message.PublishMessage) error {
	return m.p.Retained(topic, msgs)
}

// Close manager
func (m *Manager) Close() error {
	return m.p.Close()
}
