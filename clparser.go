package clparser

import (
	"strings"
)

type BasicError struct {
	What string
}

func (err *BasicError) Error() string {
	return err.What
}

// -t 30 -i -g something
// -i -g something
// --whatever=something -t 30 -i
// --something -i -g
// something somethingelse whatever
func Parse(args []string) (map[string]string, error) {
	if len(args) == 0 {
		return map[string]string{}, &BasicError{"empty arguments"}
	}
	resultMap := make(map[string]string)
	for currentPosition := 0; currentPosition < len(args); currentPosition++ {
		currentElement := args[currentPosition]
		if currentElement[0] == '-' {
			// cases like -t 30 ...
			// or --something -i -g
			// or --whatever=something ...
			// or -i -g ....
			if len(currentElement) < 2 {
				return map[string]string{}, &BasicError{"empty dash"}
			}
			if currentElement[1] == '-' {
				// case like --whatever=something
				// or --something -i -g
				k, v, err := parseDash(currentElement)
				if err != nil {
					return map[string]string{}, err
				}
				resultMap[k] = v
			} else {
				// cases like -i -g ....
				// or -t 30 ...
				switch nextPosition := currentPosition + 1; true {
				case nextPosition >= len(args):
					// edge case of -i without anything after
					fallthrough
				case args[nextPosition][0] == '-':
					// case when current element is -i -g with next element is -g
					resultMap[currentElement] = ""
				default:
					// case when current element is -t 30 with next element is 30
					resultMap[currentElement] = args[nextPosition]
					// increment current position to skip next element that was added
					currentPosition++
				}
			}
		} else {
			// case of something somethingelse whatever
			resultMap[currentElement] = ""
		}
	}
	return resultMap, nil
}

// case like --whatever=something
// or --something -i -g
func parseDash(element string) (string, string, error) {
	if len(element) == 0 {
		return "", "", &BasicError{"empty element"}
	}
	if strings.ContainsAny(element, "=") {
		// case like --whatever=something
		splitElement := strings.Split(element, "=")
		if len(splitElement) != 2 {
			return "", "", &BasicError{"more than one = in -- option"}
		}
		return splitElement[0], splitElement[1], nil
	}
	// case like --something ....
	return element, "", nil
}
