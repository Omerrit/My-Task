package utils

import (
	"bytes"
	"fmt"
	"gerrit-share.lan/go/errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
)

func ParseMultiForm(request *http.Request) (values url.Values, file *multipart.Part, err error) {
	contentTypeValue := request.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(contentTypeValue)
	if err != nil {
		return nil, nil, fmt.Errorf("incorrect multipart header: %s", err.Error())
	}
	boundary, ok := params["boundary"]
	if !ok {
		return nil, nil, fmt.Errorf("boundary is missed")
	}

	values = make(url.Values, 0)
	reader := multipart.NewReader(request.Body, boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil, fmt.Errorf("file should be the last part of multipart form")
			}
			return nil, nil, fmt.Errorf("failed to parse parts: %s", err.Error())
		}
		name := part.FormName()
		if len(part.FileName()) == 0 {
			if len(name) == 0 {
				continue
			}
			buffer := bytes.Buffer{}
			_, err := buffer.ReadFrom(part)
			if err != nil {
				return nil, nil, err
			}
			values.Add(part.FormName(), string(buffer.Bytes()))
			continue
		}
		file = part
		break
	}
	return values, file, nil
}
