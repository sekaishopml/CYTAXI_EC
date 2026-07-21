export interface AIConfig {
  provider: "rule_based" | "ml_service";
  mlEndpoint?: string;
  mlApiKey?: string;
}

export interface Location {
  lat: number;
  lng: number;
}

export interface DriverCandidate {
  id: string;
  name: string;
  lat: number;
  lng: number;
  rating: number;
  acceptRate: number;
  vehicle: string;
  trustScore: number;
  tier: "bronze" | "silver" | "gold" | "platinum" | "elite";
  distanceKm: number;
  etaSeconds: number;
  zoneBonus: number;
}

export interface MatchingRequest {
  pickup: Location;
  destination: Location;
  passengerId: string;
  preferences?: {
    vehicleType?: string;
    minRating?: number;
    maxEtaSeconds?: number;
  };
}

export interface MatchingResult {
  selected: DriverCandidate[];
  alternatives: DriverCandidate[];
  scores: Record<string, number>;
  strategy: string;
  etaSeconds: number;
}

export interface ETAPrediction {
  etaSeconds: number;
  confidence: number;
  trafficFactor: number;
  distanceKm: number;
  routePolyline?: string;
}

export interface DynamicPrice {
  baseFare: number;
  distanceFare: number;
  timeFare: number;
  surgeMultiplier: number;
  demandLevel: number;
  supplyLevel: number;
  nightSurcharge: number;
  airportSurcharge: number;
  total: number;
  currency: string;
}

export interface FraudAssessment {
  riskLevel: "low" | "medium" | "high" | "critical";
  riskScore: number;
  flags: FraudFlag[];
}

export interface FraudFlag {
  type: string;
  severity: number;
  description: string;
}

export interface AIService {
  matchDrivers(request: MatchingRequest, candidates: DriverCandidate[]): MatchingResult;
  predictETA(origin: Location, dest: Location, timestamp?: number): ETAPrediction;
  calculateDynamicPrice(
    baseParams: { distanceKm: number; durationSec: number },
    demandLevel: number,
    supplyLevel: number,
    hour: number,
    zone: string,
  ): DynamicPrice;
  assessFraud(tripData: {
    passengerId: string;
    driverId?: string;
    distanceKm: number;
    fare: number;
    paymentMethod: string;
    hour: number;
    zone: string;
  }): FraudAssessment;
}

export class RuleBasedAIService implements AIService {
  private config: AIConfig;

  constructor(config?: Partial<AIConfig>) {
    this.config = { provider: "rule_based", ...config };
  }

  matchDrivers(request: MatchingRequest, candidates: DriverCandidate[]): MatchingResult {
    const scores: Record<string, number> = {};
    const weights = { distance: 0.35, eta: 0.25, rating: 0.15, acceptRate: 0.10, trustScore: 0.10, zone: 0.05 };
    const maxDistance = Math.max(...candidates.map(c => c.distanceKm), 1);
    const maxEta = Math.max(...candidates.map(c => c.etaSeconds), 1);

    const scored = candidates.map(c => {
      const distanceScore = 1 - (c.distanceKm / maxDistance);
      const etaScore = 1 - Math.min(c.etaSeconds / 600, 1);
      const ratingScore = c.rating / 5;
      const acceptRateScore = c.acceptRate / 100;
      const trustScore = c.trustScore / 100;
      const zoneScore = c.zoneBonus;

      const score =
        weights.distance * distanceScore +
        weights.eta * etaScore +
        weights.rating * ratingScore +
        weights.acceptRate * acceptRateScore +
        weights.trustScore * trustScore +
        weights.zone * zoneScore;

      scores[c.id] = score;
      return { ...c, score };
    });

    scored.sort((a, b) => b.score - a.score);

    const selected = scored.slice(0, 1);
    const alternatives = scored.slice(1, 4);

    return {
      selected: selected.map(({ score: _, ...rest }) => rest),
      alternatives: alternatives.map(({ score: _, ...rest }) => rest),
      scores,
      strategy: "weighted_multi_factor",
      etaSeconds: selected[0]?.etaSeconds ?? 300,
    };
  }

