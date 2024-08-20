package marathon

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type FFmpegConfig struct {
	FFmpegArgs string
	FFmpegPath string
}

type FFmpegBroadcaster struct {
	Config *FFmpegConfig

	ffmpegCmd *exec.Cmd
}

func DefaultFFmpegConfig() (s *FFmpegConfig) {
	return &FFmpegConfig{
		FFmpegArgs: FFmpegArgs,
		FFmpegPath: "",
	}
}

func NewFFmpegBroadcaster(config *FFmpegConfig) (IStreamingBroadcaster, error) {
	broadcaster := &FFmpegBroadcaster{
		Config: config,
	}
	return broadcaster, nil
}

func (s *FFmpegBroadcaster) CommandArgs(path string, streamingUrl string) []string {
	var args = []string{
		"-re",               // Real time
		"-stream_loop", "1", // Infinite Loop
		"-f", "concat", // Concat all videos
		"-i", path,
	}
	if s.Config.FFmpegArgs != "" {
		args = append(args, strings.Split(s.Config.FFmpegArgs, " ")...)
	}
	args = append(args, []string{
		"-f", "flv",
		streamingUrl,
	}...)

	return args
}

// Start Stream with ffmpeg
func (s *FFmpegBroadcaster) Play(path string, streamingUrl string) {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Warnf("file %s does not exist", path)
		return
	}

	// Kill stream
	if s.ffmpegCmd != nil {
		log.Debugf("killing process - PID %d", s.ffmpegCmd.Process.Pid)
		s.ffmpegCmd.Process.Kill()
		s.ffmpegCmd.Wait()
	}

	// Command Args
	args := s.CommandArgs(path, streamingUrl)
	cmd := exec.Command(s.Config.FFmpegPath, args...)
	s.ffmpegCmd = cmd
	if log.GetLevel() >= log.ErrorLevel {
		cmd.Stderr = os.Stderr
	}
	if log.GetLevel() >= log.DebugLevel {
		cmd.Stdout = os.Stdout
	}

	log.Debugf("executing command - %s", strings.Replace(cmd.String(), streamingUrl, "[hidden]", 1))
	err := cmd.Start()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("command running - PID %d", cmd.Process.Pid)
}

func (s *FFmpegBroadcaster) Prepare() error {
	var ffmpeg string = s.Config.FFmpegPath
	// Check if exists
	if ffmpeg != "" {
		if _, err := os.Stat(ffmpeg); os.IsNotExist(err) {
			log.Warnf("file not found for FFmpeg Path: %s. Trying to look up in binary PATH", ffmpeg)
			ffmpeg = ""
		}
	}
	// Look in PATH if not found, or configured
	if ffmpeg == "" {
		path, err := exec.LookPath("ffmpeg")
		if err != nil {
			log.Errorf("FFmpeg not found. %v", err)
			return err
		}
		ffmpeg = path
	}
	// Check version
	cmd := exec.Command(ffmpeg, "-version")
	if log.GetLevel() >= log.ErrorLevel {
		cmd.Stderr = os.Stderr
	}
	if log.GetLevel() >= log.DebugLevel {
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Run(); err != nil {
		log.Errorf("FFmpeg version not supported. %v", err)
		return err
	}
	s.Config.FFmpegPath = ffmpeg
	log.Infof("using FFmpeg at %s", s.Config.FFmpegPath)
	return nil
}

func (s *FFmpegBroadcaster) Shutdown() {
	if s.ffmpegCmd != nil {
		s.ffmpegCmd.Process.Kill()
	}
}
