import React, { useState, useEffect, useMemo } from 'react'
import {
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  ScatterChart,
  Scatter,
  Area,
  AreaChart
} from 'recharts'
import { 
  ExclamationTriangleIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
  ArrowRightIcon,
  InformationCircleIcon,
  CogIcon,
  ChartBarIcon,
  CalculatorIcon
} from '@heroicons/react/24/outline'
import { cn } from '../lib/utils'
import { typography, colors, spacing, effects, buildCardStyles, goldTheme } from '../utils/designSystem'

const PortfolioRiskAnalysis = ({ 
  portfolio = null,
  className = "",
  onPortfolioChange 
}) => {
  const [activeTab, setActiveTab] = useState('overview')
  const [riskMetrics, setRiskMetrics] = useState(null)
  const [loading, setLoading] = useState(false)
  const [simulationParams, setSimulationParams] = useState({
    simulations: 10000,
    timeHorizon: 252, // 1 year in trading days
    confidence: 0.95
  })

  // Sample portfolio data if none provided
  const samplePortfolio = {
    totalValue: 100000,
    assets: [
      {
        symbol: 'BTC',
        allocation: 40,
        currentPrice: 45000,
        volatility: 65,
        expectedReturn: 25
      },
      {
        symbol: 'ETH',
        allocation: 30,
        currentPrice: 3200,
        volatility: 70,
        expectedReturn: 30
      },
      {
        symbol: 'ADA',
        allocation: 15,
        currentPrice: 0.45,
        volatility: 85,
        expectedReturn: 35
      },
      {
        symbol: 'SOL',
        allocation: 10,
        currentPrice: 95,
        volatility: 90,
        expectedReturn: 40
      },
      {
        symbol: 'AAPL',
        allocation: 5,
        currentPrice: 190,
        volatility: 25,
        expectedReturn: 12
      }
    ],
    lastUpdated: new Date()
  }

  const currentPortfolio = portfolio || samplePortfolio

  // Mock risk metrics (in production, this would come from your backend)
  const mockRiskMetrics = useMemo(() => ({
    overallRiskScore: 68.5,
    varAnalysis: {
      var1Day: 2850,
      var1Week: 6450,
      var1Month: 12800,
      confidence: 0.95,
      method: 'Monte Carlo',
      timestamp: new Date()
    },
    monteCarloAnalysis: {
      simulations: simulationParams.simulations,
      timeHorizon: simulationParams.timeHorizon,
      expectedReturn: 24.5,
      expectedVolatility: 58.2,
      var95: 18.3,
      var99: 26.7,
      maxDrawdown: 45.2,
      sharpeRatio: 1.35,
      percentiles: {
        '5th': -23.5,
        '10th': -18.7,
        '25th': -9.2,
        '50th': 8.5,
        '75th': 28.3,
        '90th': 45.1,
        '95th': 58.9
      },
      paths: generateSamplePaths(20, simulationParams.timeHorizon),
      timestamp: new Date()
    },
    positionSizing: [
      {
        symbol: 'BTC',
        currentAllocation: 40,
        recommendedAllocation: 35,
        maxPosition: 45,
        riskScore: 0.65,
        reasoning: 'High volatility - consider reducing position',
        confidence: 0.78
      },
      {
        symbol: 'ETH',
        currentAllocation: 30,
        recommendedAllocation: 32,
        maxPosition: 35,
        riskScore: 0.70,
        reasoning: 'Good risk-reward ratio',
        confidence: 0.72
      },
      {
        symbol: 'ADA',
        currentAllocation: 15,
        recommendedAllocation: 12,
        maxPosition: 18,
        riskScore: 0.85,
        reasoning: 'Very high volatility - reduce exposure',
        confidence: 0.85
      },
      {
        symbol: 'SOL',
        currentAllocation: 10,
        recommendedAllocation: 8,
        maxPosition: 12,
        riskScore: 0.90,
        reasoning: 'Extreme volatility - limit position',
        confidence: 0.88
      },
      {
        symbol: 'AAPL',
        currentAllocation: 5,
        recommendedAllocation: 13,
        maxPosition: 20,
        riskScore: 0.25,
        reasoning: 'Low risk - consider increasing for diversification',
        confidence: 0.65
      }
    ],
    correlationAnalysis: {
      assetPairs: {
        'BTC-ETH': 0.78,
        'BTC-ADA': 0.65,
        'BTC-SOL': 0.62,
        'BTC-AAPL': 0.15,
        'ETH-ADA': 0.72,
        'ETH-SOL': 0.69,
        'ETH-AAPL': 0.12,
        'ADA-SOL': 0.81,
        'ADA-AAPL': 0.08,
        'SOL-AAPL': 0.05
      },
      diversificationScore: 45.8,
      clusterAnalysis: {
        'Cryptocurrency': ['BTC', 'ETH', 'ADA', 'SOL'],
        'Stocks': ['AAPL']
      }
    },
    riskFactors: [
      'High crypto correlation (78% BTC-ETH)',
      'Significant maximum drawdown risk (45%)',
      'High portfolio volatility (58%)',
      'Concentration in crypto assets (95%)'
    ],
    recommendations: [
      'Diversify into traditional assets to reduce correlation',
      'Consider implementing stop-loss strategies',
      'Rebalance portfolio based on risk-adjusted recommendations',
      'Add defensive assets during high volatility periods'
    ]
  }), [simulationParams])

  useEffect(() => {
    // Simulate loading risk metrics
    setLoading(true)
    setTimeout(() => {
      setRiskMetrics(mockRiskMetrics)
      setLoading(false)
    }, 1000)
  }, [mockRiskMetrics])

  const tabs = [
    { id: 'overview', label: 'Overview', icon: ChartBarIcon },
    { id: 'monte-carlo', label: 'Monte Carlo', icon: CalculatorIcon },
    { id: 'correlations', label: 'Correlations', icon: InformationCircleIcon },
    { id: 'recommendations', label: 'Recommendations', icon: ExclamationTriangleIcon }
  ]

  const getRiskColor = (score) => {
    if (score >= 80) return 'text-red-400'
    if (score >= 60) return 'text-orange-400'
    if (score >= 40) return 'text-yellow-400'
    return 'text-green-400'
  }

  const getRiskLabel = (score) => {
    if (score >= 80) return 'Very High'
    if (score >= 60) return 'High'
    if (score >= 40) return 'Medium'
    return 'Low'
  }

  const COLORS = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6']

  const pieData = currentPortfolio.assets.map((asset, index) => ({
    name: asset.symbol,
    value: asset.allocation,
    color: COLORS[index % COLORS.length],
    amount: (currentPortfolio.totalValue * asset.allocation / 100).toLocaleString()
  }))

  if (loading) {
    return (
      <div className={cn(buildCardStyles(true, 'default'), "p-6", className)}>
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-gray-700 rounded w-1/3"></div>
          <div className="h-4 bg-gray-700 rounded w-2/3"></div>
          <div className="h-64 bg-gray-700 rounded"></div>
        </div>
      </div>
    )
  }

  return (
    <div className={cn(buildCardStyles(true, 'default'), "p-6", className)}>
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className={cn(typography.cardTitle, colors.textPrimary)}>
            Portfolio Risk Analysis
          </h2>
          <p className={cn(typography.bodySmall, colors.textSecondary)}>
            Comprehensive risk assessment and optimization
          </p>
        </div>
        <div className={cn(
          "px-4 py-2 rounded-lg",
          goldTheme.backgroundSubtle,
          "border border-amber-500/30"
        )}>
          <div className="flex items-center space-x-2">
            <span className={cn(typography.bodySmall, colors.textSecondary)}>
              Risk Score:
            </span>
            <span className={cn(
              typography.bodyMedium,
              "font-semibold",
              getRiskColor(riskMetrics?.overallRiskScore || 0)
            )}>
              {riskMetrics?.overallRiskScore.toFixed(1) || 0} / 100
            </span>
            <span className={cn(typography.bodySmall, colors.textSecondary)}>
              ({getRiskLabel(riskMetrics?.overallRiskScore || 0)})
            </span>
          </div>
        </div>
      </div>

      {/* Tab Navigation */}
      <div className="flex space-x-1 mb-6 bg-gray-800/50 p-1 rounded-lg">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={cn(
              "flex items-center space-x-2 px-4 py-2 rounded-md transition-colors text-sm",
              activeTab === tab.id
                ? "bg-blue-600 text-white"
                : "text-gray-400 hover:text-gray-200 hover:bg-gray-700/50"
            )}
          >
            <tab.icon className="w-4 h-4" />
            <span>{tab.label}</span>
          </button>
        ))}
      </div>

      {/* Tab Content */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          {/* Portfolio Composition */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div>
              <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
                Portfolio Composition
              </h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={pieData}
                      cx="50%"
                      cy="50%"
                      innerRadius={60}
                      outerRadius={100}
                      paddingAngle={2}
                      dataKey="value"
                    >
                      {pieData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip
                      formatter={(value, name) => [
                        `${value}% ($${pieData.find(d => d.name === name)?.amount})`,
                        name
                      ]}
                      contentStyle={{
                        backgroundColor: '#1F2937',
                        border: '1px solid #374151',
                        borderRadius: '8px'
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div>
              <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
                Key Risk Metrics
              </h3>
              <div className="space-y-4">
                <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                  <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                    Expected Annual Return
                  </span>
                  <span className={cn(typography.bodyMedium, "text-green-400 font-semibold")}>
                    {riskMetrics?.monteCarloAnalysis.expectedReturn.toFixed(1)}%
                  </span>
                </div>
                <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                  <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                    Portfolio Volatility
                  </span>
                  <span className={cn(typography.bodyMedium, "text-orange-400 font-semibold")}>
                    {riskMetrics?.monteCarloAnalysis.expectedVolatility.toFixed(1)}%
                  </span>
                </div>
                <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                  <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                    Value at Risk (95%)
                  </span>
                  <span className={cn(typography.bodyMedium, "text-red-400 font-semibold")}>
                    ${riskMetrics?.varAnalysis.var1Month.toLocaleString()}
                  </span>
                </div>
                <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                  <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                    Sharpe Ratio
                  </span>
                  <span className={cn(
                    typography.bodyMedium,
                    riskMetrics?.monteCarloAnalysis.sharpeRatio > 1 ? "text-green-400" : "text-yellow-400",
                    "font-semibold"
                  )}>
                    {riskMetrics?.monteCarloAnalysis.sharpeRatio.toFixed(2)}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Position Sizing Recommendations */}
          <div>
            <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
              Position Sizing Recommendations
            </h3>
            <div className="space-y-3">
              {riskMetrics?.positionSizing.map((rec) => (
                <div key={rec.symbol} className="p-4 bg-gray-800/30 rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center space-x-3">
                      <span className={cn(typography.bodyMedium, colors.textPrimary, "font-semibold")}>
                        {rec.symbol}
                      </span>
                      <span className={cn(
                        typography.bodySmall,
                        "px-2 py-1 rounded",
                        rec.riskScore > 0.8 ? "bg-red-500/20 text-red-400" :
                        rec.riskScore > 0.6 ? "bg-orange-500/20 text-orange-400" :
                        rec.riskScore > 0.4 ? "bg-yellow-500/20 text-yellow-400" :
                        "bg-green-500/20 text-green-400"
                      )}>
                        Risk: {(rec.riskScore * 100).toFixed(0)}%
                      </span>
                    </div>
                    <div className="flex items-center space-x-4">
                      <div className="text-center">
                        <div className={cn(typography.bodySmall, colors.textSecondary)}>Current</div>
                        <div className={cn(typography.bodyMedium, colors.textPrimary)}>{rec.currentAllocation}%</div>
                      </div>
                      <ArrowRightIcon className="w-4 h-4 text-gray-500" />
                      <div className="text-center">
                        <div className={cn(typography.bodySmall, colors.textSecondary)}>Recommended</div>
                        <div className={cn(
                          typography.bodyMedium,
                          rec.recommendedAllocation > rec.currentAllocation ? "text-green-400" : "text-orange-400"
                        )}>
                          {rec.recommendedAllocation}%
                        </div>
                      </div>
                    </div>
                  </div>
                  <p className={cn(typography.bodySmall, colors.textSecondary)}>
                    {rec.reasoning}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'monte-carlo' && (
        <div className="space-y-6">
          {/* Simulation Parameters */}
          <div className="p-4 bg-gray-800/30 rounded-lg">
            <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
              Simulation Parameters
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className={cn(typography.bodySmall, colors.textSecondary, "block mb-2")}>
                  Number of Simulations
                </label>
                <input
                  type="number"
                  value={simulationParams.simulations}
                  onChange={(e) => setSimulationParams(prev => ({
                    ...prev,
                    simulations: parseInt(e.target.value)
                  }))}
                  className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white"
                />
              </div>
              <div>
                <label className={cn(typography.bodySmall, colors.textSecondary, "block mb-2")}>
                  Time Horizon (Days)
                </label>
                <input
                  type="number"
                  value={simulationParams.timeHorizon}
                  onChange={(e) => setSimulationParams(prev => ({
                    ...prev,
                    timeHorizon: parseInt(e.target.value)
                  }))}
                  className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white"
                />
              </div>
              <div>
                <label className={cn(typography.bodySmall, colors.textSecondary, "block mb-2")}>
                  Confidence Level
                </label>
                <select
                  value={simulationParams.confidence}
                  onChange={(e) => setSimulationParams(prev => ({
                    ...prev,
                    confidence: parseFloat(e.target.value)
                  }))}
                  className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white"
                >
                  <option value={0.90}>90%</option>
                  <option value={0.95}>95%</option>
                  <option value={0.99}>99%</option>
                </select>
              </div>
            </div>
          </div>

          {/* Monte Carlo Results */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div>
              <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
                Simulation Paths
              </h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={generateTimeAxis(riskMetrics?.monteCarloAnalysis.timeHorizon)}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
                    <XAxis 
                      dataKey="day" 
                      stroke="#9CA3AF"
                      fontSize={12}
                    />
                    <YAxis 
                      stroke="#9CA3AF"
                      fontSize={12}
                      domain={['dataMin * 0.8', 'dataMax * 1.2']}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: '#1F2937',
                        border: '1px solid #374151',
                        borderRadius: '8px'
                      }}
                    />
                    {riskMetrics?.monteCarloAnalysis.paths.slice(0, 5).map((path, index) => (
                      <Line
                        key={index}
                        type="monotone"
                        dataKey={`path${index}`}
                        stroke={COLORS[index]}
                        strokeWidth={1}
                        dot={false}
                        opacity={0.6}
                      />
                    ))}
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div>
              <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
                Return Distribution
              </h3>
              <div className="space-y-3">
                {Object.entries(riskMetrics?.monteCarloAnalysis.percentiles || {}).map(([percentile, value]) => (
                  <div key={percentile} className="flex justify-between items-center p-2 bg-gray-800/20 rounded">
                    <span className={cn(typography.bodySmall, colors.textSecondary)}>
                      {percentile} Percentile
                    </span>
                    <span className={cn(
                      typography.bodySmall,
                      value >= 0 ? "text-green-400" : "text-red-400",
                      "font-mono"
                    )}>
                      {value >= 0 ? '+' : ''}{value.toFixed(1)}%
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'correlations' && (
        <div className="space-y-6">
          {/* Correlation Matrix */}
          <div>
            <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
              Asset Correlation Matrix
            </h3>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr>
                    <th className="text-left p-2"></th>
                    {currentPortfolio.assets.map(asset => (
                      <th key={asset.symbol} className={cn(typography.bodySmall, colors.textSecondary, "text-center p-2")}>
                        {asset.symbol}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {currentPortfolio.assets.map((asset1, i) => (
                    <tr key={asset1.symbol}>
                      <td className={cn(typography.bodySmall, colors.textSecondary, "p-2 font-semibold")}>
                        {asset1.symbol}
                      </td>
                      {currentPortfolio.assets.map((asset2, j) => {
                        const correlation =  i === j ? 1 : 
                          riskMetrics?.correlationAnalysis.assetPairs[`${asset1.symbol}-${asset2.symbol}`] ||
                          riskMetrics?.correlationAnalysis.assetPairs[`${asset2.symbol}-${asset1.symbol}`] ||
                          0

                        return (
                          <td key={asset2.symbol} className="text-center p-2">
                            <span className={cn(
                              typography.bodySmall,
                              "px-2 py-1 rounded font-mono",
                              correlation >= 0.7 ? "bg-red-500/20 text-red-400" :
                              correlation >= 0.4 ? "bg-yellow-500/20 text-yellow-400" :
                              correlation >= 0 ? "bg-green-500/20 text-green-400" :
                              "bg-blue-500/20 text-blue-400"
                            )}>
                              {correlation.toFixed(2)}
                            </span>
                          </td>
                        )
                      })}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Diversification Score */}
          <div className={cn(
            "p-6 rounded-lg",
            riskMetrics?.correlationAnalysis.diversificationScore > 70 ? goldTheme.backgroundSubtle :
            riskMetrics?.correlationAnalysis.diversificationScore > 40 ? "bg-yellow-500/10" :
            "bg-red-500/10"
          )}>
            <div className="flex items-center justify-between">
              <div>
                <h3 className={cn(typography.bodyLarge, colors.textPrimary)}>
                  Diversification Score
                </h3>
                <p className={cn(typography.bodySmall, colors.textSecondary)}>
                  Lower correlation = better diversification
                </p>
              </div>
              <div className="text-right">
                <div className={cn(
                  typography.displayMedium,
                  "font-bold",
                  riskMetrics?.correlationAnalysis.diversificationScore > 70 ? "text-amber-400" :
                  riskMetrics?.correlationAnalysis.diversificationScore > 40 ? "text-yellow-400" :
                  "text-red-400"
                )}>
                  {riskMetrics?.correlationAnalysis.diversificationScore.toFixed(1)}
                </div>
                <div className={cn(typography.bodySmall, colors.textSecondary)}>
                  / 100
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'recommendations' && (
        <div className="space-y-6">
          {/* Risk Factors */}
          <div>
            <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
              Risk Factors
            </h3>
            <div className="space-y-3">
              {riskMetrics?.riskFactors.map((factor, index) => (
                <div key={index} className="flex items-start space-x-3 p-3 bg-red-500/10 border border-red-500/30 rounded-lg">
                  <ExclamationTriangleIcon className="w-5 h-5 text-red-400 mt-0.5" />
                  <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                    {factor}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Recommendations */}
          <div>
            <h3 className={cn(typography.bodyLarge, colors.textPrimary, "mb-4")}>
              Optimization Recommendations
            </h3>
            <div className="space-y-3">
              {riskMetrics?.recommendations.map((recommendation, index) => (
                <div key={index} className="flex items-start space-x-3 p-3 bg-blue-500/10 border border-blue-500/30 rounded-lg">
                  <InformationCircleIcon className="w-5 h-5 text-blue-400 mt-0.5" />
                  <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                    {recommendation}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Action Items */}
          <div className={cn(goldTheme.backgroundSubtle, "p-6 rounded-lg border border-amber-500/30")}>
            <h3 className={cn(typography.bodyLarge, "text-amber-400 mb-4")}>
              Immediate Action Items
            </h3>
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <input type="checkbox" className="rounded border-gray-600 bg-gray-700 text-amber-600" />
                <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                  Reduce SOL allocation from 10% to 8%
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <input type="checkbox" className="rounded border-gray-600 bg-gray-700 text-amber-600" />
                <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                  Increase AAPL allocation from 5% to 13%
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <input type="checkbox" className="rounded border-gray-600 bg-gray-700 text-amber-600" />
                <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                  Set stop-loss orders at 15% below current levels
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <input type="checkbox" className="rounded border-gray-600 bg-gray-700 text-amber-600" />
                <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                  Review portfolio weekly during high volatility periods
                </span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

// Helper functions
function generateSamplePaths(numPaths, timeHorizon) {
  const paths = []
  for (let i = 0; i < numPaths; i++) {
    const path = []
    let value = 100000 // Starting portfolio value
    for (let day = 0; day <= timeHorizon; day++) {
      if (day === 0) {
        path.push(value)
      } else {
        // Simple geometric brownian motion simulation
        const drift = 0.0008 // Daily drift
        const volatility = 0.025 // Daily volatility
        const randomShock = (Math.random() - 0.5) * 2 * Math.sqrt(3) // Normalized random
        value *= Math.exp(drift + volatility * randomShock)
        path.push(value)
      }
    }
    paths.push(path)
  }
  return paths
}

function generateTimeAxis(timeHorizon) {
  const data = []
  const paths = generateSamplePaths(5, timeHorizon)
  
  for (let day = 0; day <= timeHorizon; day += Math.floor(timeHorizon / 50)) {
    const dataPoint = { day }
    paths.forEach((path, index) => {
      dataPoint[`path${index}`] = path[day]
    })
    data.push(dataPoint)
  }
  
  return data
}

export default PortfolioRiskAnalysis