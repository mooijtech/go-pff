// This file is part of go-pff (https://github.com/mooijtech/go-pff)
// Copyright (C) 2021 Marten Mooij (https://www.mooijtech.com/)
package pff

import (
	"encoding/binary"
	"errors"
)

// BTreeNode represents a branch- or leaf node in the b-tree.
type BTreeNode struct {
	StartOffset int
}

// NewBTreeNode is a constructor for b-tree nodes.
func NewBTreeNode(btreeNodeStartOffset int) BTreeNode {
	return BTreeNode{
		StartOffset: btreeNodeStartOffset,
	}
}

// BTreeNodeEntry represents a node entry.
type BTreeNodeEntry struct {
	Identifier      int
	Data            []byte
}

// NewBTreeNodeEntry is a constructor for b-tree node entries.
func NewBTreeNodeEntry(identifier int, data []byte) BTreeNodeEntry {
	return BTreeNodeEntry {
		Identifier:      identifier,
		Data:            data,
	}
}

// GetNodeBTree returns the Node B-Tree (NBT).
//
// References "2.3. The 32-bit header data" and "2.4. The 64-bit header data", "5. The index b-tree":
// An index b-tree consists of:
// - branch nodes that point to branch or leaf nodes
// - leaf nodes that contain the index data
func (pff *PFF) GetNodeBTree(formatType string) (BTreeNode, error) {
	var btreeStartOffset int

	if formatType == FormatType64 || formatType == FormatType64With4k {
		offset, err := pff.Read(8, 224)

		if err != nil {
			return BTreeNode{}, err
		}

		btreeStartOffset = int(binary.LittleEndian.Uint64(offset))
	} else if formatType == FormatType32 {
		offset, err := pff.Read(4, 188)

		if err != nil {
			return BTreeNode{}, err
		}

		btreeStartOffset = int(binary.LittleEndian.Uint32(offset))
	} else {
		return BTreeNode{}, errors.New("unsupported format type")
	}

	return NewBTreeNode(btreeStartOffset), nil
}

// GetBlockBTree returns the Block B-Tree (BBT).
//
// References "2.3. The 32-bit header data" and "2.4. The 64-bit header data", "5. The index b-tree":
// An index b-tree consists of:
// - branch nodes that point to branch or leaf nodes
// - leaf nodes that contain the index data
func (pff *PFF) GetBlockBTree(formatType string) (BTreeNode, error) {
	var btreeStartOffset int

	if formatType == FormatType64 || formatType == FormatType64With4k {
		offset, err := pff.Read(8, 240)

		if err != nil {
			return BTreeNode{}, err
		}

		btreeStartOffset = int(binary.LittleEndian.Uint64(offset))
	} else if formatType == FormatType32 {
		offset, err := pff.Read(4, 196)

		if err != nil {
			return BTreeNode{}, err
		}

		btreeStartOffset = int(binary.LittleEndian.Uint32(offset))
	} else {
		return BTreeNode{}, errors.New("unsupported format type")
	}

	return NewBTreeNode(btreeStartOffset), nil
}

// GetBTreeNodeEntryCount returns the amount of entries in this node.
//
// References "5. The index b-tree".
func (pff *PFF) GetBTreeNodeEntryCount(formatType string, btreeNode BTreeNode) (int, error) {
	if formatType == FormatType64 {
		entryCount, err := pff.Read(1, btreeNode.StartOffset + 488)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{entryCount[0], 0})), nil
	} else if formatType == FormatType64With4k {
		entryCount, err := pff.Read(2, btreeNode.StartOffset + 4056)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16(entryCount)), nil
	} else if formatType == FormatType32 {
		entryCount, err := pff.Read(1, btreeNode.StartOffset + 496)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{entryCount[0], 0})), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

