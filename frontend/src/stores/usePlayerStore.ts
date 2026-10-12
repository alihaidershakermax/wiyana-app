import { defineStore } from "pinia";
import { reactive } from "vue";

const PLAYER_STORAGE_KEY = "wiyana_player_state";

interface PlayerStorage {
  resume: Record<string, number>;
  completed: Record<string, boolean>;
}

function loadStorage(): PlayerStorage {
  try {
    const raw = localStorage.getItem(PLAYER_STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw) as Partial<PlayerStorage>;
      return {
        resume: parsed.resume || {},
        completed: parsed.completed || {},
      };
    }
  } catch {
    // Ignore corrupted storage.
  }
  return { resume: {}, completed: {} };
}

function persist(state: PlayerStorage): void {
  try {
    localStorage.setItem(PLAYER_STORAGE_KEY, JSON.stringify(state));
  } catch {
    // Storage unavailable — state stays in memory only.
  }
}

export const usePlayerStore = defineStore("player", () => {
  const initial = loadStorage();
  const state = reactive<PlayerStorage>({
    resume: { ...initial.resume },
    completed: { ...initial.completed },
  });

  function getResume(lessonId: string): number {
    return state.resume[lessonId] ?? 0;
  }

  function setResume(lessonId: string, seconds: number): void {
    if (!Number.isFinite(seconds) || seconds <= 0) return;
    state.resume[lessonId] = seconds;
    persist(state);
  }

  function isCompleted(lessonId: string): boolean {
    return state.completed[lessonId] === true;
  }

  function markCompleted(lessonId: string): void {
    state.completed[lessonId] = true;
    persist(state);
  }

  function completedCount(lessonIds: string[]): number {
    return lessonIds.filter((id) => isCompleted(id)).length;
  }

  return {
    getResume,
    setResume,
    isCompleted,
    markCompleted,
    completedCount,
  };
});
