package schema

type GetProteinReq struct {
	Code string `json:"code"`
}

type ProteinData struct {
	ProteinID string  `json:"protein_id"`
	Data      Protein `json:"data"`
}

type Protein struct {
	PrimaryAccession     string               `json:"primary_accession"`
	RecommendedName      string               `json:"recommended_name"`
	Organism             Organism             `json:"organism"`
	EntryAudit           EntryAudit           `json:"entry_audit"`
	Functions            []string             `json:"functions"`
	SubunitStructure     []string             `json:"subunit_structure"`
	SubcellularLocations []string             `json:"subcellular_locations"`
	DiseaseAssociations  []DiseaseAssociation `json:"disease_associations"`
	Isoforms             []Isoform            `json:"isoforms"`
	Features             []Feature            `json:"features"`
	PDBIDs               []string             `json:"pdb_ids"`
	PDBLink              string               `json:"pdb_link"`
	Sequence             string               `json:"sequence"`
	FileHash             string               `json:"file_hash"`
	IPFSCid              string               `json:"ipfs_cid"`
}

func (p *Protein) GenerateContractArgs(functionName string) []string {
	switch functionName {
	case "StoreMetadata":
	case "QueryMetadata":
	case "UpdateMetadata":
	case "MetadataExists":
	}
	return []string{}
}

type Organism struct {
	ScientificName string `json:"scientific_name"`
	CommonName     string `json:"common_name"`
}

type EntryAudit struct {
	FirstPublicDate          string `json:"first_public_date"`
	LastAnnotationUpdateDate string `json:"last_annotation_update_date"`
	SequenceVersion          int    `json:"sequence_version"`
	EntryVersion             int    `json:"entry_version"`
}

type DiseaseAssociation struct {
	DiseaseName    string `json:"disease_name"`
	Acronym        string `json:"acronym"`
	CrossReference string `json:"cross_reference"`
}

type Isoform struct {
	IsoformName    string `json:"isoform_name"`
	SequenceStatus string `json:"sequence_status"`
}

type Feature struct {
	Type        string `json:"type"`
	Location    string `json:"location"`
	Description string `json:"description"`
}