// GetBTreeNodeMaxEntryCount returns the maximum amount of entries in this node.
//
// References "5. The index b-tree".
func (pff *PFF) GetBTreeNodeMaxEntryCount(formatType string, btreeNode BTreeNode) (int, error) {
	if formatType == FormatType64 {
		maxEntryCount, err := pff.Read(1, btreeNode.StartOffset + 489)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{maxEntryCount[0], 0})), nil
	} else if formatType == FormatType64With4k {
		maxEntryCount, err := pff.Read(2, btreeNode.StartOffset + 4058)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16(maxEntryCount)), nil
	} else if formatType == FormatType32 {
		maxEntryCount, err := pff.Read(1, btreeNode.StartOffset + 497)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{maxEntryCount[0], 0})), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

// GetBTreeNodeEntrySize returns the entry size of a node entry.
//
// References "5. The index b-tree":
func (pff *PFF) GetBTreeNodeEntrySize(formatType string, btreeNode BTreeNode) (int, error) {
	if formatType == FormatType64 {
		entrySize, err := pff.Read(1, btreeNode.StartOffset+490)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{entrySize[0], 0})), nil
	} else if formatType == FormatType64With4k {
		entrySize, err := pff.Read(1, btreeNode.StartOffset+4060)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{entrySize[0], 0})), nil
	} else if formatType == FormatType32 {
		entrySize, err := pff.Read(2, btreeNode.StartOffset+498)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16(entrySize)), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

// GetBTreeNodeLevel returns a zero value representing a leaf node or a value greater than zero representing branch nodes.
//
// References "5. The index b-tree"
func (pff *PFF) GetBTreeNodeLevel(formatType string, btreeNode BTreeNode) (int, error) {
	if formatType == FormatType64 {
		nodeLevel, err := pff.Read(1, btreeNode.StartOffset+491)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{nodeLevel[0], 0})), nil
	} else if formatType == FormatType64With4k {
		nodeLevel, err := pff.Read(1, btreeNode.StartOffset+4061)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{nodeLevel[0], 0})), nil
	} else if formatType == FormatType32 {
		nodeLevel, err := pff.Read(1, btreeNode.StartOffset+499)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{nodeLevel[0], 0})), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

