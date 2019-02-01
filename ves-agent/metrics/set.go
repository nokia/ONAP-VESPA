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

package metrics

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"github.com/nokia/onap-vespa/govel"

	log "github.com/sirupsen/logrus"
)

// MeasKeys is a map of keys / values fore metrics structure
// used in filtering. Some kind of composite key
type MeasKeys map[string]string

func (keys MeasKeys) String() string {
	bld := strings.Builder{}
	for k, v := range keys {
		if _, err := fmt.Fprintf(&bld, "%s: %s,", k, v); err != nil {
			log.Panic(err)
		}
	}
	s := bld.String()
	return s[:len(s)-1]
}

// JSONObjectKeys transform this set of key/values pairs into a list
// of `govel.Key`
func (keys MeasKeys) JSONObjectKeys() []govel.Key {
	jKeys := []govel.Key{}
	for k, v := range keys {
		s := v // Copy string before getting a reference to it
		jKeys = append(jKeys, govel.Key{KeyName: k, KeyValue: &s})
	}
	return jKeys
}

// MatchJSONObjectKeys checks if this set of key/value pairs match exactly  with
// the given list of `govel.Key`
func (keys MeasKeys) MatchJSONObjectKeys(jKeys []govel.Key) bool {
	if len(jKeys) != len(keys) {
		return false
	}
	for _, k := range jKeys {
		if v, ok := keys[k.KeyName]; !ok || v != *k.KeyValue {
			return false
		}
	}
	return true
}

// EventMeasurementSet is a list of VES Measurements events
type EventMeasurementSet []*govel.EventMeasurements

// Batch converts the measurement set into a batch of events
func (set EventMeasurementSet) Batch() govel.Batch {
	batch := make([]govel.Event, len(set))
	for i, evt := range set {
		batch[i] = evt
	}
	return batch
}

// MeasurementFactory is used to create new measurement events
type MeasurementFactory interface {
	Create(vmID string, timestamp time.Time) (*govel.EventMeasurements, error)
}

// MeasurementFactoryFunc is a function which can be used as a MeasurementFactory
type MeasurementFactoryFunc func(vmID string, timestamp time.Time) (*govel.EventMeasurements, error)

// Create is called to create a new `govel.EventMeasurements`
func (f MeasurementFactoryFunc) Create(vmID string, timestamp time.Time) (*govel.EventMeasurements, error) {
	return f(vmID, timestamp)
}

// EventMeasurementSetBuilder is an utility to construct a set of measurement events
type EventMeasurementSetBuilder struct {
	set      EventMeasurementSet // The measurement set under construction
	provider MeasurementFactory  // A provider to create new govel.EventMeasurements
	valid    bool                // Can this builder still be used ?
}

// NewEventMeasurementSetBuilder creates a new EventMeasurementSetBuilder using the provided
// measurement factory. If an error occurs during value insertion, the full set will be
// invalidated, and trying to use the EventMeasurementSetBuilder after that will cause a panic
func NewEventMeasurementSetBuilder(fact MeasurementFactory) EventMeasurementSetBuilder {
	return EventMeasurementSetBuilder{
		set:      EventMeasurementSet{},
		provider: fact,
		valid:    true,
	}
}

// checkValid panics if the EventMeasurementSetBuilder is not valid anymore
func (bld *EventMeasurementSetBuilder) checkValid() {
	if !bld.valid {
		log.Panic("EventMeasurementSetBuilder has an illegal state")
	}
}

func (bld *EventMeasurementSetBuilder) invalidate() {
	// In case of error, invalidate this builder
	// It must need to br created again
	bld.set = EventMeasurementSet{}
	bld.valid = false
}

// find returns a pointer to the EventMeasurement in the set for VM with id vmID with the provided timestamp.
// If does not exist, creates it and store it before returning it
func (bld *EventMeasurementSetBuilder) find(vmID string, timestamp time.Time) (*govel.EventMeasurements, error) {
	t := timestamp.UnixNano() / 1000
	for _, evt := range bld.set {
		if evt.LastEpochMicrosec == t && evt.SourceName == vmID {
			return evt, nil
		}
	}
	evt, err := bld.provider.Create(vmID, timestamp)
	if err != nil {
		return nil, err
	}
	bld.set = append(bld.set, evt)
	// evt can be returned since it's a pointer (bld.set holds pointer data)
	return evt, nil
}

// Measurements returns the built EventMeasurementSet
func (bld *EventMeasurementSetBuilder) Measurements() EventMeasurementSet {
	bld.checkValid()
	return bld.set
}

// Set sets the metric given by "fields" for VM with ID "vmID" at time "timestamp" with value "values".
// "keys" are used when and array is encountered to select the entry (or set the entry when missing)
//
// Returns an error or nil. If an error occurs, the EventMeasurementSetBuilder must not be used anymore
// Panics if an error occurred in the previous call
func (bld *EventMeasurementSetBuilder) Set(fields string, vmID string, timestamp time.Time, value float64, keys MeasKeys) error {
	bld.checkValid()
	if len(fields) == 0 {
		return errors.New("Fields cannot be empty")
	}
	if len(vmID) == 0 {
		return errors.New("VmID cannot be empty")
	}
	// Get or create the govel.EventMeasurements
	metric, err := bld.find(vmID, timestamp)
	if err != nil {
		return err
	}
	subfields := strings.Split(fields, ".")
	// Use reflection to set the structure field
	err = setField(reflect.ValueOf(metric), subfields, value, keys)
	if err != nil {
		bld.invalidate()
	}
	return err
}

