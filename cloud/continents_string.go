// generated by stringer -type=Continent; DO NOT EDIT

package cloud

import "fmt"

const _Continent_name = "AfricaAntarcticaAsiaEuropeNorthAmericaOceaniaSouthAmerica"

var _Continent_index = [...]uint8{0, 6, 16, 20, 26, 38, 45, 57}

func (i Continent) String() string {
	if i >= Continent(len(_Continent_index)-1) {
		return fmt.Sprintf("Continent(%d)", i)
	}
	return _Continent_name[_Continent_index[i]:_Continent_index[i+1]]
}