// GetBTreeNodePageType returns the page type.
//
// References "5. The index b-tree", "3.4. Page types".
func (pff *PFF) GetBTreeNodePageType(formatType string, btreeNode BTreeNode) (int, error) {
	if formatType == FormatType64 {
		pageType, err := pff.Read(1, btreeNode.StartOffset+496)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{pageType[0], 0})), nil
	} else if formatType == FormatType64With4k {
		pageType, err := pff.Read(1, btreeNode.StartOffset+4072)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{pageType[0], 0})), nil
	} else if formatType == FormatType32 {
		pageType, err := pff.Read(1, btreeNode.StartOffset+500)

		if err != nil {
			return -1, err
		}

		return int(binary.LittleEndian.Uint16([]byte{pageType[0], 0})), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

// GetBTreeBranchNodeEntryOffset returns the offset of the b-tree node entry.
//
// References "5.1. The 32-bit index b-tree node", "5.2. The 64-bit index b-tree node"
func (pff *PFF) GetBTreeBranchNodeEntryOffset(formatType string, nodeEntry []byte) (int, error) {
	if formatType == FormatType64 {
		return int(binary.LittleEndian.Uint64(nodeEntry[16:24])), nil
	} else if formatType == FormatType64With4k {
		return -1, nil
	} else if formatType == FormatType32 {
		return -1, nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

func (btreeNodeEntry *BTreeNodeEntry) GetLocalDescriptorsOffset(formatType string) (int, error) {
	if formatType == FormatType64 {
		return int(binary.LittleEndian.Uint64(btreeNodeEntry.Data[16:24])), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

func (btreeNodeEntry *BTreeNodeEntry) GetDataOffset(formatType string) (int, error) {
	if formatType == FormatType64 {
		return int(binary.LittleEndian.Uint64(btreeNodeEntry.Data[8:16])), nil
	} else {
		return -1, errors.New("unsupported format type")
	}
}

// GetBTreeNodeEntries returns an array of b-tree nodes for a given b-tree node.
//
// References "5. The index b-tree".
func (pff *PFF) GetBTreeNodeEntries(formatType string, btreeNode BTreeNode) ([]BTreeNodeEntry, error) {
	var nodeEntries []byte
	var err error

	if formatType == FormatType64 {
		nodeEntries, err = pff.Read(488, btreeNode.StartOffset)
	} else if formatType == FormatType64With4k {
		nodeEntries, err = pff.Read(4056, btreeNode.StartOffset)
	} else if formatType == FormatType32 {
		nodeEntries, err = pff.Read(496, btreeNode.StartOffset)
	} else {
		return nil, errors.New("unsupported format type")
	}

	if err != nil {
		return []BTreeNodeEntry{}, err
	}

	nodeEntryCount, err := pff.GetBTreeNodeEntryCount(formatType, btreeNode)

	if err != nil {
		return []BTreeNodeEntry{}, err
	}

	nodeEntrySize, err := pff.GetBTreeNodeEntrySize(formatType, btreeNode)

	if err != nil {
		return []BTreeNodeEntry{}, err
	}

	nodeLevel, err := pff.GetBTreeNodeLevel(formatType, btreeNode)

	if err != nil {
		return []BTreeNodeEntry{}, err
	}

	_, err = pff.GetBTreeNodePageType(formatType, btreeNode)

	if err != nil {
		return []BTreeNodeEntry{}, err
	}

	// Node entries
	// (number of entries x entry size)
	entries := make([]BTreeNodeEntry, nodeEntryCount)

	for i := 0; i < nodeEntryCount; i++ {
		nodeEntry := nodeEntries[(i * nodeEntrySize) : (i * nodeEntrySize) + nodeEntrySize]

		if nodeLevel > 0 {
			// Branch node
			identifier := binary.LittleEndian.Uint32(nodeEntry[:8])

			entries[i] = NewBTreeNodeEntry(int(identifier), nodeEntry)
		} else {
			// Leaf node
			identifier := binary.LittleEndian.Uint32(nodeEntry[:8])

			entries[i] = NewBTreeNodeEntry(int(identifier), nodeEntry)
		}
	}

	return entries, nil
}

// FindBTreeNode walks the b-tree and finds the node with the given identifier.
func (pff *PFF) FindBTreeNode(formatType string, btreeNode BTreeNode, identifier int) (BTreeNodeEntry, error) {
	btreeNodeEntries, err := pff.GetBTreeNodeEntries(formatType, btreeNode)

	if err != nil {
		return BTreeNodeEntry{}, err
	}

	btreeNodeLevel, err := pff.GetBTreeNodeLevel(formatType, btreeNode)

	if err != nil {
		return BTreeNodeEntry{}, err
	}

	if btreeNodeLevel > 0 {
		// Branch node entries
		// Branch node entries point to other branch nodes.

		for i := 0; i < len(btreeNodeEntries); i++ {
			btreeNodeEntry := btreeNodeEntries[i]

			if btreeNodeEntry.Identifier == identifier {
				return btreeNodeEntry, nil
			}

			btreeNodeEntryOffset, err := pff.GetBTreeBranchNodeEntryOffset(formatType, btreeNodeEntry.Data)

			if err != nil {
				return BTreeNodeEntry{}, err
			}

			btreeNodeEntryNode := NewBTreeNode(btreeNodeEntryOffset)

			// Recursively walk through the branch node entries.
			btreeNodeEntry, err = pff.FindBTreeNode(formatType, btreeNodeEntryNode, identifier)

			if err != nil {
				return BTreeNodeEntry{}, nil
			}

			if btreeNodeEntry.Identifier == identifier {
				return btreeNodeEntry, nil
			}
		}
	} else {
		// Leaf node entries
		// Leaf node entries point to data and the local descriptors.

		for i := 0; i < len(btreeNodeEntries); i++ {
			btreeNodeEntry := btreeNodeEntries[i]

			if btreeNodeEntry.Identifier == identifier {
				return btreeNodeEntry, nil
			}
		}
	}

	return BTreeNodeEntry{}, nil
}