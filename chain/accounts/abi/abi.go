// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Bokerchain/Boker/chain/log"
)

// The ABI holds information about a contract's context and available
// invokable methods. It will allow you to type check function calls and
// packs data accordingly.
type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
}

// JSON returns a parsed ABI interface and error if it failed.
func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

//用符合ABI格式打包给定的方法。 方法调用的数据将由method_id，args0，arg1，... argN组成。 方法ID包含4个字节和参数都是32个字节。
//方法ID是从哈希的前4个字节创建的方法字符串签名。 （signature = baz（uint32，string32））
func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error) {

	// Fetch the ABI of the requested method
	var method Method

	if name == "" {
		method = abi.Constructor
	} else {
		m, exist := abi.Methods[name]
		if !exist {
			return nil, fmt.Errorf("method '%s' not found", name)
		}
		method = m
	}
	arguments, err := method.pack(args...)
	if err != nil {
		return nil, err
	}
	// Pack up the method ID too if not a constructor and return
	if name == "" {
		return arguments, nil
	}
	return append(method.Id(), arguments...), nil
}

func (abi ABI) Unpack(v interface{}, name string, output []byte) (err error) {

	if err = bytesAreProper(output); err != nil {
		return err
	}

	var unpack unpacker
	if method, ok := abi.Methods[name]; ok {
		unpack = method
	} else if event, ok := abi.Events[name]; ok {
		unpack = event
	} else {
		return fmt.Errorf("abi: could not locate named method or event.")
	}

	if unpack.isTupleReturn() {
		return unpack.tupleUnpack(v, output)
	}

	log.Info("Unpack", "output", output)
	return unpack.singleUnpack(v, output)
}

func (abi ABI) InputUnpack(v []interface{}, name string, input []byte) (err error) {

	//判断输入数据是否正确
	if len(input) == 0 {
		return errors.New("abi: unmarshalling empty input")
	} else if len(input)%32 != 0 {
		return errors.New("abi: improperly formatted input")
	}

	//得到abi的方法信息
	if method, ok := abi.Methods[name]; ok {

		//判断输入接口是否和参数数量一致
		if len(v) != len(abi.Methods[name].Inputs) {
			return errors.New("abi: methods count not equal to interface")
		}

		//根据参数数量判断解码方式
		if len(method.Inputs) <= 1 {
			//单参数解码
			return method.singleInputUnpack(v[0], input)
		} else {
			//多参数解码
			return method.multInputUnpack(v, input)
		}
	}
	return errors.New("abi: could not locate named method")
}

func (abi *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type      string
		Name      string
		Constant  bool
		Indexed   bool
		Anonymous bool
		Inputs    []Argument
		Outputs   []Argument
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	abi.Methods = make(map[string]Method)
	abi.Events = make(map[string]Event)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			abi.Constructor = Method{
				Inputs: field.Inputs,
			}
		// empty defaults to function according to the abi spec
		case "function", "":
			abi.Methods[field.Name] = Method{
				Name:    field.Name,
				Const:   field.Constant,
				Inputs:  field.Inputs,
				Outputs: field.Outputs,
			}
		case "event":
			abi.Events[field.Name] = Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs,
			}
		}
	}

	return nil
}
