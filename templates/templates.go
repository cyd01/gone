package templates

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDir struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDir) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDir{fs: _escLocal, name: name}
	}
	return _escDir{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/.edit.html.swp": {
		local:   ".edit.html.swp",
		size:    12288,
		modtime: 1442683616,
		compressed: `
H4sIAAAJbogA/+yaTW/cRBjHp0gVqpoAKlyRJkYIiLBnvaFNs/VGSpOiFqU0kA2kqoo0a8+uR7E9lj0b
bxItHKCIb4DEF+DOBa70iMSJG+LEmRM3Dkj8Z3azvKhqbhSk+Uk/zavHzzMzOUTrfuu9W7fpatAm4DlC
fvji972b3z94+uELhMSqSFRFziQ7LOtw9Yrfvuqv+eHllSARiYy5FkmQqlqLUaVKESTkw+mCrFHVQV3y
WLChYnUVs6HU6agfxCpng3FRoLsQTIu8zLBIzWrNtYwZFtVBqvPs7IgcDsffGOmBf3WRrLTDlmm+5C3R
5y/tPemoHA6Hw+FwOBwOx7+ILs+Rj1A+NWu/MSvP/aN0OBwOh8PhcDgcDofD8f+FJ4R0LxLy2gKxv/+f
/v//7bOEfAofwLtwH74DV+EV+CK8CM/Dn58h5CH8En4OP4OfwI/hMZRQwATuww3YhRG8BM/D3xYJ+QX+
CL+D38Cv4QfwLfgm7MAWfBW+Ain8FXH/BL+CGm7D67AD2zBcmOa2uPAkd9nhcDgcDofD4fhvEZ0JM5/e
rmMm66vkaJ1ciFgiD1GiMlBVbmoXIlmUI031USm6Xj3q51J79JBnIzR3+aHwKJvO02KseSU4LXiOsVgV
WhSYW6mm7nrtlrd+chLMeieTiJ3Ot+8zr6O50KlKut7Ond2eR3mspSq6Hp4quU4nE89EiABpnPEaS8ZY
SFS2d8n3aap12WHm0+L4QB2KapCpxn6AzFk7vLwWtlZC6vvrJJomi6xTwRObda2PMmHiYMv0/v11usxQ
n5isBojXr+Wx6NCwXepr874Bz2V21KGbqqgV4nmdetujWCac7vKipj1sV1NJE9+fI9O5Aj25KnisbKns
d9PXiFm55Ekii2GHBqsit+9qZKJTvLvVetm2+0prlduOcmx7KjlMdYe2bEOr0lSpqWdioE8b87M5oRQ5
lqqWZnOprKksUmECTUzWp2mb+T7P5LDo0Ok2m+UDU60EPZnPy3k1lIUJZxqN2dpHDM8imqc3a/fV2Oyt
7eqrKhGVjy6biJ+rY/9x443oH0j9mCnLiMMeaLR0b3Nro7dxzx5rZM96dp1Nmiyua9whjGipcQvm1w03
1HaQaHpPUOLPhY7zrMDlm922pmmCZiVQ1ZCFa2trbGzmmOUeMaH37nQ4ZFu9rVnV1xWuiz0MngWJTvBs
tLR1Z7N3d+cGtS/c2bu+fWuTej5j769sMvMw3b/Zu71Nw6BFe395nrEbb3vkjwAAAP//IzQ4mAAwAAA=
`,
	},

	"/edit.html": {
		local:   "edit.html",
		size:    983,
		modtime: 1442683616,
		compressed: `
H4sIAAAJbogA/3xT0W7bOBB8lr5iQ+BegpNoJTgcHMsGcnaAK5A2QaOiLYI80BJtERFFgVpbdoz8e5ek
47ZAkRdb3J1ZzSxH+dnibl58v7+BGnUD91/+u/0wB5Zw/vVyzvmiWMC3/4uPt5ClIyisaHuFyrSi4fzm
E4tZjdhdcT4MQzpcpsauefGZ79yszJGPjwn+wkwrrNgsjnP/xp1u2n76hznZeDwO9ACWoprFUY4KGzk7
HNJOYP36mvNQiKnV476RgPtOThnKHfKy74kbRfwc8rPH+eK6uH6Ec06VczjQb5QMcvmsMFmaXdKrF9Wu
r2BpbCWtK008RJuX9/rvtDpRVb4+moA7a2HXqn07vpJkYlf7IOWtmY06z/bttJQtWhkQzlIiGrUmlKtL
e8K5lrBSwAGAzHYm7BpUD6qtpVUoq2A8auQKT4rQdKdnq9a17wRXiEY7NUc50aAqrH3hr9/Npf9KPXEi
opVpMVkJrZr9FcxN25tG9H8Du92UqhLwQBGAgm5ncHos+9kJWEkVTQEpjf83fSdKOTkNpiVLev9Fh8G2
v9anp5n3lXN/+bM45yEoce5W6wJzliRwTFePonw2W2lXjRnS0mgu+EX2zzgbXWaQJA5dqS2UJJsiGXbs
A5SvjNWgJdammrL7u4eCgSjdiqfsFEWPpIC+3UUrNAWxJO00iIE1Aw29GDEX3mPV5/eID2zVdhs8Zrjf
LLUi5lY0Gzo+iK1kwL0e7gQ5vZwEO9fBrbNPH8ws/hEAAP//8+j85dcDAAA=
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}