package sharepoint

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

//goland:noinspection GoSnakeCaseUsage
const (
	MICROSOFT_GRAPH_DEFAULT_SCOPE                        = "https://graph.microsoft.com/.default"
	SHAREPOINT_ONLINE_CLIENT_ID                          = "00000003-0000-0ff1-ce00-000000000000"
	ERROR_NOT_AUTHENTICATED                              = "not authenticated"
	ERROR_SITE_URL_NOT_SET                               = "site URL not set"
	ERROR_USERNAME_PASSWORD_AUTHENTICATION_NOT_SUPPORTED = "username/password authentication not supported in current implementation"
	ERROR_NO_VALID_CREDENTIALS_PROVIDED                  = "no valid credentials provided"
	ERROR_GRAPH_SERVICE_CLIENT_NOT_INITIALIZED           = "graph service client not initialized"
	ERROR_UPLOAD_FILE_NOT_IMPLEMENTED                    = "UploadFile method not fully implemented with official SDK"
	ERROR_DOWNLOAD_FILE_NOT_IMPLEMENTED                  = "DownloadFile method not fully implemented with official SDK"
	ERROR_LIST_FILES_NOT_IMPLEMENTED                     = "ListFiles method not fully implemented with official SDK"
	ERROR_CREATE_FOLDER_NOT_IMPLEMENTED                  = "CreateFolder method not fully implemented with official SDK"
	ERROR_DELETE_FILE_NOT_IMPLEMENTED                    = "DeleteFile method not fully implemented with official SDK"
	ERROR_GET_LIST_ITEMS_NOT_IMPLEMENTED                 = "GetListItems method not fully implemented with official SDK"
	ERROR_ADD_LIST_ITEM_NOT_IMPLEMENTED                  = "AddListItem method not fully implemented with official SDK"
	ERROR_UPDATE_LIST_ITEM_NOT_IMPLEMENTED               = "UpdateListItem method not fully implemented with official SDK"
	ERROR_DELETE_LIST_ITEM_NOT_IMPLEMENTED               = "DeleteListItem method not fully implemented with official SDK"
	HTTP_STATUS_OK                                       = 200
	HTTP_STATUS_CREATED                                  = 201
	HTTP_STATUS_NO_CONTENT                               = 204
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
		proxyURL     string
		client       *msgraphsdk.GraphServiceClient
		ctx          context.Context
	}
)

// Default authentication scopes
var DEFAULT_SCOPES = []string{
	MICROSOFT_GRAPH_DEFAULT_SCOPE,
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func New(siteURL string, clientID string, clientSecret string) *SHARE_POINT {
	sp := &SHARE_POINT{
		siteURL:      siteURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		ctx:          context.Background(),
	}
	sp.Authenticate()
	return sp
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func NewWithUser(siteURL string, username string, password string) *SHARE_POINT {
	sp := &SHARE_POINT{
		siteURL:  siteURL,
		username: username,
		password: password,
		ctx:      context.Background(),
	}
	sp.Authenticate()
	return sp
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) Authenticate() error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.clientID != "" && s.clientSecret != "" {
		options := &azidentity.ClientSecretCredentialOptions{}
		credential, err := azidentity.NewClientSecretCredential(
			"common",
			s.clientID,
			s.clientSecret,
			options,
		)
		if err == nil {
			client, clientErr := msgraphsdk.NewGraphServiceClientWithCredentials(credential, DEFAULT_SCOPES)
			if clientErr != nil {
				err = fmt.Errorf("failed to create graph service client: %v", clientErr)
			} else {
				s.client = client
			}
		}
	} else if s.username != "" && s.password != "" {
		err = errors.New(ERROR_USERNAME_PASSWORD_AUTHENTICATION_NOT_SUPPORTED)
	} else {
		err = errors.New(ERROR_NO_VALID_CREDENTIALS_PROVIDED)
	}

	return err
}

func (s *SHARE_POINT) GetSiteURL() (string, error) {
	result := ""
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.siteURL != "" {
		result = s.siteURL
	} else {
		err = errors.New(ERROR_SITE_URL_NOT_SET)
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

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) UploadFile(localPath, remotePath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		var file *os.File
		if file, err = os.Open(localPath); err == nil {
			defer file.Close()
			if _, err = file.Stat(); err == nil {
				if s.client == nil {
					err = errors.New(ERROR_GRAPH_SERVICE_CLIENT_NOT_INITIALIZED)
				} else {
					err = errors.New(ERROR_UPLOAD_FILE_NOT_IMPLEMENTED)
				}
			}
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) DownloadFile(remotePath, localPath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		localDirectory := filepath.Dir(localPath)
		if err = os.MkdirAll(localDirectory, 0755); err == nil {
			err = errors.New(ERROR_DOWNLOAD_FILE_NOT_IMPLEMENTED)
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
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_LIST_FILES_NOT_IMPLEMENTED)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) CreateFolder(folderPath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_CREATE_FOLDER_NOT_IMPLEMENTED)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) DeleteFile(filePath string) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_DELETE_FILE_NOT_IMPLEMENTED)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) GetListItems(listName string) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_GET_LIST_ITEMS_NOT_IMPLEMENTED)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) AddListItem(listName string, item map[string]interface{}) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_ADD_LIST_ITEM_NOT_IMPLEMENTED)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) UpdateListItem(listName string, itemID int, item map[string]interface{}) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_UPDATE_LIST_ITEM_NOT_IMPLEMENTED)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (s *SHARE_POINT) DeleteListItem(listName string, itemID int) error {
	err := error(nil)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.client == nil {
		err = errors.New(ERROR_NOT_AUTHENTICATED)
	} else {
		err = errors.New(ERROR_DELETE_LIST_ITEM_NOT_IMPLEMENTED)
	}
	return err
}
