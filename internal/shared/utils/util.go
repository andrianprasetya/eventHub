package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Slugify(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}

func GenerateDomainName(s string) string {
	dateTimePrefix := time.Now().Format("20060102150405")
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	randomNumber := r.Intn(900) + 100

	return fmt.Sprintf("%s-%s%d", Slugify(s), dateTimePrefix, randomNumber)
}

func ToJSONString(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error marshalling feature:", err)
		return "{}"
	}
	return string(bytes)
}

func ToStringJSON(v string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(v), &data)
	return data, err
}
