'use client'

import React from 'react'

interface Badge {
  id: number
  name: string
  description: string
  iconUrl: string
  rarity: string
}

interface Props {
  badges: Badge[]
}

const Badges: React.FC<Props> = ({ badges }) => {
  if (!badges || badges.length === 0) {
    return <p>No badges earned yet.</p>
  }
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-bold mb-4">Earned Badges</h3>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {badges.map((badge) => (
          <div key={badge.id} className="border border-gray-200 rounded-lg p-4 flex flex-col items-center">
            <img src={badge.iconUrl} alt={badge.name} className="w-16 h-16 mb-2" />
            <h4 className="font-semibold">{badge.name}</h4>
            <p className="text-sm text-gray-500 text-center">{badge.description}</p>
            <p className="mt-1 text-xs text-gray-400 capitalize">{badge.rarity}</p>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Badges
