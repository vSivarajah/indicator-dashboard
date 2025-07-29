// Enhanced Design System Constants and Utilities

export const colors = {
  // Background colors - consistent across all components
  cardBg: 'bg-card/50 backdrop-blur-sm',
  cardBgHover: 'hover:bg-card/70',
  gradientBg: 'bg-gradient-to-br from-background via-background to-muted/20',
  
  // Professional border colors
  border: 'border-border',
  borderAccent: 'border-accent/30',
  borderHover: 'hover:border-accent/50',
  
  // Consistent text colors using CSS variables
  textPrimary: 'text-foreground',
  textSecondary: 'text-muted-foreground',
  textMuted: 'text-muted-foreground/70',
  textSuccess: 'text-green-400',
  textDanger: 'text-red-400',
  textWarning: 'text-yellow-400',
  textAccent: 'text-accent',
  textGold: 'text-yellow-500',
  
  // Professional accent palette
  accent: 'text-accent',
  accentLight: 'text-accent/80',
  accentDark: 'text-accent',
  accentSubtle: 'text-accent/60',
  
  // Background accents (very subtle)
  bgAccent: 'bg-accent',
  bgAccentLight: 'bg-accent/80',
  bgAccentSubtle: 'bg-accent/5',
  
  // Consistent button styling
  buttonPrimary: 'bg-primary text-primary-foreground',
  buttonSecondary: 'bg-secondary text-secondary-foreground',
  buttonOutline: 'border border-accent/50 text-accent bg-background',
}

export const spacing = {
  // Standardized padding system
  cardPadding: 'p-6',
  cardPaddingSmall: 'p-4',
  cardPaddingLarge: 'p-8',
  sidebarPadding: 'p-6',
  headerPadding: 'px-6 py-4',
  
  // Standardized margins
  sectionMargin: 'mb-8',
  elementMargin: 'mb-4',
  smallMargin: 'mb-2',
  
  // Grid and layout gaps
  gridGap: 'gap-6',
  gridGapSmall: 'gap-4',
  gridGapLarge: 'gap-8',
  
  // Consistent spacing between elements
  spaceBetweenItems: 'space-y-4',
  spaceBetweenItemsSmall: 'space-y-2',
  spaceBetweenItemsLarge: 'space-y-6',
}

export const effects = {
  // Subtle professional shadows
  cardShadow: 'hover:shadow-lg hover:shadow-black/20',
  subtleShadow: 'shadow-sm shadow-black/10',
  
  // Minimal transforms
  scaleUp: 'hover:scale-[1.01]',
  scaleUpSmall: 'hover:scale-[1.005]',
  
  // Smooth transitions
  transition: 'transition-all duration-200',
  transitionFast: 'transition-all duration-150',
  transitionSlow: 'transition-all duration-300',
  
  // Professional border radius
  rounded: 'rounded-lg',
  roundedSmall: 'rounded-md',
  roundedLarge: 'rounded-xl',
  
  // Professional border effects
  professionalBorder: 'border border-gray-700/50',
  accentBorder: 'border border-orange-500/30',
  accentBorderHover: 'hover:border-orange-500/50 hover:shadow-orange-500/10',
  focusRing: 'focus:ring-2 focus:ring-orange-500/30 focus:border-orange-500/50',
  
  // Gold theme specific effects
  goldGlowHover: 'hover:shadow-lg hover:shadow-yellow-500/20 hover:border-yellow-500/30',
}

export const typography = {
  // Standardized heading hierarchy
  heroTitle: 'text-4xl md:text-5xl font-bold leading-tight',
  pageTitle: 'text-3xl md:text-4xl font-bold leading-tight',
  sectionTitle: 'text-2xl md:text-3xl font-semibold leading-tight',
  cardTitle: 'text-lg font-semibold leading-tight',
  cardSubtitle: 'text-base font-medium leading-tight',
  
  // Consistent body text
  bodyLarge: 'text-base font-normal leading-relaxed',
  body: 'text-sm font-normal leading-relaxed',
  bodySmall: 'text-xs font-normal leading-relaxed',
  caption: 'text-xs font-medium leading-relaxed',
  
  // Standardized font weights
  fontBold: 'font-bold',
  fontSemibold: 'font-semibold',
  fontMedium: 'font-medium',
  fontNormal: 'font-normal',
  
  // Value displays (for metrics, prices, etc.)
  valueDisplay: 'text-2xl md:text-3xl font-bold leading-none',
  valueDisplayLarge: 'text-3xl md:text-4xl font-bold leading-none',
  valueChange: 'text-sm font-medium leading-none',
}

// Enhanced component style builders
export const buildCardStyles = (hover = true, variant = 'default', size = 'default') => {
  const baseStyles = `${colors.cardBg} ${colors.border} ${effects.rounded} ${effects.transition} group relative`
  
  const hoverEffects = hover ? `${colors.cardBgHover} ${colors.borderHover} ${effects.cardShadow} ${effects.scaleUpSmall}` : ''
  
  const variants = {
    default: '',
    accent: `${colors.borderAccent}`,
    professional: `${colors.border}`,
    elevated: 'shadow-lg',
  }
  
  const sizes = {
    small: spacing.cardPaddingSmall,
    default: spacing.cardPadding,
    large: spacing.cardPaddingLarge,
  }
  
  return `${baseStyles} ${sizes[size]} ${hoverEffects} ${variants[variant]}`
}

// Standardized card height system for consistent layouts
export const cardHeights = {
  indicatorCard: 'min-h-[180px]', // Fixed minimum height for all indicator cards
  chartCard: 'h-[400px]',         // Standard chart container height
  portfolioCard: 'min-h-[200px]', // Portfolio card minimum height
  dcaSidebarCard: 'min-h-[160px]', // DCA sidebar card height
}

