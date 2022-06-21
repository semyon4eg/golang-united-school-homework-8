package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id string `json:"id"`
	Email string `json:"email"`
	Age int `json:"age"`
}

type List []Item

var l List

var i Item

var errNoFile = errors.New("-fileName flag has to be specified")
var errNoOperation = errors.New("-operation flag has to be specified")
var errNoItem = errors.New("-item flag has to be specified")
var errNoId = errors.New("-id flag has to be specified")

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Perform(args Arguments, writer io.Writer) error {

	if fileName, ok := args["fileName"]; ok && fileName == "" {
		return errNoFile
	} else if ok {
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		checkError(err)

		defer file.Close()

		fileInfo, err := os.Stat(fileName)
		checkError(err)
	
		fileSize := fileInfo.Size()
		if fileSize == 0 {
			err = ioutil.WriteFile(fileName, []byte("[]"), 0755)
			checkError(err)
		}
	}

	if op, ok := args["operation"]; ok && op == "" {
		return errNoOperation
	} else if ok {
		switch op {
		case "list":
			list(args, writer)
		case "add":
			return add(args, writer)
		case "findById":
			return findById(args, writer)
		case "remove":
			return remove(args, writer)
		default:
			return fmt.Errorf("Operation %v not allowed!", op)
		}
	}

	return nil
}

func parseArgs() Arguments {

	var args Arguments = make(Arguments, 0)

	var idFlag = flag.String("id", "", "ID")
	var itemFlag = flag.String("item", "", "item")
	var operationFlag = flag.String("operation", "", "operation")
	var fileNameFlag = flag.String("fileName", "", "file name")
	flag.Parse()

	args["id"] = *idFlag
	args["item"] = *itemFlag
	args["operation"] = *operationFlag
	args["fileName"] = *fileNameFlag

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func list(args Arguments, writer io.Writer) {

	f, err := ioutil.ReadFile(args["fileName"])
	checkError(err)

	writer.Write(f)
}

func add(args Arguments, writer io.Writer) error {

	if item, ok := args["item"]; ok && item == "" {
		return errNoItem
	} else if ok {

		err := json.Unmarshal([]byte(item), &i)
		checkError(err)

		f, err := os.ReadFile(args["fileName"])
		checkError(err)

		err = json.Unmarshal(f, &l)
		checkError(err)

		for _, val := range l {
			if val.Id == i.Id {
				writer.Write([]byte("Item with id " + i.Id + " already exists"))
				return nil
			}
		}

		l = append(l, i)
		b, err := json.Marshal(l)
		checkError(err)

		err = ioutil.WriteFile(args["fileName"], b, 0755)
		checkError(err)
	}

	return nil
}

func findById(args Arguments, writer io.Writer) error {

	if id, ok := args["id"]; ok && id == "" {
		return errNoId
	} else if ok {

		f, err := os.ReadFile(args["fileName"])
		checkError(err)

		err = json.Unmarshal(f, &l)
		checkError(err)

		for _, i := range l {
			if i.Id == id {

				b, err := json.Marshal(i)
				checkError(err)

				writer.Write(b)
			}
		}

		writer.Write([]byte(""))
	}

	return nil
}

func remove(args Arguments, writer io.Writer) error {

	if id, ok := args["id"]; ok && id == "" {
		return errNoId
	} else if ok {
		
		f, err := os.ReadFile(args["fileName"])
		checkError(err)

		err = json.Unmarshal(f, &l)
		checkError(err)

		for i, v := range l {
			if v.Id == id {

				l = append(l[:i], l[i+1:]... )

				b, err := json.Marshal(l)
				checkError(err)

				err = ioutil.WriteFile(args["fileName"], b, 0755)
				checkError(err)

				return nil
			}
		}
		
		writer.Write([]byte(fmt.Sprintf("Item with id %v not found", args["id"])))
	}

	return nil
}
