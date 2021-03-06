// Code generated by go-bindata.
// sources:
// storage/migrations/mysql/20170523194156-init.sql
// storage/migrations/mysql/20170529135133-custom_capabilities.sql
// storage/migrations/mysql/20170607184935-keys.sql
// storage/migrations/mysql/20170705192449-alter_node.sql
// storage/migrations/mysql/20180417152820-add_unique_key.sql
// storage/migrations/mysql/bindata.go
// DO NOT EDIT!

package mysql

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _storageMigrationsMysql20170523194156InitSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xd1\x4a\xc3\x30\x14\x86\xef\xf3\x14\xe7\x6e\x2b\x3a\x10\xa1\x20\x94\x5d\x74\x6b\xd4\x62\xc9\xb4\x6b\x05\xef\x7a\xba\x64\x35\xb0\x26\x21\x49\x1d\x7b\x7b\xe9\x90\xd1\x89\xba\x7a\x17\xc8\xf7\x27\x87\xef\xfc\x64\x36\x83\xab\x56\x36\x16\xbd\x80\xd2\x90\x65\x4e\xe3\x82\x42\x11\x2f\x32\x0a\x95\xd2\x5c\x54\x30\x25\x00\x95\x3f\x18\x51\xc1\x07\xda\xcd\x3b\xda\xe9\xed\x4d\x00\x6c\x55\x00\x2b\xb3\x0c\x12\x7a\x1f\x97\x59\x01\x93\xc9\x75\x4f\x3a\x8f\xbe\x73\xe3\x58\xe4\xdc\x0a\x37\x12\x76\xc2\x39\xa9\x55\xca\x07\x78\x18\xfe\xce\x77\x86\xa3\x17\xbc\x82\x5a\x36\x52\xf9\xb3\xb7\x8f\x80\x15\x8d\x74\xde\xfe\x81\x94\x2c\x7d\x29\x29\x3c\xd1\xb7\xc1\xac\xd3\xd3\x31\x20\x01\x50\xf6\x90\x32\x3a\x4f\x95\xd2\xc9\xe2\x34\xc2\xf2\x31\xce\xd7\xb4\x98\x77\x7e\x7b\x17\x91\x6f\x5a\x37\x68\xb0\x96\x3b\xe9\xa5\x70\x5f\x7a\x7b\xd3\xf1\x7f\x64\xd4\x56\xef\x9d\xb0\x0c\xdb\xc1\x5a\xc2\xcb\x81\x57\x61\x7b\x89\xe3\x32\x66\x87\x7e\xab\x6d\x3b\xfe\x17\x83\x8d\xc8\x34\xf2\xb5\xef\x1b\xd5\x1c\x2e\xa6\xc6\x1a\x1c\xf6\x34\xd1\x7b\x45\x92\x7c\xf5\xfc\xa3\xd0\xe8\xec\xea\x58\xe1\x88\x7c\x06\x00\x00\xff\xff\xd8\x81\xbe\x43\xe8\x02\x00\x00")

func storageMigrationsMysql20170523194156InitSqlBytes() ([]byte, error) {
	return bindataRead(
		_storageMigrationsMysql20170523194156InitSql,
		"storage/migrations/mysql/20170523194156-init.sql",
	)
}

func storageMigrationsMysql20170523194156InitSql() (*asset, error) {
	bytes, err := storageMigrationsMysql20170523194156InitSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage/migrations/mysql/20170523194156-init.sql", size: 744, mode: os.FileMode(420), modTime: time.Unix(1572206971, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageMigrationsMysql20170529135133Custom_capabilitiesSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x94\x51\x6b\xdb\x30\x14\x85\xdf\xf5\x2b\x2e\x79\x71\xcc\x5c\x08\x63\x85\x81\x97\x81\xea\xa8\xa9\xa8\x2a\xaf\x96\x1c\xd8\x53\xa5\xcc\x4a\x66\x88\x2d\x63\xa9\xcb\xf6\xef\x87\xe7\x78\x4b\xdd\x10\xb6\xbd\x94\xfa\xcd\xe2\x9c\x73\xf5\x5d\xec\x83\x2e\x2e\xe0\x4d\x55\x6e\x5b\xed\x0d\xe4\x0d\x42\x49\x46\xb0\x24\x20\xf1\x15\x23\xa0\xbe\xe8\x46\xaf\xcb\x5d\xe9\x4b\xe3\x1e\x5a\xb3\x51\x30\x45\x00\xaa\xb6\x85\xc1\x45\xd1\x1a\xe7\x14\xac\x70\x96\xdc\xe0\x6c\xfa\x76\x16\x02\x4f\x25\xf0\x9c\x31\x58\x90\x6b\x9c\x33\x09\x41\x10\x75\x06\x67\x3c\x2d\x14\xf4\xcf\x60\xb8\x3c\x63\xa8\x75\x65\x06\xfd\x5f\x19\xbe\xe9\xdd\xa3\xf9\x87\x09\x39\xa7\xf7\x39\x81\x5b\xf2\x19\x94\xee\x51\x84\xf1\xfc\xd7\xd8\xe9\x13\xc0\x68\xb8\x7e\x74\xb8\x56\x88\x42\x20\x7c\x49\x39\x99\xd3\xba\xb6\x8b\x2b\xe8\x86\x09\x22\xe7\x8f\x7e\xf3\x3e\x46\x88\x72\x41\x32\x09\x94\xcb\xf4\xe4\x0a\xcf\xc7\x47\x03\x4c\x88\x00\xa6\x82\x30\x92\x48\x18\x59\x92\x94\x27\x58\x3e\x0b\x5a\xb7\x76\xef\x4c\xcb\x0f\x31\x87\xd7\x95\x69\x5d\x69\xeb\xee\xa4\xd9\x69\xbf\xb1\x6d\x35\x28\x1a\xbd\x35\xcc\xea\x42\xf8\xee\x03\xd8\xfe\x50\x61\x04\x93\xa3\x94\xc9\x28\x14\xb4\xfb\xb3\xea\xeb\x2c\xbd\x83\x63\x3c\xb8\xc1\x2b\xca\x97\xbf\x15\x1f\x3e\x42\x10\x84\xfd\xb2\x53\xfe\xb2\x34\x07\xdb\xe4\x79\xd0\x2b\x64\x3a\x16\x4d\xc6\xa6\xd7\xc8\x33\x3a\x9c\x9c\x12\xfe\x0f\x57\x8c\xd0\x22\x4b\x3f\x9d\xaa\x33\x15\xa3\x8c\x70\x7c\x77\xa6\xeb\xc6\xbf\xaf\x8a\x11\xca\x39\x4b\x93\xdb\xde\x23\x62\xf4\xa4\x3f\x17\x76\x5f\x23\x41\x97\x1c\x33\x10\xf7\x4c\xc8\xae\x49\x83\x77\x97\xb3\xd9\x2c\x00\x41\x24\x54\xc6\x39\xbd\x35\x0f\xde\x7c\xf7\x30\x87\x80\x56\x8d\x75\xae\x5c\xef\x0c\x14\x76\x5f\x83\xff\x5a\x3a\xe8\xd3\x4a\x5b\x07\xf1\xcf\x00\x00\x00\xff\xff\x42\xc8\x60\xda\x9e\x05\x00\x00")

