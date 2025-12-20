// Package ftpserver
// File:        ftpserver.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/ftp/server/ftpserver.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: FTP server functionality for Go applications
// --------------------------------------------------------------------------------
package ftpserver

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"time"

	ftpserverlib "github.com/fclairamb/ftpserverlib"
	"github.com/spf13/afero"
	"github.com/xiang-tai-duo/go-bootstrap/logger"
)

//goland:noinspection GoSnakeCaseUsage
const (
	BANNER           = "Welcome to Go-Boost FTP Server"
	CLIENT_CONNECTED = "Client connected"
)

//goland:noinspection SpellCheckingInspection
var (
	ConnectionTimeout   = 30
	Debug               = false
	IdleTimeout         = 300
	ListenAddress       = ":21"
	Directory           = "./ftp_files"
	PasvHost            = ""
	PasvPortMin         = 30000
	PasvPortMax         = 31000
	PasvPortMaxAttempts = 100
)

//goland:noinspection GoNameStartsWithPackageName,SpellCheckingInspection,GoSnakeCaseUsage
type (
	FTP_SERVER      struct{}
	PASV_PORT_RANGE struct {
		minPort     int
		maxPort     int
		current     int
		maxAttempts int
	}
)

func New() *FTP_SERVER {
	return &FTP_SERVER{}
}

//goland:noinspection GoBoolExpressions,GoUnusedParameter
func (f *FTP_SERVER) AuthUser(ftpClient ftpserverlib.ClientContext, user, pass string) (ftpserverlib.ClientDriver, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("AuthUser: user=%s", user))
	}
	return f, nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) ClientConnected(ftpClient ftpserverlib.ClientContext) (string, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("ClientConnected: from %s (ID: %d)", ftpClient.RemoteAddr(), ftpClient.ID()))
	}
	return CLIENT_CONNECTED, nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) ClientDisconnected(ftpClient ftpserverlib.ClientContext) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("ClientDisconnected: ID: %d", ftpClient.ID()))
	}
}

//goland:noinspection GoUnusedExportedFunction
func (f *FTP_SERVER) ListenAsync() {
	go func() {
		if err := os.MkdirAll(Directory, 0755); err == nil {
			logger.Logger.Debug(fmt.Sprintf("Starting FTP server on %s", ListenAddress))
			logger.Logger.Debug("Anonymous login enabled")
			logger.Logger.Debug(fmt.Sprintf("File storage directory: %s", Directory))
			err = ftpserverlib.NewFtpServer(&FTP_SERVER{}).ListenAndServe()
		} else {
			err = fmt.Errorf("failed to create FTP directory: %v", err)
		}
	}()
}

func (f *FTP_SERVER) GetSettings() (*ftpserverlib.Settings, error) {
	return f.loadConfig()
}

func (f *FTP_SERVER) GetTLSConfig() (*tls.Config, error) {
	return nil, fmt.Errorf("TLS is not configured on this server")
}

//goland:noinspection GoBoolExpressions,GoUnusedParameter
func (f *FTP_SERVER) WelcomeUser(ftpClient ftpserverlib.ClientContext) (string, error) {
	if Debug {
		logger.Logger.Debug("WelcomeUser called")
	}
	return BANNER, nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) UserLeft(ftpClient ftpserverlib.ClientContext) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("UserLeft: ID: %d", ftpClient.ID()))
	}
}

