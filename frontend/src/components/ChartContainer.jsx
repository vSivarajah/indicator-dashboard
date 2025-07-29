import { useState, memo } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Button } from './ui/button'
import { Tabs, TabsList, TabsTrigger } from './ui/tabs'
import { Separator } from './ui/separator'
import { cn } from '@/lib/utils'
import { 
  buildCardStyles, 
  cardHeights, 
  colors, 
  typography, 
  spacing, 
  effects,
  layoutStyles
} from '../utils/designSystem'

const ChartContainer = memo(function ChartContainer({ title, description, timeframe = "1M", loading = false }) {
  const [selectedTimeframe, setSelectedTimeframe] = useState(timeframe)
  const timeframes = ["1D", "1W", "1M", "3M", "1Y"]
  
  return (
    <Card className={cn(
      buildCardStyles(true, 'professional'),
      cardHeights.chartCard,
      "group flex flex-col"
    )}>
      <CardHeader className={cn(spacing.cardPaddingSmall, "pb-4")}>
        <div className={layoutStyles.cardHeaderLayout}>
          <div className="space-y-1 flex-1">
            <CardTitle className={cn(
              typography.cardTitle,
              colors.textPrimary,
              "group-hover:text-accent transition-colors"
            )}>
              {title}
            </CardTitle>
            <CardDescription className={cn(
              typography.bodySmall,
              colors.textSecondary
            )}>
              {description}
            </CardDescription>
          </div>
          
          <Tabs value={selectedTimeframe} onValueChange={setSelectedTimeframe}>
            <TabsList className="bg-muted/50">
              {timeframes.map((tf) => (
                <TabsTrigger
                  key={tf}
                  value={tf}
                  className="text-xs data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
                >
                  {tf}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
        </div>
      </CardHeader>

      <CardContent className={cn(spacing.cardPaddingSmall, "space-y-4 flex-1 flex flex-col")}>
        {/* Chart placeholder with consistent height */}
        <div className="flex-1 bg-gradient-to-br from-muted/50 to-muted/30 rounded-lg flex items-center justify-center border border-border/50 relative overflow-hidden group-hover:border-accent/30 transition-all duration-300">
          {/* Animated background pattern */}
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-accent/5 opacity-50 group-hover:opacity-100 transition-opacity duration-300"></div>
          
          {loading ? (
            /* Loading skeleton */
            <div className="w-full h-full p-4">
              <div className="flex items-end justify-between h-full space-x-2">
                {Array.from({ length: 12 }, (_, i) => (
                  <div
                    key={i}
                    className="bg-muted rounded-t animate-pulse"
                    style={{
                      height: `${Math.random() * 60 + 20}%`,
                      animationDelay: `${i * 100}ms`,
                      animationDuration: '2s'
                    }}
                  ></div>
                ))}
              </div>
            </div>
          ) : (
            <div className="relative z-10 text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-gradient-to-r from-primary to-accent rounded-full flex items-center justify-center opacity-40 group-hover:opacity-60 transition-opacity duration-300 group-hover:scale-110">
                <svg className="w-8 h-8 text-white transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                </svg>
              </div>
              <p className="text-sm text-muted-foreground group-hover:text-foreground transition-colors duration-300">
                Advanced charting coming soon
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                Selected: {selectedTimeframe}
              </p>
            </div>
          )}
        </div>
        
        {/* Chart metadata */}
        <Separator />
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4 text-xs text-muted-foreground">
            <span>Last updated: 2 min ago</span>
            <span>•</span>
            <span>Source: Multiple APIs</span>
          </div>
          <Button variant="ghost" size="sm" className="text-xs text-primary hover:text-primary/80 h-auto p-0 font-medium">
            View Details →
          </Button>
        </div>
      </CardContent>
    </Card>
  )
})

export default ChartContainer