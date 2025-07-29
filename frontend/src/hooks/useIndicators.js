import { useState, useEffect } from 'react'
import apiService from '../services/api'

export function useIndicators() {
  const [data, setData] = useState({
    crypto: {},
    macro: {},
    portfolio: {},
    loading: true,
    error: null
  })

  const [lastUpdated, setLastUpdated] = useState(null)

  const fetchData = async () => {
    try {
      setData(prev => ({ ...prev, loading: true, error: null }))
      
      // Fetch market data and indicators
      const [marketData, cryptoPrices, bitcoinDominance] = await Promise.all([
        apiService.fetchMarketData().catch(() => null),
        apiService.fetchCryptoPrices().catch(() => null),
        apiService.fetchBitcoinDominance().catch(() => null)
      ])
      
      // Try to fetch traditional indicators (may fail if not implemented)
      let indicators = {
        crypto: {},
        macro: {}
      }
      
      try {
        indicators = await apiService.fetchAllIndicators()
      } catch (error) {
        console.warn('Traditional indicators not available:', error.message)
      }
      
      // Try to fetch portfolio risk (may fail if not implemented)
      let portfolioRisk = {}
      try {
        portfolioRisk = await apiService.fetchPortfolioRisk()
      } catch (error) {
        console.warn('Portfolio risk not available:', error.message)
      }
      
      // Combine real market data with existing indicator structure
      const combinedData = {
        crypto: {
          ...indicators.crypto,
          prices: cryptoPrices?.data || null,
          dominance: bitcoinDominance?.data || indicators.crypto.dominance || null
        },
        macro: indicators.macro,
        market: marketData?.data || null,
        portfolio: portfolioRisk,
        loading: false,
        error: null
      }
      
      setData(combinedData)
      setLastUpdated(new Date())
    } catch (error) {
      console.error('Error fetching indicators:', error)
      setData(prev => ({
        ...prev,
        loading: false,
        error: error.message || 'Failed to fetch data'
      }))
    }
  }

  // Initial fetch
  useEffect(() => {
    fetchData()
  }, [])

  // Auto-refresh every 5 minutes
  useEffect(() => {
    const interval = setInterval(fetchData, 5 * 60 * 1000) // 5 minutes
    return () => clearInterval(interval)
  }, [])

  return {
    ...data,
    lastUpdated,
    refresh: fetchData
  }
}

export function useMVRVChart() {
  const [chartData, setChartData] = useState({
    data: null,
    loading: true,
    error: null
  })

  const fetchChartData = async () => {
    try {
      setChartData(prev => ({ ...prev, loading: true, error: null }))
      
      const data = await apiService.fetchChartData('mvrv')
      
      setChartData({
        data,
        loading: false,
        error: null
      })
    } catch (error) {
      console.error('Error fetching MVRV chart data:', error)
      setChartData({
        data: null,
        loading: false,
        error: error.message || 'Failed to fetch chart data'
      })
    }
  }

  useEffect(() => {
    fetchChartData()
  }, [])

  return {
    ...chartData,
    refresh: fetchChartData
  }
}

export function useBackendHealth() {
  const [isHealthy, setIsHealthy] = useState(false)
  const [isChecking, setIsChecking] = useState(true)

  const checkHealth = async () => {
    try {
      setIsChecking(true)
      await apiService.checkHealth()
      setIsHealthy(true)
    } catch (error) {
      console.error('Backend health check failed:', error)
      setIsHealthy(false)
    } finally {
      setIsChecking(false)
    }
  }

  useEffect(() => {
    checkHealth()
    
    // Check health every 30 seconds
    const interval = setInterval(checkHealth, 30 * 1000)
    return () => clearInterval(interval)
  }, [])

  return {
    isHealthy,
    isChecking,
    checkHealth
  }
}