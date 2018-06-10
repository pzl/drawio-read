package main

import (
	"fmt"
	"os"

	// byte array operations
	"bytes"
	"io"
	"io/ioutil"

	//byte parsing
	"encoding/binary"

	// drawio-related parsing
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"net/url"
)

const pngHead = "\x89PNG\r\n\x1a\n"
const ztxtHead = "zTXt"
const endHead = "IEND"

type Chunk struct {
	Length uint32
	Type   string
	Data   []byte
	Crc32  []byte
}

type ZtChunk struct {
	Name string
	Text string
}

type MXFile struct {
	Diagram string `xml:"diagram"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func deflate(data []byte) (string, error) {
	reader := flate.NewReader(bytes.NewReader(data))
	defer reader.Close()

	uncompressed, err := ioutil.ReadAll(reader)
	return string(uncompressed), err
}

func (c Chunk) parse_ztxt() ZtChunk {
	var z ZtChunk
	nul_pos := bytes.Index(c.Data, []byte{0})
	z.Name = string(c.Data[:nul_pos])
	// there are two NUL bytes. One for separation, other to mark compression
	data, err := deflate(c.Data[nul_pos+2:])
	check(err)
	z.Text = data

	return z
}

func valid_png(f *os.File) (bool, error) {
	head := make([]byte, 8)
	_, err := io.ReadFull(f, head)
	if err != nil {
		return false, err
	}
	return string(head) == pngHead, nil
}

func read_section(f *os.File) Chunk {
	var c Chunk
	buf4 := make([]byte, 4)

	io.ReadFull(f, buf4)
	c.Length = binary.BigEndian.Uint32(buf4)

	io.ReadFull(f, buf4)
	c.Type = string(buf4)

	databuf := make([]byte, c.Length)
	io.ReadFull(f, databuf)
	c.Data = databuf

	io.ReadFull(f, buf4)
	c.Crc32 = buf4

	return c
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Argument required: path to PNG file")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	valid, err := valid_png(file)
	check(err)
	if !valid {
		fmt.Println("Not valid PNG file")
		os.Exit(1)
	}

	for {
		chunk := read_section(file)
		if chunk.Type == ztxtHead {
			z := chunk.parse_ztxt()
			if z.Name == "mxGraphModel" {
				data := z.Text
				data, _ = url.QueryUnescape(data)
				var mxfile MXFile
				xml.Unmarshal([]byte(data), &mxfile)
				decoded, err := base64.StdEncoding.DecodeString(mxfile.Diagram)
				check(err)

				uncompressed, err := deflate(decoded)
				check(err)
				final, _ := url.QueryUnescape(uncompressed)
				fmt.Println(final)
				os.Exit(0)
			}
			break
		} else if chunk.Type == endHead {
			fmt.Println("reached the end")
			break
		}
	}
}
