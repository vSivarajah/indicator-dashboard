import { useState, useRef, useEffect } from 'react'

function DCAChart({ chartData, title = "DCA Performance" }) {
  const [activeChart, setActiveChart] = useState('value')
  const [hoveredPoint, setHoveredPoint] = useState(null)
  const [tooltip, setTooltip] = useState({ show: false, x: 0, y: 0, data: null })
  const chartRef = useRef(null)
  const [chartAnimation, setChartAnimation] = useState(false)

  useEffect(() => {
    if (chartData) {
      setChartAnimation(true)
      const timer = setTimeout(() => setChartAnimation(false), 1000)
      return () => clearTimeout(timer)
    }
  }, [chartData, activeChart])

  // Loading skeleton component
  const LoadingSkeleton = () => (
    <div className="bg-gradient-to-br from-gray-800/50 to-gray-900/30 rounded-2xl p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="h-6 bg-gradient-to-r from-gray-700 to-gray-600 rounded-lg w-48 animate-pulse"></div>
        <div className="flex space-x-2">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="h-8 w-16 bg-gradient-to-r from-gray-700 to-gray-600 rounded-lg animate-pulse" style={{ animationDelay: `${i * 0.1}s` }}></div>
          ))}
        </div>
      </div>
      <div className="h-64 bg-gradient-to-br from-gray-800/70 to-gray-700/50 rounded-xl relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-transparent via-gray-600/20 to-transparent animate-shimmer"></div>
        {/* Skeleton chart lines */}
        <div className="absolute bottom-8 left-4 right-4 space-y-2">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="h-1 bg-gradient-to-r from-blue-500/30 to-purple-500/30 rounded-full animate-pulse" style={{ width: `${60 + i * 10}%`, animationDelay: `${i * 0.2}s` }}></div>
          ))}
        </div>
      </div>
    </div>
  )

  if (!chartData) {
    return <LoadingSkeleton />
  }

  const chartTypes = [
    { key: 'value', label: 'Portfolio Value', color: 'text-blue-400', gradient: 'from-blue-500 to-cyan-400' },
    { key: 'investment', label: 'Total Invested', color: 'text-green-400', gradient: 'from-green-500 to-emerald-400' },
    { key: 'price', label: 'Asset Price', color: 'text-yellow-400', gradient: 'from-yellow-500 to-orange-400' },
    { key: 'pnl', label: 'P&L', color: 'text-purple-400', gradient: 'from-purple-500 to-pink-400' }
  ]

  const getChartData = (type) => {
    switch (type) {
      case 'value':
        return chartData.value_history || []
      case 'investment':
        return chartData.investment_history || []
      case 'price':
        return chartData.price_history || []
      case 'pnl':
        return chartData.pnl_history || []
      default:
        return []
    }
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    })
  }

  const formatValue = (value, type) => {
    if (type === 'price') {
      return `$${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}`
    }
    return `$${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
  }

  const getChartColor = (type) => {
    const colors = {
      value: '#60A5FA',
      investment: '#34D399',
      price: '#FBBF24',
      pnl: '#A78BFA'
    }
    return colors[type] || '#60A5FA'
  }

  const handleMouseMove = (event, point, index) => {
    if (!chartRef.current) return
    
    const rect = chartRef.current.getBoundingClientRect()
    const x = event.clientX - rect.left
    const y = event.clientY - rect.top
    
    setTooltip({
      show: true,
      x: Math.min(x + 10, rect.width - 200),
      y: Math.max(y - 10, 0),
      data: {
        date: point.date,
        value: point.value,
        type: activeChart
      }
    })
    setHoveredPoint(index)
  }

  const handleMouseLeave = () => {
    setTooltip({ show: false, x: 0, y: 0, data: null })
    setHoveredPoint(null)
  }

  const data = getChartData(activeChart)
  const maxValue = Math.max(...data.map(d => d.value))
  const minValue = Math.min(...data.map(d => d.value))
  const range = maxValue - minValue

  return (
    <div className="bg-gradient-to-br from-gray-800/50 to-gray-900/30 rounded-2xl p-6 backdrop-blur-sm border border-gray-700/30 shadow-xl">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">{title}</h3>
        <div className="flex space-x-1 bg-gray-800/50 rounded-xl p-1">
          {chartTypes.map((chart) => (
            <button
              key={chart.key}
              onClick={() => setActiveChart(chart.key)}
              className={`relative px-4 py-2 text-sm font-medium rounded-lg transition-all duration-300 transform hover:scale-105 ${
                activeChart === chart.key
                  ? `bg-gradient-to-r ${chart.gradient} text-white shadow-lg`
                  : 'text-gray-400 hover:text-white hover:bg-gray-700/50'
              }`}
            >
              {activeChart === chart.key && (
                <div className={`absolute inset-0 bg-gradient-to-r ${chart.gradient} opacity-20 rounded-lg animate-pulse`}></div>
              )}
              <span className="relative z-10">{chart.label}</span>
            </button>
          ))}
        </div>
      </div>

      {data.length > 0 ? (
        <div className="relative" ref={chartRef} onMouseLeave={handleMouseLeave}>
          {/* Enhanced SVG Chart */}
          <div className="h-80 bg-gradient-to-br from-gray-900/70 to-gray-800/50 rounded-xl p-6 overflow-hidden border border-gray-700/30 shadow-inner">
            <svg
              width="100%"
              height="100%"
              viewBox="0 0 800 250"
              className="overflow-visible"
            >
              {/* Enhanced Grid and Gradients */}
              <defs>
                <pattern id="modernGrid" width="40" height="25" patternUnits="userSpaceOnUse">
                  <path d="M 40 0 L 0 0 0 25" fill="none" stroke="#374151" strokeWidth="0.3" opacity="0.5"/>
                </pattern>
                <linearGradient id="chartGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stopColor={getChartColor(activeChart)} stopOpacity="0.3"/>
                  <stop offset="100%" stopColor={getChartColor(activeChart)} stopOpacity="0.05"/>
                </linearGradient>
                <filter id="glow">
                  <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
                  <feMerge>
                    <feMergeNode in="coloredBlur"/>
                    <feMergeNode in="SourceGraphic"/>
                  </feMerge>
                </filter>
              </defs>
              <rect width="100%" height="100%" fill="url(#modernGrid)" />

              {/* Gradient Fill Area */}
              {data.length > 1 && (
                <polygon
                  fill="url(#chartGradient)"
                  points={[
                    `0,250`,
                    ...data.map((point, index) => {
                      const x = (index / (data.length - 1)) * 800
                      const y = 250 - ((point.value - minValue) / range) * 200
                      return `${x},${y}`
                    }),
                    `800,250`
                  ].join(' ')}
                  className={chartAnimation ? 'animate-pulse' : ''}
                />
              )}
              
              {/* Enhanced Chart Line */}
              {data.length > 1 && (
                <polyline
                  fill="none"
                  stroke={getChartColor(activeChart)}
                  strokeWidth="3"
                  filter="url(#glow)"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  points={data.map((point, index) => {
                    const x = (index / (data.length - 1)) * 800
                    const y = 250 - ((point.value - minValue) / range) * 200
                    return `${x},${y}`
                  }).join(' ')}
                  className={chartAnimation ? 'animate-pulse' : ''}
                  style={{
                    strokeDasharray: chartAnimation ? '10 5' : 'none',
                    animation: chartAnimation ? 'dashArray 2s ease-in-out' : 'none'
                  }}
                />
              )}

              {/* Interactive Data Points */}
              {data.map((point, index) => {
                const x = (index / (data.length - 1)) * 800
                const y = 250 - ((point.value - minValue) / range) * 200
                const isHovered = hoveredPoint === index
                return (
                  <g key={index}>
                    {/* Hover area */}
                    <circle
                      cx={x}
                      cy={y}
                      r="12"
                      fill="transparent"
                      className="cursor-pointer"
                      onMouseEnter={(e) => handleMouseMove(e, point, index)}
                    />
                    
                    {/* Glow effect for hovered point */}
                    {isHovered && (
                      <circle
                        cx={x}
                        cy={y}
                        r="8"
                        fill={getChartColor(activeChart)}
                        opacity="0.3"
                        className="animate-ping"
                      />
                    )}
                    
                    {/* Data point */}
                    <circle
                      cx={x}
                      cy={y}
                      r={isHovered ? "6" : "4"}
                      fill={getChartColor(activeChart)}
                      stroke="white"
                      strokeWidth="2"
                      className="transition-all duration-200 cursor-pointer filter drop-shadow-lg"
                      style={{
                        transform: isHovered ? 'scale(1.3)' : 'scale(1)',
                        transformOrigin: `${x}px ${y}px`
                      }}
                    />
                  </g>
                )
              })}
            </svg>
            
            {/* Interactive Tooltip */}
            {tooltip.show && tooltip.data && (
              <div 
                className="absolute z-20 bg-gray-900/95 backdrop-blur-sm border border-gray-600/50 rounded-xl p-4 shadow-2xl pointer-events-none transform -translate-x-1/2 -translate-y-full"
                style={{ left: tooltip.x, top: tooltip.y }}
              >
                <div className="text-xs text-gray-400 mb-1">{formatDate(tooltip.data.date)}</div>
                <div className="text-lg font-bold text-white">{formatValue(tooltip.data.value, tooltip.data.type)}</div>
                <div className="text-xs text-gray-400 capitalize">{tooltip.data.type} Value</div>
                <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-full">
                  <div className="w-2 h-2 bg-gray-900 rotate-45 border-r border-b border-gray-600/50"></div>
                </div>
              </div>
            )}
          </div>

          {/* Enhanced Chart Legend */}
          <div className="mt-6 flex items-center justify-between text-sm">
            <div className="flex items-center space-x-2">
              <div className="w-3 h-3 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
              <span className="text-gray-400">{formatDate(data[0]?.date)}</span>
            </div>
            <div className="flex items-center space-x-6 bg-gray-800/50 rounded-lg px-4 py-2">
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 rounded-full bg-green-400"></div>
                <span className="text-gray-300">Min: <span className="font-medium text-white">{formatValue(minValue, activeChart)}</span></span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 rounded-full bg-red-400"></div>
                <span className="text-gray-300">Max: <span className="font-medium text-white">{formatValue(maxValue, activeChart)}</span></span>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-gray-400">{formatDate(data[data.length - 1]?.date)}</span>
              <div className="w-3 h-3 rounded-full bg-gradient-to-r from-gray-500 to-gray-600"></div>
            </div>
          </div>
        </div>
      ) : (
        <div className="flex items-center justify-center h-80 text-gray-400">
          <div className="text-center">
            <div className="w-16 h-16 mx-auto mb-4 opacity-30">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M3 3v18h18V3H3zm16 16H5V5h14v14zm-8-2l3-4 2 3h-6l1-1z"/>
              </svg>
            </div>
            <p className="text-lg font-medium text-gray-300 mb-2">No Data Available</p>
            <p className="text-sm text-gray-500">No data points available for the selected chart type</p>
          </div>
        </div>
      )}
    </div>
  )
}

// Add custom CSS animations
const style = document.createElement('style')
style.textContent = `
  @keyframes dashArray {
    0% { stroke-dasharray: 0 1000; }
    100% { stroke-dasharray: 1000 0; }
  }
  
  @keyframes shimmer {
    0% { transform: translateX(-100%); }
    100% { transform: translateX(100%); }
  }
  
  .animate-shimmer {
    animation: shimmer 2s infinite;
  }
`

if (typeof document !== 'undefined' && !document.querySelector('#dca-chart-styles')) {
  style.id = 'dca-chart-styles'
  document.head.appendChild(style)
}

export default DCAChart