func storageMigrationsMysql20170529135133Custom_capabilitiesSqlBytes() ([]byte, error) {
	return bindataRead(
		_storageMigrationsMysql20170529135133Custom_capabilitiesSql,
		"storage/migrations/mysql/20170529135133-custom_capabilities.sql",
	)
}

func storageMigrationsMysql20170529135133Custom_capabilitiesSql() (*asset, error) {
	bytes, err := storageMigrationsMysql20170529135133Custom_capabilitiesSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage/migrations/mysql/20170529135133-custom_capabilities.sql", size: 1438, mode: os.FileMode(420), modTime: time.Unix(1572206971, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageMigrationsMysql20170607184935KeysSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\xcd\x6a\xe3\x30\x18\x45\xf7\x7e\x8a\x0f\x66\xe1\x84\x49\xc0\x33\xb4\x2b\xd3\x85\x12\xab\x41\x54\x91\x13\x4b\x86\x66\x65\x29\x95\x9a\x0a\x62\xcb\x44\x0a\xe9\xe3\x17\x37\x3f\x60\xb7\xdd\x75\xab\xab\x7b\xcf\x11\x28\x9a\x4e\xe1\x6f\x6d\x77\x07\x15\x0c\x94\x6d\x14\x21\x2a\x70\x01\x02\xcd\x28\x06\xd9\x38\x6d\x24\xa0\x2c\x83\x79\x4e\xcb\x25\x03\x69\xb5\x84\x19\x59\x10\x26\x46\xff\x93\x31\xb0\x5c\x00\x2b\x29\x05\x54\x8a\xbc\x22\x6c\x5e\xe0\x25\x66\x02\x56\x05\x59\xa2\x62\x03\x4f\x78\x03\x8f\xa4\xe0\x22\xed\x0f\xbf\xa8\x56\x6d\xed\xde\x06\x6b\xfc\x6f\x01\xe6\x05\x46\x02\x03\x61\x19\x7e\x06\xe9\x83\x0a\x47\x5f\xb9\xd7\xea\xd8\x6a\x15\x8c\x96\x90\xb3\xeb\x8b\x46\x97\x58\x4e\xe4\x35\x1d\x0f\x07\x94\xd6\x07\xe3\xfd\xb9\xd6\xf7\x1d\x7d\xce\xa0\xcb\x85\x2f\x4d\x6f\xbc\xb7\xae\x21\x43\xe4\xed\x78\x9c\x46\x59\x91\xaf\x06\x24\x6e\x02\x53\xb5\xf9\x06\x78\x03\x94\x8c\xac\xcb\xa1\xe1\xcf\xbd\x4e\xb4\x4b\x26\xd0\x13\x9e\x74\x8a\xe1\xec\xf1\x07\x38\xc5\x78\x05\xff\xd2\xa8\xf7\x15\x32\x77\x6a\x22\x4e\x16\x0c\x51\xe0\x6b\xca\x45\x87\x8f\xef\xee\x93\x24\x89\x81\x63\x01\xb5\xf1\x5e\xed\x4c\x15\xcc\x7b\x80\x07\x88\x49\xdd\x3a\xef\xed\x76\x6f\x40\xbb\x53\x03\xe1\xcd\x7a\x38\xaf\x59\xd7\xc4\xe9\x47\x00\x00\x00\xff\xff\x23\x9f\x25\xdf\x69\x02\x00\x00")

