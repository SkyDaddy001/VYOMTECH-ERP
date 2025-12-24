-- ============================================================
-- ATTRIBUTION ANALYTICS QUERIES
-- ============================================================
-- These queries are analytics/BI-ready (Superset compatible)
-- Designed for marketing ROI, conversion funnel analysis
-- All support multi-tenant isolation via tenant_id filter

-- ============================================================
-- 1. FIRST-TOUCH ATTRIBUTION REPORT
-- ============================================================
-- Use case: Marketing attribution, campaign ROI calculation
-- Question: "Which campaign brought the most leads?"

SELECT
  ae.source,
  ae.sub_source,
  ae.campaign,
  ae.medium,
  COUNT(DISTINCT ae.lead_id) AS leads_attributed,
  COUNT(DISTINCT CASE WHEN l.status = 'won' THEN ae.lead_id END) AS conversions,
  ROUND(
    COUNT(DISTINCT CASE WHEN l.status = 'won' THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS conversion_rate,
  SUM(COALESCE(res.sale_amount, 0)) AS total_revenue,
  ROUND(AVG(COALESCE(res.sale_amount, 0)), 2) AS avg_deal_size,
  DATEDIFF(NOW(), MIN(ae.occurred_at)) AS days_since_first
FROM
  attribution_event ae
  LEFT JOIN sales_lead l ON ae.tenant_id = l.tenant_id AND ae.lead_id = l.id
  LEFT JOIN real_estate_sale res ON l.id = res.lead_id AND res.status = 'closed'
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
  AND ae.touch_order = 1  -- FIRST TOUCH ONLY
  AND ae.occurred_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
GROUP BY
  ae.source,
  ae.sub_source,
  ae.campaign,
  ae.medium
ORDER BY
  conversions DESC,
  leads_attributed DESC;

-- ============================================================
-- 2. LAST-TOUCH ATTRIBUTION REPORT
-- ============================================================
-- Use case: Sales optimization, channel effectiveness
-- Question: "Which channel closes more deals?"

SELECT
  ae.source,
  ae.campaign,
  ae.medium,
  COUNT(DISTINCT ae.lead_id) AS touches,
  COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) AS closed_deals,
  ROUND(
    COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS close_rate,
  SUM(COALESCE(res.sale_amount, 0)) AS revenue,
  ROUND(AVG(COALESCE(res.sale_amount, 0)), 2) AS avg_deal_size,
  AVG(DATEDIFF(res.sale_date, ae.occurred_at)) AS avg_days_to_close
FROM
  attribution_event ae
  INNER JOIN (
    -- Get the LAST touch for each lead
    SELECT
      tenant_id,
      lead_id,
      MAX(touch_order) AS max_order
    FROM
      attribution_event
    WHERE
      tenant_id = 'TENANT_ID'
      AND deleted_at IS NULL
    GROUP BY
      tenant_id,
      lead_id
  ) last_touch ON ae.tenant_id = last_touch.tenant_id
    AND ae.lead_id = last_touch.lead_id
    AND ae.touch_order = last_touch.max_order
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
  AND ae.occurred_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
GROUP BY
  ae.source,
  ae.campaign,
  ae.medium
ORDER BY
  closed_deals DESC,
  revenue DESC;

-- ============================================================
-- 3. ASSISTED ATTRIBUTION
-- ============================================================
-- Use case: Understanding multi-touch journeys
-- Question: "Which channels support (but don't close) deals?"

SELECT
  ae_assist.source AS assisting_channel,
  ae_assist.campaign AS assisting_campaign,
  COUNT(DISTINCT ae_assist.lead_id) AS assisted_touches,
  COUNT(DISTINCT CASE
    WHEN res.status = 'closed' THEN ae_assist.lead_id
  END) AS conversions_assisted,
  SUM(COALESCE(res.sale_amount, 0)) AS revenue_assisted,
  ROUND(
    SUM(COALESCE(res.sale_amount, 0)) / NULLIF(COUNT(DISTINCT ae_assist.lead_id), 0),
    2
  ) AS avg_deal_value
FROM
  attribution_event ae_assist
  LEFT JOIN real_estate_sale res ON ae_assist.lead_id = res.lead_id
WHERE
  ae_assist.tenant_id = 'TENANT_ID'
  AND ae_assist.deleted_at IS NULL
  AND ae_assist.touch_order > 1  -- NOT first or last
  AND ae_assist.touch_order < (
    SELECT MAX(touch_order)
    FROM attribution_event ae_last
    WHERE
      ae_last.tenant_id = ae_assist.tenant_id
      AND ae_last.lead_id = ae_assist.lead_id
  )
  AND ae_assist.occurred_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
GROUP BY
  ae_assist.source,
  ae_assist.campaign
ORDER BY
  revenue_assisted DESC;

-- ============================================================
-- 4. MARKETING FUNNEL - LEAD SOURCE TO CONVERSION
-- ============================================================
-- Use case: End-to-end funnel analysis
-- Question: "How many leads at each stage per campaign?"

SELECT
  mc.name AS campaign_name,
  COUNT(DISTINCT ae.lead_id) AS initial_leads,
  COUNT(DISTINCT CASE WHEN l.status IN ('qualified', 'negotiation') THEN ae.lead_id END) AS qualified_leads,
  COUNT(DISTINCT sv.id) AS site_visits,
  COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) AS closed_deals,
  ROUND(
    COUNT(DISTINCT CASE WHEN l.status IN ('qualified', 'negotiation') THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS qualification_rate,
  ROUND(
    COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS conversion_rate,
  SUM(COALESCE(res.sale_amount, 0)) AS total_revenue
FROM
  attribution_event ae
  LEFT JOIN marketing_campaign mc ON ae.campaign = mc.id
  LEFT JOIN sales_lead l ON ae.tenant_id = l.tenant_id AND ae.lead_id = l.id
  LEFT JOIN site_visit sv ON l.id = sv.lead_id
  LEFT JOIN real_estate_sale res ON l.id = res.lead_id AND res.status = 'closed'
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
  AND ae.touch_order = 1  -- First touch per lead
  AND ae.occurred_at >= DATE_SUB(NOW(), INTERVAL 180 DAY)
GROUP BY
  mc.name
ORDER BY
  initial_leads DESC;

-- ============================================================
-- 5. TIME-TO-CONVERSION ANALYSIS
-- ============================================================
-- Use case: Sales cycle optimization
-- Question: "How long does it take from first touch to conversion?"

SELECT
  las.first_source,
  las.last_source,
  COUNT(DISTINCT las.lead_id) AS conversions,
  AVG(las.days_to_conversion) AS avg_days_to_conversion,
  PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY las.days_to_conversion) AS median_days,
  PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY las.days_to_conversion) AS p90_days,
  MIN(las.days_to_conversion) AS min_days,
  MAX(las.days_to_conversion) AS max_days,
  COUNT(CASE WHEN las.days_to_conversion <= 7 THEN 1 END) AS conversions_under_7_days
FROM
  lead_attribution_snapshot las
WHERE
  las.tenant_id = 'TENANT_ID'
  AND las.last_touch_at IS NOT NULL
GROUP BY
  las.first_source,
  las.last_source
ORDER BY
  conversions DESC;

-- ============================================================
-- 6. CHANNEL ATTRIBUTION COMPARISON
-- ============================================================
-- Use case: Multi-model attribution comparison
-- Question: "What if we credit last-touch vs first-touch?"

WITH first_touch_attribution AS (
  SELECT
    ae.source,
    ae.campaign,
    COUNT(DISTINCT ae.lead_id) AS leads,
    ROUND(SUM(COALESCE(res.sale_amount, 0)) / COUNT(DISTINCT ae.lead_id), 2) AS avg_value
  FROM
    attribution_event ae
    LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
  WHERE
    ae.tenant_id = 'TENANT_ID'
    AND ae.deleted_at IS NULL
    AND ae.touch_order = 1
  GROUP BY
    ae.source,
    ae.campaign
),
last_touch_attribution AS (
  SELECT
    ae.source,
    ae.campaign,
    COUNT(DISTINCT ae.lead_id) AS leads,
    ROUND(SUM(COALESCE(res.sale_amount, 0)) / COUNT(DISTINCT ae.lead_id), 2) AS avg_value
  FROM
    attribution_event ae
    INNER JOIN (
      SELECT
        tenant_id,
        lead_id,
        MAX(touch_order) AS max_order
      FROM attribution_event
      WHERE tenant_id = 'TENANT_ID'
      GROUP BY tenant_id, lead_id
    ) lt ON ae.tenant_id = lt.tenant_id
      AND ae.lead_id = lt.lead_id
      AND ae.touch_order = lt.max_order
    LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
  WHERE
    ae.tenant_id = 'TENANT_ID'
    AND ae.deleted_at IS NULL
  GROUP BY
    ae.source,
    ae.campaign
)
SELECT
  COALESCE(ft.source, lt.source) AS source,
  COALESCE(ft.campaign, lt.campaign) AS campaign,
  COALESCE(ft.leads, 0) AS first_touch_leads,
  COALESCE(ft.avg_value, 0) AS first_touch_value,
  COALESCE(lt.leads, 0) AS last_touch_leads,
  COALESCE(lt.avg_value, 0) AS last_touch_value,
  ABS(COALESCE(ft.leads, 0) - COALESCE(lt.leads, 0)) AS lead_difference,
  ROUND(
    (COALESCE(ft.avg_value, 0) - COALESCE(lt.avg_value, 0)) * 100.0 / NULLIF(COALESCE(lt.avg_value, 1), 0),
    2
  ) AS value_difference_pct
FROM
  first_touch_attribution ft
  FULL OUTER JOIN last_touch_attribution lt ON ft.source = lt.source
    AND ft.campaign = lt.campaign
ORDER BY
  first_touch_leads DESC;

-- ============================================================
-- 7. DEVICE & GEO ATTRIBUTION
-- ============================================================
-- Use case: Channel optimization by device/location
-- Question: "Which devices convert best per channel?"

SELECT
  ae.source,
  ae.device,
  ae.country,
  COUNT(DISTINCT ae.lead_id) AS leads,
  COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) AS conversions,
  ROUND(
    COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS conversion_rate,
  SUM(COALESCE(res.sale_amount, 0)) AS revenue,
  ROUND(AVG(COALESCE(res.sale_amount, 0)), 2) AS avg_deal_size
FROM
  attribution_event ae
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
  AND ae.touch_order = 1
GROUP BY
  ae.source,
  ae.device,
  ae.country
ORDER BY
  revenue DESC;

-- ============================================================
-- 8. CUSTOM PAYLOAD ANALYSIS (JSON)
-- ============================================================
-- Use case: WhatsApp, Telephony, QR code attribution
-- Question: "What's the conversion rate for WhatsApp leads?"

SELECT
  ae.source,
  JSON_EXTRACT(ae.custom_payload, '$.whatsapp_msg_id') AS whatsapp_id,
  JSON_EXTRACT(ae.custom_payload, '$.call_duration') AS call_duration,
  COUNT(DISTINCT ae.lead_id) AS leads,
  COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) AS conversions,
  ROUND(
    COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS conversion_rate,
  SUM(COALESCE(res.sale_amount, 0)) AS revenue
FROM
  attribution_event ae
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.source IN ('whatsapp', 'phone', 'qr')
  AND ae.custom_payload IS NOT NULL
GROUP BY
  ae.source,
  JSON_EXTRACT(ae.custom_payload, '$.whatsapp_msg_id')
ORDER BY
  conversions DESC;

-- ============================================================
-- 9. REVENUE ATTRIBUTION (MULTI-MODEL COMPARISON)
-- ============================================================
-- Use case: Finance/accounting reconciliation
-- Question: "How should we credit revenue across channels?"

SELECT
  'first-touch' AS attribution_model,
  ae.source,
  ae.campaign,
  COUNT(DISTINCT ae.lead_id) AS attributed_conversions,
  SUM(COALESCE(res.sale_amount, 0)) AS attributed_revenue,
  ROUND(SUM(COALESCE(res.sale_amount, 0)) / COUNT(DISTINCT ae.lead_id), 2) AS cpa,
  ROUND(SUM(COALESCE(mc.budget, 0)), 2) AS spend,
  ROUND(SUM(COALESCE(res.sale_amount, 0)) / NULLIF(SUM(COALESCE(mc.budget, 0)), 0), 2) AS roi
FROM
  attribution_event ae
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
  LEFT JOIN marketing_campaign mc ON ae.campaign = mc.id
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
  AND ae.touch_order = 1
GROUP BY
  ae.source,
  ae.campaign

UNION ALL

SELECT
  'last-touch' AS attribution_model,
  ae.source,
  ae.campaign,
  COUNT(DISTINCT ae.lead_id) AS attributed_conversions,
  SUM(COALESCE(res.sale_amount, 0)) AS attributed_revenue,
  ROUND(SUM(COALESCE(res.sale_amount, 0)) / COUNT(DISTINCT ae.lead_id), 2) AS cpa,
  ROUND(SUM(COALESCE(mc.budget, 0)), 2) AS spend,
  ROUND(SUM(COALESCE(res.sale_amount, 0)) / NULLIF(SUM(COALESCE(mc.budget, 0)), 0), 2) AS roi
FROM
  attribution_event ae
  INNER JOIN (
    SELECT tenant_id, lead_id, MAX(touch_order) AS max_order
    FROM attribution_event
    WHERE tenant_id = 'TENANT_ID'
    GROUP BY tenant_id, lead_id
  ) lt ON ae.tenant_id = lt.tenant_id
    AND ae.lead_id = lt.lead_id
    AND ae.touch_order = lt.max_order
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
  LEFT JOIN marketing_campaign mc ON ae.campaign = mc.id
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
GROUP BY
  ae.source,
  ae.campaign

ORDER BY
  attribution_model,
  attributed_revenue DESC;

-- ============================================================
-- 10. LEAD QUALITY BY SOURCE
-- ============================================================
-- Use case: Lead quality assessment
-- Question: "Which sources produce the highest-quality leads?"

SELECT
  las.first_source,
  COUNT(DISTINCT las.lead_id) AS total_leads,
  COUNT(DISTINCT CASE WHEN l.status IN ('qualified', 'negotiation') THEN las.lead_id END) AS qualified_leads,
  COUNT(DISTINCT sv.id) AS site_visits,
  COUNT(DISTINCT res.id) AS conversions,
  ROUND(
    COUNT(DISTINCT CASE WHEN res.id IS NOT NULL THEN las.lead_id END) * 100.0 /
    COUNT(DISTINCT las.lead_id),
    2
  ) AS conversion_rate,
  AVG(las.days_to_conversion) AS avg_sales_cycle,
  AVG(COALESCE(res.sale_amount, 0)) AS avg_deal_value,
  COUNT(DISTINCT res.id) / COUNT(DISTINCT las.lead_id) AS deals_per_lead_ratio
FROM
  lead_attribution_snapshot las
  LEFT JOIN sales_lead l ON las.lead_id = l.id
  LEFT JOIN site_visit sv ON las.lead_id = sv.lead_id
  LEFT JOIN real_estate_sale res ON las.lead_id = res.lead_id AND res.status = 'closed'
WHERE
  las.tenant_id = 'TENANT_ID'
GROUP BY
  las.first_source
ORDER BY
  conversion_rate DESC,
  conversions DESC;

-- ============================================================
-- 11. TOUCHPOINT SEQUENCE ANALYSIS
-- ============================================================
-- Use case: Customer journey optimization
-- Question: "What's the most common path to conversion?"

SELECT
  las.touch_sequence,
  COUNT(DISTINCT las.lead_id) AS frequency,
  SUM(CASE WHEN res.status = 'closed' THEN 1 ELSE 0 END) AS conversions,
  ROUND(SUM(CASE WHEN res.status = 'closed' THEN 1 ELSE 0 END) * 100.0 /
    COUNT(DISTINCT las.lead_id), 2) AS conversion_rate,
  AVG(las.days_to_conversion) AS avg_cycle_days,
  SUM(COALESCE(res.sale_amount, 0)) AS total_revenue
FROM
  lead_attribution_snapshot las
  LEFT JOIN real_estate_sale res ON las.lead_id = res.lead_id AND res.status = 'closed'
WHERE
  las.tenant_id = 'TENANT_ID'
  AND las.touch_sequence IS NOT NULL
GROUP BY
  las.touch_sequence
ORDER BY
  conversions DESC,
  frequency DESC
LIMIT 20;

-- ============================================================
-- 12. DAILY ATTRIBUTION TRENDS
-- ============================================================
-- Use case: Time-series analysis, trend detection
-- Question: "Is attribution changing week-over-week?"

SELECT
  DATE(ae.occurred_at) AS date,
  ae.source,
  COUNT(DISTINCT ae.lead_id) AS leads,
  COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) AS conversions,
  SUM(COALESCE(res.sale_amount, 0)) AS revenue,
  ROUND(
    COUNT(DISTINCT CASE WHEN res.status = 'closed' THEN ae.lead_id END) * 100.0 /
    COUNT(DISTINCT ae.lead_id),
    2
  ) AS conversion_rate
FROM
  attribution_event ae
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id AND res.status = 'closed'
WHERE
  ae.tenant_id = 'TENANT_ID'
  AND ae.deleted_at IS NULL
  AND ae.touch_order = 1
  AND ae.occurred_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
GROUP BY
  DATE(ae.occurred_at),
  ae.source
ORDER BY
  date DESC,
  revenue DESC;

-- ============================================================
-- MATERIALIZED VIEWS (for performance optimization)
-- ============================================================

-- Create a view for first-touch attribution (high-traffic query)
CREATE OR REPLACE VIEW v_attribution_first_touch AS
SELECT
  ae.id,
  ae.tenant_id,
  ae.lead_id,
  ae.source,
  ae.sub_source,
  ae.campaign,
  ae.medium,
  ae.device,
  ae.country,
  ae.utm_campaign,
  ae.utm_source,
  ae.utm_medium,
  ae.occurred_at,
  res.status,
  res.sale_amount
FROM
  attribution_event ae
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id
WHERE
  ae.touch_order = 1
  AND ae.deleted_at IS NULL;

-- Create a view for last-touch attribution
CREATE OR REPLACE VIEW v_attribution_last_touch AS
SELECT
  ae.id,
  ae.tenant_id,
  ae.lead_id,
  ae.source,
  ae.sub_source,
  ae.campaign,
  ae.medium,
  ae.device,
  ae.country,
  ae.occurred_at,
  res.status,
  res.sale_amount,
  ROW_NUMBER() OVER (PARTITION BY ae.tenant_id, ae.lead_id ORDER BY ae.touch_order DESC) AS rn
FROM
  attribution_event ae
  LEFT JOIN real_estate_sale res ON ae.lead_id = res.lead_id
WHERE
  ae.deleted_at IS NULL
HAVING
  rn = 1;
