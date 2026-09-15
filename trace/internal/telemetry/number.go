// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"encoding/json"
	"math"
	"strconv"
)

// protoInt64 represents the protobuf encoding of integers which can be either
// strings or integers.
type protoInt64 int64

// Int64 returns the protoInt64 as an int64.
func (i *protoInt64) Int64() int64 { return int64(*i) }

// UnmarshalJSON decodes both strings and integers.
func (i *protoInt64) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		parsedInt, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*i = protoInt64(parsedInt)
	} else {
		var parsedInt int64
		if err := json.Unmarshal(data, &parsedInt); err != nil {
			return err
		}
		*i = protoInt64(parsedInt)
	}
	return nil
}

// protoUint64 represents the protobuf encoding of integers which can be either
// strings or integers.
type protoUint64 uint64

// Int64 returns the protoUint64 as a uint64.
func (i *protoUint64) Uint64() uint64 { return uint64(*i) }

// UnmarshalJSON decodes both strings and integers.
func (i *protoUint64) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		parsedUint, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*i = protoUint64(parsedUint)
	} else {
		var parsedUint uint64
		if err := json.Unmarshal(data, &parsedUint); err != nil {
			return err
		}
		*i = protoUint64(parsedUint)
	}
	return nil
}

// protoFloat64 represents the protobuf JSON encoding of doubles.
type protoFloat64 float64

// Float64 returns the protoFloat64 as a float64.
func (f *protoFloat64) Float64() float64 { return float64(*f) }

// MarshalJSON encodes finite values as JSON numbers and non-finite values as
// strings as required by protobuf JSON.
func (f protoFloat64) MarshalJSON() ([]byte, error) {
	v := float64(f)
	switch {
	case math.IsNaN(v):
		return json.Marshal("NaN")
	case math.IsInf(v, 1):
		return json.Marshal("Infinity")
	case math.IsInf(v, -1):
		return json.Marshal("-Infinity")
	default:
		return json.Marshal(v)
	}
}

// UnmarshalJSON decodes both JSON numbers and protobuf JSON strings.
func (f *protoFloat64) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		parsedFloat, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*f = protoFloat64(parsedFloat)
	} else {
		var parsedFloat float64
		if err := json.Unmarshal(data, &parsedFloat); err != nil {
			return err
		}
		*f = protoFloat64(parsedFloat)
	}
	return nil
}
