package narwhal_lib

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func Breakdown() (map[string][]byte, error) {
	b, err := ioutil.ReadFile("./compose_test/sample.yml")
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	var ret = make(map[string][]byte)
	for k, v := range data {

		out, err := yaml.Marshal(&v)
		if err != nil {
			return nil, err
		}
		ret[k] = out
	}
	return ret, nil
}

func TestParseConfig_both_context_and_file(t *testing.T) {
	m, err := Breakdown()
	if err != nil {
		t.Error(err)
	}
	subj := m["rocket"]

	expected := Builds{
		Context: "random",
		File:    "df",
	}
	actual, err := parseConfig(subj)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("not equal")
	}
}

func TestParseConfig_string(t *testing.T) {
	m, err := Breakdown()
	if err != nil {
		t.Error(err)
	}
	subj := m["golang"]

	expected := Builds{
		Context: "welp",
		File:    "Dockerfile",
	}
	actual, err := parseConfig(subj)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("not equal")
	}
}

func TestParseConfig_only_file(t *testing.T) {
	m, err := Breakdown()
	if err != nil {
		t.Error(err)
	}
	subj := m["dotnet"]

	expected := Builds{
		Context: ".",
		File:    "RandomFile",
	}
	actual, err := parseConfig(subj)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("not equal")
	}
}

func TestParseConfig_only_context(t *testing.T) {
	m, err := Breakdown()
	if err != nil {
		t.Error(err)
	}
	subj := m["node"]

	expected := Builds{
		Context: "golang",
		File:    "Dockerfile",
	}
	actual, err := parseConfig(subj)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("not equal")
	}
}

func TestParseConfig_empty(t *testing.T) {
	m, err := Breakdown()
	if err != nil {
		t.Error(err)
	}
	subj := m["ror"]

	expected := Builds{
		Context: ".",
		File:    "Dockerfile",
	}
	actual, err := parseConfig(subj)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("not equal")
	}
}
