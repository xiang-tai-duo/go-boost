// Package ftpserver
// File:        ftpserver.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ftp/server/ftpserver.go
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
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
type (
	FTP_SERVER      struct{}
	PASV_PORT_RANGE struct {
		current     int
		maxAttempts int
		maxPort     int
		minPort     int
	}
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection,GoUnusedConst
const (
	BANNER                        = "Welcome to Go-Boost FTP Server"
	CLIENT_CONNECTED              = "Client connected"
	DEFAULT_CONNECTION_TIMEOUT    = 30
	DEFAULT_DEBUG                 = false
	DEFAULT_DIRECTORY             = "./ftp_files"
	DEFAULT_IDLE_TIMEOUT          = 300
	DEFAULT_LISTEN_ADDRESS        = ":21"
	DEFAULT_PASV_HOST             = ""
	DEFAULT_PASV_PORT_MAX         = 31000
	DEFAULT_PASV_PORT_MAX_ATTEMPS = 100
	DEFAULT_PASV_PORT_MIN         = 30000
	DIR_MODE                      = 0755
	FILE_MODE                     = 0644
	MODULE_NAME_FTPSERVER         = "ftpserver"
	SERVER_NAME                   = "go-boost File Server"
	TLS_NOT_CONFIGURED            = "TLS is not configured on this server"
)

//goland:noinspection SpellCheckingInspection
var (
	ConnectionTimeout   = DEFAULT_CONNECTION_TIMEOUT
	Debug               = DEFAULT_DEBUG
	Directory           = DEFAULT_DIRECTORY
	IdleTimeout         = DEFAULT_IDLE_TIMEOUT
	ListenAddress       = DEFAULT_LISTEN_ADDRESS
	PasvHost            = DEFAULT_PASV_HOST
	PasvPortMax         = DEFAULT_PASV_PORT_MAX
	PasvPortMaxAttempts = DEFAULT_PASV_PORT_MAX_ATTEMPS
	PasvPortMin         = DEFAULT_PASV_PORT_MIN
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_FTPSERVER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_FTPSERVER, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_FTPSERVER, logger.SKIP_STACK_FRAMES_BASE)
}

func New() *FTP_SERVER {
	result := &FTP_SERVER{}
	return result
}

//goland:noinspection GoBoolExpressions,GoUnusedParameter
func (f *FTP_SERVER) AuthUser(ftpClient ftpserverlib.ClientContext, user, pass string) (ftpserverlib.ClientDriver, error) {
	result := ftpserverlib.ClientDriver(f)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("AuthUser: user=%s", user))
	}
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) CanAllocate(size int) (bool, error) {
	result := true
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("CanAllocate: size=%d", size))
	}
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Chmod(path string, mode os.FileMode) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Chmod: %s, mode=%o", path, mode))
	}
	err = os.Chmod(filepath.Join(Directory, path), mode)
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Chown(path string, uid, gid int) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Chown: %s, uid=%d, gid=%d", path, uid, gid))
	}
	return err
}

//goland:noinspection GoBoolExpressions,SpellCheckingInspection
func (f *FTP_SERVER) Chtimes(path string, atime time.Time, mtime time.Time) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Chtimes: %s, atime=%v, mtime=%v", path, atime, mtime))
	}
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) ClientConnected(ftpClient ftpserverlib.ClientContext) (string, error) {
	result := CLIENT_CONNECTED
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("ClientConnected: from %s (ID: %d)", ftpClient.RemoteAddr(), ftpClient.ID()))
	}
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) ClientDisconnected(ftpClient ftpserverlib.ClientContext) {
	if Debug {
		__debug(fmt.Sprintf("ClientDisconnected: ID: %d", ftpClient.ID()))
	}
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Create(path string) (afero.File, error) {
	result := afero.File(nil)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Create: %s", path))
	}
	result, err = os.Create(filepath.Join(Directory, path))
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) DeleteDir(path string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("DeleteDir: %s", path))
	}
	err = os.RemoveAll(filepath.Join(Directory, path))
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) DeleteFile(path string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("DeleteFile: %s", path))
	}
	err = os.Remove(filepath.Join(Directory, path))
	return err
}

//goland:noinspection SpellCheckingInspection
func (p *PASV_PORT_RANGE) FetchNext() (int, int, bool) {
	result := 0
	if p.current < p.minPort || p.current >= p.maxPort {
		p.current = p.minPort
	} else {
		p.current++
	}
	result = p.current
	return result, result, true
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) GetFile(path string) ([]byte, error) {
	result := make([]byte, 0)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("GetFile: %s", path))
	}
	result, err = os.ReadFile(filepath.Join(Directory, path))
	return result, err
}

func (f *FTP_SERVER) GetSettings() (*ftpserverlib.Settings, error) {
	result, err := f.loadConfig()
	return result, err
}

