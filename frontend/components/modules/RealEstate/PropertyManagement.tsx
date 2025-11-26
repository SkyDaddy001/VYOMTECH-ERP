import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Plus, Edit2, Trash2, Eye, MapPin, Home, Users } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface PropertyProject {
  id: string
  project_name: string
  project_code: string
  location: string
  city: string
  total_units: number
  project_type: string
  status: string
  launch_date: string
  expected_completion: string
  developer_name: string
  noc_status: string
}

interface PropertyUnit {
  id: string
  unit_number: string
  floor: number
  unit_type: string
  facing: string
  carpet_area: number
  sbua: number
  status: string
}

export function PropertyManagement() {
  const [projects, setProjects] = useState<PropertyProject[]>([])
  const [units, setUnits] = useState<PropertyUnit[]>([])
  const [selectedProject, setSelectedProject] = useState<string | null>(null)
  const [showProjectForm, setShowProjectForm] = useState(false)
  const [showUnitForm, setShowUnitForm] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')

  const [projectForm, setProjectForm] = useState({
    project_name: '',
    project_code: '',
    location: '',
    city: '',
    total_units: 0,
    project_type: 'residential',
    status: 'planning',
    launch_date: '',
    expected_completion: '',
    developer_name: '',
    architect_name: ''
  })

  const [unitForm, setUnitForm] = useState({
    unit_number: '',
    floor: 0,
    unit_type: '2BHK',
    facing: 'north',
    carpet_area: 0,
    carpet_area_with_balcony: 0,
    utility_area: 0,
    plinth_area: 0,
    sbua: 0,
    uds_sqft: 0
  })

  useEffect(() => {
    fetchProjects()
  }, [])

  useEffect(() => {
    if (selectedProject) {
      fetchUnits(selectedProject)
    }
  }, [selectedProject])

  const fetchProjects = async () => {
    try {
      const response = await fetch('/api/v1/real-estate/projects', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        }
      })
      const data = await response.json()
      setProjects(data || [])
      if (data && data.length > 0 && !selectedProject) {
        setSelectedProject(data[0].id)
      }
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    }
  }

  const fetchUnits = async (projectId: string) => {
    try {
      const response = await fetch(`/api/v1/real-estate/projects/${projectId}/units`, {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        }
      })
      const data = await response.json()
      setUnits(data || [])
    } catch (error) {
      console.error('Failed to fetch units:', error)
    }
  }

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const response = await fetch('/api/v1/real-estate/projects', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: JSON.stringify(projectForm)
      })

      if (response.ok) {
        fetchProjects()
        setShowProjectForm(false)
        setProjectForm({
          project_name: '',
          project_code: '',
          location: '',
          city: '',
          total_units: 0,
          project_type: 'residential',
          status: 'planning',
          launch_date: '',
          expected_completion: '',
          developer_name: '',
          architect_name: ''
        })
      }
    } catch (error) {
      console.error('Failed to create project:', error)
    }
  }

  const handleCreateUnit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedProject) return

    try {
      const response = await fetch('/api/v1/real-estate/units', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: JSON.stringify({
          ...unitForm,
          project_id: selectedProject,
          block_id: 'default-block'
        })
      })

      if (response.ok) {
        fetchUnits(selectedProject)
        setShowUnitForm(false)
        setUnitForm({
          unit_number: '',
          floor: 0,
          unit_type: '2BHK',
          facing: 'north',
          carpet_area: 0,
          carpet_area_with_balcony: 0,
          utility_area: 0,
          plinth_area: 0,
          sbua: 0,
          uds_sqft: 0
        })
      }
    } catch (error) {
      console.error('Failed to create unit:', error)
    }
  }

  const filteredUnits = units.filter(unit => {
    const matchesSearch = unit.unit_number.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesStatus = statusFilter === 'all' || unit.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const unitStats = {
    total: units.length,
    available: units.filter(u => u.status === 'available').length,
    booked: units.filter(u => u.status === 'booked').length,
    sold: units.filter(u => u.status === 'sold').length
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold text-gray-900">Property Management</h2>
        <Button onClick={() => setShowProjectForm(true)} className="bg-blue-600 hover:bg-blue-700">
          <Plus className="w-4 h-4 mr-2" /> New Project
        </Button>
      </div>

      {/* Projects Section */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {projects.map(project => (
          <Card 
            key={project.id}
            className={`cursor-pointer transition ${selectedProject === project.id ? 'border-blue-600 border-2' : 'hover:shadow-lg'}`}
            onClick={() => setSelectedProject(project.id)}
          >
            <CardContent className="pt-6">
              <h3 className="font-semibold text-sm mb-2">{project.project_name}</h3>
              <div className="space-y-1 text-xs text-gray-600">
                <p><strong>Code:</strong> {project.project_code}</p>
                <p><strong>Units:</strong> {project.total_units}</p>
                <p><strong>Type:</strong> {project.project_type}</p>
                <p className="pt-2">
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    project.status === 'planning' ? 'bg-yellow-100 text-yellow-800' :
                    project.status === 'under_construction' ? 'bg-blue-100 text-blue-800' :
                    project.status === 'ready' ? 'bg-green-100 text-green-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {project.status}
                  </span>
                </p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Project Form */}
      {showProjectForm && (
        <Card className="border-blue-300 bg-blue-50">
          <CardHeader>
            <CardTitle>New Property Project</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreateProject} className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-2">Project Name</label>
                <Input
                  value={projectForm.project_name}
                  onChange={(e) => setProjectForm({...projectForm, project_name: e.target.value})}
                  placeholder="Project name"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Project Code</label>
                <Input
                  value={projectForm.project_code}
                  onChange={(e) => setProjectForm({...projectForm, project_code: e.target.value})}
                  placeholder="Project code"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Location</label>
                <Input
                  value={projectForm.location}
                  onChange={(e) => setProjectForm({...projectForm, location: e.target.value})}
                  placeholder="Location"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">City</label>
                <Input
                  value={projectForm.city}
                  onChange={(e) => setProjectForm({...projectForm, city: e.target.value})}
                  placeholder="City"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Total Units</label>
                <Input
                  type="number"
                  value={projectForm.total_units}
                  onChange={(e) => setProjectForm({...projectForm, total_units: parseInt(e.target.value)})}
                  placeholder="Total units"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Type</label>
                <Select value={projectForm.project_type} onValueChange={(value) => setProjectForm({...projectForm, project_type: value})}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="residential">Residential</SelectItem>
                    <SelectItem value="commercial">Commercial</SelectItem>
                    <SelectItem value="mixed">Mixed</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Developer</label>
                <Input
                  value={projectForm.developer_name}
                  onChange={(e) => setProjectForm({...projectForm, developer_name: e.target.value})}
                  placeholder="Developer name"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Architect</label>
                <Input
                  value={projectForm.architect_name}
                  onChange={(e) => setProjectForm({...projectForm, architect_name: e.target.value})}
                  placeholder="Architect name"
                />
              </div>

              <div className="col-span-2 flex gap-2">
                <Button type="submit" className="bg-green-600 hover:bg-green-700 flex-1">Create Project</Button>
                <Button type="button" onClick={() => setShowProjectForm(false)} className="bg-gray-500 hover:bg-gray-600 flex-1">Cancel</Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Units Section */}
      {selectedProject && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-2xl font-bold">Property Units</h3>
            <Button onClick={() => setShowUnitForm(true)} className="bg-green-600 hover:bg-green-700">
              <Plus className="w-4 h-4 mr-2" /> Add Unit
            </Button>
          </div>

          {/* Unit Stats */}
          <div className="grid grid-cols-4 gap-4">
            <Card>
              <CardContent className="pt-6 text-center">
                <Home className="w-8 h-8 mx-auto mb-2 text-blue-600" />
                <p className="text-3xl font-bold">{unitStats.total}</p>
                <p className="text-sm text-gray-600">Total Units</p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6 text-center">
                <MapPin className="w-8 h-8 mx-auto mb-2 text-green-600" />
                <p className="text-3xl font-bold">{unitStats.available}</p>
                <p className="text-sm text-gray-600">Available</p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6 text-center">
                <Users className="w-8 h-8 mx-auto mb-2 text-yellow-600" />
                <p className="text-3xl font-bold">{unitStats.booked}</p>
                <p className="text-sm text-gray-600">Booked</p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6 text-center">
                <Home className="w-8 h-8 mx-auto mb-2 text-red-600" />
                <p className="text-3xl font-bold">{unitStats.sold}</p>
                <p className="text-sm text-gray-600">Sold</p>
              </CardContent>
            </Card>
          </div>

          {/* Unit Form */}
          {showUnitForm && (
            <Card className="border-green-300 bg-green-50">
              <CardHeader>
                <CardTitle>Add New Unit</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleCreateUnit} className="grid grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium mb-2">Unit Number</label>
                    <Input
                      value={unitForm.unit_number}
                      onChange={(e) => setUnitForm({...unitForm, unit_number: e.target.value})}
                      placeholder="e.g., A-101"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Floor</label>
                    <Input
                      type="number"
                      value={unitForm.floor}
                      onChange={(e) => setUnitForm({...unitForm, floor: parseInt(e.target.value)})}
                      placeholder="Floor number"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Type</label>
                    <Select value={unitForm.unit_type} onValueChange={(value) => setUnitForm({...unitForm, unit_type: value})}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="1BHK">1BHK</SelectItem>
                        <SelectItem value="2BHK">2BHK</SelectItem>
                        <SelectItem value="3BHK">3BHK</SelectItem>
                        <SelectItem value="shop">Shop</SelectItem>
                        <SelectItem value="office">Office</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Facing</label>
                    <Select value={unitForm.facing} onValueChange={(value) => setUnitForm({...unitForm, facing: value})}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="north">North</SelectItem>
                        <SelectItem value="south">South</SelectItem>
                        <SelectItem value="east">East</SelectItem>
                        <SelectItem value="west">West</SelectItem>
                        <SelectItem value="corner">Corner</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Carpet Area (Sq.Ft)</label>
                    <Input
                      type="number"
                      step="0.01"
                      value={unitForm.carpet_area}
                      onChange={(e) => setUnitForm({...unitForm, carpet_area: parseFloat(e.target.value)})}
                      placeholder="Carpet area"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">SBUA (Sq.Ft)</label>
                    <Input
                      type="number"
                      step="0.01"
                      value={unitForm.sbua}
                      onChange={(e) => setUnitForm({...unitForm, sbua: parseFloat(e.target.value)})}
                      placeholder="SBUA"
                    />
                  </div>

                  <div className="col-span-3 flex gap-2">
                    <Button type="submit" className="bg-green-600 hover:bg-green-700 flex-1">Create Unit</Button>
                    <Button type="button" onClick={() => setShowUnitForm(false)} className="bg-gray-500 hover:bg-gray-600 flex-1">Cancel</Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}

          {/* Units Table */}
          <Card>
            <CardHeader>
              <div className="flex justify-between items-center">
                <CardTitle>Units List</CardTitle>
                <div className="flex gap-2">
                  <Input
                    placeholder="Search unit number..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="w-48"
                  />
                  <Select value={statusFilter} onValueChange={setStatusFilter}>
                    <SelectTrigger className="w-40">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Status</SelectItem>
                      <SelectItem value="available">Available</SelectItem>
                      <SelectItem value="booked">Booked</SelectItem>
                      <SelectItem value="sold">Sold</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Unit No</TableHead>
                      <TableHead>Floor</TableHead>
                      <TableHead>Type</TableHead>
                      <TableHead>Facing</TableHead>
                      <TableHead>Carpet Area</TableHead>
                      <TableHead>SBUA</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredUnits.map(unit => (
                      <TableRow key={unit.id}>
                        <TableCell className="font-semibold">{unit.unit_number}</TableCell>
                        <TableCell>{unit.floor}</TableCell>
                        <TableCell>{unit.unit_type}</TableCell>
                        <TableCell>{unit.facing}</TableCell>
                        <TableCell>{unit.carpet_area.toLocaleString()} sq.ft</TableCell>
                        <TableCell>{unit.sbua.toLocaleString()} sq.ft</TableCell>
                        <TableCell>
                          <span className={`px-2 py-1 rounded text-xs font-medium ${
                            unit.status === 'available' ? 'bg-green-100 text-green-800' :
                            unit.status === 'booked' ? 'bg-yellow-100 text-yellow-800' :
                            'bg-blue-100 text-blue-800'
                          }`}>
                            {unit.status}
                          </span>
                        </TableCell>
                        <TableCell>
                          <div className="flex gap-2">
                            <button className="p-1 hover:bg-blue-50 rounded">
                              <Eye className="w-4 h-4 text-blue-600" />
                            </button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}

export default PropertyManagement