func (f *FTP_SERVER) loadConfig() (*ftpserverlib.Settings, error) {
	return &ftpserverlib.Settings{
		ListenAddr:               ListenAddress,
		PublicHost:               PasvHost,
		Banner:                   BANNER,
		IdleTimeout:              IdleTimeout,
		ConnectionTimeout:        ConnectionTimeout,
		ActiveTransferPortNon20:  false,
		DisableMLSD:              false,
		DisableMLST:              false,
		DisableMFMT:              false,
		TLSRequired:              ftpserverlib.ClearOrEncrypted,
		DisableLISTArgs:          false,
		DisableSite:              false,
		DisableActiveMode:        false,
		EnableHASH:               false,
		DisableSTAT:              false,
		DisableSYST:              false,
		EnableCOMB:               false,
		DeflateCompressionLevel:  0,
		DefaultTransferType:      ftpserverlib.TransferTypeBinary,
		ActiveConnectionsCheck:   ftpserverlib.IPMatchRequired,
		PasvConnectionsCheck:     ftpserverlib.IPMatchRequired,
		PassiveTransferPortRange: newPasvPortRange(PasvPortMin, PasvPortMax, PasvPortMaxAttempts),
	}, nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) GetFile(path string) ([]byte, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("GetFile: %s", path))
	}
	return os.ReadFile(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) PutFile(path string, data []byte) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("PutFile: %s, size=%d", path, len(data)))
	}
	return os.WriteFile(filepath.Join(Directory, path), data, 0644)
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) DeleteFile(path string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("DeleteFile: %s", path))
	}
	return os.Remove(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) RenameFile(from, to string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("RenameFile: %s -> %s", from, to))
	}
	return os.Rename(filepath.Join(Directory, from), filepath.Join(Directory, to))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) MakeDir(path string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("MakeDir: %s", path))
	}
	return os.MkdirAll(filepath.Join(Directory, path), 0755)
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) DeleteDir(path string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("DeleteDir: %s", path))
	}
	return os.RemoveAll(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) ListDir(path string) ([]os.FileInfo, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("ListDir: %s", path))
	}
	result := make([]os.FileInfo, 0)
	if entries, err := os.ReadDir(filepath.Join(Directory, path)); err == nil {
		result := make([]os.FileInfo, 0, len(entries))
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			result = append(result, info)
		}
	} else {
		result = nil
	}
	return result, nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Stat(path string) (os.FileInfo, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Stat: %s", path))
	}
	return os.Stat(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) CanAllocate(size int) (bool, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("CanAllocate: size=%d", size))
	}
	return true, nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Chmod(path string, mode os.FileMode) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Chmod: %s, mode=%o", path, mode))
	}
	return os.Chmod(filepath.Join(Directory, path), mode)
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) SetTime(path string, mtime, atime int64) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("SetTime: %s, mtime=%d, atime=%d", path, mtime, atime))
	}
	return nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Chown(path string, uid, gid int) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Chown: %s, uid=%d, gid=%d", path, uid, gid))
	}
	return nil
}

//goland:noinspection GoBoolExpressions,SpellCheckingInspection
func (f *FTP_SERVER) Chtimes(path string, atime time.Time, mtime time.Time) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Chtimes: %s, atime=%v, mtime=%v", path, atime, mtime))
	}
	return nil
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Create(path string) (afero.File, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Create: %s", path))
	}
	return os.Create(filepath.Join(Directory, path))
}

func (p *PASV_PORT_RANGE) FetchNext() (int, int, bool) {
	if p.current < p.minPort || p.current >= p.maxPort {
		p.current = p.minPort
	} else {
		p.current++
	}
	return p.current, p.current, true
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Mkdir(path string, perm os.FileMode) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Mkdir: %s, perm=%o", path, perm))
	}
	return os.MkdirAll(filepath.Join(Directory, path), perm)
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) MkdirAll(path string, perm os.FileMode) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("MkdirAll: %s, perm=%o", path, perm))
	}
	return os.MkdirAll(filepath.Join(Directory, path), perm)
}

func (f *FTP_SERVER) Name() string {
	return "go-boost File Server"
}

func (p *PASV_PORT_RANGE) NumberAttempts() int {
	return p.maxAttempts
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Open(path string) (afero.File, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Open: %s", path))
	}
	return os.Open(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("OpenFile: %s, flag=%d, perm=%o", path, flag, perm))
	}
	return os.OpenFile(filepath.Join(Directory, path), flag, perm)
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Remove(path string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Remove: %s", path))
	}
	return os.Remove(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) RemoveAll(path string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("RemoveAll: %s", path))
	}
	return os.RemoveAll(filepath.Join(Directory, path))
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Rename(oldPath string, newPath string) error {
	if Debug {
		logger.Logger.Debug(fmt.Sprintf("Rename: %s -> %s", oldPath, newPath))
	}
	return os.Rename(filepath.Join(Directory, oldPath), filepath.Join(Directory, newPath))
}

//goland:noinspection SpellCheckingInspection
func newPasvPortRange(minPort, maxPort, maxAttempts int) *PASV_PORT_RANGE {
	return &PASV_PORT_RANGE{
		minPort:     minPort,
		maxPort:     maxPort,
		current:     minPort - 1,
		maxAttempts: maxAttempts,
	}
}
