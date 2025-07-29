import { useState, memo } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { buildCardStyles, colors, typography, effects } from '../utils/designSystem'

const HoldingsComparisonChart = memo(function HoldingsComparisonChart({ 
  holdings,
  loading = false,
  title = "Holdings Comparison"
}) {
  const [sortBy, setSortBy] = useState('value') // value, pnl, pnlPercent
  const [sortOrder, setSortOrder] = useState('desc') // asc, desc

  // Process and sort holdings data
  const processHoldingsData = () => {
    if (!holdings || holdings.length === 0) return []
    
    return holdings
      .map(holding => ({
        symbol: holding.symbol,
        name: holding.symbol, // In real app, this would be the full name
        value: holding.value || 0,
        costBasis: (holding.amount * holding.average_price) || 0,
        pnl: holding.pnl || 0,
        pnlPercent: holding.pnl_percent || 0,
        amount: holding.amount || 0,
        currentPrice: holding.current_price || 0,
        averagePrice: holding.average_price || 0
      }))
      .sort((a, b) => {
        const modifier = sortOrder === 'desc' ? -1 : 1
        return (a[sortBy] - b[sortBy]) * modifier
      })
  }

  const chartData = processHoldingsData()

  // Get bar color based on P&L
  const getBarColor = (pnl) => {
    if (pnl > 0) return '#10B981' // Green for profits
    if (pnl < 0) return '#EF4444' // Red for losses
    return '#6B7280' // Gray for break-even
  }

  // Custom tooltip component
  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      const isPositive = data.pnl >= 0
      
      return (
        <div className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-3 shadow-lg`}>
          <div className="text-white font-medium mb-2">{data.symbol}</div>
          <div className="space-y-1">
            <div className="text-sm text-gray-300">
              Amount: {data.amount.toFixed(8)}
            </div>
            <div className="text-sm text-gray-300">
              Current Price: ${data.currentPrice.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className="text-sm text-gray-300">
              Avg Price: ${data.averagePrice.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className="text-sm text-gray-300">
              Current Value: ${data.value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </div>
            <div className="text-sm text-gray-400">
              Cost Basis: ${data.costBasis.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
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

  // Format Y-axis labels based on current sort
  const formatYAxisLabel = (value) => {
    if (sortBy === 'pnlPercent') {
      return `${value.toFixed(0)}%`
    } else if (value >= 1000000) {
      return `$${(value / 1000000).toFixed(1)}M`
    } else if (value >= 1000) {
      return `$${(value / 1000).toFixed(1)}K`
    }
    return `$${value.toFixed(0)}`
  }

  // Handle sort option change
  const handleSortChange = (newSortBy) => {
    if (sortBy === newSortBy) {
      setSortOrder(sortOrder === 'desc' ? 'asc' : 'desc')
    } else {
      setSortBy(newSortBy)
      setSortOrder('desc')
    }
  }

  // Sort options
  const sortOptions = [
    { key: 'value', label: 'Value', icon: '💰' },
    { key: 'pnl', label: 'P&L ($)', icon: '📈' },
    { key: 'pnlPercent', label: 'P&L (%)', icon: '📊' }
  ]

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
  if (!holdings || holdings.length === 0) {
    return (
      <div className={buildCardStyles(false)}>
        <h3 className="text-lg font-semibold text-white mb-4">{title}</h3>
        <div className="h-80 flex flex-col items-center justify-center text-gray-400">
          <div className="w-16 h-16 rounded-full bg-gray-700/50 flex items-center justify-center mb-4">
            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
            </svg>
          </div>
          <p className="text-center">No holdings to compare</p>
          <p className="text-sm text-gray-500 mt-1">Add some holdings to see comparison</p>
        </div>
      </div>
    )
  }

  return (
    <div className={buildCardStyles(false)}>
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-white">{title}</h3>
          <div className="text-sm text-gray-400 mt-1">
            Sorted by {sortOptions.find(opt => opt.key === sortBy)?.label} ({sortOrder === 'desc' ? 'highest first' : 'lowest first'})
          </div>
        </div>
        
        {/* Sort options */}
        <div className="flex space-x-1 bg-gray-900/50 rounded-lg p-1">
          {sortOptions.map((option) => (
            <button
              key={option.key}
              onClick={() => handleSortChange(option.key)}
              className={`px-3 py-1 text-xs font-medium rounded-md transition-all duration-300 flex items-center space-x-1 ${
                sortBy === option.key
                  ? 'bg-gradient-to-r from-orange-400 to-purple-500 text-white shadow-md scale-105'
                  : 'text-gray-400 hover:text-white hover:bg-gray-700/50 hover:scale-105'
              }`}
            >
              <span>{option.icon}</span>
              <span>{option.label}</span>
              {sortBy === option.key && (
                <span className="text-xs">
                  {sortOrder === 'desc' ? '↓' : '↑'}
                </span>
              )}
            </button>
          ))}
        </div>
      </div>
      
      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart 
            data={chartData}
            margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke="rgba(75, 85, 99, 0.3)" />
            <XAxis 
              dataKey="symbol"
              stroke="#9CA3AF"
              fontSize={12}
            />
            <YAxis 
              tickFormatter={formatYAxisLabel}
              stroke="#9CA3AF"
              fontSize={12}
            />
            <Tooltip content={<CustomTooltip />} />
            <Bar 
              dataKey={sortBy}
              radius={[4, 4, 0, 0]}
              animationDuration={800}
            >
              {chartData.map((entry, index) => (
                <Cell 
                  key={`cell-${index}`}
                  fill={getBarColor(entry.pnl)}
                  className="hover:opacity-80 transition-opacity cursor-pointer"
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Holdings summary */}
      <div className="mt-4 pt-4 border-t border-gray-700/30">
        <div className="grid grid-cols-3 gap-4 text-sm">
          <div>
            <span className="text-gray-400">Total Holdings:</span>
            <span className="text-white ml-2 font-medium">{chartData.length}</span>
          </div>
          <div>
            <span className="text-gray-400">Profitable:</span>
            <span className="text-green-400 ml-2 font-medium">
              {chartData.filter(h => h.pnl > 0).length}
            </span>
          </div>
          <div>
            <span className="text-gray-400">Losing:</span>
            <span className="text-red-400 ml-2 font-medium">
              {chartData.filter(h => h.pnl < 0).length}
            </span>
          </div>
        </div>
      </div>
    </div>
  )
})

export default HoldingsComparisonChart