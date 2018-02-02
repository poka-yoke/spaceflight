package trek

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// RFC1123HostnameRegExp : Hostname Regular expression according with
// RFC1123HostnameRegExp.
const RFC1123HostnameRegExp = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]|\\*)\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])$"

// SimpleURL : URL Regular expression.
const SimpleURL = "^http(|s)://(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]|\\*)\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])(|/.*)$"

// IsURL returns true if argument can be considered a URL.
func IsURL(url string) bool {
	m, err := regexp.MatchString(SimpleURL, url)
	if err != nil {
		return false
	}
	return m
}

// IsHostname returns true if argument can be considered a hostname.
func IsHostname(name string) bool {
	m, err := regexp.MatchString(RFC1123HostnameRegExp, name)
	if err != nil {
		return false
	}
	return m
}

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

// Add adds a new redirect line for original to final, and returns it.
func Add(initialContent, original, final string) (
	result string, err error,
) {
	if original == "" || final == "" {
		err = fmt.Errorf("You must specify both original and final")
		return
	}
	newContent := ""
	if IsHostname(original) && IsURL(final) {
		newContent = fmt.Sprintf(
			"\nserver {\n"+
				"\tlisten 80;\n"+
				"\tserver_name\t%s;\n"+
				"\treturn 301\t%s;\n"+
				"}\n",
			original,
			final,
		)
	} else {
		newContent = fmt.Sprintf("%s %s;\n", original, final)
	}
	result = fmt.Sprintf("%s%s", initialContent, newContent)
	return
}
