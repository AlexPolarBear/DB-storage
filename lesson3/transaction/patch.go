package transaction

import (
	"encoding/json"
	"log"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func transaction() {
	sourceData := []byte(`{"name": "Alex", "age": 23, "city": "Saint-Peterburg"}`)
	targetData := make(map[string]interface{})

	patch, err := jsonpatch.CreatePatch(sourceData, targetData)
	if err != nil {
		log.Fatal("Error creating JSON Patch:", err)
	}

	patchedData, err := patch.Apply(targetData)
	if err != nil {
		log.Fatal("Error using JSON Patch:", err)
	}

	patchedJSON, err := json.Marshal(patchedData)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}

	log.Println(string(patchedJSON))
}