func storageMigrationsMysql20170607184935KeysSqlBytes() ([]byte, error) {
	return bindataRead(
		_storageMigrationsMysql20170607184935KeysSql,
		"storage/migrations/mysql/20170607184935-keys.sql",
	)
}

func storageMigrationsMysql20170607184935KeysSql() (*asset, error) {
	bytes, err := storageMigrationsMysql20170607184935KeysSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage/migrations/mysql/20170607184935-keys.sql", size: 617, mode: os.FileMode(420), modTime: time.Unix(1572206971, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageMigrationsMysql20170705192449Alter_nodeSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe2\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x48\xc8\xcb\x4f\x49\x4d\x50\xf0\xf5\x77\xf1\x74\x8b\x54\x70\xf6\xf7\x09\xf5\xf5\x53\x48\x48\x4c\x49\x29\x4a\x2d\x2e\x4e\x50\x28\x4b\x2c\x4a\xce\x48\x2c\xd2\x30\x32\x35\xd5\x54\xf0\xf3\x0f\x51\xf0\x0b\xf5\xf1\x51\x70\x71\x75\x73\x0c\xf5\x09\x51\x50\x57\xb7\xe6\x42\x31\xdc\x25\xbf\x3c\x8f\x2c\xe3\x0d\xb0\x9b\x0e\x08\x00\x00\xff\xff\x05\x76\xd0\xcb\xba\x00\x00\x00")

func storageMigrationsMysql20170705192449Alter_nodeSqlBytes() ([]byte, error) {
	return bindataRead(
		_storageMigrationsMysql20170705192449Alter_nodeSql,
		"storage/migrations/mysql/20170705192449-alter_node.sql",
	)
}

func storageMigrationsMysql20170705192449Alter_nodeSql() (*asset, error) {
	bytes, err := storageMigrationsMysql20170705192449Alter_nodeSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage/migrations/mysql/20170705192449-alter_node.sql", size: 186, mode: os.FileMode(420), modTime: time.Unix(1572206971, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageMigrationsMysql20180417152820Add_unique_keySql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8f\xc1\x4b\xc3\x30\x14\x87\xef\xfd\x2b\x7e\xb7\x6e\xc8\x60\x88\x3d\x15\x0f\xcf\x36\x6e\xc5\x98\xb2\x36\x15\x6f\x26\x5b\xc3\x16\x5c\x9b\xd2\x04\xa7\xff\xbd\xd4\xa9\xd0\xe3\x7b\xf0\xbe\xef\x7b\xd1\x6a\x85\x9b\xce\x1e\x47\x1d\x0c\x9a\x21\x22\x2e\x59\x05\x49\x0f\x9c\x41\xf5\xae\x35\x0a\x94\xe7\xc8\x4a\xde\x3c\x0b\xa8\x77\xf3\xa5\xf0\xa1\xc7\xc3\x49\x8f\x8b\xdb\x24\x59\x42\x94\x12\xa2\xe1\x1c\xf4\x38\x5d\x2a\xdb\xaa\x74\x4e\x39\xe8\x41\xef\xed\xd9\x06\x6b\xbc\x42\xb6\x25\xb1\x61\xff\xc0\x49\x41\x6d\x3b\x1a\xef\xd5\x75\x7a\x9a\x14\x2f\x54\x65\x5b\xaa\xe6\x8a\x34\xca\x2a\x46\x92\xa1\x11\xc5\xae\x61\x28\x44\xce\x5e\x7f\x9b\x4a\xf1\x97\xbb\xf8\x59\x2c\xd3\x68\xf6\x5a\xee\x2e\x7d\x54\x17\x1b\x41\x1c\xf5\x8e\xd7\x72\xe2\xc4\x77\xc9\x7a\xbd\x8e\x51\x33\x89\xce\x78\xaf\x8f\xe6\x2d\x98\xcf\x80\x7b\xc4\x45\x37\x38\xef\xed\xfe\x6c\xd0\xba\x4b\x8f\x70\xb2\x1e\x57\x9a\x75\x7d\x9c\x7e\x07\x00\x00\xff\xff\xf2\x9d\x9b\xd9\x39\x01\x00\x00")

func storageMigrationsMysql20180417152820Add_unique_keySqlBytes() ([]byte, error) {
	return bindataRead(
		_storageMigrationsMysql20180417152820Add_unique_keySql,
		"storage/migrations/mysql/20180417152820-add_unique_key.sql",
	)
}

