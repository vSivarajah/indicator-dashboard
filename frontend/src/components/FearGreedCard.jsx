import { useState, useEffect } from 'react'
import { InformationCircleIcon, ChartBarIcon, ArrowTrendingUpIcon, ArrowTrendingDownIcon } from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Button } from './ui/button'
import { Separator } from './ui/separator'
import { cn } from '@/lib/utils'

function FearGreedCard({ fearGreedData, indicator }) {
  const [details, setDetails] = useState(null)
  const [showDetails, setShowDetails] = useState(false)

  // Use either the real API data or fall back to indicator data
  const data = fearGreedData || indicator

  useEffect(() => {
    if (fearGreedData?.details) {
      setDetails(fearGreedData.details)
    }
  }, [fearGreedData])

  if (!data) {
    return (
      <Card className="animate-pulse">
        <CardHeader>
          <div className="h-4 bg-muted rounded w-3/4 mb-4"></div>
          <div className="h-8 bg-muted rounded w-1/2 mb-2"></div>
          <div className="h-4 bg-muted rounded w-full"></div>
        </CardHeader>
      </Card>
    )
  }

  const getRiskVariant = (riskLevel) => {
    switch (riskLevel?.toLowerCase()) {
      case 'low': return 'success'
      case 'medium': return 'warning'
      case 'high': return 'destructive'
      default: return 'secondary'
    }
  }

  const getChangeIcon = (change) => {
    const changeValue = parseInt(change)
    if (changeValue > 0) return <ArrowTrendingUpIcon className="w-4 h-4 text-green-400" />
    if (changeValue < 0) return <ArrowTrendingDownIcon className="w-4 h-4 text-red-400" />
    return <div className="w-4 h-4" />
  }

  const getValueColor = (value) => {
    const numValue = parseInt(value)
    if (numValue >= 75) return 'text-red-400'
    if (numValue >= 55) return 'text-orange-400'
    if (numValue >= 45) return 'text-yellow-400'
    if (numValue >= 25) return 'text-blue-400'
    return 'text-green-400'
  }

  const formatTradingAdvice = (recommendation) => {
    if (!recommendation) return null
    
    const actionColors = {
      'buy': 'text-green-400',
      'sell': 'text-red-400',
      'hold': 'text-yellow-400',
      'caution': 'text-orange-400'
    }

    return (
      <div className="p-4 bg-muted/30 rounded-lg border space-y-3">
        <div className="flex items-center justify-between">
          <span className={cn("font-semibold", actionColors[recommendation.action] || 'text-muted-foreground')}>
            {recommendation.action?.toUpperCase()}
          </span>
          <Badge variant="outline" className="text-xs">
            {recommendation.confidence}% confidence
          </Badge>
        </div>
        <div className="text-sm text-muted-foreground">
          Time horizon: {recommendation.time_horizon}
        </div>
        {recommendation.reasoning && (
          <ul className="text-xs text-muted-foreground space-y-1">
            {recommendation.reasoning.map((reason, index) => (
              <li key={index} className="flex items-start space-x-2">
                <div className="w-1 h-1 bg-muted-foreground rounded-full mt-2 flex-shrink-0"></div>
                <span>{reason}</span>
              </li>
            ))}
          </ul>
        )}
      </div>
    )
  }

  return (
    <Card className="group relative transition-all duration-300 hover:shadow-lg hover:shadow-primary/10 hover:border-primary/50 hover:scale-[1.02] bg-card/50 backdrop-blur-sm">
      {/* Subtle gradient overlay on hover */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-accent/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
      
      <CardHeader className="relative pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <ChartBarIcon className="w-5 h-5 text-blue-400" />
            <CardTitle className="text-lg font-semibold text-foreground group-hover:text-primary transition-colors">
              Fear & Greed Index
            </CardTitle>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setShowDetails(!showDetails)}
            className="h-8 w-8 text-muted-foreground hover:text-foreground"
          >
            <InformationCircleIcon className="w-4 h-4" />
          </Button>
        </div>
        <CardDescription className="text-sm text-muted-foreground">
          {data.description || 'Market sentiment indicator (0-100)'}
        </CardDescription>
      </CardHeader>

      <CardContent className="relative space-y-4">
        {/* Main Value */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <span className={cn("text-3xl font-bold", getValueColor(data.value))}>
              {data.value}
            </span>
            <div className="flex items-center space-x-1">
              {getChangeIcon(data.change)}
              <span className={data.change?.startsWith('+') ? 'text-green-400' : 'text-red-400'}>
                {data.change}
              </span>
            </div>
          </div>
          <Badge variant={getRiskVariant(data.risk_level)} className="capitalize">
            {details?.classification || data.risk_level}
          </Badge>
        </div>

        {/* Status */}
        <p className="text-sm text-muted-foreground">
          {data.status}
        </p>

        {/* Source Attribution */}
        {details?.data_source && (
          <div className="text-xs text-muted-foreground">
            Source: <a 
              href={details.data_source.website}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-400 hover:text-blue-300 transition-colors underline"
            >
              {details.data_source.provider}
            </a>
          </div>
        )}

        {/* Trading Recommendation */}
        {showDetails && details?.trading_recommendation && (
          <>
            <Separator />
            {formatTradingAdvice(details.trading_recommendation)}
          </>
        )}

        {/* Component Breakdown */}
        {showDetails && details?.components && (
          <>
            <Separator />
            <div className="space-y-3">
              <h4 className="text-sm font-medium text-foreground">Component Analysis</h4>
              <div className="grid grid-cols-2 gap-3 text-xs">
                {Object.entries(details.components).map(([key, component]) => (
                  <div key={key} className="space-y-1">
                    <div className="flex items-center justify-between">
                      <span className="text-muted-foreground capitalize">{key.replace('_', ' ')}</span>
                      <span className="text-foreground font-medium">{component.value}</span>
                    </div>
                    <div className="text-muted-foreground text-xs">
                      Weight: {component.weight}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </>
        )}

        {/* Last Updated */}
        <Separator />
        <div className="text-xs text-muted-foreground">
          Last updated: {new Date(data.timestamp).toLocaleString()}
          {details?.next_update && (
            <span className="ml-2">• Next update: {details.next_update}</span>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export default FearGreedCard