package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetVideoMetadata(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "video_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name           string
		fileContent    string
		expectedResult VideoMetadata
		expectError    bool
	}{
		{
			name:        "Valid H264 Video",
			fileContent: "dummy h264 video content",
			expectedResult: VideoMetadata{
				Duration: 120.5,
				IsH264:   true,
				Width:    1920,
				Height:   1080,
				FPS:      30,
			},
			expectError: false,
		},
		{
			name:        "Non-H264 Video",
			fileContent: "dummy video content",
			expectedResult: VideoMetadata{
				Duration: 60.0,
				IsH264:   false,
				Width:    1280,
				Height:   720,
				FPS:      24,
			},
			expectError: false,
		},
		{
			name:        "Invalid Video File",
			fileContent: "invalid content",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, "test_video.mp4")
			err := os.WriteFile(filePath, []byte(tc.fileContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := GetVideoMetadata(filePath)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result.Duration != tc.expectedResult.Duration {
					t.Errorf("Expected duration %f, but got %f", tc.expectedResult.Duration, result.Duration)
				}

				if result.IsH264 != tc.expectedResult.IsH264 {
					t.Errorf("Expected IsH264 %v, but got %v", tc.expectedResult.IsH264, result.IsH264)
				}

				if result.Width != tc.expectedResult.Width {
					t.Errorf("Expected width %d, but got %d", tc.expectedResult.Width, result.Width)
				}

				if result.Height != tc.expectedResult.Height {
					t.Errorf("Expected height %d, but got %d", tc.expectedResult.Height, result.Height)
				}

				if result.FPS != tc.expectedResult.FPS {
					t.Errorf("Expected FPS %f, but got %f", tc.expectedResult.FPS, result.FPS)
				}
			}
		})
	}
}

func TestGetVideoMetadataWithEmptyFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "video_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	emptyFilePath := filepath.Join(tempDir, "empty_video.mp4")
	err = os.WriteFile(emptyFilePath, []byte{}, 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	_, err = GetVideoMetadata(emptyFilePath)
	if err == nil {
		t.Errorf("Expected an error for empty file, but got none")
	}
}

func TestGetVideoMetadataWithNonExistentFile(t *testing.T) {
	nonExistentFilePath := "C:\\Users\\EDY\\Downloads\\system_uploads_2024062004483933.mp4"
	_, err := GetVideoMetadata(nonExistentFilePath)
	if err == nil {
		t.Errorf("Expected an error for non-existent file, but got none")
	}
}
