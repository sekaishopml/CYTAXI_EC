import React, { useState, useEffect, useCallback } from "react";
import { useTrip } from "@/contexts/trip";
import { useAuth } from "@/contexts/auth";
import { acceptTrip, startTrip, updateLocation, finishTrip } from "@/services/tracking";

type DriverTripPhase = "accepted" | "arrived" | "driving" | "completed";

interface TripInfo {
  id: string;
  passenger: string;
  passengerRating: number;
  pickup: string;
  destination: string;
  fare: string;
  distance: string;
  etaPickup: number;
  etaDest: number;
}

export default function CurrentTripPage() {
  const { current, startTrip: ctxStart, completeTrip } = useTrip();
  const { driver } = useAuth();
  const driverId = driver?.id || "drv_1000";
  const [phase, setPhase] = useState<DriverTripPhase>("accepted");
  const [eta, setEta] = useState(current?.eta || 300);
  const [distanceKm, setDistanceKm] = useState(5.5);
  const [events, setEvents] = useState<string[]>([]);

  const mockTrip: TripInfo = {
    id: current?.id || "trip_001",
    passenger: "María G.",
    passengerRating: 4.9,
    pickup: current?.pickup || "Av. 9 de Octubre 1200",
    destination: current?.destination || "Mall del Sol, Km 4.5",
    fare: current?.fare || "$8.50",
    distance: "4.2 km",
    etaPickup: 5,
    etaDest: 12,
  };

  const addEvent = useCallback((msg: string) => {
    setEvents((prev) => [...prev, `${new Date().toLocaleTimeString()}: ${msg}`]);
  }, []);

  useEffect(() => {
    if (phase !== "driving") return;
    let watchId: number | undefined;
    const sendLocation = () => {
      if (typeof navigator !== "undefined" && navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
          async (pos) => {
            try {
              await updateLocation(mockTrip.id, driverId, pos.coords.latitude, pos.coords.longitude);
              addEvent(`📍 ${pos.coords.latitude.toFixed(4)}, ${pos.coords.longitude.toFixed(4)}`);
            } catch {}
          },
          () => {},
          { enableHighAccuracy: true }
        );
      }
    };
    sendLocation();
    const interval = setInterval(() => {
      sendLocation();
      setDistanceKm((prev) => {
        const next = Math.max(0, prev - 0.15);
        if (next <= 0.1) {
          handleComplete();
          return 0;
        }
        return next;
      });
      setEta((prev) => Math.max(0, prev - 3));
    }, 3000);
    return () => {
      clearInterval(interval);
      if (watchId !== undefined && typeof navigator !== "undefined" && navigator.geolocation) {
        navigator.geolocation.clearWatch(watchId);
      }
    };
  }, [phase, mockTrip.id, driverId]);

  const handleAccept = useCallback(async () => {
    setPhase("accepted");
    try {
      await acceptTrip(mockTrip.id, driverId);
    } catch {}
    addEvent("Trip accepted");
  }, [mockTrip.id, driverId]);

  const handleArrive = useCallback(async () => {
    setPhase("arrived");
    try {
      await startTrip(mockTrip.id, driverId);
    } catch {}
    addEvent("Arrived at pickup");
  }, [mockTrip.id, driverId]);

  const handleStartTrip = useCallback(async () => {
    setPhase("driving");
    ctxStart();
    addEvent("Trip started");
  }, [ctxStart]);

  const handleComplete = useCallback(async () => {
    setPhase("completed");
    completeTrip();
    try {
      await finishTrip(mockTrip.id);
    } catch {}
    addEvent("Trip completed");
  }, [completeTrip, mockTrip.id]);

  const handleReject = useCallback(() => {
    addEvent("Trip declined");
    if (current) completeTrip();
  }, [current, completeTrip]);

  if (!current && phase !== "completed") {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-6">
        <div className="w-16 h-16 rounded-full bg-brand-green/10 flex items-center justify-center mb-4">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#006c49" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>
          </svg>
        </div>
        <h2 className="text-xl font-bold text-gray-900 mb-2">Esperando viajes</h2>
        <p className="text-sm text-gray-500">Mantente en línea para recibir solicitudes</p>
      </div>
    );
  }

  const phaseColors: Record<DriverTripPhase, { bg: string; text: string; dot: string }> = {
    accepted: { bg: "bg-brand-green/5", text: "text-brand-green", dot: "bg-brand-green" },
    arrived: { bg: "bg-blue-50", text: "text-blue-600", dot: "bg-blue-500" },
    driving: { bg: "bg-brand-green/5", text: "text-brand-green", dot: "bg-brand-green" },
    completed: { bg: "bg-gray-50", text: "text-gray-600", dot: "bg-gray-400" },
  };

  const pc = phaseColors[phase];
  const showActions = phase === "accepted" || phase === "arrived";

  return (
    <div className="space-y-4 pb-8">
      {/* Header */}
      <div className={`${pc.bg} rounded-2xl p-5 border border-gray-100`}>
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <div className={`w-3 h-3 rounded-full ${pc.dot} animate-pulse`} />
            <span className={`text-sm font-semibold ${pc.text}`}>
              {phase === "accepted" && "En camino al pasajero"}
              {phase === "arrived" && "Esperando pasajero"}
              {phase === "driving" && "Viaje en curso"}
              {phase === "completed" && "Viaje finalizado"}
            </span>
          </div>
          <span className="text-xs text-gray-400 font-mono">{mockTrip.id}</span>
        </div>

        <div className="space-y-3">
          <div className="flex items-start gap-3">
            <div className="flex flex-col items-center gap-1 mt-1">
              <div className="w-3 h-3 rounded-full bg-brand-green" />
              <div className="w-0.5 h-8 bg-gray-200" />
              <div className="w-3 h-3 rounded-full bg-blue-500" />
            </div>
            <div className="flex-1 min-w-0">
              <div className="mb-4">
                <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Recoger en</p>
                <p className="text-sm font-semibold text-gray-900 truncate">{mockTrip.pickup}</p>
              </div>
              <div>
                <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Destino</p>
                <p className="text-sm font-semibold text-gray-900 truncate">{mockTrip.destination}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-3">
        <div className="bg-white rounded-xl p-3 text-center border border-gray-100">
          <p className="text-xs text-gray-400 mb-1">Distancia</p>
          <p className="text-lg font-bold text-gray-900">{phase === "driving" ? distanceKm.toFixed(1) : mockTrip.distance}</p>
          <p className="text-xs text-gray-400">km</p>
        </div>
        <div className="bg-white rounded-xl p-3 text-center border border-gray-100">
          <p className="text-xs text-gray-400 mb-1">ETA</p>
          <p className="text-lg font-bold text-gray-900">{Math.ceil((phase === "driving" ? eta : (phase === "accepted" ? mockTrip.etaPickup : mockTrip.etaDest)) / 60) || 1}</p>
          <p className="text-xs text-gray-400">min</p>
        </div>
        <div className="bg-white rounded-xl p-3 text-center border border-gray-100">
          <p className="text-xs text-gray-400 mb-1">Tarifa</p>
          <p className="text-lg font-bold text-brand-green">{mockTrip.fare}</p>
          <p className="text-xs text-gray-400">estimado</p>
        </div>
      </div>

      {/* Passenger info */}
      <div className="bg-white rounded-xl p-4 border border-gray-100 flex items-center gap-3">
        <div className="w-12 h-12 rounded-full bg-brand-green/10 flex items-center justify-center text-xl flex-shrink-0">
          👤
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-semibold text-gray-900">{mockTrip.passenger}</p>
          <div className="flex items-center gap-2 mt-0.5">
            <span className="text-yellow-500 text-xs">★★★★</span>
            <span className="text-xs font-medium text-gray-500">{mockTrip.passengerRating}</span>
          </div>
        </div>
        {phase === "accepted" && (
          <div className="flex gap-1">
            <button className="p-2 rounded-lg bg-gray-100 hover:bg-gray-200">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-gray-600"><path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72 12.84 12.84 0 00.7 2.81 2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45 12.84 12.84 0 002.81.7A2 2 0 0122 16.92z"/></svg>
            </button>
            <button className="p-2 rounded-lg bg-gray-100 hover:bg-gray-200">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-gray-600"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/></svg>
            </button>
          </div>
        )}
      </div>

      {/* Actions */}
      {showActions && (
        <div className="flex gap-4">
          {phase === "accepted" && (
            <>
              <button onClick={handleArrive}
                className="flex-1 h-16 rounded-2xl bg-brand-green text-white font-bold text-lg shadow-lg shadow-brand-green/20 active:scale-[0.97] transition-transform flex items-center justify-center gap-2">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                Llegué al pickup
              </button>
              <button onClick={handleReject}
                className="h-16 px-8 rounded-2xl bg-red-500 text-white font-bold text-lg shadow-lg shadow-red-500/20 active:scale-[0.97] transition-transform">
                ✕
              </button>
            </>
          )}
          {phase === "arrived" && (
            <button onClick={handleStartTrip}
              className="flex-1 h-16 rounded-2xl bg-blue-500 text-white font-bold text-lg shadow-lg shadow-blue-500/20 active:scale-[0.97] transition-transform">
              Iniciar viaje
            </button>
          )}
        </div>
      )}

      {/* Active trip controls */}
      {phase === "driving" && (
        <button onClick={handleComplete}
          className="w-full h-16 rounded-2xl bg-red-500 text-white font-bold text-lg shadow-lg shadow-red-500/20 active:scale-[0.97] transition-transform">
          Finalizar viaje
        </button>
      )}

      {/* Completed */}
      {phase === "completed" && (
        <div className="bg-white rounded-2xl p-6 text-center border border-gray-100">
          <div className="text-4xl mb-3">✅</div>
          <h3 className="text-lg font-bold text-gray-900 mb-1">Viaje completado</h3>
          <p className="text-sm text-gray-500">Tarifa: {mockTrip.fare}</p>
          <button onClick={() => setPhase("accepted")}
            className="mt-4 px-8 h-12 rounded-xl bg-brand-green text-white font-semibold shadow-md">
            Nuevo viaje simulado
          </button>
        </div>
      )}

      {/* Event log (debug) */}
      <details className="bg-white rounded-xl border border-gray-100">
        <summary className="px-4 py-3 text-sm font-semibold text-gray-500 cursor-pointer">
          Event Log ({events.length})
        </summary>
        <div className="px-4 pb-3 space-y-1 max-h-32 overflow-y-auto">
          {events.map((e, i) => (
            <p key={i} className="text-xs font-mono text-gray-400">{e}</p>
          ))}
        </div>
      </details>
    </div>
  );
}
