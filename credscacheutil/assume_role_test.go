// Copyright (c) 2023 Aton-Kish
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package credscacheutil

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestAssumeRoleCacheKeyGenerator_String(t *testing.T) {
	type expected struct {
		res string
	}

	tests := []struct {
		name      string
		generator AssumeRoleCacheKeyGenerator
		expected  expected
	}{
		{
			name: "positive case: with RoleARN",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN: "role_arn",
			},
			expected: expected{
				res: `{"RoleArn": "role_arn"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
			},
			expected: expected{
				res: `{"RoleArn": "role_arn", "RoleSessionName": "role_session_name"}`,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:    "role_arn",
				ExternalID: aws.String("external_id"),
			},
			expected: expected{
				res: `{"ExternalId": "external_id", "RoleArn": "role_arn"}`,
			},
		},
		{
			name: "positive case: with RoleARN, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				SerialNumber: aws.String("mfa_serial"),
			},
			expected: expected{
				res: `{"RoleArn": "role_arn", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:  "role_arn",
				Duration: time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "RoleArn": "role_arn"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
			},
			expected: expected{
				res: `{"ExternalId": "external_id", "RoleArn": "role_arn", "RoleSessionName": "role_session_name"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				SerialNumber:    aws.String("mfa_serial"),
			},
			expected: expected{
				res: `{"RoleArn": "role_arn", "RoleSessionName": "role_session_name", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "RoleArn": "role_arn", "RoleSessionName": "role_session_name"}`,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				ExternalID:   aws.String("external_id"),
				SerialNumber: aws.String("mfa_serial"),
			},
			expected: expected{
				res: `{"ExternalId": "external_id", "RoleArn": "role_arn", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:    "role_arn",
				ExternalID: aws.String("external_id"),
				Duration:   time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "ExternalId": "external_id", "RoleArn": "role_arn"}`,
			},
		},
		{
			name: "positive case: with RoleARN, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				SerialNumber: aws.String("mfa_serial"),
				Duration:     time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "RoleArn": "role_arn", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
				SerialNumber:    aws.String("mfa_serial"),
			},
			expected: expected{
				res: `{"ExternalId": "external_id", "RoleArn": "role_arn", "RoleSessionName": "role_session_name", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "ExternalId": "external_id", "RoleArn": "role_arn", "RoleSessionName": "role_session_name"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				SerialNumber:    aws.String("mfa_serial"),
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "RoleArn": "role_arn", "RoleSessionName": "role_session_name", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				ExternalID:   aws.String("external_id"),
				SerialNumber: aws.String("mfa_serial"),
				Duration:     time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "ExternalId": "external_id", "RoleArn": "role_arn", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
				SerialNumber:    aws.String("mfa_serial"),
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "ExternalId": "external_id", "RoleArn": "role_arn", "RoleSessionName": "role_session_name", "SerialNumber": "mfa_serial"}`,
			},
		},
		{
			name: "positive case: with RoleARN, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:  "role_arn",
				Duration: time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: `{"DurationSeconds": 3600, "RoleArn": "role_arn"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.generator.String()

			assert.Equal(t, tt.expected.res, actual)
		})
	}
}

func TestAssumeRoleCacheKeyGenerator_CacheKey(t *testing.T) {
	type expected struct {
		res string
		err error
	}

	tests := []struct {
		name      string
		generator AssumeRoleCacheKeyGenerator
		expected  expected
	}{
		{
			name: "positive case: with RoleARN",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN: "role_arn",
			},
			expected: expected{
				res: "de1969e7a880d858c9bef3ba110acf78869d4527",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
			},
			expected: expected{
				res: "0c4a0c26056d87adcb5be96d428ab967ceb2daef",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:    "role_arn",
				ExternalID: aws.String("external_id"),
			},
			expected: expected{
				res: "542f374592ab812af01e95f77d5b3a3ff52aac10",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				SerialNumber: aws.String("mfa_serial"),
			},
			expected: expected{
				res: "3b5ffe7e787c6aca89ff03530dc3c413d3af1611",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:  "role_arn",
				Duration: time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "191aa88b0bb6e3b4f1a2d40d88eb9f22c2fc8fa4",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
			},
			expected: expected{
				res: "71eca3af61a2ca310d4de2243f6e9b3f243dd4f7",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				SerialNumber:    aws.String("mfa_serial"),
			},
			expected: expected{
				res: "cda918cacd9e1d1c71d510d187e90c5817e04b97",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "d4c06aa9a54967371823e504843d43af217bd80a",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				ExternalID:   aws.String("external_id"),
				SerialNumber: aws.String("mfa_serial"),
			},
			expected: expected{
				res: "c10285ba7a3c3718ff67fe7cedc11dff9b1f9373",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:    "role_arn",
				ExternalID: aws.String("external_id"),
				Duration:   time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "c4d752519e46952275e7f7975a4b8f6c70165853",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				SerialNumber: aws.String("mfa_serial"),
				Duration:     time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "38b73eece8246fa77d02085b7ac27289afd1cfc6",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID, SerialNumber",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
				SerialNumber:    aws.String("mfa_serial"),
			},
			expected: expected{
				res: "696d8943794ca230409bccc0c6e4158934275eed",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "98d5d14fdb662800ed57cf48cd4f3c2ec9ff46a0",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				SerialNumber:    aws.String("mfa_serial"),
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "4af18d294c15620b13617e9a7f44a793a0b65e50",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, ExternalID, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:      "role_arn",
				ExternalID:   aws.String("external_id"),
				SerialNumber: aws.String("mfa_serial"),
				Duration:     time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "9f707334eb3cc733e2f264fc7d39d4fb0b2b04f1",
				err: nil,
			},
		},
		{
			name: "positive case: with RoleARN, RoleSessionName, ExternalID, SerialNumber, Duration(=3600s)",
			generator: AssumeRoleCacheKeyGenerator{
				RoleARN:         "role_arn",
				RoleSessionName: "role_session_name",
				ExternalID:      aws.String("external_id"),
				SerialNumber:    aws.String("mfa_serial"),
				Duration:        time.Duration(3600) * time.Second,
			},
			expected: expected{
				res: "9460fb8eea6805c12dbd1c9f1426735b38fc5fb4",
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.generator.CacheKey()

			if tt.expected.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.res, actual)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			}
		})
	}
}
