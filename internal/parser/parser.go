package parser

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/sirupsen/logrus" //nolint:depguard
)

// ReadBinary Parse binary request from cargo
func ReadBinary(r io.Reader) (json []byte, crate []byte, err error) {
	var size uint32
	// Todo: Test against bad input
	if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
		logrus.Error(fmt.Errorf("read JSON size failed: %s", err))
	}

	json = make([]byte, size)
	if err = binary.Read(r, binary.LittleEndian, &json); err != nil {
		logrus.Error(fmt.Errorf("read JSON failed: %s", err))
	}

	if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
		logrus.Error(fmt.Errorf("read .crate size failed: %s", err))
	}

	crate = make([]byte, size)
	if err = binary.Read(r, binary.LittleEndian, &crate); err != nil {
		logrus.Error(fmt.Errorf("read .crate failed: %s", err))
	}

	return json, crate, nil
}
