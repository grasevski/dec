package dec

import (
	"database/sql/driver"
	"github.com/shopspring/decimal"
)

type Decimal int64

const factor, digits, base = 1000000000, 10, 10

var digitsDecimal = decimal.New(1, digits)

func Avg(first Decimal, rest ...Decimal) Decimal { return Sum(first, rest...) / Decimal(len(rest)+1) }

func Max(first Decimal, rest ...Decimal) Decimal {
	for _, x := range rest {
		if x > first {
			first = x
		}
	}
	return first
}

func Min(first Decimal, rest ...Decimal) Decimal {
	for _, x := range rest {
		if x < first {
			first = x
		}
	}
	return first
}

func New(intPart int64, decimalPart int32) Decimal {
	return Decimal(intPart*factor + int64(decimalPart))
}

func NewFromFloat(value float64) Decimal   { return Decimal(value * factor) }
func NewFromFloat32(value float32) Decimal { return NewFromFloat(float64(value)) }

func NewFromString(value string) (Decimal, error) {
	var d, err = decimal.NewFromString(value)
	return NewFromDecimal(d), err
}

func RequireFromString(value string) Decimal { return NewFromDecimal(decimal.RequireFromString(value)) }

func Sum(first Decimal, rest ...Decimal) (sum Decimal) {
	sum = first
	for _, x := range rest {
		sum += x
	}
	return
}

func (d Decimal) Abs() Decimal {
	if d < 0 {
		return -d
	}
	return d
}

func (d Decimal) Cmp(d2 Decimal) int {
	if d < d2 {
		return -1
	} else if d > d2 {
		return 1
	}
	return 0
}

func (d Decimal) Div(d2 Decimal) Decimal                      { return d * factor / d2 }
func (d Decimal) DivRound(d2 Decimal, precision int8) Decimal { return d.Div(d2).Round(precision) }

func (d Decimal) Float64() (f float64, exact bool) {
	f = float64(d) / factor
	exact = NewFromFloat(f) == d
	return
}

func (d Decimal) Floor() Decimal {
	if d < 0 && d%factor != 0 {
		d -= factor
	}
	return (d / factor) * factor
}

func (d *Decimal) GobDecode(data []byte) (err error) {
	var tmp = d.Decimal()
	err = tmp.GobDecode(data)
	*d = NewFromDecimal(tmp)
	return
}

func (d Decimal) GobEncode() ([]byte, error)     { return d.Decimal().GobEncode() }
func (d Decimal) IntPart() int64                 { return int64(d / factor) }
func (d Decimal) MarshalBinary() ([]byte, error) { return d.Decimal().MarshalBinary() }
func (d Decimal) MarshalJSON() ([]byte, error)   { return d.Decimal().MarshalJSON() }
func (d Decimal) MarshalText() ([]byte, error)   { return d.Decimal().MarshalText() }
func (d Decimal) Mul(d2 Decimal) Decimal         { return d * d2 / factor }
func (d Decimal) Round(places int8) Decimal      { return NewFromDecimal(d.Decimal().Round(int32(places))) }

func (d Decimal) RoundBank(places int8) Decimal {
	return NewFromDecimal(d.Decimal().RoundBank(int32(places)))
}
func (d Decimal) RoundCash(interval uint8) Decimal {
	return NewFromDecimal(d.Decimal().RoundCash(interval))
}

func (d *Decimal) Scan(value interface{}) (err error) {
	var tmp = d.Decimal()
	err = tmp.Scan(value)
	*d = NewFromDecimal(tmp)
	return
}

func (d Decimal) Shift(shift int8) Decimal {
	if shift >= 0 {
		for i := int8(0); i < shift; i++ {
			d *= base
		}
	} else {
		for i := int8(0); i < -shift; i++ {
			d /= base
		}
	}
	return d
}

func (d Decimal) Sign() int {
	if d > 0 {
		return 1
	} else if d < 0 {
		return -1
	}
	return 0
}

func (d Decimal) String() string                 { return d.Decimal().String() }
func (d Decimal) StringFixed(places int8) string { return d.Decimal().StringFixed(int32(places)) }

func (d Decimal) StringFixedBank(places int8) string {
	return d.Decimal().StringFixedBank(int32(places))
}

func (d Decimal) StringFixedCash(interval uint8) string { return d.Decimal().StringFixedCash(interval) }

func (d Decimal) Truncate(precision int8) Decimal {
	var p Decimal = factor
	for i := int8(0); i < precision; i++ {
		p /= base
	}
	return p * (d / p)
}

func (d *Decimal) UnmarshalBinary(data []byte) (err error) {
	var tmp = d.Decimal()
	err = tmp.UnmarshalBinary(data)
	*d = NewFromDecimal(tmp)
	return
}

func (d *Decimal) UnmarshalJSON(decimalBytes []byte) (err error) {
	var tmp = d.Decimal()
	err = tmp.UnmarshalJSON(decimalBytes)
	*d = NewFromDecimal(tmp)
	return
}

func (d *Decimal) UnmarshalText(text []byte) (err error) {
	var tmp = d.Decimal()
	err = tmp.UnmarshalText(text)
	*d = NewFromDecimal(tmp)
	return
}

func (d Decimal) Value() (driver.Value, error) { return d.Decimal().Value() }

type NullDecimal struct {
	Decimal
	Valid bool
}

func (d NullDecimal) MarshalJSON() ([]byte, error) { return d.NullDecimal().MarshalJSON() }

func (d *NullDecimal) Scan(value interface{}) (err error) {
	var tmp = d.NullDecimal()
	err = tmp.Scan(value)
	*d = NewFromNullDecimal(tmp)
	return
}

func (d *NullDecimal) UnmarshalJSON(decimalBytes []byte) (err error) {
	var tmp = d.NullDecimal()
	err = tmp.UnmarshalJSON(decimalBytes)
	*d = NewFromNullDecimal(tmp)
	return
}

func (d NullDecimal) Value() (driver.Value, error) { return d.NullDecimal().Value() }

func (d Decimal) Decimal() decimal.Decimal     { return decimal.New(int64(d), -digits) }
func NewFromDecimal(d decimal.Decimal) Decimal { return Decimal(d.Mul(digitsDecimal).IntPart()) }

func (d NullDecimal) NullDecimal() decimal.NullDecimal {
	return decimal.NullDecimal{Decimal: d.Decimal.Decimal(), Valid: d.Valid}
}

func NewFromNullDecimal(d decimal.NullDecimal) NullDecimal {
	return NullDecimal{NewFromDecimal(d.Decimal), d.Valid}
}
