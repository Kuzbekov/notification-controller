/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package notifier

import (
	"github.com/fluxcd/pkg/runtime/events"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewSentry(t *testing.T) {
	s, err := NewSentry("https://test@localhost/1")
	require.NoError(t, err)
	assert.Equal(t, s.Client.Options().Dsn, "https://test@localhost/1")
}

func TestToSentryEvent(t *testing.T) {
	// Construct test event
	e := events.Event{
		InvolvedObject: corev1.ObjectReference{
			Kind:      "GitRepository",
			Namespace: "flux-system",
			Name:      "test-app",
		},
		Severity:  "info",
		Timestamp: metav1.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
		Message:   "message",
		Metadata: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
		ReportingController: "source-controller",
	}

	// Map to Sentry event
	s := toSentryEvent(e)

	// Assertions
	require.Equal(t, time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), s.Timestamp)
	require.Equal(t, sentry.LevelInfo, s.Level)
	require.Equal(t, "source-controller", s.ServerName)
	require.Equal(t, "GitRepository: flux-system/test-app", s.Transaction)
	require.Equal(t, map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}, s.Extra)
	require.Equal(t, "message", s.Message)
}
