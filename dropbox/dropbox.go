package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	client = &http.Client{}
	// https://dropbox.github.io/dropbox-api-v2-explorer/#files_list_folder
	OptionRecursive        = droboxOption{"recursive", true}
	OptionIncludeMediaInfo = droboxOption{"include_media_info", true}
)

const (
	listDirURL      = "https://api.dropboxapi.com/2/files/list_folder"
	listContinueURL = "https://api.dropboxapi.com/2/files/list_folder/continue"
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

func makeRequest(method string, url string, headers map[string]string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Request failed: [%d] %s", resp.StatusCode, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func ListDropboxDir(path string, token string, dropboxOptions ...droboxOption) ([]DropboxFile, error) {
	commonHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	var fileList response
	{
		rawBody := map[string]interface{}{
			"path": path,
		}
		for _, option := range dropboxOptions {
			rawBody[option.Key] = option.Value
		}

		respBytes, err := makeRequest(http.MethodPost, listDirURL, commonHeaders, rawBody)
		err = json.Unmarshal(respBytes, &fileList)
		if err != nil {
			return []DropboxFile{}, err
		}
	}

	output := []DropboxFile{}
	for _, entry := range fileList.Entries {
		output = append(output, entry)
	}

	hasMore := fileList.HasMore
	for hasMore {
		rawBody := map[string]string{
			"cursor": fileList.Cursor,
		}
		respBytes, err := makeRequest(http.MethodPost, listContinueURL, commonHeaders, rawBody)
		if err != nil {
			return []DropboxFile{}, err
		}
		err = json.Unmarshal(respBytes, &fileList)
		if err != nil {
			return []DropboxFile{}, err
		}
		for _, entry := range fileList.Entries {
			output = append(output, entry)
		}
		hasMore = fileList.HasMore
	}

	return output, nil
}

type droboxOption struct {
	Key   string
	Value interface{}
}
