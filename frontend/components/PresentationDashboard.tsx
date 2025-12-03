'use client'

import React, { useState } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'

export interface Slide {
  id: string
  title: string
  subtitle?: string
  content: React.ReactNode
  backgroundColor?: string
  textColor?: string
}

interface PresentationDashboardProps {
  slides: Slide[]
  title?: string
  currentSlideIndex?: number
  onSlideChange?: (index: number) => void
  showSlideNumbers?: boolean
  autoPlay?: boolean
  autoPlayInterval?: number
}

export default function PresentationDashboard({
  slides,
  title = 'Dashboard Presentation',
  currentSlideIndex = 0,
  onSlideChange,
  showSlideNumbers = true,
  autoPlay = false,
  autoPlayInterval = 5000
}: PresentationDashboardProps) {
  const [currentIndex, setCurrentIndex] = useState(currentSlideIndex)
  const [isAnimating, setIsAnimating] = useState(false)

  React.useEffect(() => {
    if (!autoPlay) return

    const interval = setInterval(() => {
      goToNextSlide()
    }, autoPlayInterval)

    return () => clearInterval(interval)
  }, [autoPlay, autoPlayInterval, currentIndex])

  const currentSlide = slides[currentIndex]

  const goToPreviousSlide = () => {
    if (currentIndex > 0) {
      setIsAnimating(true)
      setTimeout(() => {
        const newIndex = currentIndex - 1
        setCurrentIndex(newIndex)
        onSlideChange?.(newIndex)
        setIsAnimating(false)
      }, 300)
    }
  }

  const goToNextSlide = () => {
    if (currentIndex < slides.length - 1) {
      setIsAnimating(true)
      setTimeout(() => {
        const newIndex = currentIndex + 1
        setCurrentIndex(newIndex)
        onSlideChange?.(newIndex)
        setIsAnimating(false)
      }, 300)
    }
  }

  const goToSlide = (index: number) => {
    setIsAnimating(true)
    setTimeout(() => {
      setCurrentIndex(index)
      onSlideChange?.(index)
      setIsAnimating(false)
    }, 300)
  }

  return (
    <div className="w-full h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-black flex flex-col">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-900 to-blue-800 text-white px-8 py-4 shadow-lg border-b border-blue-700">
        <h1 className="text-3xl font-bold">{title}</h1>
        <p className="text-blue-200 text-sm mt-1">
          Slide {currentIndex + 1} of {slides.length}
        </p>
      </div>

      {/* Main Slide Area */}
      <div className="flex-1 flex items-center justify-center p-8 overflow-hidden">
        <div
          className={`w-full h-full max-w-6xl transition-all duration-300 ${
            isAnimating ? 'opacity-50 scale-95' : 'opacity-100 scale-100'
          }`}
        >
          {/* Slide Content - PowerPoint Style */}
          <div
            className={`w-full h-full rounded-xl shadow-2xl overflow-hidden flex flex-col bg-gradient-to-br ${
              currentSlide.backgroundColor || 'from-white via-blue-50 to-blue-100'
            }`}
          >
            {/* Slide Header */}
            <div className="bg-gradient-to-r from-blue-600 to-blue-700 text-white px-12 py-8 border-b-4 border-blue-800">
              <h2 className="text-4xl font-bold mb-2">{currentSlide.title}</h2>
              {currentSlide.subtitle && (
                <p className="text-blue-100 text-xl">{currentSlide.subtitle}</p>
              )}
            </div>

            {/* Slide Content */}
            <div className={`flex-1 overflow-auto px-12 py-8 text-gray-800`}>
              {currentSlide.content}
            </div>

            {/* Slide Footer */}
            <div className="bg-gray-100 border-t border-gray-300 px-12 py-4 flex justify-between items-center">
              <div className="text-sm text-gray-600">
                {new Date().toLocaleDateString('en-IN', {
                  weekday: 'long',
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric'
                })}
              </div>
              {showSlideNumbers && (
                <div className="text-sm font-semibold text-gray-700">
                  Slide {currentIndex + 1}/{slides.length}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Navigation Controls */}
      <div className="bg-gray-800 border-t border-gray-700 px-8 py-6 flex items-center justify-between">
        {/* Previous Button */}
        <button
          onClick={goToPreviousSlide}
          disabled={currentIndex === 0}
          className={`flex items-center gap-2 px-6 py-3 rounded-lg font-semibold transition-all ${
            currentIndex === 0
              ? 'bg-gray-700 text-gray-500 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-700 active:scale-95'
          }`}
        >
          <ChevronLeft size={20} />
          Previous
        </button>

        {/* Slide Thumbnails/Dots */}
        <div className="flex gap-2 mx-4 flex-wrap justify-center">
          {slides.map((slide, idx) => (
            <button
              key={slide.id}
              onClick={() => goToSlide(idx)}
              className={`w-3 h-3 rounded-full transition-all ${
                idx === currentIndex
                  ? 'bg-blue-500 w-8'
                  : 'bg-gray-600 hover:bg-gray-500'
              }`}
              title={`Go to slide ${idx + 1}`}
            />
          ))}
        </div>

        {/* Next Button */}
        <button
          onClick={goToNextSlide}
          disabled={currentIndex === slides.length - 1}
          className={`flex items-center gap-2 px-6 py-3 rounded-lg font-semibold transition-all ${
            currentIndex === slides.length - 1
              ? 'bg-gray-700 text-gray-500 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-700 active:scale-95'
          }`}
        >
          Next
          <ChevronRight size={20} />
        </button>
      </div>

      {/* Keyboard Navigation Hint */}
      <div className="bg-gray-900 text-gray-400 text-xs px-8 py-2 text-center border-t border-gray-700">
        Use ← → arrows or click to navigate • ESC to exit presentation mode
      </div>
    </div>
  )
}
