package k8s

import (
	"bytes"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func Manifests(r io.ReadCloser) ([]runtime.Object, error) {
	defer r.Close()

	dd, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	fn := scheme.Codecs.UniversalDeserializer().Decode

	bb := bytes.Split(dd, []byte("---\n"))
	mm := make([]runtime.Object, 0, len(bb))
	for _, b := range bb {
		m, _, err := fn(b, nil, nil)
		if err != nil {
			return nil, err
		}
		mm = append(mm, m)
	}
	return mm, nil
}
