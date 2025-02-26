package restream

import (
	"testing"
	"time"

	"github.com/datarhei/core/ffmpeg"
	"github.com/datarhei/core/net"
	"github.com/datarhei/core/restream/app"
	"github.com/stretchr/testify/require"
)

func getDummyRestreamer(portrange net.Portranger) Restreamer {
	ffmpeg, _ := ffmpeg.New(ffmpeg.Config{
		Binary:    "ffmpeg",
		Portrange: portrange,
	})

	rs, _ := New(Config{
		FFmpeg: ffmpeg,
	})

	return rs
}

func getDummyProcess() *app.Config {
	return &app.Config{
		ID: "process",
		Input: []app.ConfigIO{
			{
				ID:      "in",
				Address: "testsrc=size=1280x720:rate=25",
				Options: []string{
					"-f",
					"lavfi",
					"-re",
				},
			},
		},
		Output: []app.ConfigIO{
			{
				ID:      "out",
				Address: "-",
				Options: []string{
					"-codec",
					"copy",
					"-f",
					"null",
				},
			},
		},
		Options: []string{
			"-loglevel",
			"info",
		},
		Reconnect:      true,
		ReconnectDelay: 10,
		Autostart:      false,
		StaleTimeout:   0,
	}
}

func TestAddProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	_, err = rs.GetProcess(process.ID)
	require.NotEqual(t, nil, err, "Unset process found (%s)", process.ID)

	err = rs.AddProcess(process)
	require.Equal(t, nil, err, "Failed to add process (%s)", err)

	_, err = rs.GetProcess(process.ID)
	require.Equal(t, nil, err, "Set process not found (%s)", process.ID)

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")
}

func TestAutostartProcess(t *testing.T) {
	rs := getDummyRestreamer(nil)
	process := getDummyProcess()
	process.Autostart = true

	rs.AddProcess(process)

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestAddInvalidProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)

	// Invalid process ID
	process := getDummyProcess()
	process.ID = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process without ID")

	// Invalid input ID
	process = getDummyProcess()
	process.Input[0].ID = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process input without ID")

	// Invalid input address
	process = getDummyProcess()
	process.Input[0].Address = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process input without address")

	// Duplicate input ID
	process = getDummyProcess()
	process.Input = append(process.Input, process.Input[0])

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process input with duplicate ID")

	// No inputs
	process = getDummyProcess()
	process.Input = nil

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process without inputs")

	// Invalid output ID
	process = getDummyProcess()
	process.Output[0].ID = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process output without ID")

	// Invalid output address
	process = getDummyProcess()
	process.Output[0].Address = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process output without address")

	// Duplicate output ID
	process = getDummyProcess()
	process.Output = append(process.Output, process.Output[0])

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process output with duplicate ID")

	// No outputs
	process = getDummyProcess()
	process.Output = nil

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process without outputs")
}

func TestRemoveProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	err = rs.AddProcess(process)
	require.Equal(t, nil, err, "Failed to add process (%s)", err)

	err = rs.DeleteProcess(process.ID)
	require.Equal(t, nil, err, "Set process not found (%s)", process.ID)

	_, err = rs.GetProcess(process.ID)
	require.NotEqual(t, nil, err, "Unset process found (%s)", process.ID)
}

func TestGetProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)

	_, err = rs.GetProcess(process.ID)
	require.Equal(t, nil, err, "Process not found (%s)", process.ID)

	list := rs.GetProcessIDs()
	require.Len(t, list, 1, "expected 1 process")
	require.Equal(t, process.ID, list[0], "expected same process ID")
}

func TestStartProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)

	err = rs.StartProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to start non-existing process")

	err = rs.StartProcess(process.ID)
	require.Equal(t, nil, err, "should be able to start existing process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	err = rs.StartProcess(process.ID)
	require.Equal(t, nil, err, "should be able to start already running process")

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestStopProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)
	rs.StartProcess(process.ID)

	err = rs.StopProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to stop non-existing process")

	err = rs.StopProcess(process.ID)
	require.Equal(t, nil, err, "should be able to stop existing running process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")

	err = rs.StopProcess(process.ID)
	require.Equal(t, nil, err, "should be able to stop already stopped process")

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")
}

func TestRestartProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)

	err = rs.RestartProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to restart non-existing process")

	err = rs.RestartProcess(process.ID)
	require.Equal(t, nil, err, "should be able to restart existing stopped process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")

	rs.StartProcess(process.ID)

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestReloadProcess(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)

	err = rs.ReloadProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to reload non-existing process")

	err = rs.ReloadProcess(process.ID)
	require.Equal(t, nil, err, "should be able to reload existing stopped process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")

	rs.StartProcess(process.ID)

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	err = rs.ReloadProcess(process.ID)
	require.Equal(t, nil, err, "should be able to reload existing process")

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestProcessData(t *testing.T) {
	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)

	data, _ := rs.GetProcessMetadata(process.ID, "foobar")
	require.Equal(t, nil, data, "nothing should be stored under the key")

	rs.SetProcessMetadata(process.ID, "foobar", process)

	data, _ = rs.GetProcessMetadata(process.ID, "foobar")
	require.NotEqual(t, nil, data, "there should be something stored under the key")

	p := data.(*app.Config)

	require.Equal(t, process.ID, p.ID, "failed to retrieve stored data")
}

