// API service for backend communication
const API_BASE_URL = 'http://localhost:8080/api/v1'

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL
  }

  async fetchIndicator(indicatorType) {
    try {
      const response = await fetch(`${this.baseURL}/indicators/${indicatorType}`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error(`Error fetching ${indicatorType} indicator:`, error)
      throw error
    }
  }

  async fetchMVRVIndicator() {
    return this.fetchIndicator('mvrv')
  }

  async fetchDominanceIndicator() {
    return this.fetchIndicator('dominance')
  }

  async fetchFearGreedIndicator() {
    return this.fetchIndicator('fear-greed')
  }

  async fetchBubbleRiskIndicator() {
    return this.fetchIndicator('bubble-risk')
  }

  async fetchMacroIndicator(type) {
    try {
      const response = await fetch(`${this.baseURL}/macro/${type}`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error(`Error fetching macro indicator ${type}:`, error)
      throw error
    }
  }

  async fetchInflationIndicator() {
    return this.fetchMacroIndicator('inflation')
  }

  async fetchInterestRatesIndicator() {
    return this.fetchMacroIndicator('interest-rates')
  }

  async fetchPortfolioRisk() {
    try {
      const response = await fetch(`${this.baseURL}/portfolio/risk`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('Error fetching portfolio risk:', error)
      throw error
    }
  }

  async fetchMarketCycle() {
    try {
      const response = await fetch(`${this.baseURL}/market/cycle`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('Error fetching market cycle:', error)
      throw error
    }
  }

  async fetchChartData(indicator) {
    try {
      const response = await fetch(`${this.baseURL}/charts/${indicator}`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error(`Error fetching chart data for ${indicator}:`, error)
      throw error
    }
  }

  async fetchMarketData() {
    try {
      const response = await fetch(`${this.baseURL}/market/summary`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('Error fetching market data:', error)
      throw error
    }
  }

  async fetchCryptoPrices() {
    try {
      const response = await fetch(`${this.baseURL}/market/prices`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('Error fetching crypto prices:', error)
      throw error
    }
  }

  async fetchBitcoinDominance() {
    try {
      const url = `${this.baseURL}/market/dominance`
      console.log('API: Fetching Bitcoin dominance from:', url)
      
      const response = await fetch(url)
      console.log('API: Response status:', response.status, response.statusText)
      console.log('API: Response headers:', Object.fromEntries(response.headers.entries()))
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const data = await response.json()
      console.log('API: Response data:', data)
      return data
    } catch (error) {
      console.error('Error fetching Bitcoin dominance:', error)
      throw error
    }
  }

  async fetchMarketHealth() {
    try {
      const response = await fetch(`${this.baseURL}/market/health`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('Error fetching market health:', error)
      throw error
    }
  }

  async fetchAllIndicators() {
    try {
      const [mvrv, dominance, fearGreed, bubbleRisk, inflation, interestRates, marketData] = await Promise.all([
        this.fetchMVRVIndicator(),
        this.fetchDominanceIndicator(),
        this.fetchFearGreedIndicator(),
        this.fetchBubbleRiskIndicator(),
        this.fetchInflationIndicator(),
        this.fetchInterestRatesIndicator(),
        this.fetchMarketData()
      ])

      return {
        crypto: {
          mvrv,
          dominance,
          fearGreed,
          bubbleRisk
        },
        macro: {
          inflation,
          interestRates
        },
        market: marketData
      }
    } catch (error) {
      console.error('Error fetching all indicators:', error)
      throw error
    }
  }

  // Health check
  async checkHealth() {
    try {
      const response = await fetch(`${this.baseURL.replace('/api/v1', '')}/health`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      console.error('Health check failed:', error)
      throw error
    }
  }
}

// Create singleton instance
const apiService = new ApiService()

export default apiService