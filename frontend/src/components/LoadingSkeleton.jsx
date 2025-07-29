import { effects, colors } from '../utils/designSystem'

function LoadingSkeleton({ type = 'card', count = 1, className = '' }) {
  const skeletons = Array.from({ length: count }, (_, index) => index)

  const CardSkeleton = ({ delay = 0 }) => (
    <div 
      className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-6 animate-pulse`}
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-start justify-between mb-4">
        <div className="space-y-2">
          <div className="h-4 bg-gray-700 rounded w-24"></div>
          <div className="h-3 bg-gray-700 rounded w-32"></div>
        </div>
        <div className="h-6 bg-gray-700 rounded-full w-16"></div>
      </div>
      
      <div className="flex items-baseline space-x-3 mb-3">
        <div className="h-8 bg-gray-700 rounded w-20"></div>
        <div className="h-4 bg-gray-700 rounded w-12"></div>
      </div>
      
      <div className="h-3 bg-gray-700 rounded w-full"></div>
    </div>
  )

  const ChartSkeleton = ({ delay = 0 }) => (
    <div 
      className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-6 animate-pulse`}
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-center justify-between mb-6">
        <div className="space-y-2">
          <div className="h-5 bg-gray-700 rounded w-32"></div>
          <div className="h-3 bg-gray-700 rounded w-48"></div>
        </div>
        <div className="flex space-x-1 bg-gray-900/50 rounded-lg p-1">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-6 w-8 bg-gray-700 rounded"></div>
          ))}
        </div>
      </div>
      
      <div className="h-64 bg-gray-700 rounded-lg mb-4"></div>
      
      <div className="flex items-center justify-between">
        <div className="flex space-x-4">
          <div className="h-3 bg-gray-700 rounded w-24"></div>
          <div className="h-3 bg-gray-700 rounded w-20"></div>
        </div>
        <div className="h-3 bg-gray-700 rounded w-16"></div>
      </div>
    </div>
  )

  const HeroSkeleton = () => (
    <div className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-6 animate-pulse`}>
      <div className="text-center mb-8">
        <div className="h-12 bg-gray-700 rounded w-64 mx-auto mb-4"></div>
        <div className="h-4 bg-gray-700 rounded w-96 mx-auto"></div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[1, 2, 3].map((i) => (
          <div key={i} className={`${colors.cardBg} ${colors.border} ${effects.rounded} p-6`}>
            <div className="flex items-center justify-between mb-4">
              <div className="h-12 w-12 bg-gray-700 rounded"></div>
              <div className="h-6 bg-gray-700 rounded-full w-16"></div>
            </div>
            <div className="space-y-2">
              <div className="h-4 bg-gray-700 rounded w-20"></div>
              <div className="h-8 bg-gray-700 rounded w-24"></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )

  if (type === 'hero') {
    return <HeroSkeleton />
  }

  if (type === 'chart') {
    return (
      <div className={`grid grid-cols-1 lg:grid-cols-2 gap-6 ${className}`}>
        {skeletons.map((_, index) => (
          <ChartSkeleton key={index} delay={index * 100} />
        ))}
      </div>
    )
  }

  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 ${className}`}>
      {skeletons.map((_, index) => (
        <CardSkeleton key={index} delay={index * 100} />
      ))}
    </div>
  )
}

export default LoadingSkeleton