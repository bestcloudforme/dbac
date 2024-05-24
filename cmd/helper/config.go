package helper

import (
	"encoding/json"
	"log"
	"os"
)

// func LoadConfig() string {
// 	data, err := ioutil.ReadFile("app.config")
// 	if err != nil {
// 		log.Panicf("failed reading data from file: %s", err)
// 	}
// 	data_str := string(data)
// 	res := strings.Split(data_str, "=")
// 	//fmt.Println(res)
// 	//fmt.Println(data_str)
// 	return res[1]
// }

func LoadProfile() string {
	profilePath := os.Getenv("HOME") + "/.dbac-profiles.json"
	data, err := os.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("Failed to read profile file: %v", err)
	}

	var curProfile Profiles
	if err := json.Unmarshal(data, &curProfile); err != nil {
		log.Fatalf("Failed to unmarshal profiles: %v", err)
	}

	return curProfile.Current
}
