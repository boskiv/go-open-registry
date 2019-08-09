package parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadBinary(r io.Reader) (json []byte, crate []byte, err error) {
	var size uint32
	if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, nil, fmt.Errorf("read JSON size failed: %s", err)
	}

	json = make([]byte, size)
	if err = binary.Read(r, binary.LittleEndian, &json); err != nil {
		return nil, nil, fmt.Errorf("read JSON failed: %s", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, nil, fmt.Errorf("read .crate size failed: %s", err)
	}

	crate = make([]byte, size)
	if err = binary.Read(r, binary.LittleEndian, &crate); err != nil {
		return nil, nil, fmt.Errorf("read .crate failed: %s", err)
	}

	return json, crate, nil

}
