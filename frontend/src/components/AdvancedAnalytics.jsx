import React, { useState, useEffect } from 'react'
import { 
  ChartBarIcon, 
  CpuChipIcon, 
  ExclamationTriangleIcon,
  LightBulbIcon,
  ArrowTrendingUpIcon,
  BeakerIcon
} from '@heroicons/react/24/outline'
import { cn } from '../lib/utils'
import { typography, colors, spacing, effects, buildCardStyles, goldTheme } from '../utils/designSystem'
import AdvancedChart from './AdvancedChart'
import PortfolioRiskAnalysis from './PortfolioRiskAnalysis'

const AdvancedAnalytics = ({ className = "" }) => {
  const [activeSection, setActiveSection] = useState('ai-insights')
  const [mlData, setMLData] = useState(null)
  const [loading, setLoading] = useState(false)

  // Sample ML insights data
  const mockMLInsights = {
    pricePrediction: {
      symbol: 'BTC',
      predictedPrice: 52400,
      confidence: 0.76,
      direction: 'Bullish',
      priceChange: 3.2,
      riskLevel: 'Medium',
      timestamp: new Date()
    },
    marketCycle: {
      currentStage: 'Mid Bull',
      nextStage: 'Late Bull',
      stageConfidence: 0.82,
      stageProgress: 67.5,
      estimatedDuration: 120,
      timestamp: new Date()
    },
    anomaly: {
      anomalyScore: 0.35,
      isAnomaly: false,
      anomalyType: 'Normal',
      severity: 'Low',
      description: 'Market conditions within normal ranges',
      indicators: [],
      timestamp: new Date()
    },
    sentiment: {
      overallSentiment: 0.68,
      sentimentLabel: 'Optimistic',
      sources: {
        social_media: 0.65,
        news: 0.71,
        on_chain: 0.73,
        technical: 0.62
      },
      confidence: 0.78,
      trendDirection: 'Improving',
      timestamp: new Date()
    }
  }

  useEffect(() => {
    // Simulate loading ML insights
    setLoading(true)
    setTimeout(() => {
      setMLData(mockMLInsights)
      setLoading(false)
    }, 1500)
  }, [])

  const sections = [
    { 
      id: 'ai-insights', 
      label: 'AI Insights', 
      icon: CpuChipIcon,
      description: 'Machine learning predictions and analysis'
    },
    { 
      id: 'advanced-charts', 
      label: 'Advanced Charts', 
      icon: ChartBarIcon,
      description: 'TradingView-style interactive charts'
    },
    { 
      id: 'risk-analysis', 
      label: 'Risk Analysis', 
      icon: ExclamationTriangleIcon,
      description: 'Monte Carlo simulations and portfolio risk'
    },
    { 
      id: 'market-intelligence', 
      label: 'Market Intelligence', 
      icon: LightBulbIcon,
      description: 'Comprehensive market cycle analysis'
    }
  ]

  const renderAIInsights = () => (
    <div className="space-y-6">
      {/* Header */}
      <div className="text-center">
        <div className={cn(
          "inline-flex items-center space-x-2 px-4 py-2 rounded-full",
          goldTheme.backgroundSubtle,
          "border border-amber-500/30 mb-4"
        )}>
          <CpuChipIcon className="w-5 h-5 text-amber-400" />
          <span className={cn(typography.bodyMedium, "text-amber-400 font-semibold")}>
            AI-Powered Market Intelligence
          </span>
        </div>
        <h2 className={cn(typography.displaySmall, colors.textPrimary, "mb-2")}>
          Advanced Analytics Suite
        </h2>
        <p className={cn(typography.bodyLarge, colors.textSecondary, "max-w-3xl mx-auto")}>
          Leveraging machine learning and artificial intelligence to provide sophisticated market predictions,
          anomaly detection, and sentiment analysis for informed investment decisions.
        </p>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className={cn(buildCardStyles(true, 'default'), "p-6")}>
              <div className="animate-pulse space-y-3">
                <div className="h-4 bg-gray-700 rounded w-1/3"></div>
                <div className="h-6 bg-gray-700 rounded w-2/3"></div>
                <div className="h-3 bg-gray-700 rounded w-full"></div>
                <div className="h-3 bg-gray-700 rounded w-3/4"></div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Price Prediction */}
          <div className={cn(buildCardStyles(true, 'default'), "p-6")}>
            <div className="flex items-center justify-between mb-4">
              <h3 className={cn(typography.cardTitle, colors.textPrimary)}>
                Price Prediction
              </h3>
              <div className={cn(
                "px-3 py-1 rounded-full text-xs font-semibold",
                mlData?.pricePrediction.direction === 'Bullish' 
                  ? "bg-green-500/20 text-green-400"
                  : "bg-red-500/20 text-red-400"
              )}>
                {mlData?.pricePrediction.direction}
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Predicted Price (24h)
                </span>
                <span className={cn(typography.bodyLarge, colors.textPrimary, "font-semibold")}>
                  ${mlData?.pricePrediction.predictedPrice.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Expected Change
                </span>
                <span className={cn(
                  typography.bodyMedium,
                  mlData?.pricePrediction.priceChange > 0 ? "text-green-400" : "text-red-400",
                  "font-semibold"
                )}>
                  {mlData?.pricePrediction.priceChange > 0 ? '+' : ''}{mlData?.pricePrediction.priceChange}%
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Confidence
                </span>
                <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                  {(mlData?.pricePrediction.confidence * 100).toFixed(1)}%
                </span>
              </div>
              <div className="w-full bg-gray-700 rounded-full h-2">
                <div 
                  className="bg-blue-500 h-2 rounded-full transition-all duration-1000"
                  style={{ width: `${(mlData?.pricePrediction.confidence || 0) * 100}%` }}
                ></div>
              </div>
            </div>
          </div>

          {/* Market Cycle Prediction */}
          <div className={cn(buildCardStyles(true, 'default'), "p-6")}>
            <div className="flex items-center justify-between mb-4">
              <h3 className={cn(typography.cardTitle, colors.textPrimary)}>
                Market Cycle Analysis
              </h3>
              <div className={cn(
                "px-3 py-1 rounded-full text-xs font-semibold",
                "bg-blue-500/20 text-blue-400"
              )}>
                Stage {mlData?.marketCycle.currentStage}
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Current Stage
                </span>
                <span className={cn(typography.bodyMedium, colors.textPrimary, "font-semibold")}>
                  {mlData?.marketCycle.currentStage}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Next Stage
                </span>
                <span className={cn(typography.bodyMedium, "text-orange-400")}>
                  {mlData?.marketCycle.nextStage}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Stage Progress
                </span>
                <span className={cn(typography.bodyMedium, colors.textPrimary)}>
                  {mlData?.marketCycle.stageProgress.toFixed(1)}%
                </span>
              </div>
              <div className="w-full bg-gray-700 rounded-full h-2">
                <div 
                  className="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full transition-all duration-1000"
                  style={{ width: `${mlData?.marketCycle.stageProgress || 0}%` }}
                ></div>
              </div>
              <p className={cn(typography.bodySmall, colors.textSecondary)}>
                Estimated time to next stage: ~{mlData?.marketCycle.estimatedDuration} days
              </p>
            </div>
          </div>

          {/* Anomaly Detection */}
          <div className={cn(buildCardStyles(true, 'default'), "p-6")}>
            <div className="flex items-center justify-between mb-4">
              <h3 className={cn(typography.cardTitle, colors.textPrimary)}>
                Anomaly Detection
              </h3>
              <div className={cn(
                "px-3 py-1 rounded-full text-xs font-semibold",
                mlData?.anomaly.isAnomaly 
                  ? "bg-red-500/20 text-red-400"
                  : "bg-green-500/20 text-green-400"
              )}>
                {mlData?.anomaly.isAnomaly ? 'Anomaly Detected' : 'Normal'}
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Anomaly Score
                </span>
                <span className={cn(typography.bodyMedium, colors.textPrimary, "font-mono")}>
                  {(mlData?.anomaly.anomalyScore * 100).toFixed(1)}%
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Severity
                </span>
                <span className={cn(
                  typography.bodyMedium,
                  mlData?.anomaly.severity === 'High' ? "text-red-400" :
                  mlData?.anomaly.severity === 'Medium' ? "text-yellow-400" :
                  "text-green-400"
                )}>
                  {mlData?.anomaly.severity}
                </span>
              </div>
              <div className="w-full bg-gray-700 rounded-full h-2">
                <div 
                  className={cn(
                    "h-2 rounded-full transition-all duration-1000",
                    mlData?.anomaly.anomalyScore > 0.7 ? "bg-red-500" :
                    mlData?.anomaly.anomalyScore > 0.3 ? "bg-yellow-500" :
                    "bg-green-500"
                  )}
                  style={{ width: `${(mlData?.anomaly.anomalyScore || 0) * 100}%` }}
                ></div>
              </div>
              <p className={cn(typography.bodySmall, colors.textSecondary)}>
                {mlData?.anomaly.description}
              </p>
            </div>
          </div>

          {/* Sentiment Analysis */}
          <div className={cn(buildCardStyles(true, 'default'), "p-6")}>
            <div className="flex items-center justify-between mb-4">
              <h3 className={cn(typography.cardTitle, colors.textPrimary)}>
                Market Sentiment
              </h3>
              <div className={cn(
                "px-3 py-1 rounded-full text-xs font-semibold",
                "bg-purple-500/20 text-purple-400"
              )}>
                {mlData?.sentiment.sentimentLabel}
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                  Overall Sentiment
                </span>
                <span className={cn(typography.bodyMedium, colors.textPrimary, "font-semibold")}>
                  {(mlData?.sentiment.overallSentiment * 100).toFixed(0)}/100
                </span>
              </div>
              <div className="space-y-2">
                {Object.entries(mlData?.sentiment.sources || {}).map(([source, value]) => (
                  <div key={source} className="flex justify-between items-center text-xs">
                    <span className={cn(colors.textSecondary, "capitalize")}>
                      {source.replace('_', ' ')}
                    </span>
                    <div className="flex items-center space-x-2">
                      <div className="w-16 bg-gray-700 rounded-full h-1">
                        <div 
                          className="bg-purple-500 h-1 rounded-full"
                          style={{ width: `${value * 100}%` }}
                        ></div>
                      </div>
                      <span className={colors.textSecondary}>
                        {(value * 100).toFixed(0)}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
              <div className="flex items-center space-x-2 pt-2">
                <ArrowTrendingUpIcon className="w-4 h-4 text-green-400" />
                <span className={cn(typography.bodySmall, "text-green-400")}>
                  Trend: {mlData?.sentiment.trendDirection}
                </span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )

  const renderAdvancedCharts = () => (
    <div className="space-y-6">
      <div className="text-center mb-8">
        <h2 className={cn(typography.displaySmall, colors.textPrimary, "mb-2")}>
          Advanced Charting Suite
        </h2>
        <p className={cn(typography.bodyLarge, colors.textSecondary)}>
          Professional-grade charts with technical indicators and multi-timeframe analysis
        </p>
      </div>
      
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <AdvancedChart 
          title="Bitcoin Price Analysis"
          symbol="BTC"
          height={400}
        />
        <AdvancedChart 
          title="Ethereum Technical Analysis"
          symbol="ETH"
          height={400}
        />
      </div>

      <AdvancedChart 
        title="Market Overview - Multi-Asset Analysis"
        symbol="CRYPTO"
        height={500}
        className="col-span-full"
      />
    </div>
  )

  const renderRiskAnalysis = () => (
    <div className="space-y-6">
      <div className="text-center mb-8">
        <h2 className={cn(typography.displaySmall, colors.textPrimary, "mb-2")}>
          Portfolio Risk Management
        </h2>
        <p className={cn(typography.bodyLarge, colors.textSecondary)}>
          Comprehensive risk analysis with Monte Carlo simulations and VaR calculations
        </p>
      </div>
      
      <PortfolioRiskAnalysis />
    </div>
  )

  const renderMarketIntelligence = () => (
    <div className="space-y-6">
      <div className="text-center mb-8">
        <h2 className={cn(typography.displaySmall, colors.textPrimary, "mb-2")}>
          Market Intelligence Center
        </h2>
        <p className={cn(typography.bodyLarge, colors.textSecondary)}>
          Comprehensive market cycle analysis and strategic recommendations
        </p>
      </div>

      {/* Coming Soon Message */}
      <div className={cn(
        buildCardStyles(true, 'default'),
        "p-12 text-center"
      )}>
        <BeakerIcon className="w-16 h-16 text-blue-400 mx-auto mb-4" />
        <h3 className={cn(typography.cardTitle, colors.textPrimary, "mb-2")}>
          Market Intelligence Features Coming Soon
        </h3>
        <p className={cn(typography.bodyLarge, colors.textSecondary, "mb-6")}>
          Advanced market cycle analysis, correlation studies, and strategic recommendations
          are currently in development.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 max-w-4xl mx-auto">
          {[
            'Cross-Asset Correlation Analysis',
            'Macro Economic Integration', 
            'Market Regime Detection',
            'Strategic Allocation Models',
            'Risk-Adjusted Portfolio Construction',
            'Dynamic Rebalancing Strategies'
          ].map((feature, index) => (
            <div key={index} className="p-4 bg-gray-800/30 rounded-lg">
              <span className={cn(typography.bodyMedium, colors.textSecondary)}>
                {feature}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )

  return (
    <div className={cn("max-w-7xl mx-auto px-6 py-8", className)}>
      {/* Navigation */}
      <div className="flex flex-wrap justify-center gap-2 mb-8 bg-gray-800/50 p-2 rounded-xl">
        {sections.map((section) => (
          <button
            key={section.id}
            onClick={() => setActiveSection(section.id)}
            className={cn(
              "flex items-center space-x-2 px-4 py-3 rounded-lg transition-all duration-200",
              "hover:scale-[1.02] hover:shadow-lg",
              activeSection === section.id
                ? "bg-blue-600 text-white shadow-lg"
                : "text-gray-300 hover:text-white hover:bg-gray-700/50"
            )}
          >
            <section.icon className="w-5 h-5" />
            <div className="text-left">
              <div className={cn(typography.bodyMedium, "font-semibold")}>
                {section.label}
              </div>
              <div className={cn(typography.bodySmall, "opacity-75")}>
                {section.description}
              </div>
            </div>
          </button>
        ))}
      </div>

      {/* Content */}
      <div className="min-h-[600px]">
        {activeSection === 'ai-insights' && renderAIInsights()}
        {activeSection === 'advanced-charts' && renderAdvancedCharts()}
        {activeSection === 'risk-analysis' && renderRiskAnalysis()}
        {activeSection === 'market-intelligence' && renderMarketIntelligence()}
      </div>
    </div>
  )
}

export default AdvancedAnalytics