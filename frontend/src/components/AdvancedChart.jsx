import React, { useState, useEffect, useRef, useMemo } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ComposedChart,
  Bar,
  Area,
  AreaChart,
  ScatterChart,
  Scatter,
  ReferenceLine,
  ReferenceArea,
  Brush
} from 'recharts'
import { 
  ChartBarIcon, 
  CogIcon, 
  PlusIcon, 
  MinusIcon,
  ArrowsPointingOutIcon,
  CalendarIcon
} from '@heroicons/react/24/outline'
import { cn } from '../lib/utils'
import { typography, colors, spacing, effects, buildCardStyles } from '../utils/designSystem'

// Timeframe options
const TIMEFRAMES = [
  { id: '5m', label: '5m', minutes: 5 },
  { id: '15m', label: '15m', minutes: 15 },
  { id: '1h', label: '1H', minutes: 60 },
  { id: '4h', label: '4H', minutes: 240 },
  { id: '1d', label: '1D', minutes: 1440 },
  { id: '1w', label: '1W', minutes: 10080 },
  { id: '1M', label: '1M', minutes: 43200 }
]

// Chart types
const CHART_TYPES = [
  { id: 'line', label: 'Line', icon: '📈' },
  { id: 'candle', label: 'Candles', icon: '🕯️' },
  { id: 'area', label: 'Area', icon: '🏔️' },
  { id: 'volume', label: 'Volume', icon: '📊' }
]

// Technical indicators
const INDICATORS = [
  { id: 'sma', label: 'SMA', name: 'Simple Moving Average', periods: [20, 50, 200] },
  { id: 'ema', label: 'EMA', name: 'Exponential Moving Average', periods: [12, 26, 50] },
  { id: 'bollinger', label: 'Bollinger Bands', name: 'Bollinger Bands', period: 20 },
  { id: 'rsi', label: 'RSI', name: 'Relative Strength Index', period: 14 },
  { id: 'macd', label: 'MACD', name: 'MACD', fast: 12, slow: 26, signal: 9 },
  { id: 'volume', label: 'Volume', name: 'Volume Profile' },
  { id: 'support_resistance', label: 'S/R', name: 'Support & Resistance' }
]

