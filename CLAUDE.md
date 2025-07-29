# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a comprehensive cryptocurrency indicator dashboard with a React frontend and Go backend that provides real-time analysis of market indicators including MVRV Z-Score, Bitcoin Dominance, Fear & Greed Index, Bubble Risk metrics, and Bitcoin Rainbow Chart analysis. The system uses PostgreSQL/TimescaleDB for time-series data storage, Redis for caching, and integrates multiple free data sources with consensus pricing algorithms for maximum accuracy and reliability.

## Dashboard Purpose & Functionality

**Vision & Core Mission**
This platform serves as a comprehensive cryptocurrency market analysis dashboard designed for serious investors, traders, and analysts who need data-driven insights for long-term market timing and risk management. The system combines on-chain metrics, macroeconomic indicators, and market sentiment analysis to provide actionable intelligence for portfolio management.

**Primary Use Cases**
- **Market Cycle Timing**: Identify optimal entry and exit points using multi-indicator confluence analysis
- **Risk Assessment**: Real-time portfolio risk evaluation with AI-powered recommendations
- **Long-term Investment Strategy**: Dollar-cost averaging optimization and systematic profit-taking guidance
- **Market Intelligence**: Comprehensive macro and crypto correlation analysis for informed decision-making

## Core Market Indicators

**MVRV Z-Score Analysis**
- **Purpose**: Market Value to Realized Value ratio for cycle timing and valuation assessment
- **Implementation**: Real-time calculations with historical Z-score tracking (`services/mvrv_service.go`)
- **Key Features**: Threshold-based risk levels, cycle identification, resistance/support analysis
- **Thresholds**: Extreme low (-1.5), Low (-0.5), Neutral (0.5-1.5), High (3.0), Extreme high (7.0)
- **Trading Signals**: Buy zones below -0.5, sell zones above 3.0, bubble warnings above 7.0

**Bitcoin Dominance Intelligence**
- **Purpose**: Market cap percentage analysis for alt-season detection and market cycle positioning
- **Implementation**: Real-time dominance tracking with trend analysis (`services/dominance_service.go`)
- **Alt-Season Detection**: Automated signals when dominance drops below 42% threshold
- **Key Levels**: Alt-season entry (42%), Strong alt-season (38%), Cycle top (70%), Cycle bottom (35%)
- **Trend Analysis**: 20-day and 50-day moving averages with strength classification

**Fear & Greed Index Integration**
- **Purpose**: Market sentiment analysis for contrarian investment opportunities
- **Data Source**: Alternative.me API with real-time sentiment scoring (0-100)
- **Implementation**: Cached sentiment data with component analysis (`services/fear_greed_service.go`)
- **Trading Psychology**: Extreme fear (0-25) = buying opportunity, Extreme greed (75-100) = distribution zone
- **Component Analysis**: Market volatility, momentum, social media, surveys, dominance, trends

**Bubble Risk Assessment**
- **Purpose**: Multi-factor analysis for detecting market overheating and bubble conditions
- **Algorithm**: Combines MVRV ratio, NVT signal, social sentiment, exchange flows, long-term holder behavior
- **Implementation**: Real-time risk scoring with confidence levels (`services/bubble_risk_service.go`)
- **Risk Categories**: Low (0-25), Medium (25-50), High (50-75), Extreme (75-90), Bubble warning (90-100)
- **Trading Recommendations**: Dynamic exit strategies based on risk score progression

**Bitcoin Rainbow Chart Analysis** (Implemented)
- **Purpose**: Logarithmic regression-based long-term cycle analysis and market timing
- **Implementation**: Mathematical model using days since Bitcoin genesis with color-coded risk bands (`services/rainbow_chart_service.go`)
- **Key Features**: 9-band risk assessment system from "Fire Sale" to "Maximum Bubble Territory"
- **Formula**: log10(price) = -17.01593313 + 5.84509503 * log10(days_from_genesis)
- **Risk Bands**: Fire Sale (0.8x), BUY! (1.0x), Accumulate (1.3x), Still Cheap (1.6x), HODL! (2.0x), Is This A Bubble? (2.4x), FOMO Intensifies (3.0x), Sell Seriously (4.0x), Maximum Bubble Territory (5.0x)
- **Trading Signals**: Automated buy recommendations below 1.6x regression, sell signals above 3.0x regression
- **Cycle Position**: Real-time calculation of current position within historical price cycles (0-100%)
- **Historical Analysis**: Complete data since Bitcoin genesis block for cycle duration and pattern recognition

## Advanced Chart Analysis Features

**Log Regression Bands**
- **Purpose**: Bitcoin price trend channels using logarithmic regression analysis
- **Implementation**: Historical price modeling with support/resistance identification
- **Key Features**: Upper/lower bands, trend channel analysis, breakout detection
- **Trading Applications**: Long-term trend confirmation, support/resistance levels, cycle positioning

