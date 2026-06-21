package goutils

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

func YamlToJson(yamlStr []byte) ([]byte, error) {
	jsonStr, err := yaml.YAMLToJSON(yamlStr)
	if err != nil {
		return nil, err
	}
	return jsonStr, nil
}

func JsonToYaml(jsonStr []byte) ([]byte, error) {
	yamlStr, err := yaml.JSONToYAML(jsonStr)
	if err != nil {
		return nil, err
	}
	return yamlStr, nil
}

func YamlToJsonString(yamlStr string) (string, error) {
	jsonStr, err := yaml.YAMLToJSON([]byte(yamlStr))
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func JsonToYamlString(jsonStr string) (string, error) {
	yamlStr, err := yaml.JSONToYAML([]byte(jsonStr))
	if err != nil {
		return "", err
	}
	return string(yamlStr), nil
}

func YamlToJsonIndent(yamlStr []byte, prefix, indent string) ([]byte, error) {
	jsonStr, err := yaml.YAMLToJSON(yamlStr)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(jsonStr, &out); err != nil {
		return nil, err
	}
	indentedJsonStr, err := json.MarshalIndent(out, prefix, indent)
	if err != nil {
		return nil, err
	}
	return indentedJsonStr, nil
}

func JsonToYamlIndent(jsonStr []byte, prefix, indent string) ([]byte, error) {
	var out any
	if err := json.Unmarshal(jsonStr, &out); err != nil {
		return nil, err
	}
	indentedJsonStr, err := json.MarshalIndent(out, prefix, indent)
	if err != nil {
		return nil, err
	}
	yamlStr, err := yaml.JSONToYAML(indentedJsonStr)
	if err != nil {
		return nil, err
	}
	return yamlStr, nil
}
