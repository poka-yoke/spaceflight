package fido

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Entry represents a configuration entry, key value pair
// It's an interface so we can have many kinds of entries
type Entry interface {
	Key() string
	Value() string
}

// Section represents a configuration section with entries and/or other sections
type Section struct {
	Entries  []Entry
	Sections []Section
}

// AddSection adds a section to a Section
func (s Section) AddSection(ns Section) Section {
	s.Sections = append(s.Sections, ns)
	return s
}

// AddEntry adds an Entry to a Section
func (s Section) AddEntry(ne Entry) Section {
	s.Entries = append(s.Entries, ne)
	return s
}

func (s Section) String() string {
	var result string
	for _, en := range s.Entries {
		result += fmt.Sprintf("%s: %s\n", en.Key(), en.Value())
	}
	for _, sec := range s.Sections {
		result += fmt.Sprintf("{\n%s\n}\n", sec)
	}
	return result
}

// StringEntry is an Entry formed by a string key and string value
type StringEntry struct {
	key   string
	value string
}

// NewStringEntry returns a new string with key as Key and value as Value
func NewStringEntry(key, value string) StringEntry {
	var s StringEntry
	s.key = key
	s.value = value
	return s
}

// Key returns the key as string
func (s StringEntry) Key() string { return s.key }

// Value returns the value as a string
func (s StringEntry) Value() string { return s.value }

// ListStringEntry is an Entry formed by a string key and
// an array of strings as value
type ListStringEntry struct {
	key       string
	value     []string
	Separator string
}

// Key returns the key as string
func (s ListStringEntry) Key() string { return s.key }

// Value returns the value as a string
func (s ListStringEntry) Value() string { return strings.Join(s.value, s.Separator) }

// ReadFromPipe returns a string containing stdin contents
func ReadFromPipe() (string, error) {
	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	out := ""
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))
		out += string(buf)
		if err != nil && err != io.EOF {
			return out, err
		}

		log.Println("Bytes:", nBytes, "Chunks:", nChunks)
	}
	return out, nil
}