**MVRV Historical Charts**
- **Purpose**: Long-term cycle visualization with Z-score historical analysis
- **Features**: Cycle identification, peak/trough analysis, current position relative to history
- **Time Frames**: Multi-year analysis with cycle overlay and duration estimates
- **Pattern Recognition**: Bull/bear market confirmation, cycle stage identification

**Moving Average Analysis**
- **Purpose**: 20-week and 21-week moving average analysis for long-term trend confirmation
- **Implementation**: Bull/bear market confirmation signals with trend strength analysis
- **Key Signals**: Golden cross (bullish), death cross (bearish), trend continuation patterns
- **Risk Management**: Stop-loss placement, trend following strategies, momentum confirmation

**Dominance Chart Visualization**
- **Purpose**: Bitcoin dominance trends with alt-season timing analysis
- **Features**: Historical dominance patterns, alt-season duration estimates, market cycle correlation
- **Predictive Analysis**: Expected alt-season duration, confidence levels, trigger conditions
- **Portfolio Allocation**: Dynamic BTC/alt allocation recommendations based on dominance trends

## Portfolio Management & Risk Analysis

**AI-Powered Risk Assessment**
- **Implementation**: Multi-factor risk analysis combining all indicators for portfolio evaluation
- **Risk Metrics**: Overall risk level, market cycle stage positioning, recommended actions
- **Dynamic Analysis**: Real-time risk score updates based on market condition changes
- **Confidence Scoring**: Statistical confidence levels for risk assessments and recommendations

**Exit Strategy Optimization**
- **Purpose**: Systematic profit-taking strategies based on market cycle analysis
- **Implementation**: Dynamic exit recommendations using indicator confluence
- **Strategy Types**: Percentage-based exits, risk-level triggered exits, cycle-stage exits
- **Risk Management**: Stop-loss recommendations, position sizing, portfolio rebalancing

**Dollar Cost Averaging (DCA) Calculator** (Implemented)
- **Purpose**: Comprehensive backtesting and optimization of systematic buying strategies
- **Features**: Historical performance analysis, market timing analysis, risk-adjusted recommendations
- **Implementation**: Dedicated route at `/dca` with full-page interface
- **Backend Service**: Complete DCA simulation engine (`services/dca_service.go`) with performance metrics
- **UI Features**: Modern slider controls, frequency toggles, interactive charts, debounced inputs
- **Analysis Types**: Performance metrics, market timing scores, purchase rank evaluation, strategic recommendations

**Portfolio Risk Testing** (Planned)
- **Purpose**: Monte Carlo simulations and stress testing for portfolio optimization
- **Features**: Scenario analysis, correlation analysis, drawdown assessment
- **Risk Metrics**: Value at Risk (VaR), maximum drawdown, Sharpe ratio analysis
- **Optimization**: Risk-adjusted portfolio construction with dynamic allocation

## Market Cycle Intelligence

**5-Stage Cycle Classification**
- **Stages**: Bear Market, Early Bull, Mid Bull, Late Bull, Cycle Top
- **Implementation**: Multi-indicator confluence analysis (`models/indicator.go:MarketCycle`)
- **Confidence Scoring**: Statistical probability for each cycle stage
- **Duration Estimates**: Expected stage duration based on historical patterns

**Cycle Stage Indicators**
- **Bear Market**: High dominance (>65%), negative MVRV Z-score, extreme fear
- **Early Bull**: Declining dominance, MVRV Z-score rising above -0.5, fear reducing
- **Mid Bull**: Dominance 45-60%, MVRV Z-score 0.5-2.0, greed increasing
- **Late Bull**: Low dominance (<45%), MVRV Z-score 2.0-5.0, extreme greed
- **Cycle Top**: Dominance reversal, MVRV Z-score >5.0, bubble risk extreme

**Predictive Analysis**
- **Pattern Recognition**: Historical cycle analysis for duration and pattern matching
- **Confluence Analysis**: Multiple indicator confirmation for cycle stage transitions
- **Risk Adjustment**: Dynamic risk assessment based on cycle progression
- **Timing Optimization**: Entry/exit timing based on cycle stage probabilities

## Macro Economic Integration

**Traditional Market Correlation**
- **US Inflation Rate**: Consumer Price Index (CPI) year-over-year analysis
- **Federal Reserve Policy**: Interest rate tracking with crypto correlation analysis
- **Implementation**: Real-time macro data integration (`services/data_fetcher.go`)
- **Impact Analysis**: Macro event impact on crypto market cycles and risk assessment

