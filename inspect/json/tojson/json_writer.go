package tojson

import "strconv"
import "encoding/base64"

const (
	objectStart       = '{'
	objectEnd         = '}'
	arrayStart        = '['
	arrayEnd          = ']'
	delimiter         = ','
	keyDelimiter      = ':'
	nullString        = "null"
	trueString        = "true"
	falseString       = "false"
	possibleFPFormats = "eEfgG"
	defaultFPFormat   = 'g'
)

var (
	fpFormats = func() [256]byte {
		var formats [256]byte
		for i := range formats {
			formats[i] = defaultFPFormat
		}
		for _, format := range ([]byte)(possibleFPFormats) {
			formats[int(format)] = format
		}
		return formats
	}()

	nullSlice  = []byte(nullString)
	trueSlice  = []byte(trueString)
	falseSlice = []byte(falseString)
)

type jsonWriter struct {
	output          []byte
	hadFirstElement bool
}

func (j *jsonWriter) Clear() {
	j.output = j.output[:0]
	j.hadFirstElement = false
}

func (j *jsonWriter) appendDelimiter() {
	if j.hadFirstElement {
		j.output = append(j.output, delimiter)
	} else {
		j.hadFirstElement = true
	}
}

func (j *jsonWriter) startObject() {
	j.output = append(j.output, objectStart)
	j.hadFirstElement = false
}

func (j *jsonWriter) endObject() {
	j.output = append(j.output, objectEnd)
	j.hadFirstElement = true
}

func (j *jsonWriter) startArray() {
	j.output = append(j.output, arrayStart)
	j.hadFirstElement = false
}

func (j *jsonWriter) endArray() {
	j.output = append(j.output, arrayEnd)
	j.hadFirstElement = true
}

func (j *jsonWriter) appendKey(key string) {
	j.appendString(key)
	j.output = append(j.output, keyDelimiter)
}

func (j *jsonWriter) appendString(value string) {
	j.output = strconv.AppendQuote(j.output, value)
}

func (j *jsonWriter) appendBool(value bool) {
	if value {
		j.output = append(j.output, trueSlice...)
	} else {
		j.output = append(j.output, falseSlice...)
	}
}

func (j *jsonWriter) appendInt(value int64) {
	j.output = strconv.AppendInt(j.output, value, 10)
}

func (j *jsonWriter) appendFloat32(value float32, format byte, precision int) {
	j.output = strconv.AppendFloat(j.output, float64(value), fpFormats[format], precision, 32)
}

func (j *jsonWriter) appendFloat64(value float64, format byte, precision int) {
	j.output = strconv.AppendFloat(j.output, value, fpFormats[format], precision, 64)
}

func (j *jsonWriter) appendBytes(value []byte) {
	j.output = strconv.AppendQuote(j.output, base64.RawURLEncoding.EncodeToString(value))
}

func (j *jsonWriter) appendNull() {
	j.output = append(j.output, nullSlice...)
}

func (j *jsonWriter) Output() []byte {
	return j.output
}

//to append raw value call appendDelimiter if applicable and then append directly to j.output
