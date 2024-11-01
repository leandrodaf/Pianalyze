package midi

// Definition of chords and their intervals as hashes for quick lookup.
// Intervals are relative to the root (0 represents the root).
var chordHashes = map[string]int{
	// Tríades e acordes básicos
	"Major":         hashChord([]int{0, 4, 7}),
	"Minor":         hashChord([]int{0, 3, 7}),
	"Augmented":     hashChord([]int{0, 4, 8}),
	"Diminished":    hashChord([]int{0, 3, 6}),
	"Suspended 2nd": hashChord([]int{0, 2, 7}),
	"Suspended 4th": hashChord([]int{0, 5, 7}),

	// Acordes com 6ª e 7ª
	"Major 6th":           hashChord([]int{0, 4, 7, 9}),
	"Minor 6th":           hashChord([]int{0, 3, 7, 9}),
	"Major 7th":           hashChord([]int{0, 4, 7, 11}),
	"Minor 7th":           hashChord([]int{0, 3, 7, 10}),
	"Dominant 7th":        hashChord([]int{0, 4, 7, 10}),
	"Augmented 7th":       hashChord([]int{0, 4, 8, 10}),
	"Augmented Major 7th": hashChord([]int{0, 4, 8, 11}),
	"Diminished 7th":      hashChord([]int{0, 3, 6, 9}),
	"Half-diminished":     hashChord([]int{0, 3, 6, 10}),
	"Minor Major 7th":     hashChord([]int{0, 3, 7, 11}),

	// Acordes com 9ª
	"Major 9th":            hashChord([]int{0, 4, 7, 11, 14}),
	"Minor 9th":            hashChord([]int{0, 3, 7, 10, 14}),
	"Dominant 9th":         hashChord([]int{0, 4, 7, 10, 14}),
	"Dominant 7th flat 9":  hashChord([]int{0, 4, 7, 10, 13}),
	"Dominant 7th sharp 9": hashChord([]int{0, 4, 7, 10, 15}),
	"Dominant 9th flat 5":  hashChord([]int{0, 4, 6, 10, 14}),
	"Dominant 9th sharp 5": hashChord([]int{0, 4, 8, 10, 14}),
	"Minor Major 9th":      hashChord([]int{0, 3, 7, 11, 14}),

	// Acordes com 11ª
	"Major 11th":            hashChord([]int{0, 4, 7, 11, 14, 17}),
	"Minor 11th":            hashChord([]int{0, 3, 7, 10, 14, 17}),
	"Dominant 11th":         hashChord([]int{0, 4, 7, 10, 14, 17}),
	"Dominant 7th sharp 11": hashChord([]int{0, 4, 7, 10, 18}),
	"Minor 11th flat 5":     hashChord([]int{0, 3, 6, 10, 17}),
	"Minor 11th sharp 5":    hashChord([]int{0, 3, 8, 10, 17}),

	// Acordes com 13ª
	"Major 13th":            hashChord([]int{0, 4, 7, 11, 14, 17, 21}),
	"Minor 13th":            hashChord([]int{0, 3, 7, 10, 14, 17, 21}),
	"Dominant 13th":         hashChord([]int{0, 4, 7, 10, 14, 17, 21}),
	"Dominant 13th flat 9":  hashChord([]int{0, 4, 7, 10, 13, 17, 21}),
	"Dominant 13th sharp 9": hashChord([]int{0, 4, 7, 10, 15, 17, 21}),

	// Acordes adicionais com tensões específicas e variações
	"Minor 6/9":                    hashChord([]int{0, 3, 7, 9, 14}),
	"6/9":                          hashChord([]int{0, 4, 7, 9, 14}),
	"Minor 7th flat 5":             hashChord([]int{0, 3, 6, 10}),
	"Major 7th sharp 5":            hashChord([]int{0, 4, 8, 11}),
	"Dominant 7th flat 9 flat 5":   hashChord([]int{0, 4, 6, 10, 13}),
	"Dominant 7th sharp 9 sharp 5": hashChord([]int{0, 4, 8, 10, 15}),

	// Outras variações avançadas
	"Suspended 4th add 9":           hashChord([]int{0, 5, 7, 14}),
	"Minor 9th flat 13":             hashChord([]int{0, 3, 7, 10, 13, 20}),
	"Dominant 7th flat 13":          hashChord([]int{0, 4, 7, 10, 20}),
	"Dominant 7th sharp 13":         hashChord([]int{0, 4, 7, 10, 22}),
	"Add 9":                         hashChord([]int{0, 4, 7, 14}),
	"Minor Add 9":                   hashChord([]int{0, 3, 7, 14}),
	"Dominant 13th flat 9 sharp 11": hashChord([]int{0, 4, 7, 10, 13, 18}),
	"Dominant 9th flat 13":          hashChord([]int{0, 4, 7, 10, 14, 20}),
	"Minor 11th add 13":             hashChord([]int{0, 3, 7, 10, 14, 21}),
	"Dominant 7th flat 9 sharp 13":  hashChord([]int{0, 4, 7, 10, 13, 22}),
	"Major 9th add 13":              hashChord([]int{0, 4, 7, 11, 14, 21}),
	"Minor 9th flat 11":             hashChord([]int{0, 3, 7, 10, 13, 17}),
	"Minor 13th sharp 11":           hashChord([]int{0, 3, 7, 10, 14, 18}),
	"Dominant 9th add sharp 11":     hashChord([]int{0, 4, 7, 10, 14, 18}),
	"Dominant 11th sharp 9":         hashChord([]int{0, 4, 7, 10, 15, 17}),
	"Suspended 4th add 13":          hashChord([]int{0, 5, 7, 21}),
	"Minor 9th add 13":              hashChord([]int{0, 3, 7, 10, 14, 21}),
	"Add 9 sharp 11":                hashChord([]int{0, 4, 7, 14, 18}),
	"Minor Add 9 sharp 11":          hashChord([]int{0, 3, 7, 14, 18}),
	"Dominant 7th flat 9 sharp 11":  hashChord([]int{0, 4, 7, 10, 13, 18}),
	"Dominant 7th sharp 9 sharp 11": hashChord([]int{0, 4, 7, 10, 15, 18}),
	"Dominant 13th sharp 9 flat 11": hashChord([]int{0, 4, 7, 10, 15, 16, 21}),
	"Minor 13th add flat 9":         hashChord([]int{0, 3, 7, 10, 13, 21}),
	"Minor 13th sharp 9":            hashChord([]int{0, 3, 7, 10, 15, 21}),
	"Major 9th sharp 13":            hashChord([]int{0, 4, 7, 11, 14, 22}),
	"Major 13th sharp 11":           hashChord([]int{0, 4, 7, 11, 14, 18, 21}),
}

