package pdfutil

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io"
)

func MergePdfByBytes(pdfBytesList [][]byte) (pdfByte []byte, err error) {
	allReader := make([]io.ReadSeeker, 0, len(pdfBytesList))
	for i := 0; i < len(pdfBytesList); i++ {
		allReader = append(allReader, bytes.NewReader(pdfBytesList[i]))
	}
	var b bytes.Buffer
	bufWriter := bufio.NewWriter(&b)
	err = api.Merge(allReader, bufWriter, nil)
	if err != nil {
		err = errors.New("merge pdf fail, final stage: " + err.Error())
		return
	}
	pdfByte = b.Bytes()
	return
}
