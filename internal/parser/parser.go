package parser

import (
	"encoding/binary"
	"go-open-registry/internal/log"
	"io"
)

// ReadBinary Parse binary request from cargo
func ReadBinary(r io.Reader) (json []byte, crate []byte, err error) {
	var size uint32
	// Todo: Test against bad input
	if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
		log.ErrorWithFields("read JSON size failed", log.Fields{
			"err": err,
		})
	}

	json = make([]byte, size)
	if err = binary.Read(r, binary.LittleEndian, &json); err != nil {
		log.ErrorWithFields("read JSON failed", log.Fields{
			"err": err,
		})
	}

	if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
		log.ErrorWithFields("read .crate size failed", log.Fields{
			"err": err,
		})
	}

	crate = make([]byte, size)
	if err = binary.Read(r, binary.LittleEndian, &crate); err != nil {
		log.ErrorWithFields("read .crate failed", log.Fields{
			"err": err,
		})
	}

	return json, crate, err
}
