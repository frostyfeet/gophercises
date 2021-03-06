package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request){
		if val, ok := pathsToUrls[r.URL.Path]; ok {
				http.Redirect(w, r, val, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	} 
	return http.HandlerFunc(fn)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	m := make (map[string]string)
	ptu, err := parseYaml(yml)
	if err != nil {
		return nil, err 
	} else {
		for _, val := range ptu {
			m[val.Path] = val.URL
		}
		return MapHandler(m, fallback), nil
	}
}

func parseYaml(yml []byte) ([]PathToURL, error) {

	var ptu []PathToURL
	err := yaml.Unmarshal(yml, &ptu)
	if err != nil {
		return nil, err
	}
	return ptu, nil
}

type PathToURL struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}
