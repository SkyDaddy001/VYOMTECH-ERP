'use client'

import React from 'react'

interface Challenge {
  id: number
  name: string
  description: string
  status: string
  progress?: number
}

interface Props {
  challenges: Challenge[]
}

const Challenges: React.FC<Props> = ({ challenges }) => {
  if (!challenges || challenges.length === 0) {
    return <p>No active challenges at the moment.</p>
  }
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-bold mb-4">Active Challenges</h3>
      <ul>
        {challenges.map((challenge) => (
          <li key={challenge.id} className="mb-4">
            <h4 className="font-semibold">{challenge.name}</h4>
            <p className="text-sm text-gray-500">{challenge.description}</p>
            {challenge.progress !== undefined && (
              <div className="w-full bg-gray-200 rounded-full h-4 mt-2">
                <div
                  className="bg-blue-600 h-4 rounded-full"
                  style={{ width: `${challenge.progress}%` }}
                ></div>
              </div>
            )}
            <p className="text-xs text-gray-400 mt-1 capitalize">{challenge.status}</p>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default Challenges
