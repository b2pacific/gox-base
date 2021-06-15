package config

import (
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/serialization"
	"io/ioutil"
	"os"
	"strings"
)

func ReadParameterizedYaml(data string, object interface{}, env string) (err error) {

	// Read all environment var and replace it in input string
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) >= 2 {
			data = strings.ReplaceAll(data, "$"+pair[0], pair[1])
		}
	}

	// Read config as map to resolve parameterized variables
	firstMap := map[string]interface{}{}
	err = serialization.ReadYamlFromStringWithEnvVar(data, &firstMap)
	if err != nil {
		return errors.Wrap(err, "could not parse yaml content ["+data+"]", err, nil)
	}

	// Process map in final result map
	newMap := map[string]interface{}{}
	newMap, _ = processMap(firstMap, env)

	yaml, err := serialization.ToYaml(newMap)
	if err != nil {
		return errors.Wrap(err, "could not parse final yaml content ["+data+"]", err, nil)
	} else {
		// fmt.Println(yaml)
	}

	return serialization.ReadYamlFromString(yaml, object)
}

func ReadParameterizedYamlFile(file string, object interface{}, env string) (err error) {
	_data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "could not open file to read ["+file+"]", err, nil)
	}
	return ReadParameterizedYaml(string(_data), object, env)
}

func processMap(input map[string]interface{}, env string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for k, v := range input {
		if val, ok := v.(string); ok {
			out[k], _ = processString(val, env)
		} else if val, ok := v.(map[string]interface{}); ok {
			out[k], _ = processMap(val, env)
		} else if val, ok := v.([]interface{}); ok {
			out[k], _ = processList(val, env)
		}
	}
	return out, nil
}

func processList(input []interface{}, env string) ([]interface{}, error) {
	out := make([]interface{}, 0)
	for _, v := range input {
		if val, ok := v.(string); ok {
			r, _ := processString(val, env)
			out = append(out, r)
		} else if val, ok := v.(map[string]interface{}); ok {
			r, _ := processMap(val, env)
			out = append(out, r)
		} else if val, ok := v.([]interface{}); ok {
			r, _ := processList(val, env)
			out = append(out, r)
		}
	}
	return out, nil
}

func processString(input string, env string) (interface{}, error) {
	if strings.HasPrefix(input, "env:string:") {
		input = strings.Replace(input, "env:string:", "env:", 1)
		p := ParameterizedString(input)
		return p.Get(env)
	} else if strings.HasPrefix(input, "env:bool:") {
		input = strings.Replace(input, "env:bool:", "env:", 1)
		p := ParameterizedBool(input)
		return p.Get(env)
	} else if strings.HasPrefix(input, "env:int:") {
		input = strings.Replace(input, "env:int:", "env:", 1)
		p := ParameterizedInt(input)
		return p.Get(env)
	} else if strings.HasPrefix(input, "env:float:") {
		input = strings.Replace(input, "env:float:", "env:", 1)
		p := ParameterizedFloat(input)
		return p.Get(env)
	} else {
		return input, nil
	}
}
