import { useState, useEffect, useRef, useCallback } from 'react'
import { ChartBarIcon, CurrencyDollarIcon, ArrowTrendingUpIcon, ClockIcon, XMarkIcon } from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Badge } from './ui/badge'
import { Separator } from './ui/separator'
import { Progress } from './ui/progress'
import { cn } from '@/lib/utils'
import DCAChart from './DCAChart'
import { 
  buildCardStyles, 
  cardHeights, 
  colors, 
  typography, 
  spacing, 
  effects,
  layoutStyles,
  iconStyles,
  brandingStyles
} from '../utils/designSystem'

// Debounce hook for smooth slider performance
function useDebounce(value, delay) {
  const [debouncedValue, setDebouncedValue] = useState(value)

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])

  return debouncedValue
}

function DCASidebar({ isOpen, onClose }) {
  const sliderRef = useRef(null)
  
  const [formData, setFormData] = useState({
    symbol: 'BTC',
    amount: 100,
    frequency: 'weekly',
    startDate: '2023-01-01',
    endDate: '2024-01-01'
  })

  const [simulation, setSimulation] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  
  // Debounce amount for smooth slider without API spam
  const debouncedAmount = useDebounce(formData.amount, 300)

  // Close sidebar with Escape key
  useEffect(() => {
    const handleEscape = (e) => {
      if (e.key === 'Escape' && isOpen) {
        onClose()
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [isOpen, onClose])

  // Prevent body scroll when sidebar is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'unset'
    }

    return () => {
      document.body.style.overflow = 'unset'
    }
  }, [isOpen])

  const handleInputChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: name === 'amount' ? parseFloat(value) || 0 : value
    }))
  }

  // Optimized slider handler with RAF for 60fps smoothness
  const handleSliderChange = useCallback((e) => {
    const value = parseFloat(e.target.value)
    setFormData(prev => ({ ...prev, amount: value }))
  }, [])

  const handleFrequencyChange = (frequency) => {
    setFormData(prev => ({ ...prev, frequency }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('http://localhost:8080/api/v1/dca/simulate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          symbol: formData.symbol,
          amount: formData.amount,
          frequency: formData.frequency,
          start_date: new Date(formData.startDate).toISOString(),
          end_date: new Date(formData.endDate).toISOString()
        })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || errorData.details || `HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      setSimulation(result)
    } catch (err) {
      console.error('DCA simulation error:', err)
      setError(err.message || 'Failed to simulate DCA strategy')
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(value)
  }

  const formatPercent = (value) => {
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  }

  const getRiskVariant = (level) => {
    switch (level?.toLowerCase()) {
      case 'low': return 'success'
      case 'medium': return 'warning'
      case 'high': return 'destructive'
      default: return 'secondary'
    }
  }

  const getRiskColor = (level) => {
    switch (level?.toLowerCase()) {
      case 'low': return 'text-green-400'
      case 'medium': return 'text-yellow-400'
      case 'high': return 'text-red-400'
      default: return 'text-gray-400'
    }
  }

  const getReturnColor = (value) => {
    return value >= 0 ? 'text-green-400' : 'text-red-400'
  }

  // Calculate slider percentage for gradient effect
  const sliderPercentage = ((formData.amount - 10) / (10000 - 10)) * 100

  return (
    <>
      {/* Backdrop Overlay */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 transition-opacity duration-300"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <div className={cn(
        "fixed top-0 left-0 h-full",
        colors.gradientBg,
        "border-r",
        colors.border,
        "shadow-2xl z-50 transition-transform duration-300 ease-out",
        "w-full md:w-[480px] lg:w-[520px] xl:w-[580px]",
        "overflow-y-auto",
        isOpen ? 'translate-x-0' : '-translate-x-full'
      )}>
        {/* Header */}
        <div className={cn(
          "sticky top-0",
          colors.cardBg,
          "border-b",
          colors.border,
          spacing.sidebarPadding,
          "z-10"
        )}>
          <div className={layoutStyles.headerLayout}>
            <div className={brandingStyles.logoContainer}>
              <CurrencyDollarIcon className={cn(iconStyles.xl, colors.accent)} />
              <div>
                <h1 className={cn(
                  typography.cardTitle,
                  colors.textPrimary
                )}>DCA Calculator</h1>
                <p className={cn(
                  typography.bodySmall,
                  colors.textSecondary
                )}>Simulate dollar cost averaging strategies</p>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className={cn(colors.textSecondary, "hover:text-foreground")}
            >
              <XMarkIcon className={iconStyles.large} />
            </Button>
          </div>
        </div>

        {/* Content */}
        <div className={cn(spacing.sidebarPadding, layoutStyles.sectionSpacing)}>
          {/* Input Form */}
          <Card className={cn(
            buildCardStyles(false, 'professional'),
            cardHeights.dcaSidebarCard
          )}>
            <CardHeader className={spacing.cardPaddingSmall}>
              <CardTitle className={cn(typography.cardTitle, colors.textPrimary)}>
                Strategy Configuration
              </CardTitle>
            </CardHeader>
            <CardContent className={spacing.cardPaddingSmall}>
            
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Cryptocurrency Selection */}
              <div className="space-y-3">
                <label className="block text-sm font-medium text-foreground">
                  Cryptocurrency
                </label>
                <select
                  name="symbol"
                  value={formData.symbol}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 bg-background border border-input rounded-lg text-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-all duration-200"
                >
                  <option value="BTC">Bitcoin (BTC)</option>
                  <option value="ETH">Ethereum (ETH)</option>
                </select>
              </div>

              {/* Investment Amount Slider */}
              <div>
                <label className="block text-sm font-medium text-foreground mb-3">
                  Investment Amount per Purchase
                </label>
                <div className="space-y-4">
                  {/* Amount Display */}
                  <div className="text-center">
                    <div className="text-3xl font-bold bg-gradient-to-r from-blue-400 via-purple-500 to-pink-500 bg-clip-text text-transparent">
                      {formatCurrency(formData.amount)}
                    </div>
                    <div className="text-sm text-muted-foreground mt-1">per purchase</div>
                  </div>
                  
                  {/* Custom Styled Slider */}
                  <div className="relative">
                    <input
                      ref={sliderRef}
                      type="range"
                      min="10"
                      max="10000"
                      step="10"
                      value={formData.amount}
                      onChange={handleSliderChange}
                      className="w-full h-3 bg-input rounded-lg appearance-none cursor-pointer slider-custom"
                      style={{
                        background: `linear-gradient(to right, 
                          #3b82f6 0%, 
                          #8b5cf6 ${sliderPercentage}%, 
                          hsl(var(--input)) ${sliderPercentage}%, 
                          hsl(var(--input)) 100%)`
                      }}
                    />
                    
                    {/* Amount Markers */}
                    <div className="flex justify-between text-xs text-muted-foreground mt-2 px-1">
                      <span>$10</span>
                      <span>$100</span>
                      <span>$500</span>
                      <span>$1K</span>
                      <span>$5K</span>
                      <span>$10K</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Frequency Toggle Buttons */}
              <div>
                <label className="block text-sm font-medium text-foreground mb-3">
                  Purchase Frequency
                </label>
                <div className="relative bg-muted/30 p-1 rounded-xl">
                  <div className="grid grid-cols-3 relative">
                    {/* Animated Background */}
                    <div 
                      className={`absolute top-1 bottom-1 bg-gradient-to-r from-primary to-accent rounded-lg transition-transform duration-300 ease-out`}
                      style={{
                        width: '33.333%',
                        transform: `translateX(${
                          formData.frequency === 'daily' ? '0%' : 
                          formData.frequency === 'weekly' ? '100%' : '200%'
                        })`
                      }}
                    />
                    
                    {/* Frequency Buttons */}
                    {['daily', 'weekly', 'monthly'].map((freq) => (
                      <button
                        key={freq}
                        type="button"
                        onClick={() => handleFrequencyChange(freq)}
                        className={`relative z-10 py-3 px-4 text-sm font-medium rounded-lg transition-all duration-200 ${
                          formData.frequency === freq
                            ? 'text-primary-foreground'
                            : 'text-muted-foreground hover:text-foreground'
                        }`}
                      >
                        {freq.charAt(0).toUpperCase() + freq.slice(1)}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              {/* Date Range */}
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-3">
                  <label className="block text-sm font-medium text-foreground">
                    Start Date
                  </label>
                  <Input
                    type="date"
                    name="startDate"
                    value={formData.startDate}
                    onChange={handleInputChange}
                  />
                </div>

                <div className="space-y-3">
                  <label className="block text-sm font-medium text-foreground">
                    End Date
                  </label>
                  <Input
                    type="date"
                    name="endDate"
                    value={formData.endDate}
                    onChange={handleInputChange}
                  />
                </div>
              </div>

              {/* Submit Button */}
              <Button
                type="submit"
                disabled={loading}
                variant="gradient"
                className="w-full px-6 py-4 font-medium transform hover:scale-[1.02] transition-all duration-200"
              >
                {loading ? (
                  <>
                    <div className="flex space-x-1 mr-3">
                      <div className="w-2 h-2 bg-current rounded-full animate-bounce"></div>
                      <div className="w-2 h-2 bg-current rounded-full animate-bounce" style={{animationDelay: '0.1s'}}></div>
                      <div className="w-2 h-2 bg-current rounded-full animate-bounce" style={{animationDelay: '0.2s'}}></div>
                    </div>
                    <span>Simulating Strategy...</span>
                  </>
                ) : (
                  <>
                    <ChartBarIcon className="w-5 h-5 mr-3" />
                    <span>Simulate DCA Strategy</span>
                  </>
                )}
              </Button>
            </form>

            {error && (
              <Card className="mt-6 border-destructive bg-destructive/10">
                <CardContent className="p-4">
                  <p className="text-destructive text-sm">{error}</p>
                </CardContent>
              </Card>
            )}
            </CardContent>
          </Card>

          {/* Results Section */}
          {simulation ? (
            <div className="space-y-6">
              {/* Summary Cards */}
              <div className="grid grid-cols-2 gap-4">
                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardContent className="p-4">
                    <div className="flex items-center space-x-2 mb-2">
                      <CurrencyDollarIcon className="w-5 h-5 text-green-400" />
                      <span className="text-xs text-muted-foreground">Total Invested</span>
                    </div>
                    <div className="text-lg font-bold text-foreground">
                      {formatCurrency(simulation.summary?.total_invested || 0)}
                    </div>
                  </CardContent>
                </Card>

                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardContent className="p-4">
                    <div className="flex items-center space-x-2 mb-2">
                      <ArrowTrendingUpIcon className="w-5 h-5 text-blue-400" />
                      <span className="text-xs text-muted-foreground">Final Value</span>
                    </div>
                    <div className="text-lg font-bold text-foreground">
                      {formatCurrency(simulation.summary?.final_value || 0)}
                    </div>
                  </CardContent>
                </Card>

                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardContent className="p-4">
                    <div className="flex items-center space-x-2 mb-2">
                      <span className="text-xs text-muted-foreground">Total Return</span>
                    </div>
                    <div className={cn("text-lg font-bold", getReturnColor(simulation.summary?.total_return || 0))}>
                      {formatCurrency(simulation.summary?.total_return || 0)}
                    </div>
                    <div className={cn("text-xs", getReturnColor(simulation.summary?.total_return_pct || 0))}>
                      {formatPercent(simulation.summary?.total_return_pct || 0)}
                    </div>
                  </CardContent>
                </Card>

                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardContent className="p-4">
                    <div className="flex items-center space-x-2 mb-2">
                      <ClockIcon className="w-5 h-5 text-purple-400" />
                      <span className="text-xs text-muted-foreground">Purchases</span>
                    </div>
                    <div className="text-lg font-bold text-foreground">
                      {simulation.summary?.purchase_count || 0}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {simulation.summary?.duration_days || 0} days
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Performance Metrics */}
              {simulation.performance && (
                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardHeader>
                    <CardTitle className="text-lg">Performance Analysis</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Annualized Return:</span>
                      <span className={`font-medium ${getReturnColor(simulation.performance.annualized_return)}`}>
                        {formatPercent(simulation.performance.annualized_return || 0)}
                      </span>
                    </div>
                    <Separator />
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Sharpe Ratio:</span>
                      <span className="font-medium text-foreground">
                        {(simulation.performance.sharpe_ratio || 0).toFixed(2)}
                      </span>
                    </div>
                    <Separator />
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Max Drawdown:</span>
                      <span className="font-medium text-destructive">
                        {formatPercent(simulation.performance.max_drawdown_pct || 0)}
                      </span>
                    </div>
                    <Separator />
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Win Rate:</span>
                      <span className="font-medium text-green-400">
                        {(simulation.performance.win_rate || 0).toFixed(1)}%
                      </span>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Market Timing Analysis */}
              {simulation.market_timing && (
                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardHeader>
                    <CardTitle className="text-lg">Market Timing Analysis</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Timing Score:</span>
                      <span className="text-lg font-bold text-blue-400">
                        {(simulation.market_timing.timing_score || 0).toFixed(0)}/100
                      </span>
                    </div>
                    <Separator />
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Avg MVRV at Purchase:</span>
                      <span className="font-medium text-foreground">
                        {(simulation.market_timing.avg_mvrv_at_purchase || 0).toFixed(2)}
                      </span>
                    </div>
                    <Separator />
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Avg Fear & Greed:</span>
                      <span className="font-medium text-foreground">
                        {(simulation.market_timing.avg_fear_greed_at_purchase || 0).toFixed(0)}
                      </span>
                    </div>
                    {simulation.market_timing.timing_analysis && (
                      <>
                        <Separator />
                        <div className="p-3 bg-muted/30 rounded-lg">
                          <p className="text-sm text-muted-foreground">{simulation.market_timing.timing_analysis}</p>
                        </div>
                      </>
                    )}
                  </CardContent>
                </Card>
              )}

              {/* Recommendations */}
              {simulation.recommendations && (
                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardHeader>
                    <CardTitle className="text-lg">Recommendations</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Market Condition:</span>
                      <span className="font-medium text-foreground">{simulation.recommendations.market_condition}</span>
                    </div>
                    <Separator />
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Risk Level:</span>
                      <Badge variant={getRiskVariant(simulation.recommendations.risk_level)} className="font-medium">
                        {simulation.recommendations.risk_level}
                      </Badge>
                    </div>
                    <Separator />
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Strategy:</span>
                      <span className="font-medium text-foreground">{simulation.recommendations.strategy}</span>
                    </div>
                    
                    {simulation.recommendations.suggestions && simulation.recommendations.suggestions.length > 0 && (
                      <>
                        <Separator />
                        <div>
                          <h4 className="text-sm font-medium text-foreground mb-3">Suggestions:</h4>
                          <ul className="space-y-2">
                            {simulation.recommendations.suggestions.map((suggestion, index) => (
                              <li key={index} className="text-sm text-muted-foreground flex items-start space-x-3">
                                <span className="text-accent mt-1">•</span>
                                <span>{suggestion}</span>
                              </li>
                            ))}
                          </ul>
                        </div>
                      </>
                    )}
                  </CardContent>
                </Card>
              )}

              {/* Chart Section */}
              {simulation?.chart_data && (
                <Card className="bg-card/50 backdrop-blur-sm">
                  <CardHeader>
                    <CardTitle className="text-lg">DCA Performance Over Time</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <DCAChart 
                      chartData={simulation.chart_data} 
                      title="Investment vs Portfolio Value"
                    />
                  </CardContent>
                </Card>
              )}
            </div>
          ) : (
            <Card className="bg-card/50 backdrop-blur-sm">
              <CardContent className="p-8 text-center">
                <ChartBarIcon className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium text-foreground mb-2">Ready to Simulate</h3>
                <p className="text-muted-foreground">Configure your DCA strategy and click simulate to see results</p>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Custom CSS for slider */}
        <style jsx>{`
          .slider-custom::-webkit-slider-thumb {
            appearance: none;
            height: 20px;
            width: 20px;
            border-radius: 50%;
            background: linear-gradient(45deg, #3b82f6, #8b5cf6);
            cursor: pointer;
            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
            border: 2px solid white;
            transition: all 0.2s ease;
          }
          
          .slider-custom::-webkit-slider-thumb:hover {
            transform: scale(1.1);
            box-shadow: 0 6px 20px rgba(59, 130, 246, 0.6);
          }
          
          .slider-custom::-moz-range-thumb {
            height: 20px;
            width: 20px;
            border-radius: 50%;
            background: linear-gradient(45deg, #3b82f6, #8b5cf6);
            cursor: pointer;
            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
            border: 2px solid white;
            transition: all 0.2s ease;
          }

          .scrollbar-thin::-webkit-scrollbar {
            width: 6px;
          }
          
          .scrollbar-track-gray-800::-webkit-scrollbar-track {
            background: #1f2937;
            border-radius: 3px;
          }
          
          .scrollbar-thumb-gray-600::-webkit-scrollbar-thumb {
            background: #4b5563;
            border-radius: 3px;
          }
          
          .scrollbar-thumb-gray-600::-webkit-scrollbar-thumb:hover {
            background: #6b7280;
          }
        `}</style>
      </div>
    </>
  )
}

export default DCASidebar