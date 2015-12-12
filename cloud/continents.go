//go:generate stringer -type=Continent
package cloud

//Continent represent  code For Continents
type Continent uint8

/*
AF	Africa
AN	Antarctica
AS	Asia
EU	Europe
NA	North america
OC	Oceania
SA	South america
*/
var isoCodeToContinente = map[string]Continent{
	"AF": Africa,
	"AN": Antarctica,
	"AS": Asia,
	"EU": Europe,
	"NA": NorthAmerica,
	"OC": Oceania,
	"SA": SouthAmerica,
}


const (
//Africa ...
	Africa Continent = iota + 1
// Antarctica ...
	Antarctica
// Asia ...
	Asia
// Europe ...
	Europe
// NorthAmerica ...
	NorthAmerica
// Oceania ...
	Oceania
// SouthAmerica ...
	SouthAmerica
//Undefined ...
	Undefined


	defaultContinent = Europe
)