// hashToChordNames maps hashes to a list of corresponding chord names.
var hashToChordNames = make(map[int][]string)

// init initializes the hashToChordNames map by populating it with chord names based on hashes.
func init() {
	for name, hash := range chordHashes {
		hashToChordNames[hash] = append(hashToChordNames[hash], name)
	}
}

// hashChord creates a unique hash for a chord pattern based on relative intervals.
// Uses bitmask where each bit represents the presence of a specific interval.
func hashChord(intervals []int) int {
	var hash int
	for _, interval := range intervals {
		// Ignores intervals outside range to prevent bit overflow
		if interval < 0 || interval >= 32 {
			continue
		}
		// Sets the corresponding bit for the interval
		hash |= 1 << interval
	}
	return hash
}

// getChordHash computes a unique hash for a set of notes.
// Notes are normalized to a single octave and represented as a bitmask.
func getChordHash(notes []int) int {
	var hash int
	for _, note := range notes {
		normNote := note % 12
		if normNote < 0 {
			normNote += 12
		}
		hash |= 1 << normNote
	}
	if hash == 0 {
		return 0
	}
	// Finds the root (lowest note present) to calculate relative intervals
	root := 0
	for root < 12 && (hash>>root)&1 == 0 {
		root++
	}
	relativeHash := 0
	for i := 0; i < 12; i++ {
		if (hash>>((root+i)%12))&1 == 1 {
			relativeHash |= 1 << i
		}
	}
	return relativeHash
}

// detectInversion checks the inversion of a chord and returns the inversion number (0 for root position).
func detectInversion(notes []int, chordHash int) int {
	var currentHash int
	for _, note := range notes {
		normNote := note % 12
		if normNote < 0 {
			normNote += 12
		}
		currentHash |= 1 << normNote
	}
	if currentHash == 0 {
		return -1
	}
	root := 0
	for root < 12 && (currentHash>>root)&1 == 0 {
		root++
	}
	for inversion := 0; inversion < 12; inversion++ {
		rotatedHash := (chordHash >> inversion) | (chordHash << (12 - inversion))
		rotatedHash &= (1 << 12) - 1
		rotatedCurrent := (currentHash >> root) | (currentHash << (12 - root))
		rotatedCurrent &= (1 << 12) - 1
		if rotatedHash == rotatedCurrent {
			return inversion
		}
	}
	return -1
}

// GetChordName checks if a set of notes matches a known chord pattern, detecting inversions and key.
// Returns the chord name, inversion, key, and a boolean indicating if a match was found.
func GetChordName(notes []int) (string, string, int, bool) {
	if len(notes) < 3 {
		return "", "", -1, false
	}
	notesHash := getChordHash(notes)
	key := -1
	for _, note := range notes {
		normNote := note % 12
		if normNote < 0 {
			normNote += 12
		}
		if key == -1 || normNote < key {
			key = normNote
		}
	}
	chordNames, exists := hashToChordNames[notesHash]
	if exists && len(chordNames) > 0 {
		for _, chordName := range chordNames {
			standardHash := chordHashes[chordName]
			inversion := detectInversion(notes, standardHash)
			if inversion != -1 {
				inversionName := ""
				switch inversion {
				case 0:
					inversionName = "Root position"
				case 1:
					inversionName = "1st inversion"
				case 2:
					inversionName = "2nd inversion"
				default:
					inversionName = "Unknown inversion"
				}
				return chordName, inversionName, key, true
			}
		}
	}
	return "", "", -1, false
}

// IsTriad checks if a chord is a triad.
// Returns true if it is a triad, false otherwise.
func IsTriad(chordName string) bool {
	hash, exists := chordHashes[chordName]
	if !exists {
		return false
	}
	count := 0
	for i := 0; i < 12; i++ {
		if (hash>>i)&1 == 1 {
			count++
		}
	}
	return count == 3
}