  predictETA(origin: Location, dest: Location, _timestamp?: number): ETAPrediction {
    const R = 6371;
    const dLat = ((dest.lat - origin.lat) * Math.PI) / 180;
    const dLng = ((dest.lng - origin.lng) * Math.PI) / 180;
    const a =
      Math.sin(dLat / 2) ** 2 +
      Math.cos((origin.lat * Math.PI) / 180) *
        Math.cos((dest.lat * Math.PI) / 180) *
        Math.sin(dLng / 2) ** 2;
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    const distanceKm = R * c;

    const avgSpeedKmph = 25;
    const trafficFactor = 1.0 + Math.random() * 0.3;
    const etaSeconds = Math.round((distanceKm / avgSpeedKmph) * 3600 * trafficFactor);

    return {
      etaSeconds,
      confidence: 0.75,
      trafficFactor,
      distanceKm: Math.round(distanceKm * 100) / 100,
    };
  }

  calculateDynamicPrice(
    baseParams: { distanceKm: number; durationSec: number },
    demandLevel: number,
    supplyLevel: number,
    hour: number,
    _zone: string,
  ): DynamicPrice {
    const baseFare = 100;
    const distanceFare = Math.round(baseParams.distanceKm * 50);
    const timeFare = Math.round((baseParams.durationSec / 60) * 3);

    const demandFactor = Math.max(0.8, Math.min(3.0, 1.0 + (demandLevel - 5) * 0.15));
    const supplyFactor = Math.max(0.8, Math.min(1.5, 1.0 + (10 - supplyLevel) * 0.05));
    const surgeMultiplier = Math.round(demandFactor * supplyFactor * 100) / 100;

    const isNight = hour >= 22 || hour <= 5;
    const nightSurcharge = isNight ? 250 : 0;
    const airportSurcharge = _zone === "airport" ? 300 : 0;

    const subtotal = Math.round((baseFare + distanceFare + timeFare) * surgeMultiplier);
    const total = subtotal + nightSurcharge + airportSurcharge;

    return {
      baseFare,
      distanceFare,
      timeFare,
      surgeMultiplier,
      demandLevel,
      supplyLevel,
      nightSurcharge,
      airportSurcharge,
      total,
      currency: "USD",
    };
  }

  assessFraud(tripData: {
    passengerId: string;
    driverId?: string;
    distanceKm: number;
    fare: number;
    paymentMethod: string;
    hour: number;
    zone: string;
  }): FraudAssessment {
    const flags: FraudFlag[] = [];
    let riskScore = 0;

    if (tripData.distanceKm < 0.3) {
      riskScore += 15;
      flags.push({ type: "suspicious_distance", severity: 15, description: "Distancia muy corta" });
    }

    if (tripData.distanceKm > 100) {
      riskScore += 10;
      flags.push({ type: "long_distance", severity: 10, description: "Distancia muy larga" });
    }

    if (tripData.fare > 10000) {
      riskScore += 20;
      flags.push({ type: "high_fare", severity: 20, description: "Tarifa anormalmente alta" });
    }

    if (tripData.hour >= 1 && tripData.hour <= 4) {
      riskScore += 5;
      flags.push({ type: "late_hour", severity: 5, description: "Viaje en madrugada" });
    }

    if (tripData.paymentMethod === "cash") {
      riskScore += 10;
      flags.push({ type: "cash_payment", severity: 10, description: "Pago en efectivo" });
    }

    let riskLevel: FraudAssessment["riskLevel"] = "low";
    if (riskScore >= 70) riskLevel = "critical";
    else if (riskScore >= 50) riskLevel = "high";
    else if (riskScore >= 30) riskLevel = "medium";

    return { riskLevel, riskScore, flags };
  }
}
