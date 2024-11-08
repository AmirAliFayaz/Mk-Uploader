package main

import (
	"archive/zip"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	
	"github.com/schollz/progressbar/v3"
)

type UploadedFile struct {
	FileCode   string `json:"file_code"`
	FileStatus string `json:"file_status"`
}

func init() {
	flag.BoolVar(&zipFlag, "z", false, "Zip file")
	flag.StringVar(&filename, "f", "", "File name")
	flag.Parse()
	
	if !flag.Parsed() {
		flag.PrintDefaults()
	}
}

var (
	zipFlag  bool
	filename string
	rex      = regexp.MustCompile(`(?im)<textarea\s+name="download_links"\s+style="width:\s+100%">\s*(.+)\s*</textarea>`)
	client   = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         "s6.uplod.ir",
				MinVersion:         tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				},
			},
		},
	}
)

func main() {
	
	if filename == "" {
		if len(os.Args) == 1 {
			log.Fatalf("Usage: %s <filename>", path.Clean(os.Args[0]))
		}
		filename = strings.Join(os.Args[1:], " ")
	}
	
	var uploadFile = filename
	
	filename = path.Base(filename)
	
	if !strings.Contains(filename, ".") {
		log.Fatalf("Invalid filename: %s", filename)
	}
	
	// Check if zip is enabled
	if zipFlag {
		zipFileName, err := createZip(filename)
		if err != nil {
			log.Fatalf("Error creating zip: %s", err)
		}
		
		filename = filename + ".zip"
		uploadFile = zipFileName
		
		defer os.Remove(zipFileName)
	}
	
	file, err := os.Open(uploadFile)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()
	
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info: %s", err)
	}
	
	bar := progressbar.DefaultBytes(
		fileInfo.Size(),
		fmt.Sprintf("Uploading %s", filename),
	)
	
	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	
	go func() {
		defer writer.Close()
		defer multipartWriter.Close()
		
		for _, field := range map[string]string{
			"sess_id":     "",
			"utype":       "anon",
			"file_descr":  "",
			"file_public": "",
			"link_rcpt":   "",
			"link_pass":   "",
			"to_folder":   "",
			"upload":      "شروع اپلود",
			"keepalive":   "1",
		} {
			_ = multipartWriter.WriteField(field, "")
		}
		
		part, err := multipartWriter.CreateFormFile("file_0", filename)
		if err != nil {
			log.Printf("Error creating form file: %s", err)
			return
		}
		
		_, err = io.Copy(io.MultiWriter(part, bar), file)
		if err != nil {
			log.Printf("Error copying file: %s", err)
		}
	}()
	
	req, err := http.NewRequest("POST", "https://s6.uplod.ir/cgi-bin/upload.cgi?upload_type=file&utype=anon", reader)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}
	
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://uplod.ir")
	req.Header.Set("Referer", "https://uplod.ir/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.1 Safari/537.36")
	
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}
	defer resp.Body.Close()
	
	var uploadedFiles []UploadedFile
	if err := json.NewDecoder(resp.Body).Decode(&uploadedFiles); err != nil {
		log.Fatalf("Error decoding response: %s", err)
	}
	
	for _, file := range uploadedFiles {
		if file.FileStatus != "OK" {
			log.Printf("Error uploading file `%s`: %s", filename, file.FileStatus)
			continue
		}
		
		downloadLink, err := getDownloadLink(file)
		if err != nil {
			log.Printf("Error getting download link: %s", err)
			continue
		}
		
		fmt.Println("Download link:", downloadLink)
	}
}

func getDownloadLink(file UploadedFile) (string, error) {
	resp, err := client.Get(fmt.Sprintf("https://uplod.ir/upload_result?st=%s&fn=%s", file.FileStatus, file.FileCode))
	if err != nil {
		return "", err
	}
	
	defer resp.Body.Close()
	
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	if !rex.MatchString(string(all)) {
		return "", fmt.Errorf("download link not found")
	}
	
	return rex.FindStringSubmatch(string(all))[1], nil
	
}

func createZip(filename string) (string, error) {
	zipFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return "", err
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	fileToZip, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fileToZip.Close()
	
	w, err := zipWriter.Create(filename)
	if err != nil {
		return "", err
	}
	
	if _, err := io.Copy(w, fileToZip); err != nil {
		return "", err
	}
	
	return zipFile.Name(), nil
}