const AdvancedChart = ({
  data = [],
  title = "Advanced Chart",
  symbol = "BTC",
  className = "",
  height = 500,
  showControls = true,
  onTimeframeChange,
  onIndicatorChange
}) => {
  const [selectedTimeframe, setSelectedTimeframe] = useState('1d')
  const [selectedChartType, setSelectedChartType] = useState('line')
  const [activeIndicators, setActiveIndicators] = useState(['sma'])
  const [showSettings, setShowSettings] = useState(false)
  const [zoomDomain, setZoomDomain] = useState(null)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const chartContainerRef = useRef(null)

  // Generate sample data if none provided
  const chartData = useMemo(() => {
    if (data.length > 0) return data
    
    // Generate realistic sample data
    const samples = 100
    const basePrice = 45000
    const sampleData = []
    
    for (let i = 0; i < samples; i++) {
      const timestamp = new Date(Date.now() - (samples - i) * 1000 * 60 * 60).getTime()
      const randomFactor = 1 + (Math.sin(i * 0.1) * 0.02) + (Math.random() - 0.5) * 0.01
      const price = basePrice * Math.pow(randomFactor, i / 10)
      const volume = 1000000 + Math.random() * 2000000
      
      sampleData.push({
        timestamp,
        time: new Date(timestamp).toLocaleDateString(),
        price: parseFloat(price.toFixed(2)),
        open: price * 0.999,
        high: price * 1.002,
        low: price * 0.998,
        close: price,
        volume: Math.round(volume),
        sma20: calculateSMA(sampleData, 20, price),
        sma50: calculateSMA(sampleData, 50, price),
        ema12: calculateEMA(sampleData, 12, price),
        rsi: calculateRSI(sampleData, 14, price)
      })
    }
    
    return sampleData
  }, [data])

  // Calculate technical indicators
  const enhancedData = useMemo(() => {
    return chartData.map((item, index) => {
      const enhanced = { ...item }
      
      // Add Bollinger Bands
      if (activeIndicators.includes('bollinger')) {
        const bb = calculateBollingerBands(chartData, index, 20)
        enhanced.bb_upper = bb.upper
        enhanced.bb_middle = bb.middle
        enhanced.bb_lower = bb.lower
      }
      
      // Add MACD
      if (activeIndicators.includes('macd')) {
        const macd = calculateMACD(chartData, index)
        enhanced.macd_line = macd.macd
        enhanced.macd_signal = macd.signal
        enhanced.macd_histogram = macd.histogram
      }
      
      return enhanced
    })
  }, [chartData, activeIndicators])

  const handleTimeframeChange = (timeframe) => {
    setSelectedTimeframe(timeframe)
    onTimeframeChange?.(timeframe)
  }

  const toggleIndicator = (indicatorId) => {
    setActiveIndicators(prev => 
      prev.includes(indicatorId)
        ? prev.filter(id => id !== indicatorId)
        : [...prev, indicatorId]
    )
  }

  const handleZoom = (zoomData) => {
    if (zoomData) {
      setZoomDomain({
        left: zoomData.startIndex,
        right: zoomData.endIndex
      })
    } else {
      setZoomDomain(null)
    }
  }

  const toggleFullscreen = () => {
    if (!isFullscreen) {
      if (chartContainerRef.current?.requestFullscreen) {
        chartContainerRef.current.requestFullscreen()
      }
    } else {
      if (document.exitFullscreen) {
        document.exitFullscreen()
      }
    }
    setIsFullscreen(!isFullscreen)
  }

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className={cn(
          buildCardStyles(true, 'default'),
          "p-4 border border-gray-700/50"
        )}>
          <p className={cn(typography.bodySmall, "text-gray-300 mb-2")}>
            {new Date(label).toLocaleString()}
          </p>
          <div className="space-y-1">
            <p className={cn(typography.bodySmall, "text-blue-400")}>
              Price: ${data.price?.toLocaleString()}
            </p>
            {data.volume && (
              <p className={cn(typography.bodySmall, "text-gray-400")}>
                Volume: {(data.volume / 1000000).toFixed(2)}M
              </p>
            )}
            {activeIndicators.includes('rsi') && data.rsi && (
              <p className={cn(typography.bodySmall, "text-purple-400")}>
                RSI: {data.rsi.toFixed(2)}
              </p>
            )}
          </div>
        </div>
      )
    }
    return null
  }

  const renderChart = () => {
    const chartProps = {
      data: enhancedData,
      margin: { top: 20, right: 30, left: 20, bottom: 60 }
    }

    switch (selectedChartType) {
      case 'area':
        return (
          <AreaChart {...chartProps}>
            <defs>
              <linearGradient id="priceGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#3B82F6" stopOpacity={0.3}/>
                <stop offset="95%" stopColor="#3B82F6" stopOpacity={0}/>
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
            <XAxis 
              dataKey="time" 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
            />
            <YAxis 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
              domain={['dataMin * 0.98', 'dataMax * 1.02']}
            />
            <Tooltip content={<CustomTooltip />} />
            <Area
              type="monotone"
              dataKey="price"
              stroke="#3B82F6"
              strokeWidth={2}
              fill="url(#priceGradient)"
            />
            {renderIndicators()}
          </AreaChart>
        )

      case 'candle':
        return (
          <ComposedChart {...chartProps}>
            <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
            <XAxis 
              dataKey="time" 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
            />
            <YAxis 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
              domain={['dataMin * 0.98', 'dataMax * 1.02']}
            />
            <Tooltip content={<CustomTooltip />} />
            {/* Simplified candlestick using bars */}
            <Bar dataKey="high" fill="#10B981" opacity={0.3} />
            <Bar dataKey="low" fill="#EF4444" opacity={0.3} />
            <Line
              type="monotone"
              dataKey="close"
              stroke="#F59E0B"
              strokeWidth={1}
              dot={false}
            />
            {renderIndicators()}
          </ComposedChart>
        )

      case 'volume':
        return (
          <ComposedChart {...chartProps}>
            <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
            <XAxis 
              dataKey="time" 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
            />
            <YAxis 
              yAxisId="price"
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
              domain={['dataMin * 0.98', 'dataMax * 1.02']}
            />
            <YAxis 
              yAxisId="volume"
              orientation="right"
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
            />
            <Tooltip content={<CustomTooltip />} />
            <Bar 
              yAxisId="volume"
              dataKey="volume" 
              fill="#6B7280" 
              opacity={0.3}
            />
            <Line
              yAxisId="price"
              type="monotone"
              dataKey="price"
              stroke="#3B82F6"
              strokeWidth={2}
              dot={false}
            />
            {renderIndicators()}
          </ComposedChart>
        )

      default:
        return (
          <LineChart {...chartProps}>
            <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
            <XAxis 
              dataKey="time" 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
            />
            <YAxis 
              stroke="#9CA3AF"
              fontSize={12}
              tickLine={false}
              domain={['dataMin * 0.98', 'dataMax * 1.02']}
            />
            <Tooltip content={<CustomTooltip />} />
            <Line
              type="monotone"
              dataKey="price"
              stroke="#3B82F6"
              strokeWidth={2}
              dot={false}
            />
            {renderIndicators()}
            {zoomDomain && (
              <Brush
                dataKey="time"
                height={30}
                stroke="#3B82F6"
                startIndex={zoomDomain.left}
                endIndex={zoomDomain.right}
              />
            )}
          </LineChart>
        )
    }
  }

  const renderIndicators = () => {
    const indicators = []

    if (activeIndicators.includes('sma')) {
      indicators.push(
        <Line
          key="sma20"
          type="monotone"
          dataKey="sma20"
          stroke="#F59E0B"
          strokeWidth={1}
          strokeDasharray="5 5"
          dot={false}
        />,
        <Line
          key="sma50"
          type="monotone"
          dataKey="sma50"
          stroke="#EF4444"
          strokeWidth={1}
          strokeDasharray="5 5"
          dot={false}
        />
      )
    }

    if (activeIndicators.includes('ema')) {
      indicators.push(
        <Line
          key="ema12"
          type="monotone"
          dataKey="ema12"
          stroke="#8B5CF6"
          strokeWidth={1}
          dot={false}
        />
      )
    }

    if (activeIndicators.includes('bollinger')) {
      indicators.push(
        <Line
          key="bb_upper"
          type="monotone"
          dataKey="bb_upper"
          stroke="#10B981"
          strokeWidth={1}
          strokeOpacity={0.6}
          dot={false}
        />,
        <Line
          key="bb_lower"
          type="monotone"
          dataKey="bb_lower"
          stroke="#10B981"
          strokeWidth={1}
          strokeOpacity={0.6}
          dot={false}
        />
      )
    }

    return indicators
  }

  return (
    <div 
      ref={chartContainerRef}
      className={cn(
        buildCardStyles(true, 'default'),
        "p-6",
        isFullscreen && "fixed inset-0 z-50 p-8",
        className
      )}
    >
      {/* Chart Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className={cn(typography.cardTitle, colors.textPrimary)}>
            {title}
          </h3>
          <p className={cn(typography.bodySmall, colors.textSecondary)}>
            {symbol} • {selectedTimeframe.toUpperCase()}
          </p>
        </div>
        
        {showControls && (
          <div className="flex items-center space-x-2">
            <button
              onClick={toggleFullscreen}
              className={cn(
                "p-2 rounded-lg",
                "hover:bg-gray-700/50 transition-colors",
                colors.textSecondary,
                "hover:text-white"
              )}
            >
              <ArrowsPointingOutIcon className="w-4 h-4" />
            </button>
            <button
              onClick={() => setShowSettings(!showSettings)}
              className={cn(
                "p-2 rounded-lg",
                "hover:bg-gray-700/50 transition-colors",
                colors.textSecondary,
                "hover:text-white"
              )}
            >
              <CogIcon className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>

      {/* Chart Controls */}
      {showControls && (
        <div className="mb-6 space-y-4">
          {/* Timeframe Selector */}
          <div className="flex items-center space-x-2">
            <span className={cn(typography.bodySmall, colors.textSecondary)}>
              Timeframe:
            </span>
            <div className="flex space-x-1">
              {TIMEFRAMES.map((tf) => (
                <button
                  key={tf.id}
                  onClick={() => handleTimeframeChange(tf.id)}
                  className={cn(
                    "px-3 py-1 rounded text-xs transition-colors",
                    selectedTimeframe === tf.id
                      ? "bg-blue-600 text-white"
                      : "bg-gray-700/50 text-gray-300 hover:bg-gray-700"
                  )}
                >
                  {tf.label}
                </button>
              ))}
            </div>
          </div>

          {/* Chart Type Selector */}
          <div className="flex items-center space-x-2">
            <span className={cn(typography.bodySmall, colors.textSecondary)}>
              Type:
            </span>
            <div className="flex space-x-1">
              {CHART_TYPES.map((type) => (
                <button
                  key={type.id}
                  onClick={() => setSelectedChartType(type.id)}
                  className={cn(
                    "px-3 py-1 rounded text-xs transition-colors flex items-center space-x-1",
                    selectedChartType === type.id
                      ? "bg-blue-600 text-white"
                      : "bg-gray-700/50 text-gray-300 hover:bg-gray-700"
                  )}
                >
                  <span>{type.icon}</span>
                  <span>{type.label}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Indicators */}
          {showSettings && (
            <div className="p-4 bg-gray-800/50 rounded-lg">
              <h4 className={cn(typography.bodyMedium, colors.textPrimary, "mb-3")}>
                Technical Indicators
              </h4>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
                {INDICATORS.map((indicator) => (
                  <label
                    key={indicator.id}
                    className="flex items-center space-x-2 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={activeIndicators.includes(indicator.id)}
                      onChange={() => toggleIndicator(indicator.id)}
                      className="rounded border-gray-600 bg-gray-700 text-blue-600"
                    />
                    <span className={cn(typography.bodySmall, colors.textSecondary)}>
                      {indicator.label}
                    </span>
                  </label>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Main Chart */}
      <div style={{ height: isFullscreen ? 'calc(100vh - 200px)' : height }}>
        <ResponsiveContainer width="100%" height="100%">
          {renderChart()}
        </ResponsiveContainer>
      </div>

      {/* Chart Legend */}
      <div className="mt-4 flex flex-wrap gap-4 text-xs">
        <div className="flex items-center space-x-1">
          <div className="w-3 h-0.5 bg-blue-500"></div>
          <span className={colors.textSecondary}>Price</span>
        </div>
        {activeIndicators.includes('sma') && (
          <>
            <div className="flex items-center space-x-1">
              <div className="w-3 h-0.5 bg-amber-500 opacity-80" style={{borderTop: '1px dashed'}}></div>
              <span className={colors.textSecondary}>SMA(20)</span>
            </div>
            <div className="flex items-center space-x-1">
              <div className="w-3 h-0.5 bg-red-500 opacity-80" style={{borderTop: '1px dashed'}}></div>
              <span className={colors.textSecondary}>SMA(50)</span>
            </div>
          </>
        )}
        {activeIndicators.includes('ema') && (
          <div className="flex items-center space-x-1">
            <div className="w-3 h-0.5 bg-purple-500"></div>
            <span className={colors.textSecondary}>EMA(12)</span>
          </div>
        )}
        {activeIndicators.includes('bollinger') && (
          <div className="flex items-center space-x-1">
            <div className="w-3 h-0.5 bg-green-500 opacity-60"></div>
            <span className={colors.textSecondary}>Bollinger Bands</span>
          </div>
        )}
      </div>
    </div>
  )
}

// Helper functions for technical indicators
function calculateSMA(data, period, currentPrice) {
  if (data.length < period) return currentPrice
  
  const slice = data.slice(-period)
  const sum = slice.reduce((acc, item) => acc + (item.price || currentPrice), 0)
  return sum / period
}

function calculateEMA(data, period, currentPrice) {
  if (data.length === 0) return currentPrice
  
  const multiplier = 2 / (period + 1)
  if (data.length < period) {
    return calculateSMA(data, Math.min(data.length, period), currentPrice)
  }
  
  const previousEMA = data[data.length - 1]?.ema12 || currentPrice
  return (currentPrice * multiplier) + (previousEMA * (1 - multiplier))
}

function calculateRSI(data, period, currentPrice) {
  if (data.length < period + 1) return 50
  
  const changes = []
  for (let i = Math.max(0, data.length - period); i < data.length; i++) {
    const change = (data[i]?.price || currentPrice) - (data[i-1]?.price || currentPrice)
    changes.push(change)
  }
  
  const gains = changes.filter(change => change > 0)
  const losses = changes.filter(change => change < 0).map(loss => Math.abs(loss))
  
  const avgGain = gains.length > 0 ? gains.reduce((a, b) => a + b, 0) / gains.length : 0
  const avgLoss = losses.length > 0 ? losses.reduce((a, b) => a + b, 0) / losses.length : 0
  
  if (avgLoss === 0) return 100
  
  const rs = avgGain / avgLoss
  return 100 - (100 / (1 + rs))
}

function calculateBollingerBands(data, index, period) {
  if (index < period - 1) {
    return { upper: data[index]?.price, middle: data[index]?.price, lower: data[index]?.price }
  }
  
  const slice = data.slice(Math.max(0, index - period + 1), index + 1)
  const prices = slice.map(item => item.price)
  const sma = prices.reduce((a, b) => a + b, 0) / prices.length
  
  const variance = prices.reduce((acc, price) => acc + Math.pow(price - sma, 2), 0) / prices.length
  const stdDev = Math.sqrt(variance)
  
  return {
    upper: sma + (stdDev * 2),
    middle: sma,
    lower: sma - (stdDev * 2)
  }
}

function calculateMACD(data, index) {
  // Simplified MACD calculation
  const ema12 = calculateEMA(data.slice(0, index + 1), 12, data[index]?.price || 0)
  const ema26 = calculateEMA(data.slice(0, index + 1), 26, data[index]?.price || 0)
  const macd = ema12 - ema26
  const signal = calculateEMA([...data.slice(0, index), { price: macd }], 9, macd)
  const histogram = macd - signal
  
  return { macd, signal, histogram }
}

export default AdvancedChart