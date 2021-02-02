// This file is part of go-pff (https://github.com/mooijtech/go-pff)
// Copyright (C) 2021 Marten Mooij (https://www.mooijtech.com/)
package pff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

// PFF represents the Personal Folder File format.
type PFF struct {
	Filepath string
	FormatType string
}

// New is a constructor for the Personal Folder File format.
func New(filePath string) PFF {
	return PFF {
		Filepath: filePath,
	}
}

// Read reads the PFF into an output buffer.
func (pff *PFF) Read(outputBufferSize int, offset int) ([]byte, error) {
	inputFile, err := os.Open(pff.Filepath)

	if err != nil {
		return nil, err
	}

	outputBuffer := make([]byte, outputBufferSize)

	_, err = inputFile.Seek(int64(offset), 0)

	if err != nil {
		return nil, err
	}

	_, err = inputFile.Read(outputBuffer)

	if err != nil {
		return nil, err
	}

	if err := inputFile.Close(); err != nil {
		return nil, err
	}

	return outputBuffer, nil
}

// GetHeader returns the file header.
//
// References "2. File header":
// The file header common to both the 64-bit and 32-bit PFF format consists of 24 bytes.
func (pff *PFF) GetHeader() ([]byte, error) {
	return pff.Read(24, 0)
}

// IsValidSignature checks if the file header contains the unique signature "!BDN".
//
// References "2. File header":
// The first 4 bytes of the file header contain the unique signature "!BDN" signifying the PFF format.
func (pff *PFF) IsValidSignature(header []byte) bool {
	return bytes.HasPrefix(header, []byte("!BDN"))
}

// Constants for identifying content types (PST, OST or PAB).
//
// References "2.1. Content types".
const (
	ContentTypePST = "PST"
	ContentTypeOST = "OST"
	ContentTypePAB = "PAB"
)

// GetContentType returns the content type which may be PST, OST or PAB.
//
// References "2. File header":
// The 9th and 10th byte contain the content type.
func (pff *PFF) GetContentType(header []byte) (string, error) {
	contentType := header[8:10]

	if bytes.Equal(contentType, []byte("SM")) {
		return ContentTypePST, nil
	} else if bytes.Equal(contentType, []byte("SO")) {
		return ContentTypeOST, nil
	} else if bytes.Equal(contentType, []byte("AB")) {
		return ContentTypePAB, nil
	} else {
		return "", errors.New("unrecognized content type")
	}
}

// Constants for identifying format types (64-bit or 32-bit).
//
// References "2.2. Format types".
const (
	FormatType32 = "32-bit"
	FormatType64 = "64-bit"
	FormatType64With4k = "64-bit-with-4k"
)

// GetFormatType returns the format type which can be either 64-bit (Unicode) or 32-bit (ANSI).
//
// References "2. File header" and "2.2. Format types":
// The 11h and 12th byte contain the format type.
func (pff *PFF) GetFormatType(header []byte) (string, error) {
	formatType := binary.LittleEndian.Uint16(header[10:12])

	if formatType == 14 || formatType  == 15 {
		return FormatType32, nil
	} else if formatType == 21 || formatType == 23 {
		return FormatType64, nil
	} else if formatType == 36 {
		return FormatType64With4k, nil
	} else {
		return "", errors.New("failed to get format type")
	}
}

// Constants for identifying encryption types.
const (
	EncryptionTypeNone = "none"
	EncryptionTypePermute = "permute"
	EncryptionTypeCyclic = "cyclic"
)

// GetEncryptionType returns the encryption type.
//
// References "2.3. The 32-bit header data", "2.4. The 64-bit header data" and "2.7. Encryption types":
// Compressible encryption (permute) is on by default with newer versions of Outlook.
func (pff *PFF) GetEncryptionType(formatType string) (string, error) {
	var encryptionType []byte
	var err error

	if formatType == FormatType64 || formatType == FormatType64With4k {
		encryptionType, err = pff.Read(1, 513)
	} else if formatType == FormatType32 {
		encryptionType, err = pff.Read(1, 461)
	} else {
		return "", errors.New("unsupported format type")
	}

	if err != nil {
		return "", err
	}

	if bytes.Equal(encryptionType, []byte{0}) {
		return EncryptionTypeNone, nil
	} else if bytes.Equal(encryptionType, []byte{1}) {
		return EncryptionTypePermute, nil
	} else if bytes.Equal(encryptionType, []byte{2}) {
		return EncryptionTypeCyclic, nil
	} else {
		return "", errors.New("unsupported encryption type")
	}
}