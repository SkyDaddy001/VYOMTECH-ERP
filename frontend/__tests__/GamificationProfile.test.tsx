import React from 'react'
import { render, screen } from '@testing-library/react'
import GamificationProfile from '@/components/dashboard/GamificationProfile'

describe('GamificationProfile component', () => {
  it('renders current level, points, and streak', () => {
    render(<GamificationProfile level={5} points={1200} streak={3} />)

    expect(screen.getByText('Current Level')).toBeTruthy()
    expect(screen.getByText('5')).toBeTruthy()

    expect(screen.getByText('Points')).toBeTruthy()
    expect(screen.getByText('1200')).toBeTruthy()

    expect(screen.getByText('Current Streak')).toBeTruthy()
    expect(screen.getByText('3')).toBeTruthy()
  })
})