// SetAdditionalObject insert a metric value in a AdditionalObjects field of a Measurement event
func (bld *EventMeasurementSetBuilder) SetAdditionalObject(vmID, objectName, objectInstance string, timestamp time.Time, value float64, keys MeasKeys) error {
	if len(vmID) == 0 || len(objectName) == 0 || len(objectInstance) == 0 {
		return errors.New("SetAdditionalObject() - Arguments cannot be empty")
	}
	metric, err := bld.find(vmID, timestamp)
	if err != nil {
		return err
	}
	var additionalObject *govel.JSONObject
	for i, obj := range metric.AdditionalObjects {
		if obj.ObjectName == objectName {
			additionalObject = &metric.AdditionalObjects[i]
			break
		}
	}
	if additionalObject == nil {
		metric.AdditionalObjects = append(metric.AdditionalObjects, govel.JSONObject{ObjectName: objectName})
		additionalObject = &metric.AdditionalObjects[len(metric.AdditionalObjects)-1]
	}
	var instance *govel.JSONObjectInstance
	for i, inst := range additionalObject.ObjectInstances {
		if keys.MatchJSONObjectKeys(inst.ObjectKeys) {
			instance = &additionalObject.ObjectInstances[i]
			break
		}
	}
	if instance == nil {
		inst := govel.JSONObjectInstance{
			ObjectKeys:     keys.JSONObjectKeys(),
			ObjectInstance: make(map[string]interface{}),
		}
		additionalObject.ObjectInstances = append(additionalObject.ObjectInstances, inst)
		instance = &additionalObject.ObjectInstances[len(additionalObject.ObjectInstances)-1]
	}
	instance.ObjectInstance[objectInstance] = value
	return nil
}

/***********************************************************
 *                                                         *
 * Functions below are a set of helpers used to manipulate *
 * measurements structures using some reflection.          *
 * It clearly miss some comments for now :(                *
 *                                                         *
 ***********************************************************/

func setField(parent reflect.Value, subfields []string, value float64, keys MeasKeys) error {
	if !parent.IsValid() {
		return errors.New("Cannot set field, parent is not valid")
	}
	// println(parent.Type().String())
	if parent.Kind() == reflect.Ptr {
		if parent.IsNil() {
			newVal := reflect.New(parent.Type().Elem())
			parent.Set(newVal)
		}
		return setField(parent.Elem(), subfields, value, keys)
	}
	if len(subfields) == 0 {
		switch parent.Kind() {
		case reflect.Float64, reflect.Float32:
			parent.SetFloat(value)
			return nil
		default:
			return fmt.Errorf("Cannot assign a float64 to %s", parent.Type().String())
		}
	} else {
		switch parent.Kind() {
		case reflect.Slice:
			return setSliceField(parent, subfields, value, keys)
		case reflect.Struct:
			field := parent.FieldByName(subfields[0])
			if !field.IsValid() {
				return fmt.Errorf("Invalid field %s on type %s", subfields[0], parent.Type().String())
			}
			return setField(field, subfields[1:], value, keys)
		default:
			return fmt.Errorf("Cannot access subfields of type %s", parent.Type().String())
		}
	}
}

func setSliceField(parent reflect.Value, fields []string, value float64, keys MeasKeys) error {
	// elemType := parent.Type().Elem()
	// f, ok := elemType.FieldByName(fields[0])
	// if !ok {
	// 	return fmt.Errorf("Field %s does not exists on %s", fields[0], parent.Type().String())
	// }
	// if f.Type.Kind() != reflect.Float64 && f.Type.Kind() != reflect.Float32 {
	// 	return fmt.Errorf("Field %s on %s is not a float (is %s)", fields[0], parent.Type().String(), f.Type.String())
	// }

	ref := sliceSearch(parent, keys)
	if !ref.IsValid() {
		ref = sliceAppendNew(parent, keys)
	}
	v := ref.FieldByName(fields[0])
	return setField(v, fields[1:], value, keys)
}

// Search in the slice the element matching the composite key.
// If not found, return a zero reflect.Value
func sliceSearch(parent reflect.Value, keys MeasKeys) reflect.Value {
searchLoop:
	for i := 0; i < parent.Len(); i++ {
		entry := parent.Index(i)
		// compare all key/value pairs from composite key
		for k, v := range keys {
			field := entry.FieldByName(k)
			if !field.IsValid() {
				// Ignore if field does not exist en struct
				continue
			} else if field.String() != v {
				continue searchLoop
			}
		}
		return entry
	}
	return reflect.Value{}
}

// Create and append a new element with proper type and fields
// initialized from composite key fields
func sliceAppendNew(parent reflect.Value, keys MeasKeys) reflect.Value {
	// Create a new Zero value with slice contained type
	obj := reflect.New(parent.Type().Elem()).Elem()
	// Set composite key fields
	for k, v := range keys {
		field := obj.FieldByName(k)
		if !field.IsValid() {
			// Ignore not found fields
			continue
		}
		field.SetString(v)
	}
	// Append new element to slice
	parent.Set(reflect.Append(parent, obj))
	// return a reference to the element. It's really a pointer
	// but a reflected value
	return parent.Index(parent.Len() - 1)
}
