package pff

import log "github.com/sirupsen/logrus"

// Parser represents a parser for PST files.
type Parser struct {}

// NewParser is a constructor for creating parsers.
func NewParser() Parser {
	return Parser{}
}

// Parse parses the given PST file.
func (parser *Parser) Parse(inputFile string) {
	pst := New(inputFile)

	log.Infof("Using Personal Folder File: %s", pst.Filepath)

	header, err := pst.GetHeader()

	if err != nil {
		log.Fatalf("Failed to get PFF header: %s", err)
	}

	if !pst.IsValidSignature(header) {
		log.Fatalf("Invalid Personal Folder File.")
	}

	contentType, err := pst.GetContentType(header)

	if err != nil {
		log.Errorf("Failed to get content type: %s", err)
	}

	log.Infof("Detected content type: %s...", contentType)

	formatType, err := pst.GetFormatType(header)

	if err != nil {
		log.Errorf("Failed to get format type: %s", formatType)
	}

	log.Infof("Detected format type: %s...", formatType)

	encryptionType, err := pst.GetEncryptionType(formatType)

	if err != nil {
		log.Errorf("Failed to get encryption type: %s", err)
	}

	log.Infof("Detected encryption type: %s...", encryptionType)

	nodeBTree, err := pst.GetNodeBTree(formatType)

	if err != nil {
		log.Errorf("Failed to get node b-tree: %s", err)
	}

	log.Infof("Node b-tree offset: %d", nodeBTree.StartOffset)

	nodeBTreeEntry, err := pst.FindBTreeNode(formatType, nodeBTree, 97)

	if err != nil {
		log.Errorf("Failed to find b-tree node entry: %s", err)
	}

	log.Debugf("Found node b-tree entry: %d", nodeBTreeEntry.Identifier)
}