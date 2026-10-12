const NOTES_STORAGE_KEY = "wiyana_lesson_notes";

export interface LessonNote {
  text: string;
  updatedAt: string;
}

function readAll(): Record<string, LessonNote> {
  try {
    const raw = localStorage.getItem(NOTES_STORAGE_KEY);
    return raw ? JSON.parse(raw) : {};
  } catch {
    return {};
  }
}

function writeAll(notes: Record<string, LessonNote>): void {
  try {
    localStorage.setItem(NOTES_STORAGE_KEY, JSON.stringify(notes));
  } catch {
    // Storage unavailable (private mode / quota) — notes stay in memory.
  }
}

export class NotesService {
  private static instance: NotesService;

  private constructor() {}

  public static getInstance(): NotesService {
    if (!NotesService.instance) {
      NotesService.instance = new NotesService();
    }
    return NotesService.instance;
  }

  public getNote(lessonId: string): LessonNote {
    return readAll()[lessonId] || { text: "", updatedAt: "" };
  }

  public saveNote(lessonId: string, text: string): LessonNote {
    const notes = readAll();
    const note: LessonNote = { text, updatedAt: new Date().toISOString() };
    notes[lessonId] = note;
    writeAll(notes);
    return note;
  }

  public clearNote(lessonId: string): void {
    const notes = readAll();
    delete notes[lessonId];
    writeAll(notes);
  }
}

export const notesService = NotesService.getInstance();
