import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { 
  HomeIcon, 
  CalculatorIcon, 
  ChartBarIcon,
  ChartPieIcon,
  Bars3Icon,
  XMarkIcon,
  SparklesIcon,
  CpuChipIcon
} from '@heroicons/react/24/outline'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { Card, CardContent } from './ui/card'
import { Separator } from './ui/separator'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import { cn } from '@/lib/utils'
import { 
  brandingStyles, 
  layoutStyles, 
  spacing, 
  colors, 
  typography,
  iconStyles,
  effects
} from '../utils/designSystem'

function MainSidebar({ onDCAToggle, isOpen, onToggle }) {
  const navigate = useNavigate()
  const location = useLocation()

  const navigationItems = [
    {
      id: 'dashboard',
      name: 'Dashboard',
      icon: HomeIcon,
      path: '/',
      action: () => navigate('/'),
      description: 'Real-time market indicators and analysis',
      shortcut: '⌘1'
    },
    {
      id: 'portfolio',
      name: 'Portfolio',
      icon: ChartPieIcon,
      path: '/portfolio',
      action: () => navigate('/portfolio'),
      description: 'Track your holdings and performance',
      shortcut: '⌘2'
    },
    {
      id: 'advanced-analytics',
      name: 'Advanced Analytics',
      icon: CpuChipIcon,
      path: '/advanced',
      action: () => navigate('/advanced'),
      description: 'AI-powered insights and risk analysis',
      shortcut: '⌘3'
    },
    {
      id: 'dca-calculator',
      name: 'DCA Calculator',
      icon: CalculatorIcon,
      action: () => onDCAToggle(),
      description: 'Simulate dollar-cost averaging strategies',
      shortcut: '⌘4'
    }
  ]

  const getActiveItem = () => {
    const currentPath = location.pathname
    const item = navigationItems.find(item => item.path === currentPath)
    return item ? item.id : 'dashboard'
  }

  return (
    <>
      {/* Mobile Overlay */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 lg:hidden"
          onClick={onToggle}
        />
      )}

      {/* Sidebar */}
      <div className={cn(
        "fixed top-0 left-0 h-full",
        colors.cardBg,
        "border-r",
        colors.border,
        "z-50 transition-transform duration-300 ease-in-out flex flex-col w-64 lg:relative lg:translate-x-0",
        isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
      )}>
        {/* Sidebar Header - Consistent with Header */}
        <div className={cn(
          "flex items-center justify-between border-b",
          colors.border,
          spacing.sidebarPadding
        )}>
          <div className={brandingStyles.logoContainer}>
            <div className={brandingStyles.logoIconBackground}>
              <ChartBarIcon className={cn(iconStyles.medium, "text-primary-foreground")} />
            </div>
            <div className={brandingStyles.sidebarLogo}>
              CryptoIQ
            </div>
          </div>
          
          {/* Close button for mobile */}
          <Button
            variant="ghost"
            size="icon"
            onClick={onToggle}
            className={cn(
              "lg:hidden",
              colors.textSecondary,
              "hover:text-foreground"
            )}
          >
            <XMarkIcon className={iconStyles.large} />
          </Button>
        </div>

        {/* Main Content - Navigation Menu */}
        <div className="flex-1 flex flex-col">
          <nav className="p-6 space-y-2 flex-1">
            <TooltipProvider>
              {navigationItems.map((item) => {
                const Icon = item.icon
                const isActive = getActiveItem() === item.id
                
                return (
                  <Tooltip key={item.id} delayDuration={300}>
                    <TooltipTrigger asChild>
                      <Button
                        variant={isActive ? "default" : "ghost"}
                        onClick={item.action}
                        className={cn(
                          "w-full justify-start space-x-4 px-5 py-4 h-auto text-left group relative overflow-hidden",
                          "transition-all duration-200 hover:scale-[1.02]",
                          isActive && "bg-gradient-to-r from-primary/20 to-accent/20 border border-primary/30 shadow-lg"
                        )}
                      >
                        {/* Animated background for active state */}
                        {isActive && (
                          <div className="absolute inset-0 bg-gradient-to-r from-primary/5 to-accent/5 animate-pulse" />
                        )}
                        
                        <Icon className={cn(
                          "w-5 h-5 transition-all duration-200",
                          isActive ? "text-primary scale-110" : "text-muted-foreground group-hover:text-foreground group-hover:scale-110"
                        )} />
                        
                        <div className="flex-1 relative z-10">
                          <span className="font-medium">{item.name}</span>
                          {item.shortcut && (
                            <span className="text-xs text-muted-foreground block">{item.shortcut}</span>
                          )}
                        </div>
                        
                        {/* Active indicator with enhanced animation */}
                        {isActive && (
                          <div className="flex items-center space-x-1">
                            <div className="w-2 h-2 bg-primary rounded-full animate-pulse" />
                            <SparklesIcon className="w-3 h-3 text-primary animate-pulse" />
                          </div>
                        )}
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent side="right" className="ml-2">
                      <div className="text-center">
                        <p className="font-medium">{item.name}</p>
                        <p className="text-xs text-muted-foreground mt-1">{item.description}</p>
                        {item.shortcut && (
                          <p className="text-xs text-accent mt-1">{item.shortcut}</p>
                        )}
                      </div>
                    </TooltipContent>
                  </Tooltip>
                )
              })}
            </TooltipProvider>
          </nav>

          {/* Enhanced Sidebar Footer */}
          <div className="p-6 border-t border-border mt-auto">
            <Card className="bg-gradient-to-br from-muted/30 to-muted/10 border-primary/20">
              <CardContent className="p-4">
                <div className="text-center space-y-2">
                  <div className="flex items-center justify-center space-x-2">
                    <div className="w-6 h-6 bg-gradient-to-r from-primary to-accent rounded flex items-center justify-center">
                      <SparklesIcon className="w-3 h-3 text-primary-foreground" />
                    </div>
                    <div className="font-medium text-foreground">CryptoIQ</div>
                  </div>
                  
                  <div className="flex items-center justify-center space-x-2">
                    <Badge variant="secondary" className="text-xs bg-primary/10 text-primary border-primary/20">
                      v1.0.0
                    </Badge>
                    <div className="w-1 h-1 bg-green-400 rounded-full animate-pulse" />
                    <span className="text-xs text-muted-foreground">Live</span>
                  </div>
                  
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    Advanced crypto market intelligence platform
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </>
  )
}

export default MainSidebar