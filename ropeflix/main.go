package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
)

func main() {

	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.
	file, err := os.Open("downloads.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// The bufio.NewScanner() function is called in which the
	// object os.File passed as its parameter and this returns a
	// object bufio.Scanner which is further used on the
	// bufio.Scanner.Split() method.
	scanner := bufio.NewScanner(file)

	// The bufio.ScanLines is used as an
	// input to the method bufio.Scanner.Split()
	// and then the scanning forwards to each
	// new line using the bufio.Scanner.Scan()
	// method.
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	// and then a loop iterates through
	// and prints each of the slice values.
	for _, eachLine := range text {
		keys := strings.Split(eachLine, ">")
		if len(keys) == 2 {
			fileName := keys[0]
			url := keys[1]
			fmt.Printf("Downloading FileUrl: %s, FileName: %s.mp4\n", url, fileName)

			referer := fmt.Sprintf("https://www.ropeflix.com/movie/%s/", fileName)
			filePath := fmt.Sprintf("/Users/e/Desktop/Ropeflix/%s.mp4", fileName)

			err := DownloadFile(filePath, referer, url)
			if err != nil {
				panic(err)
			}
			fmt.Println("Downloaded: " + url)
		}
	}
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, referer string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	// ...
	req.Header.Add("accept", `*/*`)
	req.Header.Add("accept-encoding", `identity;q=1, *;q=0`)
	req.Header.Add("accept-language", `en-US,en;q=0.9`)
	req.Header.Add("range", `bytes=0-`)
	req.Header.Add("referer", referer)
	req.Header.Add("sec-fetch-dest", `video`)
	req.Header.Add("sec-fetch-mode", `no-cors`)
	req.Header.Add("sec-fetch-site", `cross-site`)
	req.Header.Add("user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36`)

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}
