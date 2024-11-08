package goutils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JsonToProperties converts simple, flat json, {"abc":"123","def":"456"} to a string-valued map.
// Field value types in json input are assumed to be of type string. If otherwise, an error is returned.
// A map value of nil is returned with empty json.
func JsonPropertiesToMap(jsonStr []byte) (map[string]string, error) {
	// Unmarshal the json string into map[string]interface{}
	var i interface{}
	var err error

	// Return nil on no properties
	if len(jsonStr) == 0 {
		return nil, nil
	}

	if err = json.Unmarshal(jsonStr, &i); err != nil {
		return nil, fmt.Errorf("properties unmarshal error: %s", err)
	}
	// Interpret this as a case of no properties
	if i == nil {
		return nil, nil
	}

	// Named property map but we don't know what the types are
	unknownTypesMap := i.(map[string]interface{})

	// Create map[string]string to return since package properties can only be of type string. If a property
	// type is a non-string, the package definition is incorrect.
	properties := make(map[string]string)
	for k, v := range unknownTypesMap {
		switch v := v.(type) {
		case string:
			properties[k] = v
		default:
			vt := reflect.TypeOf(v).String()
			return nil, fmt.Errorf("property, %s, invalid type, %v", k, vt)
		}
	}

	return properties, nil
}
