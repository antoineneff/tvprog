package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const XMLZipURL = "https://xmltvfr.fr/xmltv/xmltv_tnt.zip"

func FetchXML() (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(XMLZipURL)
	if err != nil {
		return "", fmt.Errorf("failed to download ZIP file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read ZIP response: %w", err)
	}

	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return "", fmt.Errorf("failed to create ZIP reader: %w", err)
	}

	for _, file := range reader.File {
		if file.Name == "xmltv_tnt.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open XML file in ZIP: %w", err)
			}
			defer rc.Close()

			xmlData, err := io.ReadAll(rc)
			if err != nil {
				return "", fmt.Errorf("failed to read XML from ZIP: %w", err)
			}

			return string(xmlData), nil
		}
	}

	return "", fmt.Errorf("xmltv_tnt.xml not found in ZIP archive")
}

func ParseXML(xmlData string) (*TV, error) {
	var tv TV
	err := xml.Unmarshal([]byte(xmlData), &tv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}
	return &tv, nil
}
