package api

import (
	"net/http"
	"strings"

	"github.com/datarhei/core/http/api"
	"github.com/datarhei/core/http/handler/util"
	"github.com/datarhei/core/session"

	"github.com/labstack/echo/v4"
)

// The SessionHandler type provides handlers to retrieve session information
type SessionHandler struct {
	registry session.Registry
}

// NewSession returns a new Session type. You have to provide a session registry.
func NewSession(registry session.Registry) *SessionHandler {
	return &SessionHandler{
		registry: registry,
	}
}

// Summary returns a summary of all active and past sessions
// @Summary Get a summary of all active and past sessions
// @Description Get a summary of all active and past sessions of the given collector
// @ID session-3-summary
// @Produce json
// @Security ApiKeyAuth
// @Param collectors query string false "Comma separated list of collectors"
// @Success 200 {object} api.SessionsSummary "Sessions summary"
// @Failure 404 {object} api.Error
// @Router /api/v3/session [get]
func (s *SessionHandler) Summary(c echo.Context) error {
	collectors := strings.Split(util.DefaultQuery(c, "collectors", ""), ",")

	sessionsSummary := make(api.SessionsSummary)

	for _, name := range collectors {
		collector := s.registry.Collector(name)
		if collector == nil {
			return api.Err(http.StatusNotFound, "Unknown collector", "Registered collectors are: %s", strings.Join(s.registry.Collectors(), ", "))
		}

		summary := api.SessionSummary{}
		summary.Unmarshal(collector.Summary())

		sessionsSummary[name] = summary
	}

	return c.JSON(http.StatusOK, sessionsSummary)
}

// Active returns a list of active sessions
// @Summary Get a minimal summary of all active sessions
// @Description Get a minimal summary of all active sessions (i.e. number of sessions, bandwidth)
// @ID session-3-current
// @Produce json
// @Security ApiKeyAuth
// @Param collectors query string false "Comma separated list of collectors"
// @Success 200 {object} api.SessionsActive "Active sessions listing"
// @Failure 404 {object} api.Error
// @Router /api/v3/session/active [get]
func (s *SessionHandler) Active(c echo.Context) error {
	collectors := strings.Split(util.DefaultQuery(c, "collectors", ""), ",")

	sessionsActive := make(api.SessionsActive)

	for _, name := range collectors {
		collector := s.registry.Collector(name)
		if collector == nil {
			return api.Err(http.StatusNotFound, "Unknown collector", "Registered collectors are: %s", strings.Join(s.registry.Collectors(), ", "))
		}

		sessions := collector.Active()

		active := make([]api.Session, len(sessions))

		for i, s := range sessions {
			active[i].Unmarshal(s)
		}

		sessionsActive[name] = active
	}

	return c.JSON(http.StatusOK, sessionsActive)
}
