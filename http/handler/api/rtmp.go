package api

import (
	"net/http"

	"github.com/datarhei/core/http/api"
	"github.com/datarhei/core/rtmp"

	"github.com/labstack/echo/v4"
)

// The RTMPHandler type provides a handler for retrieving details from the RTMPHandler server
type RTMPHandler struct {
	rtmp rtmp.Server
}

// NewRTMP returns a new RTMP type. You have to provide a RTMP server instance.
func NewRTMP(rtmp rtmp.Server) *RTMPHandler {
	return &RTMPHandler{
		rtmp: rtmp,
	}
}

// ListChannels lists all currently publishing streams
// @Summary List all publishing streams
// @Description List all currently publishing streams
// @ID rtmp-3-list-channels
// @Produce json
// @Success 200 {array} api.RTMPChannel
// @Security ApiKeyAuth
// @Router /api/v3/rtmp [get]
func (rtmph *RTMPHandler) ListChannels(c echo.Context) error {
	channels := rtmph.rtmp.Channels()

	list := []api.RTMPChannel{}

	for _, c := range channels {
		list = append(list, api.RTMPChannel{
			Name: c,
		})
	}

	return c.JSON(http.StatusOK, list)
}
