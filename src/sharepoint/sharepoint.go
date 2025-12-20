package sharepoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//goland:noinspection GoSnakeCaseUsage
type (
	SHARE_POINT struct {
		mutex        sync.Mutex
		siteURL      string
		username     string
		password     string
		clientID     string
		clientSecret string
		token        string
		proxyURL     string
	}

	AUTH_RESPONSE struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	DRIVE_ITEM struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
		File struct {
			MimeType string `json:"mimeType"`
		} `json:"file"`
		Folder struct {
			ChildCount int `json:"childCount"`
		} `json:"folder"`
	}

	FILE_LIST_RESPONSE struct {
		Value []DRIVE_ITEM `json:"value"`
	}
)

//goland:noinspection GoSnakeCaseUsage
const (
	AUTHENTICATION_TOKEN_URL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	AUTHENTICATION_SCOPE     = "https://graph.microsoft.com/.default"
	DRIVE_ROOT_CONTENT_URL   = "https://graph.microsoft.com/v1.0/sites/root/drive/root:/%s:/content"
	DRIVE_ROOT_CHILDREN_URL  = "https://graph.microsoft.com/v1.0/sites/root/drive/root:/%s:/children"
	DRIVE_ROOT_ITEM_URL      = "https://graph.microsoft.com/v1.0/sites/root/drive/root:/%s"
	LIST_ITEMS_URL           = "https://graph.microsoft.com/v1.0/sites/root/lists/%s/items"
	LIST_ITEM_FIELDS_URL     = "https://graph.microsoft.com/v1.0/sites/root/lists/%s/items/%d/fields"
	LIST_ITEM_URL            = "https://graph.microsoft.com/v1.0/sites/root/lists/%s/items/%d"
	SHAREPOINT_CLIENT_ID     = "00000003-0000-0ff1-ce00-000000000000"
)

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func New(siteURL string, clientID string, clientSecret string) *SHARE_POINT {
	return &SHARE_POINT{
		siteURL:      siteURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func NewWithUser(siteURL string, username string, password string) *SHARE_POINT {
	return &SHARE_POINT{
		siteURL:  siteURL,
		username: username,
		password: password,
	}
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) Authenticate() error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var data *strings.Reader
	if s.clientID != "" && s.clientSecret != "" {
		data = strings.NewReader(
			fmt.Sprintf("client_id=%s"+
				"&client_secret=%s"+
				"&grant_type=client_credentials"+
				"&scope=%s",
				s.clientID, s.clientSecret, AUTHENTICATION_SCOPE))
	} else if s.username != "" && s.password != "" {
		data = strings.NewReader(
			fmt.Sprintf("client_id=%s"+
				"&username=%s"+
				"&password=%s"+
				"&grant_type=password"+
				"&scope=%s"+
				"&client_info=1"+
				"&response_type=token",
				SHAREPOINT_CLIENT_ID,
				s.username, s.password, AUTHENTICATION_SCOPE))
	} else {
		err = errors.New("no valid credentials provided")
	}
	if err == nil {
		request := &http.Request{}
		if request, err = http.NewRequest("POST", AUTHENTICATION_TOKEN_URL, data); err == nil {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if client := s.createHTTPClient(); client != nil {
				response := &http.Response{}
				if response, err = client.Do(request); err == nil {
					defer response.Body.Close()
					body := make([]byte, 0)
					if body, err = io.ReadAll(response.Body); err == nil {
						if response.StatusCode == http.StatusOK {
							authResp := AUTH_RESPONSE{}
							if err = json.Unmarshal(body, &authResp); err == nil {
								s.token = authResp.AccessToken
							} else {
								err = fmt.Errorf("failed to parse authentication response: %v", err)
							}
						} else {
							err = fmt.Errorf("authentication failed, status code: %d, response: %s", response.StatusCode, string(body))
						}
					} else {
						err = fmt.Errorf("failed to read authentication response: %v", err)
					}
				} else {
					err = fmt.Errorf("failed to send authentication request: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to create client: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to create authentication request: %v", err)
		}
	}
	return err
}

func (s *SHARE_POINT) GetSiteURL() (string, error) {
	err := error(nil)
	result := ""
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.siteURL != "" {
		result = s.siteURL
	} else {
		err = errors.New("site URL not set")
	}
	return result, err
}

func (s *SHARE_POINT) SetSiteURL(siteURL string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.siteURL = siteURL
	return err
}

func (s *SHARE_POINT) GetToken() (string, error) {
	err := error(nil)
	result := ""
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token != "" {
		result = s.token
	} else {
		err = errors.New("not authenticated")
	}
	return result, err
}

func (s *SHARE_POINT) SetProxy(proxyURL string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.proxyURL = proxyURL
	return err
}

func (s *SHARE_POINT) GetProxy() (string, error) {
	err := error(nil)
	result := ""
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result = s.proxyURL
	return result, err
}

func (s *SHARE_POINT) createHTTPClient() *http.Client {
	result := &http.Client{}
	if proxyURL := s.proxyURL; proxyURL != "" {
		if proxy, err := url.Parse(proxyURL); err == nil {
			result = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxy),
				},
			}
		}
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) UploadFile(localPath, remotePath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		var file *os.File
		if file, err = os.Open(localPath); err == nil {
			defer file.Close()
			var fileInfo fs.FileInfo
			if fileInfo, err = file.Stat(); err == nil {
				uploadURL := fmt.Sprintf(DRIVE_ROOT_CONTENT_URL, strings.TrimLeft(remotePath, "/"))
				request := &http.Request{}
				if request, err = http.NewRequest("PUT", uploadURL, file); err == nil {
					request.Header.Set("Authorization", "Bearer "+s.token)
					request.Header.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
					client := s.createHTTPClient()
					response := &http.Response{}
					if response, err = client.Do(request); err == nil {
						defer response.Body.Close()
						body := make([]byte, 0)
						if body, err = io.ReadAll(response.Body); err == nil {
							if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
								err = fmt.Errorf("upload failed, status code: %d, response: %s", response.StatusCode, string(body))
							}
						} else {
							err = fmt.Errorf("failed to read upload response: %v", err)
						}
					} else {
						err = fmt.Errorf("failed to send upload request: %v", err)
					}
				} else {
					err = fmt.Errorf("failed to create upload request: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to get file info: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to open local file: %v", err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) DownloadFile(remotePath, localPath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		request := &http.Request{}
		if request, err = http.NewRequest("GET", fmt.Sprintf(DRIVE_ROOT_CONTENT_URL, strings.TrimLeft(remotePath, "/")), nil); err == nil {
			request.Header.Set("Authorization", "Bearer "+s.token)
			client := s.createHTTPClient()
			response := &http.Response{}
			if response, err = client.Do(request); err == nil {
				defer response.Body.Close()
				if response.StatusCode == http.StatusOK {
					localDirectory := filepath.Dir(localPath)
					if err = os.MkdirAll(localDirectory, 0755); err == nil {
						var outputFile *os.File
						if outputFile, err = os.Create(localPath); err == nil {
							defer outputFile.Close()
							if _, err = io.Copy(outputFile, response.Body); err != nil {
								err = fmt.Errorf("failed to write file content: %v", err)
							}
						} else {
							err = fmt.Errorf("failed to create local file: %v", err)
						}
					} else {
						err = fmt.Errorf("failed to create local directory: %v", err)
					}
				} else {
					body, _ := io.ReadAll(response.Body)
					err = fmt.Errorf("download failed, status code: %d, response: %s", response.StatusCode, string(body))
				}
			} else {
				err = fmt.Errorf("failed to send download request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to create download request: %v", err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) ListFiles(folderPath string) ([]string, error) {
	result := make([]string, 0)
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		listURL := fmt.Sprintf(DRIVE_ROOT_CHILDREN_URL, strings.TrimLeft(folderPath, "/"))
		request := &http.Request{}
		if request, err = http.NewRequest("GET", listURL, nil); err == nil {
			request.Header.Set("Authorization", "Bearer "+s.token)
			request.Header.Set("Content-Type", "application/json")
			client := s.createHTTPClient()
			response := &http.Response{}
			if response, err = client.Do(request); err == nil {
				defer response.Body.Close()
				body := make([]byte, 0)
				if body, err = io.ReadAll(response.Body); err == nil {
					if response.StatusCode == http.StatusOK {
						listResp := FILE_LIST_RESPONSE{}
						if err = json.Unmarshal(body, &listResp); err == nil {
							for _, item := range listResp.Value {
								if item.File.MimeType != "" {
									result = append(result, item.Name)
								}
							}
						} else {
							err = fmt.Errorf("failed to parse file list response: %v", err)
						}
					} else {
						err = fmt.Errorf("failed to get file list, status code: %d, response: %s", response.StatusCode, string(body))
					}
				} else {
					err = fmt.Errorf("failed to read file list response: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to send file list request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to create file list request: %v", err)
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) CreateFolder(folderPath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		parentPath := filepath.Dir(folderPath)
		folderName := filepath.Base(folderPath)
		j := make([]byte, 0)
		if j, err = json.Marshal(map[string]interface{}{
			"name":                              folderName,
			"folder":                            map[string]interface{}{},
			"@microsoft.graph.conflictBehavior": "fail",
		}); err == nil {
			request := &http.Request{}
			if request, err = http.NewRequest("POST", fmt.Sprintf(DRIVE_ROOT_CHILDREN_URL, strings.TrimLeft(parentPath, "/")), bytes.NewBuffer(j)); err == nil {
				request.Header.Set("Authorization", "Bearer "+s.token)
				request.Header.Set("Content-Type", "application/json")
				client := s.createHTTPClient()
				response := &http.Response{}
				if response, err = client.Do(request); err == nil {
					defer response.Body.Close()
					body := make([]byte, 0)
					if body, err = io.ReadAll(response.Body); err == nil {
						if response.StatusCode != http.StatusCreated {
							err = fmt.Errorf("failed to create folder, status code: %d, response: %s", response.StatusCode, string(body))
						}
					} else {
						err = fmt.Errorf("failed to read folder creation response: %v", err)
					}
				} else {
					err = fmt.Errorf("failed to send folder creation request: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to create folder request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to marshal request body: %v", err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) DeleteFile(filePath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		request := &http.Request{}
		if request, err = http.NewRequest("DELETE", fmt.Sprintf(DRIVE_ROOT_ITEM_URL, strings.TrimLeft(filePath, "/")), nil); err == nil {
			request.Header.Set("Authorization", "Bearer "+s.token)
			client := s.createHTTPClient()
			response := &http.Response{}
			if response, err = client.Do(request); err == nil {
				defer response.Body.Close()
				body := make([]byte, 0)
				if body, err = io.ReadAll(response.Body); err == nil {
					if response.StatusCode != http.StatusNoContent {
						err = fmt.Errorf("failed to delete file, status code: %d, response: %s", response.StatusCode, string(body))
					}
				} else {
					err = fmt.Errorf("failed to read delete file response: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to send delete file request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to create delete file request: %v", err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) GetListItems(listName string) ([]map[string]interface{}, error) {
	err := error(nil)
	result := make([]map[string]interface{}, 0)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		request := &http.Request{}
		if request, err = http.NewRequest("GET", fmt.Sprintf(LIST_ITEMS_URL+"?expand=fields", listName), nil); err == nil {
			request.Header.Set("Authorization", "Bearer "+s.token)
			request.Header.Set("Content-Type", "application/json")
			client := s.createHTTPClient()
			response := &http.Response{}
			if response, err = client.Do(request); err == nil {
				defer response.Body.Close()
				body := make([]byte, 0)
				if body, err = io.ReadAll(response.Body); err == nil {
					if response.StatusCode == http.StatusOK {
						listResp := struct {
							Value []struct {
								Fields map[string]interface{} `json:"fields"`
							} `json:"value"`
						}{}
						if err := json.Unmarshal(body, &listResp); err == nil {
							for _, item := range listResp.Value {
								result = append(result, item.Fields)
							}
						} else {
							err = fmt.Errorf("failed to parse list items response: %v", err)
						}
					} else {
						err = fmt.Errorf("failed to get list items, status code: %d, response: %s", response.StatusCode, string(body))
					}
				} else {
					err = fmt.Errorf("failed to read get list items response: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to send get list items request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to create get list items request: %v", err)
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) AddListItem(listName string, item map[string]interface{}) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		addItemURL := fmt.Sprintf(LIST_ITEMS_URL, listName)
		j := make([]byte, 0)
		if j, err = json.Marshal(map[string]interface{}{
			"fields": item,
		}); err == nil {
			request := &http.Request{}
			if request, err = http.NewRequest("POST", addItemURL, bytes.NewBuffer(j)); err == nil {
				request.Header.Set("Authorization", "Bearer "+s.token)
				request.Header.Set("Content-Type", "application/json")
				client := s.createHTTPClient()
				response := &http.Response{}
				if response, err = client.Do(request); err == nil {
					defer response.Body.Close()
					body := make([]byte, 0)
					if body, err = io.ReadAll(response.Body); err == nil {
						if response.StatusCode != http.StatusCreated {
							err = fmt.Errorf("failed to add list item, status code: %d, response: %s", response.StatusCode, string(body))
						}
					} else {
						err = fmt.Errorf("failed to read add list item response: %v", err)
					}
				} else {
					err = fmt.Errorf("failed to send add list item request: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to create add list item request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to marshal request body: %v", err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) UpdateListItem(listName string, itemID int, item map[string]interface{}) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		updateItemURL := fmt.Sprintf(LIST_ITEM_FIELDS_URL, listName, itemID)
		j := make([]byte, 0)
		if j, err = json.Marshal(item); err == nil {
			request := &http.Request{}
			if request, err = http.NewRequest("PATCH", updateItemURL, bytes.NewBuffer(j)); err == nil {
				request.Header.Set("Authorization", "Bearer "+s.token)
				request.Header.Set("Content-Type", "application/json")
				client := s.createHTTPClient()
				response := &http.Response{}
				if response, err = client.Do(request); err == nil {
					defer response.Body.Close()
					body := make([]byte, 0)
					if body, err = io.ReadAll(response.Body); err == nil {
						if response.StatusCode != http.StatusOK {
							err = fmt.Errorf("failed to update list item, status code: %d, response: %s", response.StatusCode, string(body))
						}
					} else {
						err = fmt.Errorf("failed to read update list item response: %v", err)
					}
				} else {
					err = fmt.Errorf("failed to send update list item request: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to create update list item request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to marshal request body: %v", err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) DeleteListItem(listName string, itemID int) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token == "" {
		err = errors.New("not authenticated")
	} else {
		deleteItemURL := fmt.Sprintf(LIST_ITEM_URL, listName, itemID)
		request := &http.Request{}
		if request, err = http.NewRequest("DELETE", deleteItemURL, nil); err == nil {
			request.Header.Set("Authorization", "Bearer "+s.token)
			client := s.createHTTPClient()
			response := &http.Response{}
			if response, err = client.Do(request); err == nil {
				defer response.Body.Close()
				body := make([]byte, 0)
				if body, err = io.ReadAll(response.Body); err == nil {
					if response.StatusCode != http.StatusNoContent {
						err = fmt.Errorf("failed to delete list item, status code: %d, response: %s", response.StatusCode, string(body))
					}
				} else {
					err = fmt.Errorf("failed to read delete list item response: %v", err)
				}
			} else {
				err = fmt.Errorf("failed to send delete list item request: %v", err)
			}
		} else {
			err = fmt.Errorf("failed to create delete list item request: %v", err)
		}
	}
	return err
}
