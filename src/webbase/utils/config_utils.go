package utils

import (
	"bytes"
	"encoding/json"

	"log"
)

func PrintConfig(cfg interface{}) {
	b, _ := json.Marshal(&cfg)
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	log.Printf("\n %s \n", string(out.Bytes()))
}
