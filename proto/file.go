package proto

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Field struct {
	Repeated bool
	Optional bool
	Index    int
	Name     string
	Type     string
}

type Message struct {
	Fields map[string]Field
	Name   string
}

type Desc struct {
	Imports  []string
	Name     string
	Messages map[string]Message
}

type MessageLocation struct {
	Name string
	File string
}

var DefaultTypes = []string{
	"int32", "int64", "uint32", "uint64", "sint32", "sint64", "bool", "fixed64", "sfixed64", "double", "string", "fixed32", "sfixed32", "float",
}

var Messages = map[MessageLocation]Message{}

func GetDescriptors(path string) ([]Desc, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	descs := make([]Desc, len(files))

	for i, file := range files {
		fileP := filepath.Join(path, file.Name())
		f, err := os.Open(fileP)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(f)

		descriptor := Desc{
			Name:     file.Name(),
			Messages: make(map[string]Message),
		}

		curMessage := Message{
			Fields: make(map[string]Field),
		}

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			words := regexp.MustCompile(`\s+`).Split(line, -1)
			if len(words) == 1 {
				continue
			}

			switch words[0] {
			case "message":
				if len(curMessage.Name) > 0 {
					Messages[MessageLocation{
						Name: curMessage.Name,
						File: file.Name(),
					}] = curMessage
					descriptor.Messages[curMessage.Name] = curMessage
					curMessage = Message{
						Fields: make(map[string]Field),
					}
				}
				curMessage.Name = words[1]
			case "syntax":
			case "import":
				runes := []rune(words[1])
				descriptor.Imports = append(descriptor.Imports, string(runes[1:len(runes)-2]))
			case "repeated":
				runes := []rune(words[4])
				index, _ := strconv.Atoi(string(runes[:len(runes)-1]))
				curMessage.Fields[words[2]] = Field{
					Repeated: true,
					Optional: false,
					Name:     words[2],
					Type:     words[1],
					Index:    index,
				}
			case "required":
				runes := []rune(words[4])
				index, _ := strconv.Atoi(string(runes[:len(runes)-1]))
				curMessage.Fields[words[2]] = Field{
					Repeated: false,
					Optional: false,
					Name:     words[2],
					Type:     words[1],
					Index:    index,
				}
			case "optional":
				runes := []rune(words[4])
				index, _ := strconv.Atoi(string(runes[:len(runes)-1]))
				curMessage.Fields[words[2]] = Field{
					Repeated: false,
					Optional: false,
					Name:     words[2],
					Type:     words[1],
					Index:    index,
				}
			}
		}
		if len(curMessage.Name) > 0 {
			Messages[MessageLocation{
				Name: curMessage.Name,
				File: file.Name(),
			}] = curMessage
			descriptor.Messages[curMessage.Name] = curMessage
		}
		descs[i] = descriptor
	}

	return descs, nil
}
