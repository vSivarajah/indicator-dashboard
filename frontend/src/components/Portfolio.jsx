import { useState, useEffect, memo } from 'react'
import { PlusIcon, PencilIcon, TrashIcon, ChartPieIcon } from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/table'
import { Badge } from './ui/badge'
import { Separator } from './ui/separator'
import { cn } from '../lib/utils'
import PortfolioAllocationChart from './PortfolioAllocationChart'
import PortfolioPerformanceChart from './PortfolioPerformanceChart'
import HoldingsComparisonChart from './HoldingsComparisonChart'
import { 
  buildCardStyles, 
  cardHeights, 
  colors, 
  typography, 
  spacing, 
  effects,
  layoutStyles,
  iconStyles
} from '../utils/designSystem'

const Portfolio = memo(function Portfolio() {
  const [portfolios, setPortfolios] = useState([])
  const [selectedPortfolio, setSelectedPortfolio] = useState(null)
  const [portfolioSummary, setPortfolioSummary] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showAddHoldingModal, setShowAddHoldingModal] = useState(false)

  // Form states
  const [portfolioName, setPortfolioName] = useState('')
  const [holdingForm, setHoldingForm] = useState({
    symbol: '',
    amount: '',
    averagePrice: ''
  })

  // Fetch user portfolios on component mount
  useEffect(() => {
    fetchPortfolios()
  }, [])

  // Fetch portfolio summary when selected portfolio changes
  useEffect(() => {
    if (selectedPortfolio) {
      fetchPortfolioSummary(selectedPortfolio.id)
    }
  }, [selectedPortfolio])

  const fetchPortfolios = async () => {
    setLoading(true)
    try {
      const response = await fetch('http://localhost:8080/api/v1/portfolios?user_id=default_user')
      const data = await response.json()
      
      if (response.ok) {
        setPortfolios(data.portfolios || [])
        if (data.portfolios && data.portfolios.length > 0) {
          setSelectedPortfolio(data.portfolios[0])
        }
      } else {
        setError(data.error || 'Failed to fetch portfolios')
      }
    } catch (err) {
      setError('Network error: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const fetchPortfolioSummary = async (portfolioId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/v1/portfolios/${portfolioId}/summary`)
      const data = await response.json()
      
      if (response.ok) {
        setPortfolioSummary(data)
      } else {
        console.error('Failed to fetch portfolio summary:', data.error)
      }
    } catch (err) {
      console.error('Error fetching portfolio summary:', err)
    }
  }

  const createPortfolio = async () => {
    if (!portfolioName.trim()) return

    try {
      const response = await fetch('http://localhost:8080/api/v1/portfolios', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: 'default_user',
          name: portfolioName
        })
      })

      const data = await response.json()
      
      if (response.ok) {
        setPortfolioName('')
        setShowCreateModal(false)
        fetchPortfolios() // Refresh the list
      } else {
        setError(data.error || 'Failed to create portfolio')
      }
    } catch (err) {
      setError('Network error: ' + err.message)
    }
  }

  const addHolding = async () => {
    if (!selectedPortfolio || !holdingForm.symbol || !holdingForm.amount || !holdingForm.averagePrice) {
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/v1/portfolios/${selectedPortfolio.id}/holdings`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          symbol: holdingForm.symbol.toUpperCase(),
          amount: parseFloat(holdingForm.amount),
          average_price: parseFloat(holdingForm.averagePrice)
        })
      })

      const data = await response.json()
      
      if (response.ok) {
        setHoldingForm({ symbol: '', amount: '', averagePrice: '' })
        setShowAddHoldingModal(false)
        fetchPortfolios() // Refresh portfolios
        fetchPortfolioSummary(selectedPortfolio.id) // Refresh summary
      } else {
        setError(data.error || 'Failed to add holding')
      }
    } catch (err) {
      setError('Network error: ' + err.message)
    }
  }

  const removeHolding = async (holdingId) => {
    if (!selectedPortfolio) return

    try {
      const response = await fetch(`http://localhost:8080/api/v1/portfolios/${selectedPortfolio.id}/holdings/${holdingId}`, {
        method: 'DELETE'
      })

      if (response.ok) {
        fetchPortfolios() // Refresh portfolios
        fetchPortfolioSummary(selectedPortfolio.id) // Refresh summary
      } else {
        const data = await response.json()
        setError(data.error || 'Failed to remove holding')
      }
    } catch (err) {
      setError('Network error: ' + err.message)
    }
  }

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value || 0)
  }

  const formatPercentage = (value) => {
    return `${value >= 0 ? '+' : ''}${(value || 0).toFixed(2)}%`
  }

  const getRiskVariant = (risk) => {
    switch (risk?.toLowerCase()) {
      case 'low': return 'success'
      case 'medium': return 'warning'
      case 'high': return 'destructive'
      default: return 'secondary'
    }
  }

  if (loading && portfolios.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse space-y-8">
          <div className="h-8 bg-muted rounded w-64"></div>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map(i => (
              <Card key={i} className="h-64">
                <CardContent className="p-6">
                  <div className="space-y-4">
                    <div className="h-4 bg-muted rounded w-3/4"></div>
                    <div className="h-6 bg-muted rounded w-1/2"></div>
                    <div className="h-4 bg-muted rounded w-full"></div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={layoutStyles.pageContainer}>
      {/* Header */}
      <div className={cn(layoutStyles.headerLayout, spacing.sectionMargin)}>
        <div className={spacing.spaceBetweenItemsSmall}>
          <h1 className={cn(
            typography.heroTitle,
            "bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent"
          )}>
            Portfolio Management
          </h1>
          <p className={cn(
            typography.bodyLarge,
            colors.textSecondary
          )}>
            Track your crypto investments and analyze performance
          </p>
        </div>
        <Dialog open={showCreateModal} onOpenChange={setShowCreateModal}>
          <DialogTrigger asChild>
            <Button variant="default" className={cn(
              "flex items-center space-x-2",
              colors.buttonPrimary,
              effects.transition,
              effects.scaleUpSmall
            )}>
              <PlusIcon className={iconStyles.medium} />
              <span>New Portfolio</span>
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create New Portfolio</DialogTitle>
              <DialogDescription>
                Enter a name for your new portfolio to start tracking your investments.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <Input
                placeholder="Portfolio name"
                value={portfolioName}
                onChange={(e) => setPortfolioName(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && createPortfolio()}
              />
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowCreateModal(false)}>
                Cancel
              </Button>
              <Button onClick={createPortfolio} disabled={!portfolioName.trim()}>
                Create
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <Card className="mb-6 border-destructive bg-destructive/10">
          <CardContent className="p-4">
            <p className="text-destructive">{error}</p>
            <Button 
              variant="ghost"
              size="sm"
              onClick={() => setError(null)}
              className="text-destructive hover:text-destructive/80 h-auto p-0 mt-2"
            >
              Dismiss
            </Button>
          </CardContent>
        </Card>
      )}

      {portfolios.length === 0 ? (
        <Card className="text-center py-12">
          <CardContent className="space-y-6">
            <ChartPieIcon className="w-16 h-16 mx-auto text-muted-foreground" />
            <div className="space-y-2">
              <CardTitle className="text-xl">No Portfolios Yet</CardTitle>
              <CardDescription>Create your first portfolio to start tracking your investments</CardDescription>
            </div>
            <Button
              variant="gradient"
              onClick={() => setShowCreateModal(true)}
              className="px-6 py-3"
            >
              Create Portfolio
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className={cn(layoutStyles.portfolioGrid, spacing.gridGap)}>
          {/* Portfolio Selector */}
          <div className="lg:col-span-1">
            <h2 className={cn(
              typography.cardTitle,
              colors.textPrimary,
              spacing.elementMargin
            )}>Your Portfolios</h2>
            <div className={spacing.spaceBetweenItems}>
              {portfolios.map((portfolio) => (
                <Card
                  key={portfolio.id}
                  onClick={() => setSelectedPortfolio(portfolio)}
                  className={cn(
                    buildCardStyles(true, 'professional'),
                    cardHeights.portfolioCard,
                    "cursor-pointer",
                    selectedPortfolio?.id === portfolio.id
                      ? "bg-gradient-to-r from-primary to-accent text-primary-foreground"
                      : colors.cardBg
                  )}
                >
                  <div className={spacing.cardPaddingSmall}>
                    <h3 className={cn(
                      typography.cardSubtitle,
                      selectedPortfolio?.id === portfolio.id 
                        ? "text-primary-foreground" 
                        : colors.textPrimary
                    )}>
                      {portfolio.name}
                    </h3>
                    <p className={cn(
                      typography.body,
                      selectedPortfolio?.id === portfolio.id 
                        ? "text-primary-foreground/75" 
                        : colors.textSecondary
                    )}>
                      {formatCurrency(portfolio.total_value)}
                    </p>
                    <Badge 
                      variant={getRiskVariant(portfolio.risk_level)} 
                      className="text-xs mt-2"
                    >
                      {portfolio.risk_level?.toUpperCase()} RISK
                    </Badge>
                  </div>
                </Card>
              ))}
            </div>
          </div>

          {/* Portfolio Details */}
          <div className="lg:col-span-3">
            {selectedPortfolio ? (
              <div className="space-y-6">
                {/* Portfolio Overview */}
                {portfolioSummary && (
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <Card>
                      <CardContent className="p-4">
                        <h3 className="text-sm text-muted-foreground mb-1">Total Value</h3>
                        <p className="text-2xl font-bold text-foreground">
                          {formatCurrency(portfolioSummary.total_value)}
                        </p>
                      </CardContent>
                    </Card>
                    <Card>
                      <CardContent className="p-4">
                        <h3 className="text-sm text-muted-foreground mb-1">Total P&L</h3>
                        <p className={cn("text-2xl font-bold", portfolioSummary.total_pnl >= 0 ? 'text-green-400' : 'text-red-400')}>
                          {formatCurrency(portfolioSummary.total_pnl)}
                        </p>
                        <p className={cn("text-sm", portfolioSummary.total_pnl_percent >= 0 ? 'text-green-400' : 'text-red-400')}>
                          {formatPercentage(portfolioSummary.total_pnl_percent)}
                        </p>
                      </CardContent>
                    </Card>
                    <Card>
                      <CardContent className="p-4">
                        <h3 className="text-sm text-muted-foreground mb-1">Risk Level</h3>
                        <Badge 
                          variant={getRiskVariant(portfolioSummary.risk_metrics?.overall_risk)}
                          className="text-lg font-semibold"
                        >
                          {portfolioSummary.risk_metrics?.overall_risk?.toUpperCase()}
                        </Badge>
                      </CardContent>
                    </Card>
                    <Card>
                      <CardContent className="p-4">
                        <h3 className="text-sm text-muted-foreground mb-1">Holdings</h3>
                        <p className="text-2xl font-bold text-foreground">
                          {selectedPortfolio.holdings?.length || 0}
                        </p>
                      </CardContent>
                    </Card>
                  </div>
                )}

                {/* Portfolio Charts */}
                {portfolioSummary && (
                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {/* Portfolio Allocation Chart */}
                    <PortfolioAllocationChart 
                      allocationData={portfolioSummary.allocation_by_asset}
                      loading={loading}
                      title="Asset Allocation"
                    />
                    
                    {/* Portfolio Performance Chart */}
                    <PortfolioPerformanceChart 
                      portfolioSummary={portfolioSummary}
                      loading={loading}
                      title="Portfolio Performance"
                    />
                  </div>
                )}

                {/* Holdings Comparison Chart */}
                {selectedPortfolio.holdings && selectedPortfolio.holdings.length > 0 && (
                  <HoldingsComparisonChart 
                    holdings={selectedPortfolio.holdings}
                    loading={loading}
                    title="Holdings Comparison"
                  />
                )}

                {/* Holdings Table */}
                <Card>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-lg">Holdings</CardTitle>
                      <Dialog open={showAddHoldingModal} onOpenChange={setShowAddHoldingModal}>
                        <DialogTrigger asChild>
                          <Button variant="gradient" className="flex items-center space-x-2">
                            <PlusIcon className="w-4 h-4" />
                            <span>Add Holding</span>
                          </Button>
                        </DialogTrigger>
                        <DialogContent>
                          <DialogHeader>
                            <DialogTitle>Add Holding</DialogTitle>
                            <DialogDescription>
                              Add a new cryptocurrency holding to your portfolio.
                            </DialogDescription>
                          </DialogHeader>
                          <div className="space-y-4">
                            <Input
                              placeholder="Symbol (e.g., BTC, ETH)"
                              value={holdingForm.symbol}
                              onChange={(e) => setHoldingForm({...holdingForm, symbol: e.target.value})}
                            />
                            <Input
                              type="number"
                              placeholder="Amount"
                              step="any"
                              value={holdingForm.amount}
                              onChange={(e) => setHoldingForm({...holdingForm, amount: e.target.value})}
                            />
                            <Input
                              type="number"
                              placeholder="Average price ($)"
                              step="any"
                              value={holdingForm.averagePrice}
                              onChange={(e) => setHoldingForm({...holdingForm, averagePrice: e.target.value})}
                            />
                          </div>
                          <DialogFooter>
                            <Button variant="outline" onClick={() => setShowAddHoldingModal(false)}>
                              Cancel
                            </Button>
                            <Button onClick={addHolding} disabled={!holdingForm.symbol || !holdingForm.amount || !holdingForm.averagePrice}>
                              Add Holding
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                    </div>
                  </CardHeader>

                  <CardContent>
                    {selectedPortfolio.holdings && selectedPortfolio.holdings.length > 0 ? (
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead className="text-left">Asset</TableHead>
                            <TableHead className="text-right">Amount</TableHead>
                            <TableHead className="text-right">Avg Price</TableHead>
                            <TableHead className="text-right">Current Price</TableHead>
                            <TableHead className="text-right">Value</TableHead>
                            <TableHead className="text-right">P&L</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {selectedPortfolio.holdings.map((holding) => (
                            <TableRow key={holding.id}>
                              <TableCell>
                                <div className="font-medium">{holding.symbol}</div>
                              </TableCell>
                              <TableCell className="text-right">
                                {holding.amount.toFixed(8)}
                              </TableCell>
                              <TableCell className="text-right">
                                {formatCurrency(holding.average_price)}
                              </TableCell>
                              <TableCell className="text-right">
                                {formatCurrency(holding.current_price)}
                              </TableCell>
                              <TableCell className="text-right">
                                {formatCurrency(holding.value)}
                              </TableCell>
                              <TableCell className="text-right">
                                <div className={cn("font-medium", holding.pnl >= 0 ? 'text-green-400' : 'text-red-400')}>
                                  {formatCurrency(holding.pnl)}
                                </div>
                                <div className={cn("text-xs", holding.pnl >= 0 ? 'text-green-400' : 'text-red-400')}>
                                  {formatPercentage(holding.pnl_percent)}
                                </div>
                              </TableCell>
                              <TableCell className="text-right">
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => removeHolding(holding.id)}
                                  className="text-destructive hover:text-destructive/80 h-8 w-8 p-0"
                                >
                                  <TrashIcon className="w-4 h-4" />
                                </Button>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    ) : (
                      <div className="text-center py-8 space-y-4">
                        <p className="text-muted-foreground">No holdings in this portfolio</p>
                        <Button
                          variant="gradient"
                          onClick={() => setShowAddHoldingModal(true)}
                        >
                          Add First Holding
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>
            ) : (
              <Card className="text-center py-12">
                <CardContent>
                  <p className="text-muted-foreground">Select a portfolio to view details</p>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      )}

    </div>
  )
})

export default Portfolio