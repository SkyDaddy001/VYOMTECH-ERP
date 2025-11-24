'use client'

import React from 'react'

interface LeaderboardEntry {
  id: number
  userName: string
  rank: number
  points: number
  badgesCount: number
  streakDays: number
}

interface Props {
  leaderboard: LeaderboardEntry[]
}

const Leaderboard: React.FC<Props> = ({ leaderboard }) => {
  if (!leaderboard || leaderboard.length === 0) {
    return <p>No leaderboard data available.</p>
  }
  return (
    <div className="bg-white rounded-lg shadow-md p-6 overflow-x-auto">
      <h3 className="text-lg font-bold mb-4">Leaderboard</h3>
      <table className="min-w-full table-auto border-collapse border border-gray-200">
        <thead>
          <tr className="bg-gray-100">
            <th className="border border-gray-300 px-4 py-2 text-left">Rank</th>
            <th className="border border-gray-300 px-4 py-2 text-left">User</th>
            <th className="border border-gray-300 px-4 py-2 text-left">Points</th>
            <th className="border border-gray-300 px-4 py-2 text-left">Badges</th>
            <th className="border border-gray-300 px-4 py-2 text-left">Streak</th>
          </tr>
        </thead>
        <tbody>
          {leaderboard.map((entry) => (
            <tr key={entry.id} className="hover:bg-gray-50">
              <td className="border border-gray-300 px-4 py-2">{entry.rank}</td>
              <td className="border border-gray-300 px-4 py-2">{entry.userName}</td>
              <td className="border border-gray-300 px-4 py-2">{entry.points}</td>
              <td className="border border-gray-300 px-4 py-2">{entry.badgesCount}</td>
              <td className="border border-gray-300 px-4 py-2">{entry.streakDays} days</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Leaderboard