func TestLog(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	rs.AddProcess(process)

	_, err = rs.GetProcessLog("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to get log from non-existing process")

	log, err := rs.GetProcessLog(process.ID)
	require.Equal(t, nil, err, "should be able to get log from existing process")
	require.Equal(t, 0, len(log.Prelude))
	require.Equal(t, 0, len(log.Log))

	rs.StartProcess(process.ID)

	time.Sleep(3 * time.Second)

	log, _ = rs.GetProcessLog(process.ID)

	require.NotEqual(t, 0, len(log.Prelude))
	require.NotEqual(t, 0, len(log.Log))

	rs.StopProcess(process.ID)

	log, _ = rs.GetProcessLog(process.ID)

	require.NotEqual(t, 0, len(log.Prelude))
	require.NotEqual(t, 0, len(log.Log))
}

func TestPlayoutNoRange(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process := getDummyProcess()

	process.Input[0].Address = "playout:" + process.Input[0].Address

	rs.AddProcess(process)

	_, err = rs.GetPlayout("foobar", process.Input[0].ID)
	require.NotEqual(t, nil, err, "playout of non-existing process should error")

	_, err = rs.GetPlayout(process.ID, "foobar")
	require.NotEqual(t, nil, err, "playout of non-existing input should error")

	addr, _ := rs.GetPlayout(process.ID, process.Input[0].ID)
	require.Equal(t, 0, len(addr), "the playout address should be empty if no port range is given")
}

func TestPlayoutRange(t *testing.T) {
	var err error = nil

	portrange, _ := net.NewPortrange(3000, 3001)

	rs := getDummyRestreamer(portrange)
	process := getDummyProcess()

	process.Input[0].Address = "playout:" + process.Input[0].Address

	rs.AddProcess(process)

	_, err = rs.GetPlayout("foobar", process.Input[0].ID)
	require.NotEqual(t, nil, err, "playout of non-existing process should error")

	_, err = rs.GetPlayout(process.ID, "foobar")
	require.NotEqual(t, nil, err, "playout of non-existing input should error")

	addr, _ := rs.GetPlayout(process.ID, process.Input[0].ID)
	require.NotEqual(t, 0, len(addr), "the playout address should not be empty if a port range is given")
	require.Equal(t, "127.0.0.1:3000", addr, "the playout address should be 127.0.0.1:3000")
}

func TestAddressReference(t *testing.T) {
	var err error = nil

	rs := getDummyRestreamer(nil)
	process1 := getDummyProcess()
	process2 := getDummyProcess()

	process2.ID = "process2"

	rs.AddProcess(process1)

	process2.Input[0].Address = "#process:foobar=out"

	err = rs.AddProcess(process2)
	require.NotEqual(t, nil, err, "shouldn't resolve invalid reference")

	process2.Input[0].Address = "#process2:output=out"

	err = rs.AddProcess(process2)
	require.NotEqual(t, nil, err, "shouldn't resolve invalid reference")

	process2.Input[0].Address = "#process:output=foobar"

	err = rs.AddProcess(process2)
	require.NotEqual(t, nil, err, "shouldn't resolve invalid reference")

	process2.Input[0].Address = "#process:output=out"

	err = rs.AddProcess(process2)
	require.Equal(t, nil, err, "should resolve reference")
}

func TestOutputAddressValidation(t *testing.T) {
	rs := getDummyRestreamer(nil).(*restream)

	type res struct {
		path string
		err  bool
	}

	paths := map[string]res{
		"/dev/null":                   {"file:/dev/null", false},
		"/dev/../etc/passwd":          {"/etc/passwd", true},
		"/dev/fb0":                    {"file:/dev/fb0", false},
		"/etc/passwd":                 {"/etc/passwd", true},
		"/core/data/../../etc/passwd": {"/etc/passwd", true},
		"/core/data/./etc/passwd":     {"file:/core/data/etc/passwd", false},
		"file:/core/data/foobar":      {"file:/core/data/foobar", false},
		"http://example.com":          {"http://example.com", false},
		"-":                           {"pipe:", false},
		"tee:/core/data/foobar|http://example.com": {"tee:/core/data/foobar|http://example.com", false},
		"tee:/core/data/foobar|/etc/passwd":        {"tee:/core/data/foobar|/etc/passwd", true},
	}

	for path, r := range paths {
		path, _, err := rs.validateOutputAddress(path, "/core/data")

		if r.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, r.path, path)
	}
}
