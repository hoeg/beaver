package event

import "testing"

func TestExtractCommitSha(t *testing.T) {
	tests := []struct {
		artifactID string
		expected   string
	}{
		{
			artifactID: "nameofdeployemnt-abcdef123456-andid",
			expected:   "abcdef123456",
		},
		{
			artifactID: "anotherdeployment-with-hyphens-123456abcdef-andid2",
			expected:   "123456abcdef",
		},
	}

	for _, test := range tests {
		result := extractCommitSha(test.artifactID)
		if result != test.expected {
			t.Errorf("extractCommitSha(%s) = %s, expected %s", test.artifactID, result, test.expected)
		}
	}
}
