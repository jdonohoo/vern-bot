package council

// Tier defines a named council configuration.
type Tier struct {
	Name     string
	Display  string   // Human-readable display name
	Core     []string // Required persona IDs
	MinFill  int      // Min total count (fill remaining randomly)
	MaxFill  int      // Max total count (fill remaining randomly)
	Fixed    bool     // If true, no random fill (exact core members)
}

// AllTiers returns the predefined council tier configurations.
func AllTiers() map[string]Tier {
	return map[string]Tier{
		"hammers": {
			Name:    "hammers",
			Display: "Council of the Three Hammers",
			Core:    []string{"great", "mediocre", "ketamine"},
			Fixed:   true,
		},
		"conflict": {
			Name:    "conflict",
			Display: "Max Conflict",
			Core:    []string{"startup", "enterprise", "yolo", "paranoid", "optimist", "inverse"},
			Fixed:   true,
		},
		"inner": {
			Name:    "inner",
			Display: "The Inner Circle",
			Core:    []string{"architect", "inverse", "paranoid"},
			MinFill: 3,
			MaxFill: 5,
		},
		"round": {
			Name:    "round",
			Display: "The Round Table",
			Core:    []string{"mighty", "yolo", "startup", "academic", "enterprise"},
			MinFill: 6,
			MaxFill: 9,
		},
		"war": {
			Name:    "war",
			Display: "The War Room",
			Core:    []string{"mighty", "yolo", "startup", "academic", "enterprise", "ux", "retro", "optimist", "nyquil"},
			MinFill: 10,
			MaxFill: 13,
		},
		"full": {
			Name:    "full",
			Display: "The Full Vern Experience",
			Fixed:   true,
			// Core is empty — means ALL available
		},
		"random": {
			Name:    "random",
			Display: "Fate's Hand",
			// Random count, random selection
		},
	}
}
