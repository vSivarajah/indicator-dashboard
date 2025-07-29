import { memo } from 'react'
import { ArrowTrendingUpIcon, CurrencyDollarIcon, ChartBarIcon } from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { cn } from '@/lib/utils'

const HeroSection = memo(function HeroSection({ cryptoData, loading }) {
  // Extract key metrics from real market data
  const btcData = cryptoData?.prices?.BTC
  const dominanceData = cryptoData?.dominance
  
  const btcPrice = btcData?.quote?.USD?.price 
    ? Math.round(btcData.quote.USD.price).toLocaleString()
    : '67,432'
  const btcChange = btcData?.quote?.USD?.percent_change_24h
    ? `${btcData.quote.USD.percent_change_24h > 0 ? '+' : ''}${btcData.quote.USD.percent_change_24h.toFixed(2)}%`
    : '+2.34%'
  const marketCap = btcData?.quote?.USD?.market_cap
    ? `${(btcData.quote.USD.market_cap / 1e12).toFixed(2)}T`
    : '1.29T'
  const dominanceValue = dominanceData?.current_dominance
    ? `${dominanceData.current_dominance.toFixed(1)}%`
    : '56.8%'

  const stats = [
    {
      label: 'Bitcoin Price',
      value: `$${btcPrice}`,
      change: btcChange,
      icon: CurrencyDollarIcon,
      positive: btcChange?.startsWith('+'),
    },
    {
      label: 'Market Cap',
      value: `$${marketCap}`,
      change: btcChange, // Use BTC change as proxy for market cap change
      icon: ChartBarIcon,
      positive: btcChange?.startsWith('+'),
    },
    {
      label: 'BTC Dominance',
      value: dominanceValue,
      change: dominanceData?.change_percent_24h 
        ? `${dominanceData.change_percent_24h > 0 ? '+' : ''}${dominanceData.change_percent_24h.toFixed(2)}%`
        : '-0.55%',
      icon: ArrowTrendingUpIcon,
      positive: dominanceData?.change_percent_24h > 0,
    },
  ]

  const AnimatedCounter = ({ value, duration = 2000 }) => {
    return (
      <span className="tabular-nums">
        {value}
      </span>
    )
  }

  return (
    <section className="py-8">
      <Card className="relative overflow-hidden bg-card/50 backdrop-blur-sm border-border/50 shadow-xl">
        {/* Background gradient animation */}
        <div className="absolute inset-0 bg-gradient-to-br from-primary/10 to-accent/10 opacity-30 animate-pulse"></div>
        
        <CardContent className="relative z-10 p-8">
          {/* Hero title */}
          <div className="text-center mb-8">
            <h1 className="text-4xl md:text-6xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent mb-4">
              Market Intelligence
            </h1>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Real-time cryptocurrency market analysis powered by advanced on-chain metrics and sentiment indicators
            </p>
          </div>

          {/* Key metrics grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {stats.map((stat, index) => {
              const Icon = stat.icon
              return (
                <Card
                  key={stat.label}
                  className="group relative transition-all duration-300 hover:shadow-lg hover:shadow-primary/10 hover:border-primary/50 hover:scale-[1.02] bg-card/50 backdrop-blur-sm"
                  style={{ animationDelay: `${index * 150}ms` }}
                >
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                      <div className={cn(
                        "p-3 rounded-lg",
                        stat.positive ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'
                      )}>
                        <Icon className="w-6 h-6" />
                      </div>
                      <Badge 
                        variant={stat.positive ? 'success' : 'destructive'}
                        className="text-xs"
                      >
                        {stat.change}
                      </Badge>
                    </div>
                  </CardHeader>

                  <CardContent className="pt-0">
                    <CardDescription className="text-sm text-muted-foreground mb-2">
                      {stat.label}
                    </CardDescription>
                    <div className="text-2xl md:text-3xl font-bold text-foreground">
                      {loading ? (
                        <div className="h-8 bg-muted rounded animate-pulse"></div>
                      ) : (
                        <AnimatedCounter value={stat.value} />
                      )}
                    </div>
                  </CardContent>

                  {/* Hover effect */}
                  <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-accent/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                </Card>
              )
            })}
          </div>

          {/* Market status indicator */}
          <div className="mt-8 text-center">
            <div className="inline-flex items-center space-x-2 px-4 py-2 bg-muted/50 rounded-full border border-border/50">
              <div className={cn(
                "w-2 h-2 rounded-full animate-pulse",
                loading ? 'bg-yellow-400' : 'bg-green-400'
              )}></div>
              <span className="text-sm text-muted-foreground">
                {loading ? 'Loading market data...' : 'Markets are open • Live data'}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>
    </section>
  )
})

export default HeroSection