package config
import "encoding/json"

type MovingAverageType int
const (
    Sma MovingAverageType = iota
    Ema
)
var maTypeName = map[MovingAverageType]string{
    Sma: "S",
    Ema: "E",
}
var maTypeFromString = map[string]MovingAverageType{
    "S": Sma,
    "E": Ema,
}

func (mat MovingAverageType) String() string {
    return maTypeName[mat]
}
func (mat MovingAverageType) MarshalJSON() ([]byte, error) {
    return json.Marshal(mat.String())
}

type PricePoint int
const (
    Open PricePoint = iota
    High
    Low
    Close
)
var pricePointName = map[PricePoint]string{
    Open:  "O",
    High:  "H",
    Low:   "L",
    Close: "C",
}
var pricePointFromString = map[string]PricePoint{
    "O": Open,
    "H": High,
    "L": Low,
    "C": Close,
}

func (pp PricePoint) String() string {
    return pricePointName[pp]
}
func (pp PricePoint) MarshalJSON() ([]byte, error) {
    return json.Marshal(pp.String())
}

type MovingAverageKey struct {
    Type       MovingAverageType
    Period     int
    PricePoint PricePoint
}

// package config

// import "encoding/json"

// type MovingAverageType int

// const (
// 	Sma MovingAverageType = iota
// 	Ema
// )

// var maTypeName = map[MovingAverageType]string{
// 	Sma: "S",
// 	Ema: "E",
// }

// func (mat MovingAverageType) String() string {
// 	return maTypeName[mat]
// }

// func (mat MovingAverageType) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(mat.String())
// }

// type PricePoint int

// const (
// 	Open PricePoint = iota
// 	High
// 	Low
// 	Close
// )

// var pricePointName = map[PricePoint]string{
// 	Open:  "O",
// 	High:  "H",
// 	Low:   "L",
// 	Close: "C",
// }

// func (pp PricePoint) String() string {
// 	return pricePointName[pp]
// }

// func (pp PricePoint) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(pp.String())
// }

// type MovingAverageKey struct {
// 	Type       MovingAverageType
// 	Period     int
// 	PricePoint PricePoint
// }

// package config

// type MovingAverageType int

// const (
// 	Sma MovingAverageType = iota
// 	Ema
// )

// var maTypeName = map[MovingAverageType]string{
// 	Sma: "SMA",
// 	Ema: "EMA",
// }

// func (mat MovingAverageType) String() string {
// 	return maTypeName[mat]
// }

// type PricePoint int

// const (
// 	Open PricePoint = iota
// 	High
// 	Low
// 	Close
// )

// var pricePointName = map[PricePoint]string{
// 	Open: "Open",
// 	High: "High",
// 	Low: "Low",
// 	Close: "Close",
// }

// func (ppn PricePoint) String() string {
// 	return pricePointName[ppn]
// }

// type MovingAverageKey struct {
// 	Type MovingAverageType
// 	Period int
// 	PricePoint PricePoint
// }
