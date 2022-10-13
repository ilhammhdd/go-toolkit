package sqlkit

const stackSize = 1024 * 8

func GeneratePlaceHolder(n uint16) []byte {
	if n == 0 {
		return nil
	} else if n == 1 {
		return []byte{'(', '?', ')'}
	}
	var placeholder []byte = []byte{'('}
	var i uint16
	for i = 0; i < n; i++ {
		placeholder = append(placeholder, '?')
		if i != n-1 {
			placeholder = append(placeholder, ',')
		}
	}
	placeholder = append(placeholder, ')')
	return placeholder
}

func GenerateNPlaceHolder(nPlaceHolders, nParams uint16) []byte {
	if nParams == 0 {
		return nil
	} else if nPlaceHolders == 1 && nParams == 1 {
		return []byte{'(', '?', ')'}
	}

	var placeholders []byte

	var i uint16
	for i = 0; i < nPlaceHolders; i++ {
		placeholders = append(placeholders, GeneratePlaceHolder(nParams)...)
		if i != nPlaceHolders-1 {
			placeholders = append(placeholders, ',')
		}
	}
	return placeholders
}
