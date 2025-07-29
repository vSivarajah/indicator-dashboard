import { useBackendHealth } from '../hooks/useIndicators'
import { Bars3Icon, Cog6ToothIcon, MoonIcon, SunIcon } from '@heroicons/react/24/outline'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import { Switch } from './ui/switch'
import { cn } from '@/lib/utils'
import { useState } from 'react'
import { 
  brandingStyles, 
  layoutStyles, 
  spacing, 
  colors, 
  typography,
  iconStyles 
} from '../utils/designSystem'

function Header({ onMenuToggle }) {
  const { isHealthy, isChecking } = useBackendHealth()
  const [isDarkMode, setIsDarkMode] = useState(true)
  const [autoRefresh, setAutoRefresh] = useState(true)
  
  return (
    <header className={cn(
      "border-b",
      colors.border,
      colors.cardBg
    )}>
      <div className={cn(layoutStyles.pageContainerLarge, "py-0")}>
        <div className={cn(layoutStyles.headerLayout, spacing.headerPadding)}>
          <div className="flex items-center space-x-4">
            {/* Mobile Menu Toggle */}
            <Button
              variant="ghost"
              size="icon"
              onClick={onMenuToggle}
              className={cn(
                "lg:hidden",
                colors.textSecondary,
                "hover:text-foreground"
              )}
            >
              <Bars3Icon className={iconStyles.large} />
            </Button>
            
            {/* Consistent Logo Branding */}
            <div className={brandingStyles.logoContainer}>
              <div className={brandingStyles.headerLogo}>
                CryptoIQ
              </div>
              <span className={cn(
                "hidden sm:block",
                brandingStyles.logoSubtext
              )}>
                Market Intelligence
              </span>
            </div>
          </div>
          
          <div className="flex items-center space-x-4">
            {/* Live Data Status with Tooltip */}
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge 
                    variant={isHealthy ? 'success' : isChecking ? 'warning' : 'destructive'}
                    className={cn(
                      "px-4 py-2 bg-muted/50 backdrop-blur border cursor-pointer transition-all hover:scale-105",
                      isHealthy && "border-green-400/30",
                      isChecking && "border-yellow-400/30",
                      !isHealthy && !isChecking && "border-red-400/30"
                    )}
                  >
                    <div className="flex items-center space-x-2">
                      <div className={cn(
                        "w-2 h-2 rounded-full",
                        isChecking && "bg-yellow-400 animate-pulse",
                        isHealthy && "bg-green-400 animate-pulse",
                        !isHealthy && !isChecking && "bg-red-400"
                      )} />
                      <span className="text-sm">
                        {isChecking ? 'Checking...' : isHealthy ? 'Live Data' : 'Offline'}
                      </span>
                    </div>
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>
                  <p>
                    {isHealthy 
                      ? 'Backend connected - Data updating every 5 minutes' 
                      : isChecking 
                      ? 'Checking backend connection...' 
                      : 'Backend offline - Using fallback data'
                    }
                  </p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>

            {/* Settings Dropdown */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-foreground">
                  <Cog6ToothIcon className="w-5 h-5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel>Settings</DropdownMenuLabel>
                <DropdownMenuSeparator />
                
                <div className="flex items-center justify-between px-2 py-1.5">
                  <div className="flex items-center space-x-2">
                    {isDarkMode ? <MoonIcon className="w-4 h-4" /> : <SunIcon className="w-4 h-4" />}
                    <span className="text-sm">Dark Mode</span>
                  </div>
                  <Switch 
                    checked={isDarkMode} 
                    onCheckedChange={setIsDarkMode}
                    className="data-[state=checked]:bg-primary"
                  />
                </div>
                
                <div className="flex items-center justify-between px-2 py-1.5">
                  <span className="text-sm">Auto Refresh</span>
                  <Switch 
                    checked={autoRefresh} 
                    onCheckedChange={setAutoRefresh}
                    className="data-[state=checked]:bg-primary"
                  />
                </div>
                
                <DropdownMenuSeparator />
                <DropdownMenuItem className="text-sm">
                  <span>Refresh Rate: 5 minutes</span>
                </DropdownMenuItem>
                <DropdownMenuItem className="text-sm">
                  <span>Data Sources: 3 active</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>
    </header>
  )
}

export default Header