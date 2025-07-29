import { memo } from 'react'
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts'
import { buildCardStyles, colors, typography, effects } from '../utils/designSystem'

const PortfolioAllocationChart = memo(function PortfolioAllocationChart({ 
  allocationData, 
  loading = false,
  title = "Portfolio Allocation"
}) {
  // Default colors for different assets
  const defaultColors = [
    '#F97316', // Orange
    '#8B5CF6', // Purple  
    '#10B981', // Green
    '#EF4444', // Red
    '#3B82F6', // Blue
    '#F59E0B', // Amber
    '#EC4899', // Pink
    '#06B6D4', // Cyan
    '#84CC16', // Lime
    '#F43F5E'  // Rose
  ]

  // Format data for Recharts
  const chartData = allocationData?.map((item, index) => ({
    name: item.name || item.symbol,
    symbol: item.symbol,
    value: item.percentage,
    amount: item.value,
    color: item.color || defaultColors[index % defaultColors.length]
  })) || []

  // Custom tooltip component
  const CustomTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-3 shadow-lg`}>
          <div className="text-white font-medium">{data.name}</div>
          <div className="text-sm text-gray-400">Symbol: {data.symbol}</div>
          <div className="text-sm text-gray-300">
            Value: ${data.amount?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
          </div>
          <div className="text-sm text-gray-300">
            Percentage: {data.value.toFixed(1)}%
          </div>
        </div>
      )
    }
    return null
  }

  // Custom legend component
  const CustomLegend = ({ payload }) => {
    return (
      <div className="flex flex-wrap justify-center gap-4 mt-4">
        {payload?.map((entry, index) => (
          <div key={index} className="flex items-center space-x-2">
            <div 
              className="w-3 h-3 rounded-full"
              style={{ backgroundColor: entry.color }}
            />
            <span className="text-sm text-gray-300">
              {entry.payload.symbol} ({entry.payload.value.toFixed(1)}%)
            </span>
          </div>
        ))}
      </div>
    )
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
  if (!allocationData || allocationData.length === 0) {
    return (
      <div className={buildCardStyles(false)}>
        <h3 className="text-lg font-semibold text-white mb-4">{title}</h3>
        <div className="h-80 flex flex-col items-center justify-center text-gray-400">
          <div className="w-16 h-16 rounded-full bg-gray-700/50 flex items-center justify-center mb-4">
            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
            </svg>
          </div>
          <p className="text-center">No allocation data available</p>
          <p className="text-sm text-gray-500 mt-1">Add holdings to see portfolio allocation</p>
        </div>
      </div>
    )
  }

  return (
    <div className={buildCardStyles(false)}>
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-white">{title}</h3>
        <div className="text-sm text-gray-400">
          {chartData.length} asset{chartData.length !== 1 ? 's' : ''}
        </div>
      </div>
      
      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              innerRadius={40}
              outerRadius={100}
              paddingAngle={2}
              dataKey="value"
              animationBegin={0}
              animationDuration={800}
            >
              {chartData.map((entry, index) => (
                <Cell 
                  key={`cell-${index}`} 
                  fill={entry.color}
                  className="hover:opacity-80 transition-opacity cursor-pointer"
                />
              ))}
            </Pie>
            <Tooltip content={<CustomTooltip />} />
            <Legend content={<CustomLegend />} />
          </PieChart>
        </ResponsiveContainer>
      </div>

      {/* Portfolio summary */}
      <div className="mt-4 pt-4 border-t border-gray-700/30">
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-400">Total Value:</span>
            <span className="text-white ml-2 font-medium">
              ${allocationData.reduce((sum, item) => sum + (item.value || 0), 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </span>
          </div>
          <div>
            <span className="text-gray-400">Largest holding:</span>
            <span className="text-white ml-2 font-medium">
              {chartData.length > 0 ? chartData.reduce((max, item) => item.value > max.value ? item : max, chartData[0]).symbol : 'N/A'}
            </span>
          </div>
        </div>
      </div>
    </div>
  )
})

export default PortfolioAllocationChart