**Recession Risk Analysis** (Planned)
- **Purpose**: Traditional market recession indicators with crypto correlation
- **Indicators**: Yield curve inversion, unemployment rates, GDP growth, consumer confidence
- **Integration**: Macro risk scoring with crypto-specific adjustments
- **Portfolio Impact**: Risk adjustment recommendations based on macro environment

**Cross-Asset Analysis** (Planned)
- **Purpose**: Gold, stocks, bonds correlation with cryptocurrency markets
- **Implementation**: Multi-asset risk analysis and portfolio optimization
- **Risk Management**: Macro-aware portfolio construction and rebalancing strategies
- **Market Regime**: Detection of risk-on vs risk-off environments for crypto positioning

## Data Infrastructure & Analytics Foundation

**Multi-Source Data Aggregation System** (Implemented)
- **Purpose**: Combine multiple free data sources for maximum accuracy and reliability without paid APIs
- **Data Sources**: CoinCap API (rest.coincap.io/v3), CoinGecko API (free tier), Blockchain.com API (no authentication required)
- **Implementation**: Consensus pricing algorithm with confidence scoring (`services/market_data_aggregator.go`)
- **Key Features**: Real-time reliability assessment, standard deviation-based consensus, source validation
- **Confidence Scoring**: Multi-factor confidence calculation based on source agreement and reliability history
- **API Management**: Authentication handling, rate limiting, graceful degradation when sources unavailable
- **Fallback Mechanisms**: Automatic failover to available sources when APIs are down

**PostgreSQL Time-Series Database** (Implemented)
- **Purpose**: Efficient storage and retrieval of historical market data for ML/analytics and long-term analysis
- **Implementation**: Optimized PostgreSQL tables with TimescaleDB support (`internal/infrastructure/database/`)
- **Table Structure**: Separate tables for price data, indicators, market metrics, network metrics, rainbow chart data
- **Time-Series Optimization**: Custom indexes for timestamp-based queries, efficient data retrieval patterns
- **Data Types**: Price data (multi-asset), indicator calculations, market metrics, network statistics, rainbow chart analysis
- **Retention Policies**: Automated cleanup with configurable retention periods (1-5 years based on data type)
- **ML-Ready Schema**: Structured data format optimized for machine learning model training and backtesting
- **Performance**: Batch insertions, query optimization, table statistics monitoring

**Bitcoin Network Metrics Integration** (Implemented)
- **Purpose**: Real-time Bitcoin blockchain network health and security monitoring
- **Data Source**: Blockchain.com API (free, no authentication required)
- **Implementation**: Network statistics client (`internal/infrastructure/external/blockchain_client.go`)
- **Key Metrics**: Hash rate, mining difficulty, block height, total supply, transaction count, fees, mempool size
- **Analysis**: Network security trends, mining economics, transaction fee optimization
- **Historical Tracking**: Time-series storage for network metric evolution over time
- **Integration**: Combined with price data for comprehensive market analysis

**Advanced API Integration Architecture** (Implemented)
- **CoinCap API Integration**: Professional-grade integration with API key authentication
  - Base URL: rest.coincap.io/v3 (updated from deprecated v2 API)
  - Authentication: Bearer token support for rate limit increases
  - Features: Asset data, historical prices, market data, global statistics
  - Implementation: `internal/infrastructure/external/coincap_client.go`
- **Consensus Pricing Engine**: Multi-source price validation and confidence scoring
- **Health Monitoring**: Real-time API health checks and service availability monitoring
- **Error Handling**: Comprehensive error handling with retry mechanisms and fallback strategies

## Machine Learning & Analytics Foundation

**Time-Series Data Pipeline** (Ready for ML)
- **Structured Data Format**: Normalized schemas for price, indicator, and network data
- **Feature Engineering Ready**: Timestamp-indexed data with metadata for ML feature extraction
- **Historical Dataset**: Complete price and indicator history for model training
- **Real-Time Integration**: Live data pipeline for model inference and prediction
- **Backtesting Support**: Historical data access for strategy validation and optimization

**Advanced Analytics Capabilities**
- **Correlation Analysis**: Multi-asset and macro correlation tracking
- **Cycle Detection**: Automated market cycle identification using multiple indicators
- **Risk Assessment**: Multi-factor risk scoring with confidence intervals
- **Performance Metrics**: Comprehensive portfolio and strategy performance analysis

## Additional Memories

### Project Philosophy and Future Development
- Emphasize continuous learning and adaptation in the crypto market analysis ecosystem
- Commitment to building an intelligent, data-driven investment decision support system
- **Machine Learning Integration**: Foundation established with time-series database and structured data pipeline
- **Free Data Sources Strategy**: Successfully implemented comprehensive analysis using only free APIs
- **Scalable Architecture**: Built for production use with background data collection and automated analysis
- **Data Quality Focus**: Multi-source validation and consensus algorithms ensure data reliability