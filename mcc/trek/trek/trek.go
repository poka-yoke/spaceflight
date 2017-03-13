package trek

import "fmt"

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
func Add(redirects, original, final string) (resultingRedirects string, err error) {
	if original == "" || final == "" {
		err = fmt.Errorf("You must specify both original and final")
		return
	}
	newRedirect := fmt.Sprintf("%s %s;\n", original, final)
	resultingRedirects = fmt.Sprintf("%s%s", redirects, newRedirect)
	return
}
