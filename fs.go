package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
)

type FS struct {
}

type decoder func(v any) ([]any, bool)

var ErrUimplemented = errors.New("operation is unimplemented")

// TODO: handle offset and position.
func (fs *FS) write(dec decoder) []any {
	var r struct {
		FD     int `json:"fd"`
		Buffer string

		// options ignored: offset, length and position.
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	// forward the write syscall.
	n, err := syscall.Write(r.FD, []byte(r.Buffer))
	if err != nil {
		return []any{err.Error()}
	}

	// return the bytes written.
	return []any{nil, n}
}

func (fs *FS) chmod(dec decoder) []any {
	var r struct {
		Path string
		Mode uint32
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	// forward the chmod syscall.
	err := os.Chmod(r.Path, os.FileMode(r.Mode))
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) chown(dec decoder) []any {
	var r struct {
		Path string
		UID  int `json:"uid"`
		GID  int `json:"gid"`
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	// forward the chown syscall.
	err := os.Chown(r.Path, r.UID, r.GID)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) close(dec decoder) []any {
	var r struct {
		FD int `json:"fd"`
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	// forward the close syscall.
	err := syscall.Close(r.FD)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) fchmod(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) fchown(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) ftruncate(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) lchown(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) link(dec decoder) []any {
	var r struct {
		Path string
		Link string
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	// forward the link syscall.
	err := syscall.Link(r.Path, r.Link)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) lstat(dec decoder) []any {
	var r struct {
		Path string
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	var s syscall.Stat_t

	err := syscall.Lstat(r.Path, &s)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil, convertStat(s)}
}

func (fs *FS) mkdir(dec decoder) []any {
	var r struct {
		Path string
		Perm uint32
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	err := syscall.Mkdir(r.Path, r.Perm)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) read(dec decoder) []any {
	var r struct {
		FD     int `json:"fd"`
		Length int

		// options ignored: buffer, offset and position.
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	// create a slice with the amount of bytes we want to read.
	b := make([]byte, r.Length)

	// forward the read syscall.
	n, err := syscall.Read(r.FD, b)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil, n, b}
}

func (fs *FS) readdir(dec decoder) []any {
	var r struct {
		Path string
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	files, err := os.ReadDir(r.Path)
	if err != nil {
		return []any{err.Error()}
	}

	var ns []string

	// collect the file names.
	for _, f := range files {
		ns = append(ns, f.Name())
	}

	return []any{nil, ns}
}

func (fs *FS) readlink(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) rename(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) rmdir(dec decoder) []any {
	var r struct {
		Path string
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	err := syscall.Rmdir(r.Path)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) stat(dec decoder) []any {
	var r struct {
		Path string
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	var s syscall.Stat_t

	err := syscall.Stat(r.Path, &s)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil, convertStat(s)}
}

func (fs *FS) symlink(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) truncate(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) utimes(dec decoder) []any {
	return []any{ErrUimplemented.Error()}
}

func (fs *FS) open(dec decoder) []any {
	var r struct {
		Path  string
		Flags int
		Perm  uint32
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	fd, err := syscall.Open(r.Path, r.Flags, r.Perm)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil, fd}
}

func (fs *FS) fstat(dec decoder) []any {
	var r struct {
		FD int `json:"fd"`
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	var s syscall.Stat_t

	err := syscall.Fstat(r.FD, &s)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil, convertStat(s)}
}

func (fs *FS) unlink(dec decoder) []any {
	var r struct {
		Path string
	}

	if msg, ok := dec(&r); !ok {
		return msg
	}

	err := syscall.Unlink(r.Path)
	if err != nil {
		return []any{err.Error()}
	}

	return []any{nil}
}

func (fs *FS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op := strings.TrimPrefix(r.URL.Path, "/fs/")
	if op == "" {
		http.Error(w, "operation is required", http.StatusBadRequest)
		return
	}

	dec := func(v any) ([]any, bool) {
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			http.Error(w, "failed to decode body", http.StatusBadRequest)
			return []any{err.Error()}, false
		}

		return nil, true
	}

	var res []any

	switch op {
	case "write":
		res = fs.write(dec)

	case "chmod":
		res = fs.chmod(dec)

	case "chown":
		res = fs.chown(dec)

	case "close":
		res = fs.close(dec)

	case "fchmod":
		res = fs.fchmod(dec)

	case "fchown":
		res = fs.fchown(dec)

	case "fstat":
		res = fs.fstat(dec)

	case "ftruncate":
		res = fs.ftruncate(dec)

	case "lchown":
		res = fs.lchown(dec)

	case "link":
		res = fs.link(dec)

	case "lstat":
		res = fs.lstat(dec)

	case "mkdir":
		res = fs.mkdir(dec)

	case "open":
		res = fs.open(dec)

	case "read":
		res = fs.read(dec)

	case "readdir":
		res = fs.readdir(dec)

	case "readlink":
		res = fs.readlink(dec)

	case "rename":
		res = fs.rename(dec)

	case "rmdir":
		res = fs.rmdir(dec)

	case "stat":
		res = fs.stat(dec)

	case "symlink":
		res = fs.symlink(dec)

	case "truncate":
		res = fs.truncate(dec)

	case "utimes":
		res = fs.utimes(dec)

	case "unlink":
		res = fs.unlink(dec)

	default:
		http.Error(w, fmt.Sprintf("unhandled operation: %s", op), http.StatusBadRequest)
		return
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to write body", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
}

type JSStat struct {
	// same as `*os.Stat`.
	Dev     int32  `json:"dev"`
	Ino     uint64 `json:"ino"`
	Mode    uint16 `json:"mode"`
	Nlink   uint16 `json:"nlink"`
	Uid     uint32 `json:"uid"`
	Gid     uint32 `json:"gid"`
	Rdev    int32  `json:"rdev"`
	Size    int64  `json:"size"`
	Blksize int32  `json:"blksize"`
	Blocks  int64  `json:"blocks"`

	// specific to `fs.Stat`.
	AtimeMs int64 `json:"atimeMs"`
	MtimeMs int64 `json:"mtimeMs"`
	CtimeMs int64 `json:"ctimeMs"`

	// wether not the file is a directory. used for `isDirectory()`.
	Dir bool `json:"dir"`
}

func convertStat(s syscall.Stat_t) JSStat {
	// create a fs-compatible stat.
	// https://github.com/golang/go/blob/1cc19e5ba0a008df7baeb78e076e43f9d8e0abf2/src/syscall/fs_js.go#L165
	st := JSStat{
		Dev:     s.Dev,
		Ino:     s.Ino,
		Mode:    s.Mode,
		Nlink:   s.Nlink,
		Uid:     s.Uid,
		Gid:     s.Gid,
		Rdev:    s.Rdev,
		Size:    s.Size,
		Blksize: s.Blksize,
		Blocks:  s.Blocks,
	}

	st.AtimeMs = s.Atimespec.Sec * 1000
	st.MtimeMs = s.Mtimespec.Sec * 1000
	st.CtimeMs = s.Ctimespec.Sec * 1000

	// https://github.com/golang/go/blob/50034e9faac531e0e4d6cbf4d172462ca23c9be2/src/os/stat_darwin.go#L12-L42.
	if s.Mode&syscall.S_IFMT == syscall.S_IFDIR {
		st.Dir = true
	}

	return st
}
