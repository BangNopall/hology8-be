package aws

import "testing"

func TestBuildPublicObjectURLUsesAmazonS3Host(t *testing.T) {
	got := buildPublicObjectURL("hology8-storage-test", "us-east-1", "users/ktm/image 1.jpg")
	want := "https://hology8-storage-test.s3.us-east-1.amazonaws.com/users/ktm/image%201.jpg"

	if got != want {
		t.Fatalf("buildPublicObjectURL() = %q, want %q", got, want)
	}
}

func TestExtractKeyFromURLSupportsAmazonS3URLAndPlainKey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "amazon s3 url",
			input: "https://hology8-storage-test.s3.us-east-1.amazonaws.com/users/ktm/image%201.jpg",
			want:  "users/ktm/image 1.jpg",
		},
		{
			name:  "legacy idcloudhost url",
			input: "https://hology.is3.cloudhost.id/users/ktm/image%201.jpg",
			want:  "users/ktm/image 1.jpg",
		},
		{
			name:  "plain object key",
			input: "users/ktm/image 1.jpg",
			want:  "users/ktm/image 1.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractKeyFromURL(tt.input)
			if err != nil {
				t.Fatalf("extractKeyFromURL() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("extractKeyFromURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
