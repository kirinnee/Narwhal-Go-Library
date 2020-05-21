package narwhal_lib

import (
	"encoding/json"
	yml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Compose struct {
	Images map[string]Builds
}

type Builds struct {
	Context string
	File    string
}

func parse(file string) (b []byte, compose Compose, stack string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	m := make(map[string]interface{})

	b = nil
	compose = Compose{}
	err = nil
	stack = ""

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

	changeOut := make(map[string]string)

	if value, ok := m["stack"]; ok {
		stack = value.(string)
		delete(m, "stack")
	}

	//parse compose here
	if value, ok := m["images"]; ok {
		var images map[string][]byte
		images, err = breakdown(value)
		for k, v := range images {
			var build Builds
			build, err = parseConfig(v)

			imageName := k
			if !strings.Contains(imageName, ":") {
				versionNumber, ee := getVersionNumber(k)
				if ee != nil {
					return
				}
				imageName = k + ":" + strconv.Itoa(versionNumber)
				changeOut[k] = imageName
			}
			out[imageName] = build
		}
		delete(m, "images")
	}

	//parse compose here
	if value, ok := m["services"]; ok {
		var services map[string]map[string]interface{}
		services, err = breakdownToMap(value)
		for k, v := range services {
			if image, ok2 := v["image"]; ok2 {
				i := image.(string)
				if imageName, ok3 := changeOut[i]; ok3 {
					services[k]["image"] = imageName
				}
			}
		}
		m["services"] = services
	}

	b, err = yml.Marshal(&m)
	return b, Compose{out}, stack, nil

}

func getVersionNumber(k string) (i int, err error) {
	// enabled versioning
	var version = make(map[string]int)
	versionNumber := 0
	bytes, e := ioutil.ReadFile("./.narwhal_states")
	if e == nil {
		err = json.Unmarshal(bytes, &version)
		if err != nil {
			_, err = os.Create("./.narwhal_states")
			if err != nil {
				return
			}
		}
	}

	if vn, okok := version[k]; okok {
		versionNumber = vn + 1
	}
	version[k] = versionNumber

	b, err := json.Marshal(&version)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("./.narwhal_states", b, 0644)
	if err != nil {
		return
	}
	return versionNumber, nil
}

func breakdownToMap(a interface{}) (map[string]map[string]interface{}, error) {
	b, err := yml.Marshal(&a)
	if err != nil {
		return nil, err
	}
	var data map[string]map[string]interface{}
	err = yml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return data, nil

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