func storageMigrationsMysql20180417152820Add_unique_keySql() (*asset, error) {
	bytes, err := storageMigrationsMysql20180417152820Add_unique_keySqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage/migrations/mysql/20180417152820-add_unique_key.sql", size: 313, mode: os.FileMode(420), modTime: time.Unix(1572206971, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _storageMigrationsMysqlBindataGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x99\x5d\x6f\xdb\x46\xf6\xc6\xaf\xc5\x4f\xc1\x1a\x68\x21\xfd\xe1\xd8\x1c\xbe\xd3\x40\x2e\xfe\x4d\xbb\x40\xb0\x68\x0b\x6c\xbb\x57\x3b\x8b\x80\x2f\x33\x5e\xa2\xb2\xe4\x48\x72\x77\x9c\x20\xdf\x7d\xf1\x9b\x73\xe4\x3a\x69\x9a\xd8\x69\x50\xf4\x42\x16\x49\x71\xce\x9c\xb7\x79\x9e\x67\xc6\xe7\xe7\xe9\xb3\xed\xe4\xd2\x4b\xb7\x71\xbb\xfe\xe0\xa6\x74\xb8\x4d\x2f\xb7\x4f\x86\x79\x33\xf5\x87\xfe\x2c\x39\x3f\x4f\xf7\xdb\x9b\xdd\xe8\xf6\x17\xf1\xfa\xb0\xdd\xf5\x97\xee\xfc\x6a\xbe\xdc\xf5\x87\x79\xbb\xd9\x9f\x5f\xdd\xee\x5f\xae\xcf\xf3\xcc\x34\x59\x95\x17\xa6\x2b\x4d\x55\x3f\x99\x37\xf3\xe1\x6c\xff\x72\xfd\x90\x31\x9d\x29\x2a\x53\x14\x4f\xc6\x9b\xfd\x61\x7b\xf5\x62\xec\xaf\xfb\x61\x5e\xcf\x87\xd9\xed\x1f\x66\xa2\xce\x1a\xd3\x96\x5d\x51\x3d\xf9\xd9\xdd\x3e\x70\x4c\x93\x55\xa6\xcb\xcb\xb2\x7b\xd2\xaf\x0f\x6e\xf7\x62\xb3\x9d\xdc\x43\x46\xb6\x59\x69\x1a\x53\xe5\x6d\x9e\x3d\xe9\xa7\xe9\xc5\xcd\x66\x7e\x79\xe3\x5e\xfc\xec\x6e\x3f\x3a\xfa\x98\xd3\xcb\x2d\xaf\x7d\xf3\x43\xfa\xfd\x0f\x3f\xa5\xdf\x7e\xf3\xfc\xa7\x2f\x92\xe4\xba\x1f\x7f\xee\x2f\x5d\x1a\xdf\x4c\x92\xf9\xea\x7a\xbb\x3b\xa4\xcb\x64\x71\x32\xdc\x1e\xdc\xfe\x24\x59\x9c\x8c\xdb\xab\xeb\x9d\xdb\xef\xcf\x2f\x5f\xcd\xd7\x3c\xf0\x57\x07\xbe\xe6\xad\xfc\x3d\x9f\xb7\x37\x87\x79\xcd\xcd\x36\x0e\xb8\xee\x0f\xff\x39\xf7\xf3\xda\x71\xc1\x83\xfd\x61\x37\x6f\x2e\xe3\x6f\x87\xf9\xca\x9d\x24\xab\x24\xf1\x37\x9b\x31\x55\xcf\xfe\xe1\xfa\x69\xc9\x45\xfa\xaf\x7f\x33\xed\x69\xba\xe9\xaf\x5c\x2a\xc3\x56\xe9\xf2\xf8\xd4\xed\x76\xdb\xdd\x2a\x7d\x9d\x2c\x2e\x5f\xc5\xbb\xf4\xe2\x69\x8a\x57\x67\xdf\xbb\xff\x62\xc4\xed\x96\xd1\x6d\xee\xbf\xbe\xf1\xde\xed\xa2\xd9\xd5\x2a\x59\xcc\x3e\x0e\xf8\xe2\x69\xba\x99\xd7\x98\x58\xec\xdc\xe1\x66\xb7\xe1\xf6\x34\xf5\x57\x87\xb3\x6f\xb1\xee\x97\x27\x18\x4a\xbf\x7c\x79\x91\x7e\xf9\xcb\x89\x78\x12\xe7\x5a\x25\x8b\x37\x49\xb2\xf8\xa5\xdf\xa5\xc3\x8d\x4f\x65\x1e\x99\x24\x59\xbc\x10\x77\x9e\xa6\xf3\xf6\xec\xd9\xf6\xfa\x76\xf9\xd5\x70\xe3\x4f\xd3\xcb\x57\xab\x64\x31\xae\xbf\x3d\x7a\x7a\xf6\x6c\xbd\xdd\xbb\xe5\x2a\xf9\x5c\xfe\x60\x46\xec\xff\x8e\x21\xb7\xdb\x89\xdf\xfa\x70\xb8\xf1\x67\x5f\xe3\xfa\x72\x75\xca\x1b\xc9\x9b\x24\x39\xdc\x5e\xbb\xb4\xdf\xef\xdd\x81\x94\xdf\x8c\x07\xac\xc4\xf8\xb4\x1e\xc9\x62\xde\xf8\x6d\x9a\x6e\xf7\x67\x7f\x9b\xd7\xee\xf9\xc6\x6f\xef\xc6\x69\x09\x8f\xcf\xef\x59\x88\x35\x4c\x53\x2d\x63\xb2\xd8\xcf\xaf\xe2\xfd\xbc\x39\xd4\x65\xb2\xb8\x62\xf9\xa7\x77\x46\xbf\xdb\x4e\x2e\x3e\xfc\x69\xbe\x72\x29\x6d\x72\xc6\x15\xf3\xc4\x56\x59\xfa\xf9\xdd\xb9\x56\xe9\xf7\xfd\x95\x5b\xae\x74\x06\xe6\xd4\x28\xfd\x7c\xc6\xec\xc9\x9b\x0f\x8c\xfd\x71\x7e\xc5\xd8\xe8\xcd\xdb\x43\x71\xf4\x83\x43\xf1\x75\xb9\xba\xef\xf9\xdb\x06\x08\xed\x63\x06\x08\x6e\xb9\xfa\x35\xd0\xdf\x58\xd0\xe8\x7f\xdf\xc8\xf3\xfd\x37\xf3\x6e\xb9\x4a\x87\xed\x76\x7d\x7f\x74\xbf\xde\x7f\x24\xf2\xdb\xbd\x04\xee\x76\xbe\x1f\xdd\xeb\x37\xf7\x46\x6b\x4b\xd0\xe5\x2f\x14\x50\xbe\xbb\xc3\x93\xef\x00\x89\xb7\x01\xf7\xf9\x66\x3e\xfc\xf8\x72\x9d\x3e\xd5\x56\x59\x9e\xd8\x60\xbc\x0d\xed\x60\x43\xd6\xda\x90\x65\xef\xff\x78\x6f\x43\x57\xda\xd0\xe5\x36\x4c\xc6\x86\xb2\xb7\x61\x2c\x6c\x28\x32\x1b\x4c\x69\x43\x5b\xdb\xe0\xbc\x0d\xbe\x90\x7b\xd7\xd8\x50\x3b\x1b\xf2\xc1\x86\xa2\xb7\xc1\x64\x36\xf4\xc6\x86\x3c\x13\x3b\xd5\x64\x43\x53\xda\x50\x0f\x36\x4c\x7c\xe7\x36\x8c\x9d\x0d\x83\x3e\xcb\x2a\xb1\xd7\xf4\x36\x0c\xbd\x0d\x75\x69\x43\x51\xd9\x30\x64\x36\xe4\xb5\x0d\x39\x3e\x74\x36\x18\xec\x0c\xf2\x71\x9d\x0d\x5d\x26\xfe\xb5\x9d\x8c\x63\xbc\x69\x6c\x18\x5b\x1b\x7c\x63\x43\xde\xd8\xd0\x36\xea\xeb\xa8\x76\x6b\x1b\xda\xc2\x86\x7e\xb0\xa1\xaa\xe5\xde\xd4\x36\x0c\x93\x0d\x2d\xf6\x72\xb1\x5b\x57\x36\x94\xce\x06\x57\xd8\xd0\xe6\x36\x94\xb9\x0d\x86\x98\xbc\x0d\x45\x6e\x43\xd6\xdb\xd0\x55\xf2\x7e\x35\xda\x50\x95\x92\x9f\xbc\x92\x1c\xf2\x5b\xe1\x6d\x30\xad\x0d\x95\xb1\x61\x34\x36\x64\x8d\x0d\x13\xb9\x9c\x6c\x28\x06\xb9\x26\x0e\x37\xd9\x50\x4e\x32\xae\xc6\x96\xda\x20\x9f\x43\x61\x43\x36\xda\x60\x72\x8d\x8f\xda\x34\x36\x54\x9d\x0d\x99\xb1\xa1\x2b\x24\x97\x0d\xfe\x7a\xc9\x7f\xeb\x6d\x18\x9c\x0d\x4d\x21\xfe\x57\xad\x0d\xae\xb4\x61\x1a\xc5\xef\xa2\x51\x7b\xb5\x0d\x63\x6e\x43\xd1\xd9\xd0\x77\x32\xef\xd8\x8b\x9f\x4d\x2b\xbe\x7b\x67\xc3\x88\xad\x46\xea\xde\x17\x92\xe3\x61\x94\xbc\x54\xbd\xe4\xb0\xca\x6d\xf0\x9d\xf8\x3b\x34\xe2\x03\xf9\x34\x95\x0d\xad\xd6\x7f\x72\x62\xaf\x35\xd2\x17\xf9\x68\x43\x33\xd9\x90\x77\xf2\x29\x46\xa9\x27\xe3\xc9\x57\xcf\x7d\x21\x9f\x42\xfb\x89\x98\x2b\x7a\xb4\x96\x3a\x51\x0b\x62\x26\xf6\x63\x3d\xc8\x87\xe3\x79\x29\xf1\xf9\x5c\xc6\x13\xc7\xd4\x48\xdf\x75\xad\xc4\xd4\x38\xe9\x27\x62\xea\x8c\x0d\xb5\xd7\x98\xe8\xe9\x56\x7a\xb0\xab\xa5\x5e\xd4\xaa\xaf\x6c\x68\x32\x1b\x2a\xed\x57\xc6\xe2\x9f\x37\x36\x34\x5e\x7a\x8c\x1e\xaf\x74\x8d\x74\x93\xd8\xa0\x86\xd3\x20\x71\x61\x1f\xbf\xc6\x41\x72\x51\x31\x97\x11\x5b\xf4\x32\x35\x23\xae\xba\x96\xfe\xc5\x47\xfa\xb5\xd6\xbe\x21\x87\xf8\x4b\x1f\x93\x5b\xfa\x9d\x58\x89\x73\xd0\xb5\x64\xf0\xab\xb2\xc1\x8c\x36\xe4\x8c\xa7\xd6\xf4\x7a\x2f\xcf\xc8\x1f\x63\x26\x9d\xb7\xac\x64\xbd\x37\xfc\x56\xc9\x7a\xa1\xd6\x53\x66\x83\xa3\x7f\x46\x1b\x5c\xaf\xbd\xc4\x9a\x6b\xe5\xdd\xac\x7e\x1b\x37\xf8\x4c\xad\xc4\x45\x0f\x96\x85\x8c\xcf\xf2\xe3\x7b\x27\x47\x81\xf1\x08\x08\x53\x56\x7c\x9f\xda\x38\x72\xe7\x3d\xb5\x92\x2c\x16\x8f\xc1\xc7\xd3\x64\x81\x0c\x7a\x9c\x86\x3d\x39\x4d\x16\xab\x3b\x02\x7c\xc4\x6c\x44\xf1\x7f\x91\xd0\xef\x47\x11\x19\xfd\x4e\x36\x3d\x3e\x33\x1f\xd3\x2c\x77\x52\x23\x8a\x85\x8b\xa7\xef\x12\xcf\x6b\x28\xf9\x22\xfd\x84\x34\xa4\x30\xf2\x45\xda\x94\xe5\x69\x0a\xb7\x5e\xdc\xa7\xde\x65\x99\x67\xab\xf8\x1c\xc6\xbc\x10\x46\xfd\xe7\x66\x0e\x4b\x53\x35\x79\x9e\xd5\x5d\x63\x4e\xd3\x6c\xf5\x26\x59\xf4\x78\xf5\x55\x4c\xcc\xeb\x98\x8d\x8b\x54\x93\x82\xcb\x17\xf1\xef\x9b\xbb\x72\xf7\xa7\x0f\x67\x43\xdd\x4a\x3c\xfb\xed\x4e\xe2\x53\xc9\x91\x05\x1e\x89\xcd\x28\xa9\x0d\xf7\xc8\x11\x10\xf2\xb2\x88\x00\x70\x16\x5f\x03\x40\x03\x40\xa3\x90\x05\xf6\xeb\x42\xde\x8d\x60\xd8\xc8\x37\x0b\xac\x6f\x05\x8c\xf9\xce\xb9\xf7\x02\x40\x2c\x5a\x16\x56\x55\x08\x08\x61\x0b\x62\x8e\x20\xc1\xbb\x93\xd8\x64\x2c\xc0\xe2\x15\x80\x22\x01\x2a\xa0\x97\xf8\x3a\x09\x41\x0f\x4a\x7a\xc4\xe1\x95\x94\x00\xcd\x6e\x14\xe2\xc0\x7f\x88\x9b\x85\x0f\xd0\x10\x07\x1f\x97\x09\x59\x41\x16\x90\x7e\x35\x08\x81\x65\x93\x90\x0c\xdf\x90\x25\x84\x5d\xd6\x4a\xe2\xa5\x80\x37\x20\x09\x21\xe4\xf8\x99\x09\x40\x00\x0e\xbc\x47\x9c\xf8\x1d\xe7\xec\xc4\x57\x80\xd0\x38\x01\x4c\x88\x85\x7c\x93\xe7\x52\x09\xb2\xef\x25\x0e\xf2\x08\xb0\xf2\x1c\x40\x03\xe4\x88\x99\x7a\x40\x22\x8d\x82\x38\x04\xe8\x94\x64\x89\xb9\x51\xf2\x07\x9c\x20\x10\xc8\xdb\x67\x92\x03\x08\x03\xb0\x8b\x62\x60\xd0\x67\x90\x31\xbe\x19\xc9\x21\xa4\x0b\x00\xd6\x8d\x90\x17\x35\x60\x4e\x0f\xf1\x40\x06\xc4\xd8\xca\x6f\xb1\x36\xad\x8c\xa9\x75\x9c\x9f\x84\x48\x4c\xa7\xb9\xe8\xa4\x3e\x00\x2f\x64\x4a\xed\x98\x2f\x12\x74\x23\xa2\x86\xeb\x48\x1a\x83\x00\x3e\x63\xa9\xa1\x73\x92\x9f\xb6\x94\x3c\x03\xe0\x90\x1d\x3d\x83\x5d\x83\x50\xc8\x94\xd4\x5a\x25\xbe\x52\xae\x2b\xed\x21\xea\x46\x9d\x00\xf6\x52\xe7\xed\x3a\x21\xb8\xe1\x98\xeb\x41\xfa\x9a\xba\x21\x08\xf0\x2f\xcf\x25\xdf\x90\x3e\x76\x11\x86\x85\x93\xba\x46\xa2\xc8\x25\x67\x51\x38\xa9\xbf\xb1\x47\x4b\xed\x99\x5e\xf2\x35\x36\x36\x94\xc4\x5c\xd8\x50\x8e\x32\x36\xd6\xb9\x16\xb1\x41\xed\x21\xa9\x52\xc5\x09\x22\x88\xfb\x28\x26\x1a\xa9\x17\xf3\x66\x83\xf6\x4c\x23\x35\x66\x1d\x60\x8f\x39\x33\x2f\xb5\x25\xb7\xe4\x07\xd1\x50\x77\xd2\x77\x7c\xbb\x41\x72\xd9\x43\x8e\x9d\x3c\x1b\x10\x54\x46\x09\x98\xfe\x6d\xa5\xc7\x06\xbd\x1f\x95\x1c\xc9\x9f\x6f\x65\x3c\x31\x4c\x2a\xa0\x10\x2c\x10\x7c\x56\x8a\x60\xa3\xbe\x31\x07\x9d\xda\x2a\x45\x94\x78\x15\x7c\xd8\xc2\x0f\x04\x52\x14\xa6\x85\xf4\x0f\x7d\x44\xbd\xe8\x75\x6a\x80\x58\x03\x3b\xf0\x8f\x75\x05\x7e\xc5\xbc\xe7\xd2\x9f\xb1\x17\x10\x20\xb9\xe0\x06\xbe\x0c\x2a\x8c\xa2\x18\x1d\x24\xff\xe0\xd3\xa4\x82\x33\x8a\xe4\x5e\xd6\x05\xeb\x1c\xd1\x40\xfe\x27\x15\xd6\xac\x01\x7e\x67\x3d\x74\x2a\x4c\xb1\x4b\x5e\x11\x32\xed\xa8\xb6\x72\x59\xc3\x85\x8a\x20\xd6\x2a\x63\xf1\x93\xf8\x5b\x5d\x97\xf4\x1a\x42\x8c\x39\x88\x99\xf9\x88\x07\x3c\x68\x7b\x11\xdd\xc4\x59\x28\x26\x90\x3f\xfc\x46\x4c\x82\x23\x6c\x26\xe8\x23\xea\x15\x05\xb7\xd6\x9c\x75\xc5\xef\xa5\xe2\xab\xd1\xf5\x4b\x5e\x58\x47\xb1\x1f\xc0\x5a\x27\xbd\x4e\x8e\xf1\x23\xe6\xb3\x90\xda\x93\x7f\xfa\x0f\x1b\x60\x18\x22\x1f\xff\xf0\xe5\x58\x7b\xe2\xc6\x97\x28\x9c\x1a\xe9\x51\x7a\x0b\xf1\x1b\xd7\xd5\x24\x6b\x9e\x79\xc0\x7e\x7a\x31\x0a\xfa\xf2\x57\x3f\x99\xdb\x1f\x85\x67\x2f\x6b\x8b\x35\x07\xbe\xb3\xce\xa9\x15\x6b\x9a\x35\xf2\x2e\x17\xf1\xa1\xde\xc4\x02\xde\xb0\xa9\xe8\x9c\x6c\xaa\x1e\x23\xb8\x3e\xc8\x92\x9f\x59\x7f\x7d\x70\xae\x07\xc9\xb1\x8f\x1c\x0f\x3e\x42\x9d\x7d\xd0\x97\x3f\x2c\xd6\x1e\x92\xd5\x3f\x4b\xbb\x7d\x34\x67\x2a\xe5\x4c\x59\xb4\x7f\x4d\x2d\x77\x77\xa6\xfb\x77\x77\xfb\xc9\xe2\x0d\x41\x00\x69\x20\x78\xea\x5e\x77\x5b\x99\x10\x0a\x02\xc2\xeb\x8e\x0b\xf0\x61\x71\x23\xb0\x9c\x12\x2b\x40\x01\x89\x02\x24\x00\x76\xae\xbb\x3f\x04\x08\x40\xc8\x2e\x0d\xb0\x00\x28\xd8\x4d\x42\xba\x80\x16\x44\x89\x1d\x88\x9f\x9d\x2e\xa4\xd3\xf5\x42\x7c\xb5\xee\x04\x4b\x25\x42\xa7\x3b\x40\x80\x8c\x5d\x68\x71\x14\x12\x8d\x88\x83\x08\x26\x83\x7c\xd8\xb9\x01\x0a\xf1\x64\xa2\x15\x9b\xf1\xd4\x42\x77\xb2\x10\x16\xa0\x06\x80\x00\xbc\x00\x0f\xc4\x53\xab\x50\xc9\x8d\x88\xcc\x46\x77\xd6\x99\xe6\x05\x5b\x08\x14\x40\xb0\x68\xe5\x7d\x00\x10\xd1\x96\xab\x20\x84\x34\x98\x0b\xdf\xe3\xa9\x47\xa1\x64\x59\x49\xae\x10\x26\x10\xb4\xd1\x53\x1c\xc4\x40\x3c\x59\xd2\xdd\x3c\xc2\xb0\x1a\x7f\x3d\xe9\x20\x2f\x80\x16\xb9\x23\xf7\xc3\x28\x24\x12\x4f\x44\x54\x78\x32\x4f\xad\x82\x0d\x70\x8d\xef\x23\x8e\x72\xc9\x63\x14\x0e\x5e\xc4\x2e\xd7\x08\x12\x88\xc1\x29\x21\x45\xa1\xea\x95\x60\x54\xf8\x10\x1b\xcf\x01\xf1\x28\xe6\x06\x3d\x2d\xf2\x92\x13\xa7\xf3\x94\x3a\x6f\x8c\xd3\x48\x7c\x8d\x8e\x71\x4a\x20\xd4\x0f\x01\x03\xa8\x0f\x9d\x90\x26\x64\x8e\x10\xa0\xd7\xe2\x69\xc8\x28\x62\x14\x61\x19\xc9\x7e\x90\xbe\xc5\x7e\xa7\x22\x26\xd6\x50\x4f\xc1\x8c\xfa\x4c\x0d\xa9\xdd\x54\xcb\x35\x3d\x12\xc9\xbd\x93\x67\xf4\x6d\x3c\x21\x9b\xe4\xa4\x83\x5a\xc5\x13\x86\x4a\x4e\xad\x20\x73\x62\x22\xaf\xf4\x12\x84\x14\x4f\xd6\x0a\x99\x8b\x1c\x90\x37\x48\x1d\x5f\xa2\xb8\x32\x42\xc4\xc4\x48\xfd\x11\xfc\x99\x8a\x09\xea\x8c\x48\xe5\xb7\x63\xbe\xf1\x8f\xf8\x7a\x15\x5b\xbd\x91\x5e\xa4\x47\x21\x4c\xe2\x62\xdd\xb0\x26\xe8\x09\x84\x02\x6b\xa4\x53\x31\x4a\xdd\x79\x9f\x4d\x87\xd7\x13\x33\xfa\x6f\xac\x74\x8e\x4a\x37\x60\xb9\xd4\xcf\x6b\x5f\x23\x2c\xe9\x73\x72\x1c\x85\x87\x53\x61\x3f\x8a\x28\x76\x2a\xe2\x11\x12\xac\xf1\x56\xc5\x2a\x64\x8c\x58\xa3\x27\x89\xaf\x55\x71\x8d\xd8\xa2\xf6\x31\x46\xa3\xfd\x37\x89\x68\x30\x2a\xf8\xe2\x69\x4d\x26\x3e\x46\xf1\xac\xe2\x3d\x12\xb8\xd7\x7c\xd6\x92\xf3\x92\x5a\x0f\x92\x37\xf2\x14\x73\x32\xc9\x09\x12\xf1\x91\xe3\x4a\xfb\x65\x2c\x75\x03\xd4\xbc\x9f\xec\x11\x35\x9d\x6e\x56\xd8\x68\xb2\xe6\x1e\x77\xba\xf2\x2e\x8c\x7e\x56\x76\x7f\xd7\xf8\x03\xe8\xfc\xb7\xff\xaa\x7b\x30\x7f\xbf\x3b\xdb\x1f\x24\xec\xdf\xc9\xcc\x9f\xc3\xd0\xef\x4b\x83\x52\x72\x6d\x9a\xbf\x26\x23\xdf\xfd\xc7\xf4\xff\xef\xfe\x61\xfa\xa9\xbc\xec\x72\x59\xd7\x53\x25\x87\x0b\x60\x03\x38\x36\xea\xba\x8b\x18\xd6\x0b\xff\x94\x7a\x5a\x9d\xe9\xc1\x07\xeb\xbb\xd1\x8d\x40\xa6\x87\x2d\xd9\x24\xa7\xd2\xbc\xc3\x7d\xa3\x27\xa4\x5c\x33\x96\x4d\x24\xa2\x39\x72\xaf\xd7\x03\x8b\x49\x4f\x97\x33\x39\x00\x69\x54\x78\x47\x6c\x1a\x84\x9b\x1a\x3d\x7d\xf6\xba\x39\x8f\xa7\xa5\x95\xac\x6d\x6c\xc6\x8f\xfa\x08\xa7\x45\x9f\x27\x39\x3c\x29\x75\x33\x98\xeb\x61\x4c\xae\x87\x38\xc4\xc9\x38\xee\x27\xdd\xe8\xc6\x93\xed\x4a\xf3\x51\xaa\x4f\x85\x6e\xb4\x8c\xdc\xb3\xd1\xf5\x47\x8c\x32\xba\xa1\x32\xa2\x11\x1a\xfd\xcf\x01\xbf\xe3\x67\xa5\x27\xe7\x6c\xd2\xc0\x7f\x78\xad\xd4\x93\xf1\x69\x14\x3c\x81\x17\x8b\xf1\x7f\x01\x00\x00\xff\xff\xa4\xf5\xa4\x9a\x00\x20\x00\x00")

