package main

import "testing"

func TestDisplayVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		want    string
	}{
		{
			name:    "strips v prefix",
			version: "v0.0.1",
			commit:  "025f54f",
			want:    "0.0.1-025f54f",
		},
		{
			name:    "keeps plain semver",
			version: "1.2.3",
			commit:  "abcdef0",
			want:    "1.2.3-abcdef0",
		},
		{
			name:    "keeps semver prerelease and metadata",
			version: "v1.2.3-beta.1+build.5",
			commit:  "abcdef0",
			want:    "1.2.3-beta.1+build.5-abcdef0",
		},
		{
			name:    "defaults empty version",
			version: "",
			commit:  "abcdef0",
			want:    "0.0.0-abcdef0",
		},
		{
			name:    "defaults invalid version",
			version: "release-1",
			commit:  "abcdef0",
			want:    "0.0.0-abcdef0",
		},
		{
			name:    "defaults semver with numeric prerelease leading zero",
			version: "1.2.3-01",
			commit:  "abcdef0",
			want:    "0.0.0-abcdef0",
		},
		{
			name:    "defaults empty commit",
			version: "0.0.1",
			commit:  "",
			want:    "0.0.1-dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayVersion(tt.version, tt.commit)
			if got != tt.want {
				t.Fatalf("displayVersion(%q, %q) = %q, want %q", tt.version, tt.commit, got, tt.want)
			}
		})
	}
}
