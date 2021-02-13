package pff

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

type Folder struct {
	BTreeNodeEntry BTreeNodeEntry
}

func NewFolder(btreeNodeEntry BTreeNodeEntry) Folder {
	return Folder {
		BTreeNodeEntry: btreeNodeEntry,
	}
}

const (
	NodeBTreeIdentifierRootFolder = 290
)

// GetRootFolder returns the root folder.
func (pff *PFF) GetRootFolder(formatType string) (Folder, error) {
	nodeBTree, err := pff.GetNodeBTree(formatType)

	if err != nil {
		return Folder{}, err
	}

	rootFolderNode, err := pff.FindBTreeNode(formatType, nodeBTree, NodeBTreeIdentifierRootFolder)

	return NewFolder(rootFolderNode), nil
}

func (pff *PFF) GetSubFolders(formatType string, folder Folder) error {
	subFoldersIdentifier := folder.BTreeNodeEntry.Identifier + 11

	nodeBTree, err := pff.GetNodeBTree(formatType)

	if err != nil {
		return err
	}

	subFoldersNode, err := pff.FindBTreeNode(formatType, nodeBTree, subFoldersIdentifier)

	if err != nil {
		return err
	}

	subFoldersNodeDataIdentifier, err := subFoldersNode.GetDataIdentifier(formatType)

	if err != nil {
		return err
	}

	blockBTree, err := pff.GetBlockBTree(formatType)

	if err != nil {
		return err
	}

	subFoldersDataNode, err := pff.FindBTreeNode(formatType, blockBTree, subFoldersNodeDataIdentifier)

	if err != nil {
		return err
	}

	subFoldersDataNodeFileOffset, err := subFoldersDataNode.GetFileOffset(formatType)

	log.Debugf("Related sub folders identifier: %d", subFoldersIdentifier)
	log.Debugf("Offset: %d", subFoldersDataNodeFileOffset)

	n, err := pff.Read(1, subFoldersDataNodeFileOffset + 2)

	log.Debugf("It's: %d", binary.LittleEndian.Uint16([]byte{n[0], 0}))

	return nil
}