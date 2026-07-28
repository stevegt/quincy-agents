// Intent: Build a validated block index so assembly fails loudly on missing,
// duplicate, or ambiguous selectable module IDs. Source: DI-lorad

package module

import "fmt"

type SourceFile struct {
	ReferencePath string
	FilePath      string
}

type IndexedBlock struct {
	Reference Reference
	Source    SourceFile
	Block     Block
}

type Index struct {
	blocks map[string]IndexedBlock
}

func BuildIndex(sources []SourceFile) (*Index, error) {
	index := &Index{blocks: make(map[string]IndexedBlock)}

	for _, source := range sources {
		mod, err := Parse(source.FilePath)
		if err != nil {
			return nil, err
		}
		referencePath := NormalizeReferencePath(source.ReferencePath)
		for _, block := range mod.Blocks {
			if block.Metadata.ID == "" {
				return nil, fmt.Errorf("missing agent_module.id for heading %q in %s", block.Heading, source.FilePath)
			}
			reference := Reference{Path: referencePath, ID: block.Metadata.ID}
			key := referenceKey(reference)
			if existing, ok := index.blocks[key]; ok {
				return nil, fmt.Errorf("duplicate block reference %s in %s and %s", reference.String(), existing.Source.FilePath, source.FilePath)
			}
			index.blocks[key] = IndexedBlock{
				Reference: reference,
				Source:    source,
				Block:     block,
			}
		}
	}

	return index, nil
}

func (i *Index) SelectBlock(reference Reference) (IndexedBlock, error) {
	block, ok := i.blocks[referenceKey(reference)]
	if !ok {
		return IndexedBlock{}, fmt.Errorf("missing block reference %s", reference.String())
	}
	return block, nil
}

func referenceKey(reference Reference) string {
	return reference.Path + "#" + reference.ID
}
