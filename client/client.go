package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	es7 "github.com/elastic/go-elasticsearch/v7"
	es8 "github.com/elastic/go-elasticsearch/v8"
)

type MappData struct {
	Mappings Mappings `json:"mappings"`
}

type Mappings struct {
	Properties map[string]ExtraEarn `json:"properties"`
}
type ExtraEarn map[string]interface{}

func NewClient(cli interface{}, above int) *MappingClient {
	v7, ok := cli.(*es7.Client)
	if ok {
		return &MappingClient{
			v7:    v7,
			above: above,
		}
	}
	v8, ok := cli.(*es8.Client)
	if ok {
		return &MappingClient{
			v8:    v8,
			above: above,
		}
	}
	return nil
}

type MappingClient struct {
	v7    *es7.Client
	v8    *es8.Client
	above int
}

func (m *MappingClient) Create(index string, data interface{}) (string, error) {
	out, err := m.genMapping(data)
	d, _ := json.Marshal(out)
	if err != nil {
		return "", err
	}
	if m.v7 != nil {
		resp, err := m.v7.Indices.Create(index, m.v7.Indices.Create.WithBody(bytes.NewReader(d)))
		if err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			return string(data), err
		}
		return "", err
	}
	if m.v8 != nil {
		resp, err := m.v8.Indices.Create(index, m.v8.Indices.Create.WithBody(bytes.NewReader(d)))
		if err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			return string(data), err
		}
		return "", err
	}
	return "", fmt.Errorf("no es client")
}
func (m *MappingClient) Put(index string, data interface{}) (string, error) {
	out, err := m.genMapping(data)
	if err != nil {
		return "", err
	}
	if m.v7 != nil {
		d, _ := json.Marshal(out.Mappings)
		resp, err := m.v7.Indices.PutMapping(bytes.NewReader(d), m.v7.Indices.PutMapping.WithIndex(index))
		if err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			return string(data), err
		}
		return "", err
	}
	if m.v8 != nil {
		d, _ := json.Marshal(out.Mappings)
		resp, err := m.v8.Indices.PutMapping([]string{index}, bytes.NewReader(d))
		if err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			return string(data), err
		}
		return "", err
	}
	return "", fmt.Errorf("no es client")
}

func (m *MappingClient) genMapping(model_struct any) (out *MappData, out_err error) {
	if m.above <= 0 {
		m.above = 256
	}
	defer func() {
		if err := recover(); err != nil {
			out = nil
			out_err = fmt.Errorf("%v", err)
		}
	}()
	types := reflect.TypeOf(model_struct)
	if types.Kind() == reflect.Pointer {
		types = types.Elem()
	}
	data := &MappData{
		Mappings: Mappings{
			Properties: make(map[string]ExtraEarn),
		},
	}
	for i := 0; i < types.NumField(); i++ {
		tag := types.Field(i).Tag.Get("json")
		es_type := ""
		data_type := fmt.Sprintf("%v", types.Field(i).Type)
		switch data_type {
		case "uint8", "uint16", "uint32", "int8", "int16", "int32", "int":
			es_type = "integer"
		case "uint64", "int64":
			es_type = "long"
		case "float32":
			es_type = "float"
		case "float64":
			es_type = "double"
		case "time.Time", "*time.Time":
			es_type = "date"
		case "string":
			es_type = "keyword"
		default:
			return nil, fmt.Errorf("unknow type[%s]", data_type)
		}
		tmp := make(map[string]interface{})
		tmp["type"] = es_type
		if es_type == "keyword" {
			tmp["ignore_above"] = m.above
		}
		data.Mappings.Properties[tag] = tmp
	}
	return data, nil
}
func (m *MappingClient) GetMapping(index string) (string, error) {
	if m.v7 != nil {
		resp, err := m.v7.Indices.GetMapping(m.v7.Indices.GetMapping.WithIndex(index))
		if err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			return string(data), err
		}
		return "", err
	}
	if m.v8 != nil {
		resp, err := m.v8.Indices.GetMapping(m.v8.Indices.GetMapping.WithIndex(index))
		if err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			return string(data), err
		}
		return "", err
	}
	return "", nil
}
