export class PdfCache {
  private cache = new Map<string, string>();

  get(name: string): string | null {
    return this.cache.get(name) ?? null;
  }

  set(name: string, base64: string): string {
    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    const blob = new Blob([bytes], { type: "application/pdf" });
    const url = URL.createObjectURL(blob);

    // Revoke previous URL for this template if it exists
    const prev = this.cache.get(name);
    if (prev) URL.revokeObjectURL(prev);

    this.cache.set(name, url);
    return url;
  }

  has(name: string): boolean {
    return this.cache.has(name);
  }

  clear(): void {
    for (const url of this.cache.values()) {
      URL.revokeObjectURL(url);
    }
    this.cache.clear();
  }
}
