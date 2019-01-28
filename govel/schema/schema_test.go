package schema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

var testSchema = []byte(`
{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"type": "object",
	"properties": {
		"root": {
			"$ref": "#/definitions/root"
		}
	},
	"required": [
		"root"
	],
	"definitions": {
		"root": {
			"type": "object",
			"properties": {
				"A": {
					"type": "string"
				},
				"B": {
					"type": "integer"
				}
			},
			"required": [
				"A"
			]
		}
	}
}`)

type JSONSchemaTestSuite struct {
	suite.Suite
}

func TestJSONSchema(t *testing.T) {
	suite.Run(t, new(JSONSchemaTestSuite))
}

func (s *JSONSchemaTestSuite) TestBadSchema() {
	schema, err := NewSchemaFromBytes([]byte("{foobar"))
	s.Nil(schema)
	s.Error(err)
}

func (s *JSONSchemaTestSuite) TestValid() {
	schema, err := NewSchemaFromBytes(testSchema)
	s.NoError(err)
	s.NotNil(schema)

	data := map[string]interface{}{"root": map[string]interface{}{"A": "abc", "B": 12}}
	err = schema.Validate(data)
	s.NoError(err)
}

func (s *JSONSchemaTestSuite) TestInvalid() {
	schema, err := NewSchemaFromBytes(testSchema)
	s.NoError(err)
	s.NotNil(schema)

	data := map[string]interface{}{"root": map[string]interface{}{"A": "abc", "B": "12"}}
	err = schema.Validate(data)
	s.Error(err)
}

func (s *JSONSchemaTestSuite) TestLoadFromFile() {
	path := "test_schema.json"
	f, err := os.Create(path)
	if err != nil {
		s.FailNow("cannot create temporary schema file")
	}
	_, err = f.Write(testSchema)
	f.Close()
	if err != nil {
		s.FailNow("cannot create temporary schema file")
	}
	defer os.Remove(path)
	schema, err := LoadSchemaFromFile(path)
	if s.NoError(err) {
		data := map[string]interface{}{"root": map[string]interface{}{"A": "abc", "B": 12}}
		err = schema.Validate(data)
		s.NoError(err)
	}
}

func (s *JSONSchemaTestSuite) TestV2841() {
	s.NotPanics(func() { V2841() })
	sch := V2841()
	s.NotNil(sch)
}
