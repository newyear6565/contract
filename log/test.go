package main

import (
	"fmt"
	"time"
	"os"
	"strings"
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"compress/gzip"
	"compress/zlib"
	"compress/flate"
	"github.com/DataDog/zstd"
	"github.com/pierrec/lz4"
)

func timeS() int64 {
	return time.Now().UnixNano()
}

func timeE(start int64) int64 {
	return time.Now().UnixNano() - start
}

func lz4_func(toCompress []byte) {
	// toCompress := []byte(fileContent)
	fmt.Println("==================   lz4   ======================")
	compressed := make([]byte, len(toCompress))

	//compress
	start := timeS()
	l1, err := lz4.CompressBlockHC(toCompress, compressed, 0)
	span := timeE(start)
	if err != nil {
		panic(err)
	}
	fmt.Println("Data length")
	fmt.Println(len(toCompress))
	fmt.Println("compressed Data length")
	fmt.Println(l1)

	//decompress
	decompressed := make([]byte, len(toCompress))
	start = timeS()
	l2, err := lz4.UncompressBlock(compressed[:l1], decompressed)
	span2 := timeE(start)
	if err != nil {
		panic(err)
	}
	fmt.Println("\ndecompressed Data length")
	fmt.Println(l2)

	if isSame(toCompress, decompressed) {
		fmt.Println("OK \n\n")
	}
	record := fmt.Sprintf("%d,%d,%d,%d", len(toCompress), l1, span, span2)
	fmt.Println(record)
	fmt.Println([]byte(record))
	writeWithOs(lz4File, []byte(record))
}

func zstd_func(toCompress []byte) {
	fmt.Println("==================   zstd   ======================")
	start := timeS()
	dst,_ := zstd.Compress(nil, toCompress)
	span := timeE(start)
	fmt.Println("Data length")
	fmt.Println(len(toCompress))
	fmt.Println("compressed Data length")
	fmt.Println(len(dst))

	start = timeS()
	src,_ := zstd.Decompress(nil, dst)
	span2 := timeE(start)
	fmt.Println("\ndecompressed Data length")
	fmt.Println(len(src))

	if isSame(toCompress, src) {
		fmt.Println("OK \n\n")
	}

	record := fmt.Sprintf("%d,%d,%d,%d", len(toCompress), len(dst), span, span2)
	fmt.Println(record)
	fmt.Println([]byte(record))
	writeWithOs(zstdFile, []byte(record))
}

func zlib_func(toCompress []byte) {
	fmt.Println("==================   zlib   ======================")
	start := timeS()
	dst := DoZlibCompress(toCompress)
	span := timeE(start)
	fmt.Println("Data length")
	fmt.Println(len(toCompress))
	fmt.Println("compressed Data length")
	fmt.Println(len(dst))

	start = timeS()
	src := DoZlibUnCompress(dst)
	span2 := timeE(start)
	fmt.Println("\nuncompressed Data length")
	fmt.Println(len(src))
	if isSame(toCompress, src) {
		fmt.Println("OK \n\n")
	}
	record := fmt.Sprintf("%d,%d,%d,%d", len(toCompress), len(dst), span, span2)
	fmt.Println(record)
	fmt.Println([]byte(record))
	writeWithOs(zlibFile, []byte(record))
}
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func gzip_func(toCompress []byte) {
	fmt.Println("==================   gzip   ======================")
	start := timeS()
	dst, _ := GzipEncode(toCompress)
	span := timeE(start)
	fmt.Println("Data length")
	fmt.Println(len(toCompress))
	fmt.Println("compressed Data length")
	fmt.Println(len(dst))
	fmt.Println("compress time nano")
	fmt.Println(span)

	start = timeS()
	src, _ := GzipDecode(dst)
	span2 := timeE(start)
	fmt.Println("\nuncompressed Data length")
	fmt.Println(len(src))
	fmt.Println("Uncompress time nano")
	fmt.Println(span2)
	if isSame(toCompress, src) {
		fmt.Println("OK \n\n")
	}
	record := fmt.Sprintf("%d,%d,%d,%d", len(toCompress), len(dst), span, span2)
	fmt.Println(record)
	fmt.Println([]byte(record))
	writeWithOs(gzipFile, []byte(record))
}

func flateEncode(in []byte) ([]byte, error) {
	var (
		out []byte
		err error
	)

	buffer := bytes.NewBuffer(nil)
	writer, _ := flate.NewWriter(buffer, flate.BestCompression)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func flateDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func flate_func(toCompress []byte) {
	fmt.Println("==================   flate   ======================")
	start := timeS()
	dst, _ := flateEncode(toCompress)
	span := timeE(start)
	fmt.Println("Data length")
	fmt.Println(len(toCompress))
	fmt.Println("compressed Data length")
	fmt.Println(len(dst))

	start = timeS()
	src, _ := flateDecode(dst)
	span2 := timeE(start)
	fmt.Println("\nuncompressed Data length")
	fmt.Println(len(src))
	if isSame(toCompress, src) {
		fmt.Println("OK \n\n")
	}
	record := fmt.Sprintf("%d,%d,%d,%d", len(toCompress), len(dst), span, span2)
	fmt.Println(record)
	fmt.Println([]byte(record))
	writeWithOs(flateFile, []byte(record))
}




// keyToRoute returns hex bytes
// e.g {0xa1, 0xf2} -> {0xa, 0x1, 0xf, 0x2}
func keyToRoute(key []byte) []byte {
	l := len(key) * 2
	var route = make([]byte, l)
	for i, b := range key {
		route[i*2] = b/16 + '1'
		route[i*2+1] = b%16 + '1'
	}
	return route
}

// routeToKey returns native bytes
// e.g {0xa, 0x1, 0xf, 0x2} -> {0xa1, 0xf2}
func routeToKey(route []byte) []byte {
	l := len(route) / 2
	var key = make([]byte, l)
	for i := 0; i < l; i++ {
		key[i] = (route[i*2]-'1')<<4 + (route[i*2+1] - '1')
	}
	return key
}


func writeWithOs(name string, content []byte) {
	a := []byte{'\n'}
	fileObj, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()

	fileObj.Write(content)
	fileObj.Write(a)
}

var fileName = "longTx.log"
var gzipFile = "gzip.log"
var lz4File = "lz4.log"
var zstdFile = "zstd.log"
var zlibFile = "zlib.log"
var flateFile = "flate.log"

func read_file() {
	fmt.Println("==================   read   ======================")
	str, _ := ReadLine(fileName)
	for _, v := range str {
		target := routeToKey(routeToKey([]byte(v)))
		// lz4_func(target)
		// zstd_func(target)
		// zlib_func(target)
		// gzip_func(target)
		flate_func(target)
		//time.Sleep(time.Second * 1)
		// fmt.Println("----------------------------------")
		// fmt.Println()
		// fmt.Println("----------------------------------")
		// if isSame(toCompress, `routeToKey([]byte(v))`) {
		// 	fmt.Println("same")
		// } else {
		// 	fmt.Println("not same")
		// }
	}
}

func ReadLine(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result, nil
			}
			return nil, err
		}
		result = append(result, line)
	}
	return result, nil
}

func main() {
	// lz4_func()
	// zstd_func()
	// zlib_func()
	// gzip_func()
	read_file()
}

func isSame(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, v := range a {
		if v != b[idx] {
			return false
		}
	}
	return true
}

