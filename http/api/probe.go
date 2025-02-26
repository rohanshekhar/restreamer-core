package api

import (
	"encoding/json"

	"github.com/datarhei/core/restream/app"
)

// ProbeIO represents a stream of a probed file
type ProbeIO struct {
	// common
	Address  string      `json:"url"`
	Format   string      `json:"format"`
	Index    uint64      `json:"index"`
	Stream   uint64      `json:"stream"`
	Language string      `json:"language"`
	Type     string      `json:"type"`
	Codec    string      `json:"codec"`
	Coder    string      `json:"coder"`
	Bitrate  json.Number `json:"bitrate_kbps"`
	Duration json.Number `json:"duration_sec"`

	// video
	FPS    json.Number `json:"fps"`
	Pixfmt string      `json:"pix_fmt"`
	Width  uint64      `json:"width"`
	Height uint64      `json:"height"`

	// audio
	Sampling uint64 `json:"sampling_hz"`
	Layout   string `json:"layout"`
	Channels uint64 `json:"channels"`
}

func (i *ProbeIO) Unmarshal(io *app.ProbeIO) {
	if io == nil {
		return
	}

	i.Address = io.Address
	i.Format = io.Format
	i.Index = io.Index
	i.Stream = io.Stream
	i.Language = io.Language
	i.Type = io.Type
	i.Codec = io.Codec
	i.Coder = io.Coder
	i.Bitrate = toNumber(io.Bitrate)
	i.Duration = toNumber(io.Duration)

	i.FPS = toNumber(io.FPS)
	i.Pixfmt = io.Pixfmt
	i.Width = io.Width
	i.Height = io.Height

	i.Sampling = io.Sampling
	i.Layout = io.Layout
	i.Channels = io.Channels
}

// Probe represents the result of probing a file. It has a list of detected streams
// and a list of log lone from the probe process.
type Probe struct {
	Streams []ProbeIO `json:"streams"`
	Log     []string  `json:"log"`
}

// Unmarshal converts a restreamer Probe to a Probe in API representation
func (probe *Probe) Unmarshal(p *app.Probe) {
	if p == nil {
		return
	}

	probe.Streams = make([]ProbeIO, len(p.Streams))
	probe.Log = make([]string, len(p.Log))

	for i, io := range p.Streams {
		probe.Streams[i].Unmarshal(&io)
	}

	copy(probe.Log, p.Log)
}
