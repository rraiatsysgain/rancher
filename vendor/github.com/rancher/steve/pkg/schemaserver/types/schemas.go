package types

import (
	"strings"

	"github.com/rancher/wrangler/pkg/schemas"
	"github.com/sirupsen/logrus"
)

type APISchemas struct {
	InternalSchemas *schemas.Schemas
	Schemas         map[string]*APISchema
	index           map[string]*APISchema
}

func EmptyAPISchemas() *APISchemas {
	return &APISchemas{
		InternalSchemas: schemas.EmptySchemas(),
		Schemas:         map[string]*APISchema{},
		index:           map[string]*APISchema{},
	}
}

func (a *APISchemas) MustAddSchema(obj APISchema) *APISchemas {
	err := a.AddSchema(obj)
	if err != nil {
		logrus.Fatalf("failed to add schema: %v", err)
	}
	return a
}

func (a *APISchemas) MustImportAndCustomize(obj interface{}, f func(*APISchema)) {
	schema, err := a.InternalSchemas.Import(obj)
	if err != nil {
		panic(err)
	}
	apiSchema := &APISchema{
		Schema: schema,
	}
	a.Schemas[schema.ID] = apiSchema
	f(apiSchema)
}

func (a *APISchemas) AddSchemas(schema *APISchemas) error {
	for _, schema := range schema.Schemas {
		if err := a.AddSchema(*schema); err != nil {
			return err
		}
	}
	return nil
}

func (a *APISchemas) AddSchema(schema APISchema) error {
	if err := a.InternalSchemas.AddSchema(*schema.Schema); err != nil {
		return err
	}
	schema.Schema = a.InternalSchemas.Schema(schema.ID)
	a.Schemas[schema.ID] = &schema
	a.index[strings.ToLower(schema.ID)] = &schema
	a.index[strings.ToLower(schema.PluralName)] = &schema
	return nil
}

func (a *APISchemas) LookupSchema(name string) *APISchema {
	s, ok := a.Schemas[name]
	if ok {
		return s
	}
	return a.index[strings.ToLower(name)]
}
