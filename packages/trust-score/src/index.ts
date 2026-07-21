export type UserRole = "passenger" | "driver";

export interface TrustScoreConfig {
  weights: {
    punctuality: number;
    cancellations: number;
    fraud: number;
    complaints: number;
    payments: number;
  };
  tiers: TrustTier[];
}

export interface TrustTier {
  name: string;
  minScore: number;
  maxScore: number;
  color: string;
  icon: string;
  benefits: string[];
}

export interface RatingCategory {
  punctuality: number;
  service: number;
  comfort: number;
}

export interface TrustScoreInput {
  userId: string;
  role: UserRole;
  totalTrips: number;
  avgRating: number;
  ratingCategories: RatingCategory;
  cancellations: number;
  complaints: number;
  fraudFlags: number;
  onTimeRate: number;
  paymentReliability: number;
  accountAgeDays: number;
}

export interface TrustScoreResult {
  score: number;
  tier: TrustTier;
  breakdown: {
    rating: number;
    reliability: number;
    behavior: number;
    tenure: number;
  };
}

export const DEFAULT_CONFIG: TrustScoreConfig = {
  weights: {
    punctuality: 0.25,
    cancellations: 0.2,
    fraud: 0.25,
    complaints: 0.15,
    payments: 0.15,
  },
  tiers: [
    { name: "Elite", minScore: 90, maxScore: 100, color: "#7c4dff", icon: "👑", benefits: ["Priority matching", "Lower commission", "24/7 support"] },
    { name: "Platino", minScore: 80, maxScore: 89, color: "#00e676", icon: "💎", benefits: ["Priority matching", "Lower commission"] },
    { name: "Oro", minScore: 65, maxScore: 79, color: "#ffc107", icon: "🥇", benefits: ["Priority matching"] },
    { name: "Plata", minScore: 50, maxScore: 64, color: "#9ea5a0", icon: "🥈", benefits: ["Standard matching"] },
    { name: "Bronce", minScore: 0, maxScore: 49, color: "#ba1a1a", icon: "⚠️", benefits: ["Limited matching"] },
  ],
};

export function calculateTrustScore(input: TrustScoreInput, config: TrustScoreConfig = DEFAULT_CONFIG): TrustScoreResult {
  const { weights } = config;

  const ratingScore = normalizeRating(input.avgRating);
  const punctualityScore = input.onTimeRate * 100;
  const cancellationPenalty = Math.max(0, 100 - (input.cancellations / Math.max(1, input.totalTrips)) * 500);
  const fraudPenalty = Math.max(0, 100 - input.fraudFlags * 33.3);
  const complaintPenalty = Math.max(0, 100 - input.complaints * 20);
  const paymentScore = input.paymentReliability * 100;

  const reliability = (
    weights.punctuality * punctualityScore +
    weights.cancellations * cancellationPenalty +
    weights.payments * paymentScore
  ) / (weights.punctuality + weights.cancellations + weights.payments);

  const behavior = (
    weights.fraud * fraudPenalty +
    weights.complaints * complaintPenalty
  ) / (weights.fraud + weights.complaints);

  const tenure = Math.min(100, (input.accountAgeDays / 365) * 20 + input.totalTrips * 0.5);

  const totalRatingWeight = 0.5;
  const totalReliabilityWeight = 0.2;
  const totalBehaviorWeight = 0.2;
  const totalTenureWeight = 0.1;

  const score = Math.round(
    ratingScore * totalRatingWeight +
    reliability * totalReliabilityWeight +
    behavior * totalBehaviorWeight +
    tenure * totalTenureWeight
  );

  const tier = getTier(score, config);

  return {
    score,
    tier,
    breakdown: {
      rating: Math.round(ratingScore),
      reliability: Math.round(reliability),
      behavior: Math.round(behavior),
      tenure: Math.round(tenure),
    },
  };
}

export function getTier(score: number, config: TrustScoreConfig = DEFAULT_CONFIG): TrustTier {
  return config.tiers.find((t) => score >= t.minScore && score <= t.maxScore) ?? config.tiers[config.tiers.length - 1];
}

function normalizeRating(rating: number): number {
  return Math.min(100, (rating / 5) * 100);
}

export function calculateDriverScore(
  trips: number,
  rating: number,
  categories: RatingCategory,
  cancellations: number,
  onTimeRate: number
): number {
  const result = calculateTrustScore({
    userId: "driver",
    role: "driver",
    totalTrips: trips,
    avgRating: rating,
    ratingCategories: categories,
    cancellations,
    complaints: 0,
    fraudFlags: 0,
    onTimeRate,
    paymentReliability: 1,
    accountAgeDays: 180,
  });
  return result.score;
}

export function getTierBenefits(tierName: string, config: TrustScoreConfig = DEFAULT_CONFIG): string[] {
  const tier = config.tiers.find((t) => t.name === tierName);
  return tier?.benefits ?? [];
}
