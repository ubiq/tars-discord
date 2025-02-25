package main

import (
	"regexp"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

// Check if the name matches any of the given patterns
func isNameSpam(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, name)
		if matched {
			return true
		}
	}
	return false
}

// Check for spam user or display name and terminate member if conditions are met
func checkSpamName(s *discordgo.Session, m *discordgo.GuildMemberAdd) bool {
	patterns := []string{
		`(?i)Admin`,
		`(?i)Announcement`,
		`(?i)[CС][aа][pр]t[cс]h[aа]`,
		`(?i)FAQ`,
		`(?i)Giveaway`,
		`(?i)Helpdesk`,
		`(?i)Manager`,
		`(?i)MEE6`,
		`(?i)Support`,
		`\d{4}$`,
	}

	username := m.User.Username
	displayName := m.User.GlobalName
	if (isNameSpam(username, patterns) || isNameSpam(displayName, patterns)) && len(m.Roles) == 0 {
		go turdifyMember(s, m, "Name spam")
		return true
	}

	// Random-looking usernames with high entropy
	if hasHighEntropy(username) {
		go turdifyMember(s, m, "High entropy spam name")
		return true
	}

	return false
}

// Calculate entropy of a string to detect random-looking usernames
func hasHighEntropy(s string) bool {
	// Skip short strings
	if len(s) < 10 {
		return false
	}

	// Count character types (lowercase, uppercase, digits)
	var lowerCount, upperCount, digitCount, otherCount int
	for _, char := range s {
		switch {
		case unicode.IsLower(char):
			lowerCount++
		case unicode.IsUpper(char):
			upperCount++
		case unicode.IsDigit(char):
			digitCount++
		default:
			otherCount++
		}
	}

	// Check for high variance in character types (indicates randomness)
	if lowerCount > 0 && digitCount > 0 && (float64(digitCount)/float64(len(s)) > 0.2) && (float64(lowerCount)/float64(len(s)) > 0.3) {
		return true
	}

	// Calculate transitions between character types (more transitions = more random)
	transitions := 0
	lastType := -1
	for _, char := range s {
		currentType := -1
		switch {
		case unicode.IsLower(char):
			currentType = 0
		case unicode.IsUpper(char):
			currentType = 1
		case unicode.IsDigit(char):
			currentType = 2
		default:
			currentType = 3
		}

		if lastType != -1 && currentType != lastType {
			transitions++
		}
		lastType = currentType
	}

	// If we have many transitions between character types, it's likely random
	return transitions > len(s)/3
}
