package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/nrednav/cuid2"
	"golang.org/x/net/html/charset"
)

func main() {
	url := "http://servicos.cptec.inpe.br/XML/listaCidades"
	response, requestErr := http.Get(url)

	if requestErr != nil {
		log.Fatal("get error", requestErr)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal("response status error ->", response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")

	charsetReader, err := charset.NewReader(response.Body, contentType)

	if err != nil {
		log.Fatal("could not read body", err)
	}

	data, readErr := io.ReadAll(charsetReader)

	if readErr != nil {
		log.Fatalf("Failed to determine encoding: %v", err)
	}

	data = skipXMLDeclaration(data)

	var CityResponse CityResponse

	unmarshalError := xml.Unmarshal(data, &CityResponse)

	if err != nil {
		log.Fatal("could not unmarshal xml", unmarshalError)
	}

	output := "INSERT INTO city (id, external_id, state_code, name) VALUES "

	valuesToInsert := [][]string{}

	for _, city := range CityResponse.Cities {
		valuesToInsert = append(valuesToInsert, []string{
			formatValue(cuid2.Generate()),
			formatValue(city.Id),
			formatValue(city.StateCode),
			formatValue(city.Name),
		})
	}

	for _, values := range valuesToInsert {
		output += fmt.Sprintf("(%s),", strings.Join(values, ", "))
	}

	fmt.Println(output[:len(output)-1])
}

func formatValue(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

type CityResponse struct {
	XMLName xml.Name `xml:"cidades"`
	Cities  []City   `xml:"cidade"`
}

type City struct {
	Id        string `xml:"id"`
	Name      string `xml:"nome"`
	StateCode string `xml:"uf"`
}

func skipXMLDeclaration(xmlData []byte) []byte {
	if bytes.HasPrefix(xmlData, []byte("<?xml")) {
		closeBracketIndex := bytes.Index(xmlData, []byte(">"))
		xmlData = xmlData[closeBracketIndex+1:]
	}
	return xmlData
}

type XMLData struct {
	Data     []byte
	Encoding string
}
