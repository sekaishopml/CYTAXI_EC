export type Locale = "es" | "en" | "pt";

interface Messages {
  [key: string]: string | Messages;
}

const LOCALES: Record<Locale, Messages> = {
  es: {
    common: {
      loading: "Cargando...",
      error: "Ha ocurrido un error",
      retry: "Reintentar",
      cancel: "Cancelar",
      confirm: "Confirmar",
      save: "Guardar",
      delete: "Eliminar",
      back: "Atrás",
      next: "Siguiente",
      done: "Listo",
      search: "Buscar",
      noResults: "Sin resultados",
    },
    trip: {
      pickup: "¿Dónde te recogemos?",
      destination: "¿A dónde vamos?",
      confirmTitle: "Confirma tu viaje",
      searching: "Buscando conductor",
      driverFound: "Conductor encontrado",
      arriving: "Conductor llegando",
      arrived: "El conductor ha llegado",
      inProgress: "En viaje",
      destinationNear: "Llegando a tu destino",
      payment: "Procesando pago",
      rating: "Califica tu viaje",
      completed: "Viaje completado",
      newTrip: "Nuevo viaje",
      cancelTrip: "Cancelar viaje",
    },
    driver: {
      accept: "Aceptar",
      reject: "Rechazar",
      arrive: "Llegué al pickup",
      startTrip: "Iniciar viaje",
      finishTrip: "Finalizar viaje",
      waiting: "Esperando viajes",
      online: "En línea",
      offline: "Desconectado",
      earnings: "Ganancias",
    },
    errors: {
      network: "Error de conexión",
      timeout: "La solicitud tardó demasiado",
      server: "Error del servidor",
      unauthorized: "Sesión expirada",
      notFound: "No encontrado",
    },
  },
  en: {
    common: {
      loading: "Loading...",
      error: "An error occurred",
      retry: "Retry",
      cancel: "Cancel",
      confirm: "Confirm",
      save: "Save",
      delete: "Delete",
      back: "Back",
      next: "Next",
      done: "Done",
      search: "Search",
      noResults: "No results",
    },
    trip: {
      pickup: "Where shall we pick you up?",
      destination: "Where are you going?",
      confirmTitle: "Confirm your trip",
      searching: "Searching for driver",
      driverFound: "Driver found",
      arriving: "Driver arriving",
      arrived: "Driver has arrived",
      inProgress: "On trip",
      destinationNear: "Arriving at destination",
      payment: "Processing payment",
      rating: "Rate your trip",
      completed: "Trip completed",
      newTrip: "New trip",
      cancelTrip: "Cancel trip",
    },
    driver: {
      accept: "Accept",
      reject: "Decline",
      arrive: "Arrived at pickup",
      startTrip: "Start trip",
      finishTrip: "Finish trip",
      waiting: "Waiting for trips",
      online: "Online",
      offline: "Offline",
      earnings: "Earnings",
    },
    errors: {
      network: "Network error",
      timeout: "Request timed out",
      server: "Server error",
      unauthorized: "Session expired",
      notFound: "Not found",
    },
  },
  pt: {
    common: {
      loading: "Carregando...",
      error: "Ocorreu um erro",
      retry: "Tentar novamente",
      cancel: "Cancelar",
      confirm: "Confirmar",
      save: "Salvar",
      delete: "Excluir",
      back: "Voltar",
      next: "Próximo",
      done: "Concluído",
      search: "Pesquisar",
      noResults: "Sem resultados",
    },
    trip: {
      pickup: "Onde vamos buscar você?",
      destination: "Para onde vamos?",
      confirmTitle: "Confirme sua viagem",
      searching: "Procurando motorista",
      driverFound: "Motorista encontrado",
      arriving: "Motorista chegando",
      arrived: "O motorista chegou",
      inProgress: "Em viagem",
      destinationNear: "Chegando ao destino",
      payment: "Processando pagamento",
      rating: "Avalie sua viagem",
      completed: "Viagem concluída",
      newTrip: "Nova viagem",
      cancelTrip: "Cancelar viagem",
    },
    driver: {
      accept: "Aceitar",
      reject: "Recusar",
      arrive: "Cheguei ao local",
      startTrip: "Iniciar viagem",
      finishTrip: "Finalizar viagem",
      waiting: "Aguardando viagens",
      online: "Online",
      offline: "Offline",
      earnings: "Ganhos",
    },
    errors: {
      network: "Erro de conexão",
      timeout: "A solicitação expirou",
      server: "Erro do servidor",
      unauthorized: "Sessão expirada",
      notFound: "Não encontrado",
    },
  },
};

let _currentLocale: Locale = "es";
let _messages: Messages = LOCALES.es;

export function setLocale(locale: Locale): void {
  _currentLocale = locale;
  _messages = LOCALES[locale] || LOCALES.es;
}

export function getLocale(): Locale { return _currentLocale; }

export function t(path: string, params?: Record<string, string | number>): string {
  const parts = path.split(".");
  let value: unknown = _messages;
  for (const part of parts) {
    if (value && typeof value === "object" && part in (value as Record<string, unknown>)) {
      value = (value as Record<string, unknown>)[part];
    } else {
      return path;
    }
  }
  if (typeof value !== "string") return path;
  if (params) {
    return Object.entries(params).reduce((str, [key, val]) => str.replace(`{${key}}`, String(val)), value);
  }
  return value;
}

export const locales: { value: Locale; label: string }[] = [
  { value: "es", label: "Español" },
  { value: "en", label: "English" },
  { value: "pt", label: "Português" },
];

export function useI18n(): { t: typeof t; locale: Locale; setLocale: typeof setLocale; locales: typeof locales } {
  return { t, locale: _currentLocale, setLocale, locales };
}
