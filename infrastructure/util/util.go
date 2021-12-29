package util

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func CreateForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				logrus.Error("failed to create form")
				return "", nil, err
			}
		} else {
			_ = mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}
