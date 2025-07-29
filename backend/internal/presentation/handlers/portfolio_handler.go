package handlers

import (
	"net/http"
	"strconv"
	"crypto-indicator-dashboard/internal/application/dto"
	"crypto-indicator-dashboard/internal/application/usecases"
	"crypto-indicator-dashboard/pkg/errors"
	"crypto-indicator-dashboard/pkg/logger"
	"github.com/gin-gonic/gin"
)

// PortfolioHandler handles portfolio-related HTTP requests
type PortfolioHandler struct {
	portfolioUseCase *usecases.PortfolioUseCase
	logger           logger.Logger
}

// NewPortfolioHandler creates a new portfolio handler
func NewPortfolioHandler(portfolioUseCase *usecases.PortfolioUseCase, logger logger.Logger) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioUseCase: portfolioUseCase,
		logger:           logger.With("handler", "portfolio"),
	}
}

// CreatePortfolio creates a new portfolio
func (h *PortfolioHandler) CreatePortfolio(c *gin.Context) {
	var req dto.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.Validation("Invalid request format", err.Error()))
		return
	}
	
	portfolio, err := h.portfolioUseCase.CreatePortfolio(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	h.logger.Info("Portfolio created successfully", "portfolio_id", portfolio.ID, "user_id", req.UserID)
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Portfolio created successfully",
		"data":    portfolio,
	})
}

// GetPortfolio retrieves a portfolio by ID
func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	portfolioID, err := h.parseUintParam(c, "id")
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	portfolio, err := h.portfolioUseCase.GetPortfolio(c.Request.Context(), portfolioID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    portfolio,
	})
}

// GetUserPortfolios retrieves all portfolios for a user
func (h *PortfolioHandler) GetUserPortfolios(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = "default_user" // In production, get from JWT token
	}
	
	portfolios, err := h.portfolioUseCase.GetUserPortfolios(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    portfolios,
	})
}

// GetPortfolioSummary retrieves portfolio summary with analytics
func (h *PortfolioHandler) GetPortfolioSummary(c *gin.Context) {
	portfolioID, err := h.parseUintParam(c, "id")
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	summary, err := h.portfolioUseCase.GetPortfolioSummary(c.Request.Context(), portfolioID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// AddHolding adds a new holding to a portfolio
func (h *PortfolioHandler) AddHolding(c *gin.Context) {
	portfolioID, err := h.parseUintParam(c, "id")
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	var req dto.AddHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.Validation("Invalid request format", err.Error()))
		return
	}
	
	req.PortfolioID = portfolioID
	
	holding, err := h.portfolioUseCase.AddHolding(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	h.logger.Info("Holding added successfully", 
		"portfolio_id", portfolioID, 
		"symbol", req.Symbol,
		"amount", req.Amount,
	)
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Holding added successfully",
		"data":    holding,
	})
}

// UpdateHolding updates an existing holding
func (h *PortfolioHandler) UpdateHolding(c *gin.Context) {
	holdingID, err := h.parseUintParam(c, "holdingId")
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	var req dto.UpdateHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.Validation("Invalid request format", err.Error()))
		return
	}
	
	req.HoldingID = holdingID
	
	if err := h.portfolioUseCase.UpdateHolding(c.Request.Context(), &req); err != nil {
		h.handleError(c, err)
		return
	}
	
	h.logger.Info("Holding updated successfully", "holding_id", holdingID)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Holding updated successfully",
	})
}

// RemoveHolding removes a holding from a portfolio
func (h *PortfolioHandler) RemoveHolding(c *gin.Context) {
	holdingID, err := h.parseUintParam(c, "holdingId")
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	if err := h.portfolioUseCase.RemoveHolding(c.Request.Context(), holdingID); err != nil {
		h.handleError(c, err)
		return
	}
	
	h.logger.Info("Holding removed successfully", "holding_id", holdingID)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Holding removed successfully",
	})
}

// Helper methods

func (h *PortfolioHandler) parseUintParam(c *gin.Context, param string) (uint, error) {
	paramStr := c.Param(param)
	if paramStr == "" {
		return 0, errors.Validation("Missing parameter: " + param)
	}
	
	id, err := strconv.ParseUint(paramStr, 10, 32)
	if err != nil {
		return 0, errors.Validation("Invalid parameter format: " + param)
	}
	
	return uint(id), nil
}

func (h *PortfolioHandler) handleError(c *gin.Context, err error) {
	h.logger.Error("Request failed", "error", err, "path", c.Request.URL.Path)
	
	statusCode := errors.GetStatusCode(err)
	
	// Convert error to response format
	var errorResponse gin.H
	if appErr, ok := err.(*errors.AppError); ok {
		errorResponse = gin.H{
			"success": false,
			"error": gin.H{
				"type":    appErr.Type,
				"message": appErr.Message,
			},
		}
		if appErr.Details != "" {
			errorResponse["error"].(gin.H)["details"] = appErr.Details
		}
	} else {
		errorResponse = gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_ERROR",
				"message": "An internal error occurred",
			},
		}
	}
	
	c.JSON(statusCode, errorResponse)
}