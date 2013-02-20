package sass_file

import (
	"bytes"
	"github.com/ngmoco/falcore"
	"github.com/suapapa/go_sass"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// A falcore RequestFilter for serving static files
// from the filesystem.
type Filter struct {
	// File system base path for serving files
	BasePath string
	// Prefix in URL path
	PathPrefix string
}

var sc sass.Compiler

func (f *Filter) FilterRequest(req *falcore.Request) (res *http.Response) {
	// Clean asset path
	asset_path := filepath.Clean(filepath.FromSlash(req.HttpRequest.URL.Path))

	if filepath.Ext(asset_path) != ".css" {
		return
	}

	// Resolve PathPrefix
	if strings.HasPrefix(asset_path, f.PathPrefix) {
		asset_path = asset_path[len(f.PathPrefix):]
	} else {
		falcore.Debug("%v doesn't match prefix %v", asset_path, f.PathPrefix)
		res = falcore.SimpleResponse(req.HttpRequest, 404, nil, "Not found.")
		return
	}

	// Resolve FSBase
	if f.BasePath != "" {
		asset_path = filepath.Join(f.BasePath, asset_path)
	} else {
		falcore.Error("file_filter requires a BasePath")
		return falcore.SimpleResponse(req.HttpRequest, 500, nil, "Server Error\n")
	}

	scss_asset_path := strings.Replace(asset_path, ".css", ".scss", -1)
	sass_asset_path := strings.Replace(asset_path, ".css", ".sass", -1)
	if _, err := os.Stat(scss_asset_path); err == nil {
		asset_path = scss_asset_path
	} else if _, err := os.Stat(sass_asset_path); err == nil {
		asset_path = sass_asset_path
	} else {
		return
	}

	if file, err := os.Open(asset_path); err == nil {
		// Make sure it's an actual file
		if stat, err := file.Stat(); err == nil && stat.Mode()&os.ModeType == 0 {
			css, err := sc.CompileFile(asset_path)
			if err != nil {
				falcore.Error("%v", err)
				return
			}

			res = &http.Response{
				Request:       req.HttpRequest,
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString(css)),
				Header:        make(http.Header),
				ContentLength: int64(len(css)),
			}

			res.Header.Set("Content-Type", mime.TypeByExtension(".css"))

		} else {
			file.Close()
		}
	} else {
		falcore.Finest("Can't open %v: %v", asset_path, err)
	}

	return
}
