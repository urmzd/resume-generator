import { useState, useCallback } from "react";
import {
  GetResume,
  UpdateResume,
  SaveResumeFile,
} from "../../wailsjs/go/main/App";
import type { Resume, ValidationError } from "@/types";

interface UseResumeEditorReturn {
  resume: Resume | null;
  draft: Resume | null;
  loading: boolean;
  dirty: boolean;
  errors: ValidationError[];
  saving: boolean;
  fetchResume: () => Promise<void>;
  setDraft: (r: Resume) => void;
  applyChanges: () => Promise<boolean>;
  saveToFile: () => Promise<void>;
}

export function useResumeEditor(): UseResumeEditorReturn {
  const [resume, setResume] = useState<Resume | null>(null);
  const [draft, setDraftState] = useState<Resume | null>(null);
  const [loading, setLoading] = useState(false);
  const [dirty, setDirty] = useState(false);
  const [errors, setErrors] = useState<ValidationError[]>([]);
  const [saving, setSaving] = useState(false);

  const fetchResume = useCallback(async () => {
    setLoading(true);
    try {
      const r = await GetResume();
      setResume(r as unknown as Resume);
      setDraftState(r as unknown as Resume);
      setDirty(false);
      setErrors([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const setDraft = useCallback((r: Resume) => {
    setDraftState(r);
    setDirty(true);
  }, []);

  const applyChanges = useCallback(async (): Promise<boolean> => {
    if (!draft) return false;
    setLoading(true);
    try {
      const validationErrors = await UpdateResume(draft as never);
      if (validationErrors && validationErrors.length > 0) {
        setErrors(validationErrors as unknown as ValidationError[]);
        return false;
      }
      setResume(draft);
      setDirty(false);
      setErrors([]);
      return true;
    } finally {
      setLoading(false);
    }
  }, [draft]);

  const saveToFile = useCallback(async () => {
    setSaving(true);
    try {
      await SaveResumeFile();
    } finally {
      setSaving(false);
    }
  }, []);

  return {
    resume,
    draft,
    loading,
    dirty,
    errors,
    saving,
    fetchResume,
    setDraft,
    applyChanges,
    saveToFile,
  };
}