func (f *FTP_SERVER) GetTLSConfig() (*tls.Config, error) {
	result := (*tls.Config)(nil)
	err := fmt.Errorf(TLS_NOT_CONFIGURED)
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) ListDir(path string) ([]os.FileInfo, error) {
	result := make([]os.FileInfo, 0)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("ListDir: %s", path))
	}
	entries := make([]os.DirEntry, 0)
	if entries, err = os.ReadDir(filepath.Join(Directory, path)); err == nil {
		for _, entry := range entries {
			if info, infoErr := entry.Info(); infoErr == nil {
				result = append(result, info)
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (f *FTP_SERVER) ListenAsync() {
	go func() {
		err := error(nil)
		if err = os.MkdirAll(Directory, DIR_MODE); err == nil {
			__debug(fmt.Sprintf("Starting FTP server on %s", ListenAddress))
			__debug("Anonymous login enabled")
			__debug(fmt.Sprintf("File storage directory: %s", Directory))
			err = ftpserverlib.NewFtpServer(&FTP_SERVER{}).ListenAndServe()
		}
		if err != nil {
			__debug(fmt.Sprintf("failed to start FTP server: %v", err))
		}
	}()
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) MakeDir(path string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("MakeDir: %s", path))
	}
	err = os.MkdirAll(filepath.Join(Directory, path), DIR_MODE)
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Mkdir(path string, perm os.FileMode) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Mkdir: %s, perm=%o", path, perm))
	}
	err = os.MkdirAll(filepath.Join(Directory, path), perm)
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) MkdirAll(path string, perm os.FileMode) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("MkdirAll: %s, perm=%o", path, perm))
	}
	err = os.MkdirAll(filepath.Join(Directory, path), perm)
	return err
}

func (f *FTP_SERVER) Name() string {
	result := SERVER_NAME
	return result
}

//goland:noinspection SpellCheckingInspection
func (p *PASV_PORT_RANGE) NumberAttempts() int {
	result := p.maxAttempts
	return result
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Open(path string) (afero.File, error) {
	result := afero.File(nil)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Open: %s", path))
	}
	result, err = os.Open(filepath.Join(Directory, path))
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	result := afero.File(nil)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("OpenFile: %s, flag=%d, perm=%o", path, flag, perm))
	}
	result, err = os.OpenFile(filepath.Join(Directory, path), flag, perm)
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) PutFile(path string, data []byte) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("PutFile: %s, size=%d", path, len(data)))
	}
	err = os.WriteFile(filepath.Join(Directory, path), data, FILE_MODE)
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Remove(path string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Remove: %s", path))
	}
	err = os.Remove(filepath.Join(Directory, path))
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) RemoveAll(path string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("RemoveAll: %s", path))
	}
	err = os.RemoveAll(filepath.Join(Directory, path))
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Rename(oldPath string, newPath string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Rename: %s -> %s", oldPath, newPath))
	}
	err = os.Rename(filepath.Join(Directory, oldPath), filepath.Join(Directory, newPath))
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) RenameFile(from, to string) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("RenameFile: %s -> %s", from, to))
	}
	err = os.Rename(filepath.Join(Directory, from), filepath.Join(Directory, to))
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) SetTime(path string, mtime, atime int64) error {
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("SetTime: %s, mtime=%d, atime=%d", path, mtime, atime))
	}
	return err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) Stat(path string) (os.FileInfo, error) {
	result := os.FileInfo(nil)
	err := error(nil)
	if Debug {
		__debug(fmt.Sprintf("Stat: %s", path))
	}
	result, err = os.Stat(filepath.Join(Directory, path))
	return result, err
}

//goland:noinspection GoBoolExpressions
func (f *FTP_SERVER) UserLeft(ftpClient ftpserverlib.ClientContext) {
	if Debug {
		__debug(fmt.Sprintf("UserLeft: ID: %d", ftpClient.ID()))
	}
}

//goland:noinspection GoBoolExpressions,GoUnusedParameter
func (f *FTP_SERVER) WelcomeUser(ftpClient ftpserverlib.ClientContext) (string, error) {
	result := BANNER
	err := error(nil)
	if Debug {
		__debug("WelcomeUser called")
	}
	return result, err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_FTPSERVER, logger.SKIP_STACK_FRAMES_BASE)
}

func (f *FTP_SERVER) loadConfig() (*ftpserverlib.Settings, error) {
	result := &ftpserverlib.Settings{
		ActiveConnectionsCheck:   ftpserverlib.IPMatchRequired,
		ActiveTransferPortNon20:  false,
		Banner:                   BANNER,
		ConnectionTimeout:        ConnectionTimeout,
		DefaultTransferType:      ftpserverlib.TransferTypeBinary,
		DeflateCompressionLevel:  0,
		DisableActiveMode:        false,
		DisableLISTArgs:          false,
		DisableMFMT:              false,
		DisableMLSD:              false,
		DisableMLST:              false,
		DisableSTAT:              false,
		DisableSYST:              false,
		DisableSite:              false,
		EnableCOMB:               false,
		EnableHASH:               false,
		IdleTimeout:              IdleTimeout,
		ListenAddr:               ListenAddress,
		PassiveTransferPortRange: newPasvPortRange(PasvPortMin, PasvPortMax, PasvPortMaxAttempts),
		PasvConnectionsCheck:     ftpserverlib.IPMatchRequired,
		PublicHost:               PasvHost,
		TLSRequired:              ftpserverlib.ClearOrEncrypted,
	}
	err := error(nil)
	return result, err
}

//goland:noinspection SpellCheckingInspection
func newPasvPortRange(minPort, maxPort, maxAttempts int) *PASV_PORT_RANGE {
	result := &PASV_PORT_RANGE{
		current:     minPort - 1,
		maxAttempts: maxAttempts,
		maxPort:     maxPort,
		minPort:     minPort,
	}
	return result
}
