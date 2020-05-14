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

func Parse(file string) (b []byte, compose Compose, err error) {
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

		images := value.(map[interface{}]interface{})
		for k, v := range images {
			key := k.(string)

			var js []byte
			js, err = yml.Marshal(&v)
			if err != nil {
				return
			}

			var data map[string]string
			err = yml.Unmarshal(js, &data)
			if err != nil {
				return
			}

			build := Builds{
				Context: ".",
				File:    "Dockerfile",
			}
			if value, ok := data["context"]; ok {
				build.Context = value
			}

			if value, ok := data["file"]; ok {
				build.File = value
			}
			out[key] = build
		}
		delete(m, "images")
	}
	b, err = yml.Marshal(&m)
	return b, Compose{out}, nil

}
