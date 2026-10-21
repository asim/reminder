/** One MP3 MediaSource timeline, filled from the existing verse URLs. */
export type AudioSegment = { index: number; start: number; end: number };

export class ContinuousAudio {
  static supported() {
    return typeof MediaSource !== 'undefined' && MediaSource.isTypeSupported('audio/mpeg');
  }

  private source = new MediaSource();
  private url = URL.createObjectURL(this.source);
  private buffer?: SourceBuffer;
  private abort = new AbortController();
  private busy = false;
  private next: number;
  private segments: AudioSegment[] = [];
  private stopped = false;
  private failed = false;

  constructor(
    private audio: HTMLAudioElement,
    private urls: string[],
    first: number,
    private position: (segment: AudioSegment, elapsed: number) => void,
    private failure: (message: string) => void,
  ) {
    this.next = first;
    this.source.addEventListener('sourceopen', this.open, { once: true });
    audio.addEventListener('timeupdate', this.tick);
    audio.addEventListener('error', this.fail);
    audio.src = this.url;
    // The caller invokes play synchronously from the user's gesture. Fetching
    // and appending never pauses, replaces src, or calls play at verse boundaries.
  }

  private open = () => {
    if (this.stopped) return;
    try {
      this.buffer = this.source.addSourceBuffer('audio/mpeg');
      this.buffer.mode = 'sequence';
      void this.fill();
    } catch { this.fail(); }
  };

  private tick = () => {
    const segment = this.current();
    if (segment) this.position(segment, Math.max(0, this.audio.currentTime - segment.start));
    void this.fill();
  };

  current() {
    const time = this.audio.currentTime;
    return this.segments.find(segment => time >= segment.start && time < segment.end) ||
      (this.audio.ended ? this.segments.at(-1) : undefined);
  }

  seek(index: number, elapsed = 0) {
    const segment = this.segments.find(item => item.index === index);
    if (!segment || !this.buffer || this.failed) return false;
    const time = segment.start + Math.min(Math.max(0, elapsed), segment.end - segment.start - 0.001);
    const ranges = this.buffer.buffered;
    if (!Array.from({ length: ranges.length }, (_, i) => time >= ranges.start(i) && time < ranges.end(i)).some(Boolean)) return false;
    this.audio.currentTime = time;
    this.position(segment, time - segment.start);
    void this.fill();
    return true;
  }

  private operation(run: () => void) {
    const buffer = this.buffer!;
    return new Promise<void>((resolve, reject) => {
      const cleanup = () => {
        buffer.removeEventListener('updateend', done);
        buffer.removeEventListener('error', error);
        this.abort.signal.removeEventListener('abort', error);
      };
      const done = () => { cleanup(); resolve(); };
      const error = () => { cleanup(); reject(new Error('Audio buffer operation failed')); };
      buffer.addEventListener('updateend', done, { once: true });
      buffer.addEventListener('error', error, { once: true });
      this.abort.signal.addEventListener('abort', error, { once: true });
      try { run(); } catch (reason) { cleanup(); reject(reason); }
    });
  }

  private async fill() {
    if (this.busy || this.stopped || this.failed || !this.buffer || this.source.readyState !== 'open') return;
    this.busy = true;
    try {
      // Keep about 90 seconds ahead, rather than downloading the entire chapter.
      while (!this.stopped && this.next < this.urls.length) {
        const end = this.segments.at(-1)?.end || 0;
        if (end - this.audio.currentTime >= 90) break;
        // Retain the current and previous segment so language/back controls work.
        const currentIndex = this.current()?.index ?? this.next;
        const keepFrom = this.segments.find(item => item.index === currentIndex - 1)?.start || 0;
        if (keepFrom > 0 && this.buffer.buffered.length && this.buffer.buffered.start(0) < keepFrom) {
          await this.operation(() => this.buffer!.remove(0, keepFrom));
        }
        const timer = setTimeout(() => this.abort.abort(), 30000);
        let bytes: ArrayBuffer;
        try {
          const response = await fetch(this.urls[this.next], { signal: this.abort.signal });
          if (!response.ok) throw new Error('Audio download failed');
          bytes = await response.arrayBuffer();
        } finally { clearTimeout(timer); }
        if (this.stopped) return;
        await this.operation(() => this.buffer!.appendBuffer(bytes));
        if (this.stopped) return;
        const ranges = this.buffer.buffered;
        const appendedEnd = ranges.length ? ranges.end(ranges.length - 1) : end;
        if (appendedEnd <= end) throw new Error('Audio segment has no playable frames');
        this.segments.push({ index: this.next++, start: end, end: appendedEnd });
        const active = this.current();
        if (active) this.position(active, this.audio.currentTime - active.start);
      }
      if (!this.stopped && this.next === this.urls.length) this.source.endOfStream();
    } catch { if (!this.stopped) this.fail(); }
    finally { this.busy = false; }
  }

  private fail = () => {
    if (this.failed || this.stopped) return;
    this.failed = true;
    this.abort.abort();
    this.audio.pause();
    this.failure('Continuous audio could not be loaded. Press play to retry, or remove continuousAudio=1 from the URL to use the standard player.');
  }

  get hasFailed() { return this.failed; }

  dispose() {
    this.stopped = true;
    this.abort.abort();
    this.audio.removeEventListener('timeupdate', this.tick);
    this.audio.removeEventListener('error', this.fail);
    this.source.removeEventListener('sourceopen', this.open);
    URL.revokeObjectURL(this.url);
  }
}
