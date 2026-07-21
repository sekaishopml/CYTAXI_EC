export interface DemoConfig {
  enabled: boolean;
  driver: {
    id: string; name: string; vehicle: string; plate: string;
    rating: number; photo: string; tier: "gold" | "silver" | "platinum" | "elite"; trust_score: number;
  };
  matchingDelay: number;
  searchTimeout: number;
  simulateNoDrivers: boolean;
  passenger: { phone: string; name: string };
}

export const DEMO_CONFIG: DemoConfig = {
  enabled: true,
  driver: {
    id: "drv_demo", name: "Conductor Demo", vehicle: "Vehículo de prueba",
    plate: "DEMO-001", rating: 4.5, photo: "", tier: "silver", trust_score: 80,
  },
  matchingDelay: 3000,
  searchTimeout: 30000,
  simulateNoDrivers: false,
  passenger: { phone: "0000000000", name: "Usuario Demo" },
};
