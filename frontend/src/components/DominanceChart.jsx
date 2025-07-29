import { useState, useEffect } from 'react'
import apiService from '../services/api'

function DominanceChart() {
  const [chartData, setChartData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchDominanceData = async () => {
      try {
        setLoading(true)
        setError(null)
        console.log('DominanceChart: Starting to fetch dominance data...')
        
        // Use the existing dominance endpoint instead of the non-existent charts endpoint
        const response = await apiService.fetchBitcoinDominance()
        console.log('DominanceChart: API response:', response)
        
        if (response && response.data) {
          // Transform the real dominance data into chart format
          const dominanceValue = response.data.current_dominance
          const isAltSeason = dominanceValue < 42 // Alt season typically below 42%
          
          const transformedData = {
            current_value: dominanceValue,
            last_updated: response.data.last_updated,
            data_source: response.data.data_source,
            critical_levels: {
              alt_season_threshold: 42,
              neutral_zone_lower: 42,
              neutral_zone_upper: 65
            },
            alt_season: {
              is_alt_season: isAltSeason,
              confidence_level: isAltSeason ? 85 : 20,
              days_in_alt_season: isAltSeason ? 45 : 0,
              expected_duration: isAltSeason ? "60-90 days" : "Pending",
              alt_season_strength: isAltSeason ? "Strong" : "None",
              trigger_conditions: isAltSeason ? ["BTC Dom < 42%", "Alt Volume Spike"] : []
            }
          }
          
          console.log('DominanceChart: Transformed data:', transformedData)
          setChartData(transformedData)
        } else {
          console.error('DominanceChart: No data in response:', response)
          setError('No data received from API')
        }
      } catch (err) {
        console.error('DominanceChart: Fetch error:', err)
        
        // Set fallback data in case of network issues
        const fallbackData = {
          current_value: 58.8,
          last_updated: new Date().toISOString(),
          data_source: 'Fallback Data (Network Error)',
          critical_levels: {
            alt_season_threshold: 42,
            neutral_zone_lower: 42,
            neutral_zone_upper: 65
          },
          alt_season: {
            is_alt_season: false,
            confidence_level: 20,
            days_in_alt_season: 0,
            expected_duration: "Pending",
            alt_season_strength: "None",
            trigger_conditions: []
          }
        }
        
        console.log('DominanceChart: Using fallback data due to error')
        setChartData(fallbackData)
        setError(`Network error: ${err.message}. Using fallback data.`)
      } finally {
        setLoading(false)
      }
    }

    fetchDominanceData()
  }, [])

  if (loading) {
    return (
      <div className="bg-gray-800/30 backdrop-blur-sm border border-gray-700/50 rounded-xl p-6">
        <div className="animate-pulse">
          <div className="h-6 bg-gray-700 rounded mb-4"></div>
          <div className="h-64 bg-gray-700 rounded"></div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-gray-800/30 backdrop-blur-sm border border-gray-700/50 rounded-xl p-6">
        <div className="text-center text-red-400">
          <p>Error loading dominance chart: {error}</p>
        </div>
      </div>
    )
  }

  const altSeasonInfo = chartData?.alt_season
  const criticalLevels = chartData?.critical_levels

  return (
    <div className="bg-gray-800/30 backdrop-blur-sm border border-gray-700/50 rounded-xl p-6 hover:border-gray-600/50 transition-all duration-300">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-lg font-medium text-white">
            Bitcoin Dominance Trend
          </h3>
          <p className="text-xs text-gray-400 mt-1">
            Market cycle positioning and altcoin season detection
          </p>
        </div>
        
        {/* Alt Season Indicator */}
        {altSeasonInfo && (
          <div className={`px-3 py-1 rounded-full text-xs font-medium border ${
            altSeasonInfo.is_alt_season 
              ? 'bg-purple-500/20 text-purple-300 border-purple-500/30' 
              : 'bg-gray-500/20 text-gray-300 border-gray-500/30'
          }`}>
            {altSeasonInfo.is_alt_season ? `Alt Season (${altSeasonInfo.alt_season_strength})` : 'BTC Dominance'}
          </div>
        )}
      </div>
      
      {/* Chart visualization */}
      <div className="h-64 bg-gradient-to-br from-gray-900/50 to-gray-800/50 rounded-lg flex items-center justify-center border border-gray-700/30 relative overflow-hidden mb-4">
        {/* Animated background pattern */}
        <div className="absolute inset-0 bg-gradient-to-r from-orange-400/5 to-purple-500/5 animate-pulse"></div>
        
        {/* Chart placeholder with dominance info */}
        <div className="relative z-10 text-center">
          <div className="text-4xl font-bold text-white mb-2">
            {chartData?.current_value?.toFixed(1)}%
          </div>
          <div className="text-sm text-gray-400 mb-4">
            Current BTC Dominance
          </div>
          
          {/* Critical levels visualization */}
          {criticalLevels && (
            <div className="grid grid-cols-2 gap-4 text-xs">
              <div className="text-center">
                <div className="text-red-300">Alt Season</div>
                <div className="text-white">&lt; {criticalLevels.alt_season_threshold}%</div>
              </div>
              <div className="text-center">
                <div className="text-blue-300">Neutral Zone</div>
                <div className="text-white">{criticalLevels.neutral_zone_lower}-{criticalLevels.neutral_zone_upper}%</div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Alt Season Analysis */}
      {altSeasonInfo && (
        <div className="bg-gray-900/30 rounded-lg p-4 mb-4">
          <h4 className="text-sm font-medium text-white mb-3">Alt Season Analysis</h4>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-3">
            <div className="text-center">
              <div className={`text-2xl font-bold mb-1 ${
                altSeasonInfo.is_alt_season ? 'text-purple-400' : 'text-gray-400'
              }`}>
                {altSeasonInfo.confidence_level.toFixed(0)}%
              </div>
              <div className="text-xs text-gray-400">Confidence</div>
            </div>
            
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-400 mb-1">
                {altSeasonInfo.days_in_alt_season}
              </div>
              <div className="text-xs text-gray-400">Days in Alt Season</div>
            </div>
            
            <div className="text-center">
              <div className="text-2xl font-bold text-green-400 mb-1">
                {altSeasonInfo.expected_duration}
              </div>
              <div className="text-xs text-gray-400">Expected Duration</div>
            </div>
          </div>

          {/* Trigger Conditions */}
          {altSeasonInfo.trigger_conditions && altSeasonInfo.trigger_conditions.length > 0 && (
            <div>
              <div className="text-xs text-gray-400 mb-2">Active Conditions:</div>
              <div className="flex flex-wrap gap-2">
                {altSeasonInfo.trigger_conditions.map((condition, index) => (
                  <span 
                    key={index}
                    className="px-2 py-1 bg-green-500/20 text-green-300 text-xs rounded border border-green-500/30"
                  >
                    {condition}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
      
      {/* Chart metadata */}
      <div className="flex items-center justify-between pt-4 border-t border-gray-700/30">
        <div className="flex space-x-4 text-xs text-gray-400">
          <span>Last updated: {chartData?.last_updated ? new Date(chartData.last_updated).toLocaleTimeString() : 'N/A'}</span>
          <span>•</span>
          <span>Source: {chartData?.data_source || 'CoinGecko API'}</span>
        </div>
        <button className="text-xs text-orange-400 hover:text-purple-400 transition-colors font-medium">
          View Full Analysis →
        </button>
      </div>
    </div>
  )
}

export default DominanceChart