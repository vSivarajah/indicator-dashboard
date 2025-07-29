import { useState, memo } from 'react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, ReferenceLine } from 'recharts'
import { buildCardStyles, colors, typography, effects } from '../utils/designSystem'

const PortfolioPerformanceChart = memo(function PortfolioPerformanceChart({ 
  portfolioSummary,
  loading = false,
  title = "Portfolio Performance"
}) {
  const [selectedTimeframe, setSelectedTimeframe] = useState('1M')
  const timeframes = ['1D', '1W', '1M', '3M', '1Y', 'ALL']

  // Generate mock historical data based on current portfolio metrics
  // In real implementation, this would come from the backend
  const generateMockData = () => {
    if (!portfolioSummary) return []
    
    const currentValue = portfolioSummary.total_value || 0
    const totalPnL = portfolioSummary.total_pnl || 0
    const initialValue = currentValue - totalPnL
    
    const dataPoints = 30 // 30 days of data
    const data = []
    
    for (let i = dataPoints; i >= 0; i--) {
      const date = new Date()
      date.setDate(date.getDate() - i)
      
      // Simulate gradual value growth with some volatility
      const progress = (dataPoints - i) / dataPoints
      const baseValue = initialValue + (totalPnL * progress)
      const volatility = (Math.random() - 0.5) * 0.1 * baseValue
      const value = Math.max(0, baseValue + volatility)
      
      data.push({
        date: date.toISOString().split('T')[0],
        timestamp: date.getTime(),
        value: value,
        pnl: value - initialValue,
        pnlPercent: ((value - initialValue) / initialValue) * 100
      })
    }
    
    return data
  }

  const chartData = generateMockData()

  // Custom tooltip component
  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      const date = new Date(data.timestamp)
      const isPositive = data.pnl >= 0
      
      return (
        <div className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-3 shadow-lg`}>
          <div className="text-white font-medium mb-2">
            {date.toLocaleDateString('en-US', { 
              weekday: 'short', 
              year: 'numeric', 
              month: 'short', 
              day: 'numeric' 
            })}
          </div>
          <div className="space-y-1">
            <div className="text-sm text-gray-300">
              Value: ${data.value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className={`text-sm ${isPositive ? 'text-green-400' : 'text-red-400'}`}>
              P&L: {isPositive ? '+' : ''}${data.pnl.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className={`text-sm ${isPositive ? 'text-green-400' : 'text-red-400'}`}>
              {isPositive ? '+' : ''}{data.pnlPercent.toFixed(2)}%
            </div>
          </div>
        </div>
      )
    }
    return null
  }

  // Format X-axis labels
  const formatXAxisLabel = (tickItem) => {
    const date = new Date(tickItem)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }

  // Format Y-axis labels
  const formatYAxisLabel = (value) => {
    if (value >= 1000000) {
      return `$${(value / 1000000).toFixed(1)}M`
    } else if (value >= 1000) {
      return `$${(value / 1000).toFixed(1)}K`
    }
    return `$${value.toFixed(0)}`
  }

  // Loading skeleton
  if (loading) {
    return (
      <div className={buildCardStyles(false)}>
        <h3 className="text-lg font-semibold text-white mb-4">{title}</h3>
        <div className="h-80 flex items-center justify-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-orange-400"></div>
        </div>
      </div>
    )
  }

  // No data state
  if (!portfolioSummary || chartData.length === 0) {
    return (
      <div className={buildCardStyles(false)}>
        <h3 className="text-lg font-semibold text-white mb-4">{title}</h3>
        <div className="h-80 flex flex-col items-center justify-center text-gray-400">
          <div className="w-16 h-16 rounded-full bg-gray-700/50 flex items-center justify-center mb-4">
            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
            </svg>
          </div>
          <p className="text-center">No performance data available</p>
          <p className="text-sm text-gray-500 mt-1">Portfolio history will appear here</p>
        </div>
      </div>
    )
  }

  const currentValue = portfolioSummary.total_value || 0
  const totalPnL = portfolioSummary.total_pnl || 0
  const initialValue = currentValue - totalPnL
  const isPositive = totalPnL >= 0

  return (
    <div className={buildCardStyles(false)}>
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-white">{title}</h3>
          <div className="flex items-center space-x-4 mt-1">
            <div className="text-2xl font-bold text-white">
              ${currentValue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className={`text-sm font-medium ${isPositive ? 'text-green-400' : 'text-red-400'}`}>
              {isPositive ? '+' : ''}${totalPnL.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className={`text-sm ${isPositive ? 'text-green-400' : 'text-red-400'}`}>
              ({isPositive ? '+' : ''}{portfolioSummary.total_pnl_percent?.toFixed(2) || '0.00'}%)
            </div>
          </div>
        </div>
        
        {/* Timeframe selector */}
        <div className="flex space-x-1 bg-gray-900/50 rounded-lg p-1">
          {timeframes.map((tf) => (
            <button
              key={tf}
              onClick={() => setSelectedTimeframe(tf)}
              className={`px-3 py-1 text-xs font-medium rounded-md transition-all duration-300 ${
                tf === selectedTimeframe
                  ? 'bg-gradient-to-r from-orange-400 to-purple-500 text-white shadow-md scale-105'
                  : 'text-gray-400 hover:text-white hover:bg-gray-700/50 hover:scale-105'
              }`}
            >
              {tf}
            </button>
          ))}
        </div>
      </div>
      
      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="rgba(75, 85, 99, 0.3)" />
            <XAxis 
              dataKey="timestamp"
              type="number"
              scale="time"
              domain={['dataMin', 'dataMax']}
              tickFormatter={formatXAxisLabel}
              stroke="#9CA3AF"
              fontSize={12}
            />
            <YAxis 
              tickFormatter={formatYAxisLabel}
              stroke="#9CA3AF"
              fontSize={12}
            />
            <Tooltip content={<CustomTooltip />} />
            
            {/* Reference line at break-even point */}
            <ReferenceLine 
              y={initialValue} 
              stroke="#6B7280" 
              strokeDasharray="2 2"
              label={{ value: "Break-even", position: "right" }}
            />
            
            <Line
              type="monotone"
              dataKey="value"
              stroke={isPositive ? "#10B981" : "#EF4444"}
              strokeWidth={2}
              dot={false}
              activeDot={{ 
                r: 4, 
                fill: isPositive ? "#10B981" : "#EF4444",
                stroke: "#fff",
                strokeWidth: 2
              }}
              animationDuration={1000}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Performance summary */}
      <div className="mt-4 pt-4 border-t border-gray-700/30">
        <div className="grid grid-cols-3 gap-4 text-sm">
          <div>
            <span className="text-gray-400">Initial Value:</span>
            <span className="text-white ml-2 font-medium">
              ${initialValue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </span>
          </div>
          <div>
            <span className="text-gray-400">Current Value:</span>
            <span className="text-white ml-2 font-medium">
              ${currentValue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </span>
          </div>
          <div>
            <span className="text-gray-400">Time Period:</span>
            <span className="text-white ml-2 font-medium">{selectedTimeframe}</span>
          </div>
        </div>
      </div>
    </div>
  )
})

export default PortfolioPerformanceChart