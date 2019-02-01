/*
	Copyright 2019 Nokia

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

package schema

import (
	"errors"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// ErrSchemaInvalid is an error returned when data are invalid according to JSON schema
	ErrSchemaInvalid = errors.New("JSON validation failed")
)

// JSONSchema helps validating serializable data with a JSON schema
type JSONSchema struct {
	inner gojsonschema.Schema
}

// LoadSchemaFromFile loads a schema from a file
func LoadSchemaFromFile(filepath string) (*JSONSchema, error) {
	/* #nosec */
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Panic(err)
		}
	}()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return NewSchemaFromBytes(b)
}

// NewSchemaFromBytes parses and load a schema from the provided byte array
func NewSchemaFromBytes(bytes []byte) (*JSONSchema, error) {
	l := gojsonschema.NewBytesLoader(bytes)
	schema, err := gojsonschema.NewSchema(l)
	if err != nil {
		return nil, err
	}
	return &JSONSchema{*schema}, nil
}

// Validate the provided data with the schema
func (schema *JSONSchema) Validate(data interface{}) error {
	res, err := schema.inner.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return err
	}
	if !res.Valid() {
		log.Error("JSON is not valid:")
		for _, err := range res.Errors() {
			log.Error(" -- ", err.String())
		}
		return ErrSchemaInvalid
	}
	return nil
}
