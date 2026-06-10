// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package crypto

import "fmt"

// errDef is a minimal stand-in for the TTS error-definition framework
// (go.thethings.network/lorawan-stack/v3/pkg/errors), so this vendored crypto
// package has no external error dependency. It keeps the original call shape used
// by the LoRaWAN crypto: New() and WithAttributes(...).
type errDef struct{ msg string }

// Error implements the error interface.
func (e *errDef) Error() string { return e.msg }

// New returns the definition itself as a concrete error.
func (e *errDef) New() error { return e }

// WithAttributes attaches key/value context to the error message.
func (e *errDef) WithAttributes(pairs ...any) error {
	if len(pairs) == 0 {
		return e
	}
	return fmt.Errorf("%s (%v)", e.msg, pairs)
}

// errInvalidSize builds an invalid-size error definition for a fixed-length type.
func errInvalidSize(typeName, typeDescription, expectedSize string) *errDef {
	_ = typeName
	return &errDef{msg: fmt.Sprintf("invalid %s size, expected %s bytes", typeDescription, expectedSize)}
}
