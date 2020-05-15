package narwhal_lib

import (
	yml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Compose struct {
	Images map[string]Builds
}

type Builds struct {
	Context string
	File    string
}

func parse(file string) (b []byte, compose Compose, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	m := make(map[string]interface{})

	b = nil
	compose = Compose{}
	err = nil
	var f []byte
	f, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = yml.Unmarshal(f, &m)
	if err != nil {
		return
	}

	out := make(map[string]Builds)
	//parse compose here
	if value, ok := m["images"]; ok {
		var images map[string][]byte
		images, err = breakdown(value)
		for k, v := range images {
			var build Builds
			build, err = parseConfig(v)
			out[k] = build
		}
		delete(m, "images")
	}
	b, err = yml.Marshal(&m)
	return b, Compose{out}, nil

}

func breakdown(a interface{}) (map[string][]byte, error) {
	b, err := yml.Marshal(&a)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = yml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	var ret = make(map[string][]byte)
	for k, v := range data {

		out, err := yml.Marshal(&v)
		if err != nil {
			return nil, err
		}
		ret[k] = out
	}
	return ret, nil
}

func parseConfig(bytes []byte) (build Builds, err error) {
	build = Builds{
		Context: ".",
		File:    "Dockerfile",
	}
	var data map[string]string
	err = yml.Unmarshal(bytes, &data)
	if err == nil {
		if value, ok := data["context"]; ok {
			build.Context = value
		}

		if value, ok := data["file"]; ok {
			build.File = value
		}
		return
	}
	// here means yaml parsed error
	// we attempt to parse it as a string then
	var sData string
	err = yml.Unmarshal(bytes, &sData)
	if err != nil {
		return
	}
	build.Context = sData
	return
}
