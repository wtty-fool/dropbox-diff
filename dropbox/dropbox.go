package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// https://dropbox.github.io/dropbox-api-v2-explorer/#files_list_folder
	OptionRecursive        = droboxOption{"recursive", true}
	OptionIncludeMediaInfo = droboxOption{"include_media_info", true}
)

const (
	listDirURL = "https://api.dropboxapi.com/2/files/list_folder"
)

type response struct {
	Entries []DropboxFile `json:"entries"`
	Cursor  string        `json:"cursor"`
	HasMore bool          `json:"has_more"`
}

type DropboxFile struct {
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	PathLower   string `json:"path_lower"`
	PathDisplay string `json:"path_display"`
}

func ListDropboxDir(path string, token string, dropboxOptions ...droboxOption) ([]DropboxFile, error) {
	rawBody := map[string]interface{}{
		"path": path,
	}
	for _, option := range dropboxOptions {
		rawBody[option.Key] = option.Value
	}

	body, err := json.Marshal(rawBody)
	if err != nil {
		return []DropboxFile{}, err
	}

	req, err := http.NewRequest(http.MethodPost, listDirURL, bytes.NewBuffer(body))
	if err != nil {
		return []DropboxFile{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []DropboxFile{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []DropboxFile{}, fmt.Errorf("Request failed: [%d] %s", resp.StatusCode, resp.Status)
	}

	var data response
	{
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []DropboxFile{}, err
		}
		err = json.Unmarshal(respBytes, &data)
		if err != nil {
			return []DropboxFile{}, err
		}
	}

	return data.Entries, nil
}

type droboxOption struct {
	Key   string
	Value interface{}
}
