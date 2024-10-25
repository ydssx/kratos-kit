package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/go-kratos/kratos/v2/log"
)

type VideoMetadata struct {
	Duration  float64
	Width     int
	Height    int
	FPS       float64
	IsH264    bool
	Format    string
	CodecName string
}

type ffprobeMetadata struct {
	Streams []struct {
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		FPS       string `json:"r_frame_rate"`
		CodecName string `json:"codec_name"`
	} `json:"streams"`
}

func GetVideoMetadata(filePath string) (VideoMetadata, error) {
	metadata := VideoMetadata{}

	// 获取格式
	metadata.Format = strings.Trim(filepath.Ext(filePath), ".")

	// 获取时长和编码信息
	cmd := exec.Command("ffmpeg", "-i", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Run()
	output := stderr.String()

	// 解析时长
	durationRe := regexp.MustCompile(`Duration:\s+(\d+):(\d+):(\d+.\d+)`)
	durationMatches := durationRe.FindStringSubmatch(output)
	if len(durationMatches) > 0 {
		hours, _ := strconv.Atoi(durationMatches[1])
		minutes, _ := strconv.Atoi(durationMatches[2])
		seconds, _ := strconv.ParseFloat(durationMatches[3], 64)
		metadata.Duration = float64(hours)*3600 + float64(minutes)*60 + seconds
	}

	// 检查是否为H264编码
	metadata.IsH264 = strings.Contains(output, "h264") || strings.Contains(output, "H.264")

	// 获取分辨率和帧率
	cmd = exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,r_frame_rate,codec_name", "-of", "json", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return metadata, err
	}

	var ffprobeData ffprobeMetadata
	err = json.Unmarshal(out.Bytes(), &ffprobeData)
	if err != nil {
		return metadata, err
	}

	if len(ffprobeData.Streams) > 0 {
		metadata.Width = ffprobeData.Streams[0].Width
		metadata.Height = ffprobeData.Streams[0].Height
		metadata.CodecName = ffprobeData.Streams[0].CodecName

		fpsNumerator, fpsDenominator, _ := strings.Cut(ffprobeData.Streams[0].FPS, "/")
		fpsNum, _ := strconv.ParseFloat(fpsNumerator, 64)
		fpsDen, _ := strconv.ParseFloat(fpsDenominator, 64)
		if fpsDen != 0 {
			metadata.FPS = fpsNum / fpsDen
		}
	}

	return metadata, nil
}

// GenerateThumbnail generates a thumbnail image from a video file using ffmpeg.
// It takes the path of the video file as input and returns the path of the generated thumbnail.
// If the thumbnail generation is successful, it returns the thumbnail path along with no error.
// If there is an error during the thumbnail generation process, it returns an empty string and the error.
func GenerateThumbnail(videoPath string) (string, error) {
	// Generate the thumbnail path by replacing the video file extension with "_thumbnail.jpg"
	thumbnailPath := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + "_thumbnail.jpg"

	// Log the thumbnail path
	log.Infof("generate thumbnail: %s", thumbnailPath)

	// Run the ffmpeg command to generate the thumbnail
	// The command extracts a single frame from the video file at the specified timestamp and saves it as the thumbnail
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01.000", "-vframes", "1", thumbnailPath)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Check if the thumbnail file exists and is not empty.
	_, err := os.Stat(thumbnailPath)
	if err != nil || os.IsNotExist(err) {
		return "", fmt.Errorf("failed to generate thumbnail: %s", err.Error())
	}

	// Return the thumbnail path and no error if the thumbnail generation is successful
	return thumbnailPath, nil
}

// GetVideoDuration returns the duration of a video file in seconds.
func GetVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffmpeg", "-i", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// FFmpeg 在这种情况下会返回一个错误，因为我们没有实际处理视频流
		// 我们可以忽略这个错误，继续解析输出内容
	}

	// 使用正则表达式匹配时长信息
	re := regexp.MustCompile(`Duration:\s+(\d+):(\d+):(\d+.\d+)`)
	matches := re.FindStringSubmatch(stderr.String())

	if len(matches) == 0 {
		return 0, fmt.Errorf("could not find duration in ffmpeg output")
	}

	// 解析时长信息
	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.ParseFloat(matches[3], 64)

	// 计算总时长（秒）
	duration := float64(hours)*3600 + float64(minutes)*60 + seconds

	return duration, nil
}

// Metadata struct to unmarshal ffprobe JSON output
type Metadata struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

// GetMediaDimensions returns the dimensions of a media file
func GetMediaDimensions(filePath string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, 0, err
	}

	var metadata Metadata
	err = json.Unmarshal(out.Bytes(), &metadata)
	if err != nil {
		return 0, 0, err
	}

	if len(metadata.Streams) > 0 {
		return metadata.Streams[0].Width, metadata.Streams[0].Height, nil
	}
	return 0, 0, fmt.Errorf("no streams found")
}

// CompressFrameRate compresses a video file to a specified frame rate.
func CompressFrameRate(inputFile, outputFile string, frameRate int) error {
	// 构建ffmpeg命令
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-r", fmt.Sprintf("%d", frameRate), outputFile)

	// 运行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	_, err = os.Stat(outputFile)
	if err != nil || os.IsNotExist(err) {
		return fmt.Errorf("failed to generate thumbnail: %s", err.Error())
	}

	return nil
}

// CheckVideoEncodingH264 检查视频编码是否为H.264
func CheckVideoEncodingH264(videoFilePath string) (bool, error) {
	// FFmpeg命令来获取视频流的信息
	cmd := exec.Command("ffmpeg", "-i", videoFilePath)

	// 捕获输出
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		// FFmpeg命令通常会在输出视频信息时返回一个非零状态码，所以这里不直接判断错误
		// 而是继续解析输出信息
		// fmt.Println("Running FFmpeg command (non-zero exit code expected):", err)
	}

	// 解析输出来提取编码信息
	output := stderr.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Video:") {
			if strings.Contains(line, "h264") || strings.Contains(line, "H.264") {
				return true, nil
			}
			break
		}
	}

	return false, nil
}

// ConvertToH264 将视频转换为H.264编码
func ConvertToH264(inputFilePath, outputFilePath string) error {
	// FFmpeg命令来转换视频为H.264编码
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-preset", "slow", "-crf", "22", outputFilePath, "-y")

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting video to H.264: %v", err)
	}

	return nil
}

// CompressImage 压缩图片并保存到指定路径
//
// 参数:
//   - inputPath: 输入图片的文件路径
//   - outputPath: 压缩后图片的保存路径
//   - quality: 压缩质量，范围1-100，值越大质量越高
//   - maxWidth: 图片的最大宽度，如果为0则保持原始宽度
//   - maxHeight: 图片的最大高度，如果为0则保持原始高度
//
// 返回:
//   - error: 如果压缩过程中出现错误，返回相应的错误信息
func CompressImage(inputPath, outputPath string, quality, maxWidth, maxHeight int) error {
	// 打开原始图片
	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %w", err)
	}

	// 调整图片大小
	if maxWidth > 0 || maxHeight > 0 {
		src = imaging.Fit(src, maxWidth, maxHeight, imaging.Lanczos)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 保存压缩后的图片
	err = imaging.Save(src, outputPath, imaging.JPEGQuality(quality))
	if err != nil {
		return fmt.Errorf("保存压缩图片失败: %w", err)
	}

	return nil
}
