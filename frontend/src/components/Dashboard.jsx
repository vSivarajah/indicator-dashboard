import { useEffect, useState } from 'react'
import IndicatorCard from './IndicatorCard'
import FearGreedCard from './FearGreedCard'
import ChartContainer from './ChartContainer'
import DominanceChart from './DominanceChart'
import HeroSection from './HeroSection'
import { useIndicators } from '../hooks/useIndicators'
import { 
  typography, 
  colors, 
  spacing, 
  effects, 
  buildCardStyles, 
  buildButtonStyles, 
  buildProfessionalText, 
  buildGradientText,
  goldTheme, 
  iconStyles, 
  animations,
  layoutStyles,
  cardHeights
} from '../utils/designSystem'
import { TooltipProvider, Button, Badge, Tooltip, TooltipTrigger, TooltipContent } from './ui'
import { ArrowPathIcon, InformationCircleIcon, ClockIcon } from '@heroicons/react/24/outline'
import { cn } from '../lib/utils'

function Dashboard() {
  const { crypto, macro, portfolio, loading, error, lastUpdated, refresh } = useIndicators()

  // Transform API data to component format
  const transformIndicatorData = (apiData, id, title, description) => {
    if (!apiData) return null
    
    return {
      id,
      title,
      value: apiData.value || 'N/A',
      change: apiData.change || '+0.00',
      riskLevel: apiData.risk_level || 'medium',
      description,
      status: apiData.status || 'Loading...'
    }
  }

  // Transform real market data for Bitcoin dominance
  const transformDominanceData = (dominanceData) => {
    if (!dominanceData) return null
    
    return {
      id: 'btc-dominance',
      title: 'Bitcoin Dominance',
      value: `${dominanceData.current_dominance?.toFixed(1) || '54.5'}%`,
      change: dominanceData.change_percent_24h 
        ? `${dominanceData.change_percent_24h > 0 ? '+' : ''}${dominanceData.change_percent_24h.toFixed(2)}%`
        : '-0.55%',
      riskLevel: dominanceData.current_dominance > 60 ? 'high' : dominanceData.current_dominance < 45 ? 'low' : 'medium',
      description: 'BTC market cap percentage from TradingView',
      status: dominanceData.data_source === 'Fallback Data' 
        ? 'Using fallback data - TradingView unavailable' 
        : 'Real-time from TradingView'
    }
  }

  // Create mock data for missing indicators while using real dominance data
  const createMockMVRV = () => ({
    id: 'mvrv-zscore',
    title: 'MVRV Z-Score',
    value: '2.43',
    change: '+0.12',
    riskLevel: 'medium',
    description: 'Market Value to Realized Value ratio',
    status: 'Testing resistance at 2.5 - Watch for breakout'
  })

  const createMockFearGreed = () => ({
    id: 'fear-greed',
    title: 'Fear & Greed Index',
    value: '72',
    change: '+5',
    riskLevel: 'high',
    description: 'Market sentiment indicator (0-100)',
    status: 'Greed territory - Consider taking profits'
  })

  const createMockBubbleRisk = () => ({
    id: 'bubble-risk',
    title: 'Short-term Bubble Risk',
    value: 'Medium',
    change: 'Stable',
    riskLevel: 'medium',
    description: 'Overheating detection algorithm',
    status: 'Monitor closely for rapid changes'
  })

  // Crypto indicators combining real and mock data
  const cryptoIndicators = [
    createMockMVRV(),
    transformDominanceData(crypto.dominance), // Use real dominance data
    createMockFearGreed(),
    createMockBubbleRisk()
  ].filter(Boolean) // Remove null values

  // Macro indicators from real API data
  const macroIndicators = [
    transformIndicatorData(
      macro.inflation,
      'inflation',
      'US Inflation Rate',
      'Consumer Price Index YoY'
    ),
    transformIndicatorData(
      macro.interestRates,
      'interest-rates',
      'Fed Fund Rate',
      'Federal Reserve interest rate'
    )
  ].filter(Boolean)


  const portfolioMetrics = [
    { label: 'Overall Risk Level', value: portfolio.overall_risk || 'Medium', color: 'text-yellow-400' },
    { label: 'Market Cycle Stage', value: portfolio.cycle_stage || '2/5', color: 'text-blue-400' },
    { label: 'Recommended Action', value: portfolio.recommended_action || 'Hold', color: 'text-green-400' }
  ]

  // Error state
  if (error) {
    return (
      <main className="flex-1 max-w-7xl mx-auto px-6 py-8">
        <div className="bg-red-500/10 border border-red-500/30 rounded-xl p-6 text-center">
          <h2 className="text-xl font-semibold text-red-400 mb-2">Connection Error</h2>
          <p className="text-gray-300 mb-4">{error}</p>
          <button 
            onClick={refresh}
            className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          >
            Retry Connection
          </button>
        </div>
      </main>
    )
  }

  // Use the crypto indicators we created (combination of real dominance and mock data)
  const displayCryptoIndicators = cryptoIndicators
  const displayMacroIndicators = macroIndicators

  return (
    <TooltipProvider>
      <main className={cn(layoutStyles.pageContainerLarge, layoutStyles.sectionSpacing)}>
        {/* Enhanced Dashboard Header */}
        <div className={cn(layoutStyles.headerLayout, spacing.sectionMargin)}>
          <div>
            <h1 className={cn(
              typography.pageTitle,
              colors.textPrimary
            )}>
              Market Dashboard
            </h1>
            <p className={cn(
              typography.bodyLarge,
              colors.textSecondary,
              "mt-1"
            )}>
              Real-time cryptocurrency market intelligence
            </p>
          </div>
          
          <div className="flex items-center space-x-3">
            {lastUpdated && (
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge variant="outline" className={cn(
                    "px-3 py-1 border-orange-500/30 text-orange-500 hover:border-orange-500/50",
                    effects.transition
                  )}>
                    <ClockIcon className={cn(iconStyles.small, "text-orange-500", "mr-1")} />
                    <span className="text-xs">
                      {new Date(lastUpdated).toLocaleTimeString()}
                    </span>
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>
                  <p>Last updated: {new Date(lastUpdated).toLocaleString()}</p>
                </TooltipContent>
              </Tooltip>
            )}
            
            <Button
              onClick={refresh}
              disabled={loading}
              size="sm"
              variant="outline"
              className={cn(
                buildButtonStyles('outline', 'small'),
                "transition-all duration-200 hover:scale-[1.02]"
              )}
            >
              <ArrowPathIcon className={cn(
                iconStyles.small,
                "text-orange-500",
                "mr-2 transition-transform duration-500",
                loading && "animate-spin"
              )} />
              {loading ? 'Refreshing...' : 'Refresh'}
            </Button>
          </div>
        </div>
      {/* Hero Section */}
      <HeroSection cryptoData={crypto} loading={loading} />
      
      {/* Market Indicators */}
      <section>
        <div className={cn(layoutStyles.headerLayout, spacing.elementMargin)}>
          <h2 className={cn(
            typography.sectionTitle,
            colors.textPrimary
          )}>
            Market Indicators
            {loading && (
              <span className={cn("ml-2", typography.bodySmall, colors.textSecondary)}>
                <div className="inline-block w-4 h-4 border-2 border-accent border-t-transparent rounded-full animate-spin"></div>
              </span>
            )}
          </h2>
        </div>
        <div className={cn(layoutStyles.dashboardGrid, spacing.gridGap)}>
          {displayCryptoIndicators.map((indicator, index) => {
            const Component = indicator.id === 'fear-greed' ? FearGreedCard : IndicatorCard
            const props = indicator.id === 'fear-greed' 
              ? { fearGreedData: crypto.fearGreed?.data, indicator, isLoading: loading }
              : { ...indicator, isLoading: loading }
            
            return (
              <div
                key={indicator.id}
                className={cn(
                  animations.slideInStagger,
                  effects.scaleUpSmall,
                  effects.transition
                )}
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <Component {...props} />
              </div>
            )
          })}
        </div>
      </section>

      {/* Charts Section */}
      <section>
        <div className={cn(layoutStyles.headerLayout, spacing.elementMargin)}>
          <h2 className={cn(
            typography.sectionTitle,
            colors.textPrimary
          )}>
            Market Analysis Charts
          </h2>
          <div className={cn(typography.body, colors.textSecondary)}>
            Historical trends and patterns
          </div>
        </div>
        <div className={cn(layoutStyles.chartGrid, spacing.gridGap)}>
          <div 
            className="animate-in slide-in-from-left-4 duration-500"
            style={{ animationDelay: '0ms' }}
          >
            <ChartContainer 
              title="MVRV Z-Score Historical"
              description="Long-term market value analysis with cycle identification"
              loading={loading}
            />
          </div>
          <div 
            className="animate-in slide-in-from-right-4 duration-500"
            style={{ animationDelay: '100ms' }}
          >
            <DominanceChart />
          </div>
          <div 
            className="animate-in slide-in-from-left-4 duration-500"
            style={{ animationDelay: '200ms' }}
          >
            <ChartContainer 
              title="Log Regression Bands"
              description="Bitcoin price trend channels and support/resistance"
              loading={loading}
            />
          </div>
          <div 
            className="animate-in slide-in-from-right-4 duration-500"
            style={{ animationDelay: '300ms' }}
          >
            <ChartContainer 
              title="Moving Averages (20W/21W)"
              description="Long-term trend analysis and bull/bear confirmation"
              loading={loading}
            />
          </div>
        </div>
      </section>

      {/* Macro Economic Indicators */}
      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className={cn(
            typography.sectionTitle,
            buildGradientText('goldBright')
          )}>
            Macro Economic Indicators
          </h2>
          <div className="text-sm text-gray-400">
            Traditional market factors affecting crypto
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {displayMacroIndicators.map((indicator) => (
            <IndicatorCard key={indicator.id} {...indicator} />
          ))}
        </div>
      </section>

      {/* Portfolio Risk Assessment */}
      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className={cn(
            typography.sectionTitle,
            buildGradientText('goldBright')
          )}>
            Portfolio Risk Assessment
          </h2>
          <div className="text-sm text-gray-400">
            AI-powered risk analysis and exit strategies
          </div>
        </div>
        
        <div className={cn(
          buildCardStyles(true, 'gold'),
          "group p-8",
          effects.goldGlowHover
        )}>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {portfolioMetrics.map((metric, index) => (
              <div 
                key={index} 
                className={cn(
                  "text-center",
                  animations.slideInStagger
                )}
                style={{ animationDelay: `${index * 150}ms` }}
              >
                <div className={cn(
                  "text-4xl md:text-5xl font-bold mb-2 transition-all duration-300 group-hover:scale-105",
                  metric.color,
                  animations.countUp
                )}>
                  {metric.value}
                </div>
                <div className="text-sm text-gray-400">
                  {metric.label}
                </div>
              </div>
            ))}
          </div>
          
          <div className="mt-8 pt-6 border-t border-gray-700/30">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className={cn(
                  "text-lg font-medium mb-3",
                  colors.textGold
                )}>Exit Strategy Recommendations</h4>
                <ul className="space-y-2 text-sm text-gray-300">
                  <li className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
                    <span>Consider taking 25% profits if BTC hits $75K</span>
                  </li>
                  <li className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-yellow-400 rounded-full animate-pulse"></div>
                    <span>Set stop losses at 20% below current levels</span>
                  </li>
                  <li className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-yellow-300 rounded-full animate-pulse"></div>
                    <span>DCA strategy recommended for new positions</span>
                  </li>
                </ul>
              </div>
              
              <div>
                <h4 className={cn(
                  "text-lg font-medium mb-3",
                  colors.textGold
                )}>Risk Factors</h4>
                <ul className="space-y-2 text-sm text-gray-300">
                  <li className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-red-400 rounded-full animate-pulse"></div>
                    <span>High Fear & Greed (72) - Potential reversal risk</span>
                  </li>
                  <li className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-orange-400 rounded-full animate-pulse"></div>
                    <span>Fed policy uncertainty ahead</span>
                  </li>
                  <li className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-yellow-400 rounded-full animate-pulse"></div>
                    <span>MVRV Z-Score approaching resistance</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* DCA Sidebar */}
      </main>
    </TooltipProvider>
  )
}

export default Dashboard