func storageMigrationsMysqlBindataGoBytes() ([]byte, error) {
	return bindataRead(
		_storageMigrationsMysqlBindataGo,
		"storage/migrations/mysql/bindata.go",
	)
}

func storageMigrationsMysqlBindataGo() (*asset, error) {
	bytes, err := storageMigrationsMysqlBindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "storage/migrations/mysql/bindata.go", size: 20480, mode: os.FileMode(420), modTime: time.Unix(1572286257, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"storage/migrations/mysql/20170523194156-init.sql":                storageMigrationsMysql20170523194156InitSql,
	"storage/migrations/mysql/20170529135133-custom_capabilities.sql": storageMigrationsMysql20170529135133Custom_capabilitiesSql,
	"storage/migrations/mysql/20170607184935-keys.sql":                storageMigrationsMysql20170607184935KeysSql,
	"storage/migrations/mysql/20170705192449-alter_node.sql":          storageMigrationsMysql20170705192449Alter_nodeSql,
	"storage/migrations/mysql/20180417152820-add_unique_key.sql":      storageMigrationsMysql20180417152820Add_unique_keySql,
	"storage/migrations/mysql/bindata.go":                             storageMigrationsMysqlBindataGo,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"storage": &bintree{nil, map[string]*bintree{
		"migrations": &bintree{nil, map[string]*bintree{
			"mysql": &bintree{nil, map[string]*bintree{
				"20170523194156-init.sql":                &bintree{storageMigrationsMysql20170523194156InitSql, map[string]*bintree{}},
				"20170529135133-custom_capabilities.sql": &bintree{storageMigrationsMysql20170529135133Custom_capabilitiesSql, map[string]*bintree{}},
				"20170607184935-keys.sql":                &bintree{storageMigrationsMysql20170607184935KeysSql, map[string]*bintree{}},
				"20170705192449-alter_node.sql":          &bintree{storageMigrationsMysql20170705192449Alter_nodeSql, map[string]*bintree{}},
				"20180417152820-add_unique_key.sql":      &bintree{storageMigrationsMysql20180417152820Add_unique_keySql, map[string]*bintree{}},
				"bindata.go":                             &bintree{storageMigrationsMysqlBindataGo, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
