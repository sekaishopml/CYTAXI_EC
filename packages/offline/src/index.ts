export type QueueStatus = "pending" | "processing" | "completed" | "failed";

export interface QueueItem<T = unknown> {
  id: string;
  type: string;
  payload: T;
  status: QueueStatus;
  createdAt: string;
  retryCount: number;
  maxRetries: number;
  lastError?: string;
}

export interface OfflineConfig {
  dbName?: string;
  storeName?: string;
  maxRetries?: number;
  onSync?: (item: QueueItem) => Promise<boolean>;
  onComplete?: (item: QueueItem) => void;
  onError?: (item: QueueItem, error: Error) => void;
}

export class OfflineQueue {
  private db: IDBDatabase | null = null;
  private config: Required<OfflineConfig>;
  private isOnline = true;
  private processing = false;

  constructor(config: OfflineConfig = {}) {
    this.config = {
      dbName: config.dbName ?? "cytaxi_offline",
      storeName: config.storeName ?? "action_queue",
      maxRetries: config.maxRetries ?? 3,
      onSync: config.onSync ?? (async () => true),
      onComplete: config.onComplete ?? (() => {}),
      onError: config.onError ?? (() => {}),
    };
    this.initDB();
    this.watchNetwork();
  }

  private async initDB(): Promise<void> {
    return new Promise((resolve, reject) => {
      const req = indexedDB.open(this.config.dbName, 1);
      req.onupgradeneeded = () => {
        const db = req.result;
        if (!db.objectStoreNames.contains(this.config.storeName)) {
          db.createObjectStore(this.config.storeName, { keyPath: "id" });
        }
      };
      req.onsuccess = () => { this.db = req.result; resolve(); };
      req.onerror = () => reject(req.error);
    });
  }

  private watchNetwork(): void {
    if (typeof window === "undefined") return;
    this.isOnline = navigator.onLine;
    window.addEventListener("online", () => { this.isOnline = true; this.processQueue(); });
    window.addEventListener("offline", () => { this.isOnline = false; });
  }

  async enqueue<T>(type: string, payload: T, maxRetries?: number): Promise<string> {
    await this.ready();
    const item: QueueItem<T> = {
      id: `${type}_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
      type,
      payload,
      status: "pending",
      createdAt: new Date().toISOString(),
      retryCount: 0,
      maxRetries: maxRetries ?? this.config.maxRetries,
    };

    await this.put(item);

    if (this.isOnline && !this.processing) {
      this.processQueue();
    }

    return item.id;
  }

  async processQueue(): Promise<void> {
    if (this.processing || !this.isOnline) return;
    this.processing = true;

    try {
      const items = await this.getAll("pending");
      for (const item of items) {
        if (!this.isOnline) break;

        item.status = "processing";
        await this.put(item);

        try {
          const success = await this.config.onSync(item);
          if (success) {
            item.status = "completed";
            await this.put(item);
            this.config.onComplete(item);
          } else {
            throw new Error("Sync returned false");
          }
        } catch (err) {
          item.retryCount++;
          if (item.retryCount >= item.maxRetries) {
            item.status = "failed";
            item.lastError = err instanceof Error ? err.message : String(err);
            await this.put(item);
            this.config.onError(item, err instanceof Error ? err : new Error(String(err)));
          } else {
            item.status = "pending";
            await this.put(item);
          }
        }
      }
    } finally {
      this.processing = false;
    }
  }

  async getPending(): Promise<QueueItem[]> {
    return this.getAll("pending");
  }

  async getFailed(): Promise<QueueItem[]> {
    return this.getAll("failed");
  }

  async getAll(): Promise<QueueItem[]>;
  async getAll(status: QueueStatus): Promise<QueueItem[]>;
  async getAll(status?: QueueStatus): Promise<QueueItem[]> {
    await this.ready();
    return new Promise((resolve, reject) => {
      const tx = this.db!.transaction(this.config.storeName, "readonly");
      const store = tx.objectStore(this.config.storeName);
      const req = store.getAll();
      req.onsuccess = () => {
        let items = req.result as QueueItem[];
        if (status) items = items.filter((i) => i.status === status);
        resolve(items);
      };
      req.onerror = () => reject(req.error);
    });
  }

  async clear(): Promise<void> {
    await this.ready();
    const tx = this.db!.transaction(this.config.storeName, "readwrite");
    tx.objectStore(this.config.storeName).clear();
  }

  async retryFailed(): Promise<void> {
    const failed = await this.getFailed();
    for (const item of failed) {
      item.status = "pending";
      item.retryCount = 0;
      await this.put(item);
    }
    this.processQueue();
  }

  get online(): boolean { return this.isOnline; }

  private async put(item: QueueItem): Promise<void> {
    return new Promise((resolve, reject) => {
      const tx = this.db!.transaction(this.config.storeName, "readwrite");
      const req = tx.objectStore(this.config.storeName).put(item);
      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
  }

  private async ready(): Promise<void> {
    if (this.db) return;
    await new Promise<void>((resolve) => {
      const check = () => {
        if (this.db) resolve();
        else setTimeout(check, 50);
      };
      check();
    });
  }
}
