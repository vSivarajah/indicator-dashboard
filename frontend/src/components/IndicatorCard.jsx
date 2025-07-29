import { memo } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Skeleton } from './ui/skeleton'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import { cn } from '@/lib/utils'
import { 
  buildCardStyles, 
  buildProfessionalText, 
  goldTheme, 
  iconStyles, 
  animations, 
  effects, 
  cardHeights,
  typography,
  colors,
  spacing
} from '../utils/designSystem'

const IndicatorCard = memo(function IndicatorCard({ title, value, change, riskLevel, description, status, isLoading = false }) {
  const getRiskVariant = (level) => {
    switch (level?.toLowerCase()) {
      case 'low': return 'success'
      case 'medium': return 'warning'
      case 'high': return 'destructive'
      default: return 'secondary'
    }
  }

  const getRiskThemeClass = (level) => {
    switch (level?.toLowerCase()) {
      case 'low': return goldTheme.riskLow
      case 'medium': return goldTheme.riskMedium
      case 'high': return goldTheme.riskHigh
      case 'extreme': return goldTheme.riskExtreme
      default: return goldTheme.neutral
    }
  }

  const getChangeColor = (change) => {
    if (change?.startsWith('+')) return 'text-green-400'
    if (change?.startsWith('-')) return 'text-red-400'
    return 'text-muted-foreground'
  }

  // Loading skeleton state with consistent height
  if (isLoading) {
    return (
      <Card className={cn(
        colors.cardBg,
        colors.border,
        effects.rounded,
        cardHeights.indicatorCard,
        "flex flex-col"
      )}>
        <CardHeader className={cn(spacing.cardPaddingSmall, "pb-2")}>
          <div className="flex items-start justify-between">
            <div className="space-y-2 flex-1">
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-48" />
            </div>
            <Skeleton className="h-6 w-16" />
          </div>
        </CardHeader>
        <CardContent className={cn(spacing.cardPaddingSmall, "pt-0 flex-1 flex flex-col justify-between")}>
          <div className="flex items-baseline space-x-3 mb-3">
            <Skeleton className="h-9 w-24" />
            <Skeleton className="h-5 w-16" />
          </div>
          <Skeleton className="h-4 w-full mt-auto" />
        </CardContent>
      </Card>
    )
  }

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Card className={cn(
            buildCardStyles(true, 'professional'),
            cardHeights.indicatorCard,
            "group relative cursor-pointer professional-hover flex flex-col"
          )}>
            {/* Subtle hover overlay */}
            <div className="absolute inset-0 bg-gradient-to-br from-accent/5 to-accent/2 rounded-lg opacity-0 group-hover:opacity-100 transition-all duration-300" />
            
            <CardHeader className={cn("relative pb-2", spacing.cardPaddingSmall)}>
              <div className="flex items-start justify-between">
                <div className="space-y-1 flex-1">
                  <CardTitle className={cn(
                    typography.cardTitle,
                    colors.textPrimary,
                    "transition-colors duration-200"
                  )}>
                    {title}
                  </CardTitle>
                  <CardDescription className={cn(
                    typography.bodySmall,
                    colors.textSecondary,
                    "group-hover:text-foreground/80 transition-colors"
                  )}>
                    {description}
                  </CardDescription>
                </div>
                <Badge 
                  variant={getRiskVariant(riskLevel)} 
                  className={cn(
                    "capitalize transition-all duration-200 group-hover:scale-105 shrink-0",
                    getRiskThemeClass(riskLevel),
                    colors.borderAccent
                  )}
                >
                  {riskLevel}
                </Badge>
              </div>
            </CardHeader>

            <CardContent className={cn(
              "relative pt-0 flex-1 flex flex-col justify-between",
              spacing.cardPaddingSmall
            )}>
              <div className="flex items-baseline space-x-3 mb-3">
                <span className={cn(
                  typography.valueDisplay,
                  colors.textPrimary,
                  "transition-colors duration-200",
                  animations.countUp
                )}>
                  {value}
                </span>
                {change && (
                  <span className={cn(
                    typography.valueChange,
                    "transition-all duration-200 group-hover:scale-110", 
                    getChangeColor(change)
                  )}>
                    {change}
                  </span>
                )}
              </div>
              
              {status && (
                <p className={cn(
                  typography.bodySmall,
                  colors.textMuted,
                  "group-hover:text-foreground/80 transition-colors mt-auto"
                )}>
                  {status}
                </p>
              )}
            </CardContent>
          </Card>
        </TooltipTrigger>
        <TooltipContent>
          <div className="space-y-1">
            <p className="font-medium">{title}</p>
            <p className="text-xs text-muted-foreground">{description}</p>
            <div className="flex items-center space-x-2">
              <span className="text-xs">Risk Level:</span>
              <Badge variant={getRiskVariant(riskLevel)} className="capitalize text-xs">
                {riskLevel}
              </Badge>
            </div>
          </div>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
})

export default IndicatorCard