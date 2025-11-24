'use client'

import React from 'react'

interface Props {
  level: number
  points: number
  streak: number
}

const GamificationProfile: React.FC<Props> = ({ level, points, streak }) => {
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-bold mb-4">Gamification Profile</h3>
      <div className="flex justify-between space-x-6">
        <div className="flex-1 bg-blue-100 rounded p-4 text-center">
          <p className="text-xl font-semibold">{level}</p>
          <p className="text-gray-600">Current Level</p>
        </div>
        <div className="flex-1 bg-green-100 rounded p-4 text-center">
          <p className="text-xl font-semibold">{points}</p>
          <p className="text-gray-600">Points</p>
        </div>
        <div className="flex-1 bg-yellow-100 rounded p-4 text-center">
          <p className="text-xl font-semibold">{streak}</p>
          <p className="text-gray-600">Current Streak</p>
        </div>
      </div>
    </div>
  )
}

export default GamificationProfile