export const buildButtonStyles = (variant = 'primary', size = 'medium') => {
  const baseStyles = `inline-flex items-center justify-center font-medium ${effects.rounded} ${effects.transition}`
  
  const variants = {
    primary: `${colors.buttonPrimary} hover:bg-orange-600 ${effects.subtleShadow} ${effects.scaleUpSmall}`,
    secondary: `${colors.buttonSecondary} hover:bg-gray-600 ${colors.borderHover}`,
    outline: `${colors.buttonOutline} hover:bg-orange-500/10 ${effects.focusRing}`,
    ghost: `${colors.textAccent} hover:text-orange-400 hover:bg-orange-500/5`,
  }
  
  const sizes = {
    small: 'px-3 py-1.5 text-sm',
    medium: 'px-4 py-2 text-sm',
    large: 'px-6 py-3 text-base',
  }
  
  return `${baseStyles} ${variants[variant]} ${sizes[size]}`
}

// Professional text styling - clean, no gradients
export const buildProfessionalText = (variant = 'primary') => {
  const variants = {
    primary: colors.textPrimary,
    secondary: colors.textSecondary,
    accent: colors.textAccent,
    muted: colors.textMuted,
    heading: `${colors.textPrimary} font-semibold`,
    subheading: `${colors.textSecondary} font-medium`,
  }
  
  return variants[variant] || variants.primary
}

// Gradient text styling function
export const buildGradientText = (variant = 'default') => {
  const variants = {
    default: 'bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent',
    goldBright: 'bg-gradient-to-r from-yellow-400 via-yellow-500 to-orange-500 bg-clip-text text-transparent',
    blue: 'bg-gradient-to-r from-blue-400 to-blue-600 bg-clip-text text-transparent',
    purple: 'bg-gradient-to-r from-purple-400 to-purple-600 bg-clip-text text-transparent',
    accent: 'bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent',
  }
  
  return variants[variant] || variants.default
}

// Animation utilities
export const animations = {
  fadeIn: 'animate-in fade-in duration-500',
  slideUp: 'animate-in slide-in-from-bottom-4 duration-500',
  slideLeft: 'animate-in slide-in-from-right-4 duration-500',
  pulse: 'animate-pulse',
  spin: 'animate-spin',
  bounce: 'animate-bounce',
  // Gold-specific animations
  goldShimmer: 'animate-gold-shimmer',
  goldPulse: 'animate-gold-pulse',
  goldGlow: 'animate-glow',
  goldFloat: 'animate-float',
  countUp: 'animate-count-up',
  slideInStagger: 'animate-slide-in-stagger',
}

// Responsive breakpoint utilities
export const breakpoints = {
  mobile: 'max-w-sm',
  tablet: 'max-w-4xl',
  desktop: 'max-w-7xl',
  wide: 'max-w-screen-2xl',
}

// Gold theme specific utilities
export const goldTheme = {
  // Status indicators with gold theme
  positive: 'text-green-400 bg-green-400/10 border border-green-400/30',
  negative: 'text-red-400 bg-red-400/10 border border-red-400/30',
  neutral: 'text-yellow-400 bg-yellow-400/10 border border-yellow-400/30',
  warning: 'text-yellow-300 bg-yellow-300/10 border border-yellow-300/30',
  
  // Gold accents for different risk levels
  riskLow: 'text-green-400 border-green-400/30 bg-green-400/5',
  riskMedium: 'text-yellow-400 border-yellow-400/30 bg-yellow-400/5',
  riskHigh: 'text-orange-400 border-orange-400/30 bg-orange-400/5',
  riskExtreme: 'text-red-400 border-red-400/30 bg-red-400/5',
  
  // Market cycle indicators
  bear: 'text-red-400 bg-red-400/10',
  earlyBull: 'text-yellow-300 bg-yellow-300/10',
  midBull: 'text-yellow-400 bg-yellow-400/10',
  lateBull: 'text-orange-400 bg-orange-400/10',
  top: 'text-red-500 bg-red-500/10',
}

// Icon styling utilities
export const iconStyles = {
  small: 'w-4 h-4',
  medium: 'w-5 h-5',
  large: 'w-6 h-6',
  xl: 'w-8 h-8',
  accent: 'text-accent',
  accentHover: 'hover:text-accent/80',
}

// Standardized logo and branding system
export const brandingStyles = {
  // Logo container consistent sizing
  logoContainer: 'flex items-center space-x-3',
  logoIconSize: 'w-8 h-8',
  logoIconBackground: 'bg-gradient-to-r from-primary to-accent rounded-lg flex items-center justify-center',
  logoText: 'text-xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent',
  logoSubtext: 'text-sm text-muted-foreground',
  
  // Consistent branding across header and sidebar
  headerLogo: 'text-2xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent',
  sidebarLogo: 'text-lg font-bold text-foreground',
}

// Layout system for consistent spacing and alignment
export const layoutStyles = {
  // Page containers
  pageContainer: 'container mx-auto px-4 py-8',
  pageContainerLarge: 'max-w-7xl mx-auto px-6 py-8',
  
  // Grid systems
  dashboardGrid: 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4',
  portfolioGrid: 'grid grid-cols-1 lg:grid-cols-4',
  chartGrid: 'grid grid-cols-1 lg:grid-cols-2',
  
  // Flex layouts
  headerLayout: 'flex items-center justify-between',
  cardHeaderLayout: 'flex items-start justify-between',
  
  // Section spacing
  sectionSpacing: 'space-y-8',
  cardSpacing: spacing.gridGap,
}