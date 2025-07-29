import { useState } from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Header from './components/Header'
import Dashboard from './components/Dashboard'
import Portfolio from './components/Portfolio'
import AdvancedAnalytics from './components/AdvancedAnalytics'
import DCASidebar from './components/DCASidebar'
import MainSidebar from './components/MainSidebar'

function App() {
  const [isDCASidebarOpen, setIsDCASidebarOpen] = useState(false)
  const [isMainSidebarOpen, setIsMainSidebarOpen] = useState(false)

  const toggleDCASidebar = () => {
    setIsDCASidebarOpen(!isDCASidebarOpen)
  }

  const closeDCASidebar = () => {
    setIsDCASidebarOpen(false)
  }

  const toggleMainSidebar = () => {
    setIsMainSidebarOpen(!isMainSidebarOpen)
  }

  return (
    <Router>
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted">
        <div className="flex">
          {/* Main Navigation Sidebar */}
          <MainSidebar 
            onDCAToggle={toggleDCASidebar}
            isOpen={isMainSidebarOpen}
            onToggle={toggleMainSidebar}
          />
          
          {/* Main Content Area */}
          <div className="flex-1 flex flex-col min-h-screen lg:ml-0">
            <Header onMenuToggle={toggleMainSidebar} />
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/portfolio" element={<Portfolio />} />
              <Route path="/advanced" element={<AdvancedAnalytics />} />
            </Routes>
          </div>
        </div>
        
        {/* DCA Calculator Sidebar (overlay) */}
        <DCASidebar isOpen={isDCASidebarOpen} onClose={closeDCASidebar} />
      </div>
    </Router>
  )
}

export default App
