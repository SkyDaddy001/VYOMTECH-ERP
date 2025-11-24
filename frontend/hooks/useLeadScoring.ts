import { useState, useCallback } from 'react';

interface LeadScore {
  id: number;
  lead_id: number;
  tenant_id: string;
  source_quality_score: number;
  engagement_score: number;
  conversion_probability: number;
  urgency_score: number;
  overall_score: number;
  score_category: string; // hot, warm, cold, nurture
  previous_score?: number;
  score_change?: number;
  reason_text?: string;
  calculation_method: string;
  last_calculated: string;
  created_at: string;
  updated_at: string;
}

interface ScoreByCategoryResponse {
  category: string;
  count: number;
  leads: LeadScore[];
}

export function useLeadScoring() {
  const [score, setScore] = useState<LeadScore | null>(null);
  const [scoresByCategory, setScoresByCategory] = useState<ScoreByCategoryResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Get a single lead's score
  const getLeadScore = useCallback(async (leadId: number) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/leads/${leadId}/score`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch lead score');
      }

      const data = await response.json();
      setScore(data);
      return data;
    } catch (err: any) {
      const errorMsg = err.message || 'Error fetching lead score';
      setError(errorMsg);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Calculate/recalculate a lead's score
  const calculateScore = useCallback(async (leadId: number) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/leads/${leadId}/score/calculate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) {
        throw new Error('Failed to calculate lead score');
      }

      const data = await response.json();
      setScore(data.score);
      return data.score;
    } catch (err: any) {
      const errorMsg = err.message || 'Error calculating lead score';
      setError(errorMsg);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Get all leads in a specific category
  const getLeadsByCategory = useCallback(async (category: string, limit: number = 100) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(
        `/api/v1/leads/scores/category/${category}?limit=${limit}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'X-Tenant-ID': localStorage.getItem('tenantId') || '',
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to fetch leads by category');
      }

      const data = await response.json();
      setScoresByCategory(data);
      return data;
    } catch (err: any) {
      const errorMsg = err.message || 'Error fetching leads by category';
      setError(errorMsg);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Get hot leads (score >= 75)
  const getHotLeads = useCallback(async (limit: number = 50) => {
    return getLeadsByCategory('hot', limit);
  }, [getLeadsByCategory]);

  // Get warm leads (score 50-74)
  const getWarmLeads = useCallback(async (limit: number = 100) => {
    return getLeadsByCategory('warm', limit);
  }, [getLeadsByCategory]);

  // Get cold leads (score 25-49)
  const getColdLeads = useCallback(async (limit: number = 100) => {
    return getLeadsByCategory('cold', limit);
  }, [getLeadsByCategory]);

  // Batch calculate scores for all leads
  const batchCalculateScores = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/leads/scores/batch-calculate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) {
        throw new Error('Failed to start batch calculation');
      }

      const data = await response.json();
      return data;
    } catch (err: any) {
      const errorMsg = err.message || 'Error starting batch calculation';
      setError(errorMsg);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    score,
    scoresByCategory,
    loading,
    error,
    getLeadScore,
    calculateScore,
    getLeadsByCategory,
    getHotLeads,
    getWarmLeads,
    getColdLeads,
    batchCalculateScores,
  };
}
