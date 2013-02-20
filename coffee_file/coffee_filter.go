package coffee_file

import (
	"bytes"
	"github.com/ngmoco/falcore"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"os/exec"
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

func (f *Filter) FilterRequest(req *falcore.Request) (res *http.Response) {
	// Clean asset path
	asset_path := filepath.Clean(filepath.FromSlash(req.HttpRequest.URL.Path))

	if filepath.Ext(asset_path) != ".js" {
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

	coffee_asset_path := strings.Replace(asset_path, ".js", ".coffee", -1)
	if _, err := os.Stat(coffee_asset_path); err == nil {
		asset_path = coffee_asset_path
	} else {
		return
	}

	if file, err := os.Open(asset_path); err == nil {
		// Make sure it's an actual file
		if stat, err := file.Stat(); err == nil && stat.Mode()&os.ModeType == 0 {
			cmd := exec.Command("coffee", "-cp", asset_path)
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				falcore.Error("Can't compile CoffeeScript file : %v", asset_path)
				return
			}
			js := out.String()

			res = &http.Response{
				Request:       req.HttpRequest,
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString(js)),
				Header:        make(http.Header),
				ContentLength: int64(len(js)),
			}

			res.Header.Set("Content-Type", mime.TypeByExtension(".js"))

		} else {
			file.Close()
		}
	} else {
		falcore.Finest("Can't open %v: %v", asset_path, err)
	}

	return
}
