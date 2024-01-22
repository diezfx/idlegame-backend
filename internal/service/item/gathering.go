package item

// what kind of wood
// different tree types have
// level requirements, durations, exp gains
type TreeType string

const (
	SpruceType TreeType = "Spruce"
	BirchType  TreeType = "Birch"
	PineType   TreeType = "Pine"
)

type OreType string

const (
	StoneOreType   OreType = "Stone"
	CopperOreType  OreType = "Copper"
	IronOreType    OreType = "Iron"
	GoldOreType    OreType = "Gold"
	DiamondOreType OreType = "Diamond"
)

type CropType string

const (
	// low quality crop
	WheatCropType  CropType = "Wheat"
	CarrotCropType CropType = "Carrot"
	PotatoCropType CropType = "Potato"
)

func (t TreeType) String() string {
	return string(t)
}

func (t OreType) String() string {
	return string(t)
}

func (t CropType) String() string {
	return string(t)
}

type Item string

const (
	StoneBarType   Item = "StoneBar"
	CopperBarType  Item = "CopperBar"
	IronBarType    Item = "IronBar"
	GoldBarType    Item = "GoldBar"
	DiamondBarType Item = "DiamondBar"
)

const (
	SpruceBlankType Item = "SpruceLog"
	BirchLogType    Item = "BirchLog"
	PineLogType     Item = "PineLog"
)
