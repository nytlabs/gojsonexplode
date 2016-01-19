package gojsonexplode

// TODO: cannot have delimiter be " or \

import (
	"encoding/json"
	"errors"
	"strconv"
)

func explodeList(l []interface{}, parent string, delimiter string, depth int) (map[string]interface{}, error) {
	var err error
	var key string
	j := make(map[string]interface{})
	for k, i := range l {
		if len(parent) > 0 {
			key = parent + delimiter + strconv.Itoa(k)
		} else {
			key = strconv.Itoa(k)
		}

		if depth == 0 {
			j[key] = i
			continue
		}

		switch v := i.(type) {
		case nil:
			j[key] = v
		case int:
			j[key] = v
		case float64:
			j[key] = v
		case string:
			j[key] = v
		case bool:
			j[key] = v
		case []interface{}:
			out := make(map[string]interface{})
			out, err = explodeList(v, key, delimiter, (depth - 1))
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				j[newkey] = value
			}
		case map[string]interface{}:
			out := make(map[string]interface{})
			out, err = explodeMap(v, key, delimiter, (depth - 1))
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				j[newkey] = value
			}
		default:
			// do nothing
		}
	}
	return j, nil
}

func explodeMap(m map[string]interface{}, parent string, delimiter string, depth int) (map[string]interface{}, error) {
	var err error

	j := make(map[string]interface{})
	for k, i := range m {
		if len(parent) > 0 {
			k = parent + delimiter + k
		}

		if depth == 0 {
			j[k] = i
			continue
		}

		switch v := i.(type) {
		case nil:
			j[k] = v
		case int:
			j[k] = v
		case float64:
			j[k] = v
		case string:
			j[k] = v
		case bool:
			j[k] = v
		case []interface{}:
			out := make(map[string]interface{})
			out, err = explodeList(v, k, delimiter, (depth - 1))
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				j[key] = value
			}
		case map[string]interface{}:
			out := make(map[string]interface{})
			out, err = explodeMap(v, k, delimiter, (depth - 1))
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				j[key] = value
			}
		default:
			//nothing
		}
	}
	return j, nil
}

// Explodejson takes in a  nested JSON as a byte array and a delimiter and returns an
// exploded/flattened json byte array
func Explodejson(b []byte, d string, depth int) ([]byte, error) {
	var input interface{}
	var exploded map[string]interface{}
	var out []byte
	var err error
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err
	}
	switch t := input.(type) {
	case map[string]interface{}:
		exploded, err = explodeMap(t, "", d, depth)
		if err != nil {
			return nil, err
		}
	case []interface{}:
		exploded, err = explodeList(t, "", d, depth)
		if err != nil {
			return nil, err
		}
	default:
		// How did we get here? It is impossible!!
		return nil, errors.New("Possible error in JSON")
	}
	out, err = json.Marshal(exploded)
	if err != nil {
		return nil, err
	}
	return out, nil

}

// Explodejsonstr explodes a nested JSON string to an unnested one
// parameters to pass to the function are
// * s: the JSON string
// * d: the delimiter to use when unnesting the JSON object.
// * depth : the desired depth of nesting
//   -1 = no depth limit
//    0 = no nesting ( returns same as input )
//    1 = only first parent will be nested, e.g. company.address
//    2 = up until second child, e.g. company.adress.street
//    3 = ...
// Set to -1 if depth limit is not desired
//
// {"person":{"name":"Joe", "address":{"street":"123 Main St."}}}
// explodes to:
// {"person.name":"Joe", "person.address.street":"123 Main St."}
func Explodejsonstr(s string, d string, depth int) (string, error) {

	if depth == 0 {
		return s, nil
	}

	b := []byte(s)
	out, err := Explodejson(b, d, depth)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
